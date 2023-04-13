package instr

import (
	"context"
	"database/sql"
	"fmt"
	"sport/helpers"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// just read current instr penalty
func (ctx *Ctx) DalReadNextPenalty(id uuid.UUID, tx *sql.Tx) (int, error) {
	const q = "select queued_payout_cuts from instructors where id = $1"
	var rr *sql.Row
	if tx == nil {
		rr = ctx.Dal.Db.QueryRow(q, id)
	} else {
		rr = tx.QueryRow(q, id)
	}
	var cuts int
	err := rr.Scan(&cuts)
	if err != nil {
		return 0, err
	}
	if cuts < 0 {
		return 0, fmt.Errorf("DalGetNextPenalty: cuts was < 0")
	}
	return cuts, nil
}

// is thread safe
func (ctx *Ctx) DalGetNextPenaltySerial(id uuid.UUID, payoutValue, noMoreThan int) (int, error) {
	tx, err := ctx.Dal.Db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return 0, err
	}
	var p int
	for i := 0; i < 100; i++ {
		if p, err = ctx.DalGetNextPenalty(id, payoutValue, noMoreThan, tx); err != nil {
			tx.Rollback()
			if helpers.PgIsConcurrentUpdate(err) {
				continue
			}
			return 0, err
		}
		if err = tx.Commit(); err != nil {
			tx.Rollback()
			if helpers.PgIsConcurrentUpdate(err) {
				continue
			}
			return 0, err

		}
		return p, nil
	}
	return 0, fmt.Errorf("DalGetNextPenaltySerial: cannot serialize transaction -- too many requests")
}

// get and update penalty(decrement)
// not thread safe
func (ctx *Ctx) DalGetNextPenalty(id uuid.UUID, payoutValue, noMoreThan int, tx *sql.Tx) (int, error) {
	if payoutValue <= 0 || noMoreThan <= 0 {
		return 0, nil
	}
	if noMoreThan > 100 {
		return 0, fmt.Errorf("DalGetNextPenalty: invalid noMoreThan")
	}
	cuts, err := ctx.DalReadNextPenalty(id, tx)
	if err != nil {
		return 0, err
	}
	maxPayoutPenalty := (payoutValue * noMoreThan) / 100
	if maxPayoutPenalty <= 0 {
		return 0, nil
	}
	payoutPenalty := cuts
	if payoutPenalty > maxPayoutPenalty {
		payoutPenalty = maxPayoutPenalty
	}
	res, err := tx.Exec(`
		update instructors 
		set queued_payout_cuts = $1
		where id = $2`,
		cuts-payoutPenalty, id)
	if err != nil {
		return 0, err
	}
	if err := helpers.PgMustBeOneRow(res); err != nil {
		return 0, err
	}
	return payoutPenalty, nil
}

func (ctx *Ctx) DalUpdatePenalty(id uuid.UUID, amount int, tx *sql.Tx) error {
	const q = `update instructors set queued_payout_cuts = queued_payout_cuts + $2 where id = $1`
	var res sql.Result
	var err error
	if tx == nil {
		res, err = ctx.Dal.Db.Exec(q, id, amount)
	} else {
		res, err = tx.Exec(q, id, amount)
	}
	if err != nil {
		return err
	}
	if err := helpers.PgMustBeOneRow(res); err != nil {
		return err
	}
	return nil
}

// WARNING: race condition possible
func (ctx *Ctx) DalConditionalAddPenalty(id uuid.UUID, amount, freeShots int, tx *sql.Tx) error {
	if amount <= 0 {
		return nil
	}
	rr := tx.QueryRow("select refunds from instructors where id = $1", id)
	var s InstructorRefunds
	err := rr.Scan(&s)
	if err != nil {
		return err
	}
	if s == nil {
		s = make(InstructorRefunds)
	}
	t := time.Now().In(time.UTC)
	t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	tu := t.Unix()
	s[tu]++
	if s[tu] > freeShots {
		if err := ctx.DalUpdatePenalty(id, amount, tx); err != nil {
			return err
		}
	}
	res, err := tx.Exec(`
		update instructors set refunds = $1 where id = $2
	`, &s, id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

type ReadCardInfoKeyType int

const (
	ReadCardKey_InstructorID ReadCardInfoKeyType = 1
)

func (ctx *Ctx) DalReadCardInfo(
	instrID uuid.UUID,
) (*ProcessedCardInfo, error) {
	var q = `select 
			i.card_ref_id,
			i.card_brand,
			i.card_holder_name,
			i.card_summary
		from instructors i
		where i.id = $1`
	dbres := ctx.Dal.Db.QueryRow(q, instrID)
	res := new(ProcessedCardInfo)
	if err := dbres.Scan(
		&res.CardRefID,
		&res.CardBrand,
		&res.CardHolderName,
		&res.CardSummary,
	); err != nil {
		return nil, err
	}
	if res.CardRefID == "" {
		return res, sql.ErrNoRows
	}
	return res, nil
}

func (ctx *Ctx) DalUpdateCardInfo(
	userID uuid.UUID,
	payout *ProcessedCardInfo) (err error) {

	const q = `
		update instructors set 
		(card_ref_id, card_brand, card_holder_name, card_summary) 
		= 
		($1,$2,$3,$4) 
		where 
		user_id = $5
	`

	res, err := ctx.Dal.Db.Exec(
		q,
		payout.CardRefID,
		payout.CardBrand,
		payout.CardHolderName,
		payout.CardSummary,
		userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalDeletePayoutInfo(instructorID uuid.UUID) error {
	return ctx.DalUpdateCardInfo(instructorID, &ProcessedCardInfo{
		CardRefID:      "",
		CardBrand:      "",
		CardHolderName: "",
		CardSummary:    "",
	})
}

func (ctx *Ctx) DalCreateInstructor(instructor *Instructor) error {

	tx, err := ctx.Dal.Db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`insert into instructors (
		id,
		user_id,
		created_on,
		tags,
		year_exp,
		known_langs,
		bg_img_path,
		extra_img_paths,
		profile_sections,
		disabled,
		card_ref_id,
		card_brand,
		card_holder_name,
		card_summary,
		refunds,
		queued_payout_cuts,
		invoice_lines
	) values (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17 )`,
		instructor.ID,
		instructor.UserID,
		instructor.CreatedOn,
		pq.Array(instructor.Tags),
		instructor.YearExp,
		pq.Array(instructor.KnownLangs),
		instructor.BgImgPath,
		pq.Array(instructor.ExtraImgPaths),
		&instructor.ProfileSections,
		instructor.Disabled,
		instructor.CardInfo.CardRefID,
		instructor.CardInfo.CardBrand,
		instructor.CardInfo.CardHolderName,
		instructor.CardInfo.CardSummary,
		&instructor.Refunds,
		instructor.QueuedPayoutCuts,
		pq.StringArray(instructor.InvoiceLines),
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ctx *Ctx) DalReadInstructorID(userID uuid.UUID) (uuid.UUID, error) {
	r, err := ctx.Dal.Db.Query(`select 
			i.id
		from users u
			inner join instructors i on u.id = i.user_id
		where u.id = $1`, userID)
	if err != nil {
		return uuid.Nil, err
	}
	defer r.Close()
	if !r.Next() {
		return uuid.Nil, helpers.NewElementNotFoundErr(userID)
	}
	var ret uuid.UUID
	if err = r.Scan(&ret); err != nil {
		return uuid.Nil, err
	}
	return ret, err
}

type DalReadInstructorsRequest struct {
	IDs     []string
	UserIDs []string
}

func (ctx *Ctx) DalReadInstructors(req DalReadInstructorsRequest) ([]InstructorWithUser, error) {

	qb := strings.Builder{}
	qb.WriteString(fmt.Sprintf(`select %s from users u
		inner join instructors i on u.id = i.user_id `,
		InstructorWithUserCols()))

	args := make([]interface{}, 0, 1)

	if req.IDs != nil {
		next := len(args) + 1
		qb.WriteString(" where ")
		// qb.WriteString(fmt.Sprintf(" where i.id = any($%d) ", next))
		// args = append(args, pq.StringArray(req.IDs))
		helpers.PgUuidMatch(&qb, req.IDs, next, &args, "i.id")
	}

	if req.UserIDs != nil {
		next := len(args) + 1
		if next == 1 {
			qb.WriteString(" where ")
		} else {
			qb.WriteString(" and ")
		}
		qb.WriteString(fmt.Sprintf(" i.user_id = any($%d) ", next))
		args = append(args, pq.StringArray(req.UserIDs))
	}

	r, err := ctx.Dal.Db.Query(qb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var tmp InstructorWithUser
	var ret = make([]InstructorWithUser, 0, 4)
	for r.Next() {
		if err = r.Scan(tmp.ScanFields()...); err != nil {
			return nil, err
		}
		tmp.PubInstructorInfo.PostprocessAfterDbScan(ctx)
		tmp.UserInfo.PostprocessAfterDbScan(ctx.User)
		ret = append(ret, tmp)
	}

	return ret, err
}

type ReadInstructorInfoKeyType int

const (
	InstructorID ReadInstructorInfoKeyType = 1
	UserID       ReadInstructorInfoKeyType = 2
)

func (ctx *Ctx) DalReadInstructor(
	key interface{}, keyType ReadInstructorInfoKeyType,
) (*InstructorWithUser, error) {

	q := fmt.Sprintf(`select 
			%s
		from users u
			inner join instructors i on u.id = i.user_id `,
		InstructorWithUserCols())

	switch keyType {
	case InstructorID:
		q += " where i.id = $1"
	case UserID:
		q += " where u.id = $1"
	default:
		return nil, fmt.Errorf("invalid key type: %d", keyType)
	}

	r, err := ctx.Dal.Db.Query(q, key)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	if !r.Next() {
		return nil, helpers.NewElementNotFoundErr(key)
	}
	var ret InstructorWithUser
	if err = r.Scan(ret.ScanFields()...); err != nil {
		return nil, err
	}

	ret.PubInstructorInfo.PostprocessAfterDbScan(ctx)
	ret.UserInfo.PostprocessAfterDbScan(ctx.User)

	return &ret, err
}

func (ctx *Ctx) DalReadInstructorInfo(
	key interface{}, keyType ReadInstructorInfoKeyType,
) (*PubInstructorWithUser, error) {
	i, err := ctx.DalReadInstructor(key, keyType)
	if err != nil {
		return nil, err
	}
	return i.ToInfo(), nil
}

func (ctx *Ctx) DalPatchInstructor(
	userID uuid.UUID, instructor *InstructorRequest) error {

	tx, err := ctx.Dal.Db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		`update instructors set 
			tags = $1, 
			year_exp = $2, 
			disabled = $3,
			known_langs = $4,
			profile_sections = $5,
			invoice_lines = $6
		where user_id = $7`,
		pq.Array(instructor.Tags),
		instructor.YearExp,
		instructor.Disabled,
		pq.Array(instructor.KnownLangs),
		&instructor.ProfileSections,
		pq.Array(instructor.InvoiceLines),
		userID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := helpers.PgMustBeOneRow(res); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ctx *Ctx) DalCanDeleteInstructor(userID uuid.UUID) (bool, error) {
	instructorID, err := ctx.DalReadInstructorID(userID)
	if err != nil {
		return false, err
	}

	row := ctx.Dal.Db.QueryRow(`
	select count(1) 
	from reservations 
	where instructor_id = $1 and is_active = true`, instructorID)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	if count != 0 {
		return false, nil
	}

	row = ctx.Dal.Db.QueryRow(`
	select count(1) 
	from subs 
	where instructor_id = $1 and is_active = true`, instructorID)
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	if count != 0 {
		return false, nil
	}

	return true, nil
}

func (ctx *Ctx) DalDeleteInstructor(userID uuid.UUID, tx *sql.Tx) error {

	var err error
	rollbackTx := false

	if tx == nil {
		tx, err = ctx.Dal.Db.Begin()
		if err != nil {
			return err
		}
		rollbackTx = true
	}

	res, err := tx.Exec("delete from instructors where user_id = $1", userID)
	if err != nil {
		if rollbackTx {
			tx.Rollback()
		}
		return err
	}
	if err := helpers.PgMustBeOneRow(res); err != nil {
		if rollbackTx {
			tx.Rollback()
		}
		return err
	}

	if rollbackTx {
		return tx.Commit()
	}

	return nil
}

func (ctx *Ctx) DalInsertDeletedInstructor(di *DeletedInstructor, tx *sql.Tx) error {
	const q = `insert into deleted_instructors (
		email, refunds, queued_payout_cuts
	) values (
		$1, $2, $3
	)`
	var err error
	args := []interface{}{di.Email, &di.Refunds, di.QueuedPayoutCuts}
	if tx == nil {
		_, err = ctx.Dal.Db.Exec(q, args...)
	} else {
		_, err = tx.Exec(q, args...)
	}
	return err
}

func (ctx *Ctx) DalGetDeletedInstructor(email string) (*DeletedInstructor, error) {
	const q = `select 
			email, refunds, queued_payout_cuts
		from deleted_instructors where email = $1`
	res := new(DeletedInstructor)
	dbres := ctx.Dal.Db.QueryRow(q, email)
	return res, dbres.Scan(&res.Email, &res.Refunds, &res.QueuedPayoutCuts)
}
