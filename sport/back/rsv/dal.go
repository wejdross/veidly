package rsv

import (
	"database/sql"
	"fmt"
	"sport/adyen_sm"
	"sport/dal"
	"sport/dc"
	"sport/helpers"
	"sport/instr"
	"sport/train"
	"sport/user"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func UpdateReservationInstrDecision(
	d *dal.Ctx,
	id uuid.UUID,
	decision adyen_sm.InstructorDecision,
	instructorID uuid.UUID,
) error {
	res, err := d.Db.Exec(`
		update reservations set
			instructor_decision = $1
		where id = $2 and instructor_id = $3
	`, decision, id, instructorID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) UpdateReservationState(
	id uuid.UUID,
	previousState, nextState adyen_sm.State,
	timeout *time.Time,
	eargs map[string]interface{},
	tx *sql.Tx,
	// source StateChangeEventSource,
	// sourceData string,
	//force bool,
) error {
	var err error
	// if !force {
	// 	// ok, err := validateStateTransition(previousState, nextState)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	// if !ok {
	// 	// 	return nil
	// 	// }
	// 	if previousState == nextState {
	// 		return nil
	// 	}
	// }

	q := strings.Builder{}
	q.WriteString(`update reservations set 
	state = $1, sm_retries = 0`) // always set retries to 0 on successful state change
	args := []interface{}{nextState}
	next := 2
	if timeout != nil {
		q.WriteString(fmt.Sprintf(`,sm_timeout = $%d`, next))
		next++
		args = append(args, timeout)
	}
	for k := range eargs {
		q.WriteString(fmt.Sprintf(" , %s = $%d", k, next))
		args = append(args, eargs[k])
		next++
	}
	q.WriteString(fmt.Sprintf(` where id = $%d and state = $%d`, next, next+1))
	args = append(args, id, previousState)
	next += 2
	var res sql.Result
	if tx == nil {
		res, err = ctx.Dal.Db.Exec(q.String(), args...)
	} else {
		res, err = tx.Exec(q.String(), args...)
		// if tx is not nil then we wont log here.
		// its up to caller to log event
	}

	// var logReq *StateChangeEvent
	// x := StateChangeEvent{
	// 	ID:            uuid.New(),
	// 	ReservationID: id,
	// 	PreviousState: previousState,
	// 	NextState:     nextState,
	// 	Timestamp:     time.Now().In(time.UTC),
	// 	Success:       true,
	// 	Source:        source,
	// 	SourceData:    sourceData,
	// }
	// logReq = &x

	if err != nil {
		// logReq.Success = false
		// logReq.Error = err.Error()
		// LogStateChange(logReq, tx)
		return err
	}

	err = helpers.PgMustBeOneRow(res)

	// if err != nil {
	// 	logReq.Success = false
	// 	logReq.Error = err.Error()
	// }

	//LogStateChange(logReq, tx)

	return err
}

// There is no validation - internal usage only!
func (ctx *Ctx) DeleteReservation(id uuid.UUID) error {
	q := `delete from reservations where id = $1`
	res, err := ctx.Dal.Db.Exec(q, id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

// true if user can still register to training
// this approach can result in race condition: this is left to consideration
func (ctx *Ctx) ValidateCapacity(
	training *train.Training,
	occurrence time.Time,
	tx *sql.Tx) (bool, error) {

	q := `select count(*)
		from reservations 
		where 
			instructor_id = $3 and 
			training_id = $1 and 
			date_start = $2 and 
			is_confirmed = true`
	var err error
	args := []interface{}{training.ID, occurrence, training.InstructorID}
	var res *sql.Row
	if tx == nil {
		res = ctx.Dal.Db.QueryRow(q, args...)
	} else {
		res = tx.QueryRow(q, args...)
	}
	var count int
	err = res.Scan(&count)
	if count < training.Capacity {
		return true, err
	}
	return false, err
}

func (ctx *Ctx) ValidateInstructorIsFree(
	training *train.Training,
	occurrence time.Time,
	tx *sql.Tx,
) (bool, error) {
	q := `select 1 
		from reservations
		where 
			training_id != $1
			and instructor_id = $2
			and date_start >= $3 
			and date_end < $3
			and is_confirmed = true`
	args := []interface{}{training.ID, training.InstructorID, occurrence}
	var res *sql.Rows
	var err error
	if tx == nil {
		res, err = ctx.Dal.Db.Query(q, args...)
	} else {
		res, err = tx.Query(q, args...)
	}
	if err != nil {
		return false, err
	}
	err = res.Close()
	return !res.Next(), err
}

type ReadDDLRsvResponse struct {
	Rsv        []DDLRsvWithInstr
	Pagination helpers.PaginationResponse
}

type RsvWithInstrPagination struct {
	Rsv        []RsvWithInstr
	Pagination helpers.PaginationResponse
}

func (r *ReadDDLRsvResponse) ToRsvWithInstrPagination() *RsvWithInstrPagination {
	ret := RsvWithInstrPagination{}
	ret.Pagination = r.Pagination
	ret.Rsv = make([]RsvWithInstr, len(r.Rsv))
	for i := range r.Rsv {
		ret.Rsv[i] = RsvWithInstr{
			Rsv: r.Rsv[i].Rsv,
		}
		if r.Rsv[i].Instructor != nil {
			ret.Rsv[i].Instructor = new(instr.PubInstructorWithUser)
			ret.Rsv[i].Instructor.PubInstructorInfo = r.Rsv[i].Instructor.PubInstructorInfo
			ret.Rsv[i].Instructor.UserInfo = r.Rsv[i].Instructor.UserInfo
		}
	}
	return &ret
}

type ReadShortenedRsvsArgs struct {
	DateRange helpers.DateRange
	IDs       []string
}

func (ctx *Ctx) ReadShortenedRsvs(req ReadShortenedRsvsArgs) ([]ShortenedRsv, error) {
	qb := strings.Builder{}
	qb.WriteString(`
		select
			r.id,
			r.training_id,
			r.date_start,
			r.date_end,
			r.groups
		from reservations r
		where is_confirmed = true 
	`)

	args := make([]interface{}, 0, 3)

	if req.IDs != nil {
		next := len(args) + 1
		qb.WriteString((fmt.Sprintf(" and r.id = any($%d) ", next)))
		next++
		args = append(args, pq.StringArray(req.IDs))
	}

	if req.DateRange.IsNotZero() {
		next := len(args) + 1
		qb.WriteString(fmt.Sprintf(` and r.date_start <= $%d and r.date_end >= $%d `, next, next+1))
		args = append(
			args, req.DateRange.End.In(time.UTC), req.DateRange.Start.In(time.UTC))
	}

	r, err := ctx.Dal.Db.Query(qb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	res := make([]ShortenedRsv, 0, 10)
	tmp := ShortenedRsv{}

	for r.Next() {
		if err := r.Scan(
			&tmp.ID,
			&tmp.TrainingID,
			&tmp.DateStart,
			&tmp.DateEnd,
			&tmp.Groups,
		); err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}

	return res, nil
}

type ReadRsvsArgs struct {
	Pagination     *helpers.PaginationRequest
	DateRange      helpers.DateRange
	WithInstructor bool
	TrainingID     *uuid.UUID
	IsActive       *bool
	IsConfirmed    *bool
	InstructorID   *uuid.UUID
	//
	// RefID       *string
	UserID      *uuid.UUID
	ID          *uuid.UUID
	AccessToken *uuid.UUID
	IDs         []string
	InstrUserID *uuid.UUID
}

var EmptyRsvArgs = ReadRsvsArgs{}

func (ctx *Ctx) ReadDDLRsvs(
	args ReadRsvsArgs,
) (*ReadDDLRsvResponse, error) {
	whereFragmentBuilder := strings.Builder{}
	paginationFragment := ""
	var selectFragmentBuilder strings.Builder
	selectFragmentBuilder.WriteString(`
		select 
			r.id,
			r.instructor_id,

			r.training_id,
			r.training,
			r.occ,
			r.groups,

			r.user_id,
			r.user_data,

			r.date_start,
			r.date_end,

			r.is_confirmed,
			r.instructor_decision,
			r.state,
			r.sm_timeout,
			r.created_on,

			r.link_id,
			r.link_url,

			r.sm_retries,
			r.is_active,

			r.order_id,

			r.sm_cache,

			r.processing_fee,
			r.split_income_fee,
			r.split_payout,
			r.refund_amount,

			r.access_token,
			r.qr_confirmed,

			r.user_contact_data,
			r.dc,
			r.dc_rollback,
			
			r.use_user_acc,

			uu.id,
			coalesce(uu.user_data, '{}'::jsonb),
			coalesce(uu.avatar_relpath, ''),
			coalesce(uu.contact_data, '{}'::jsonb)
		`)

	if args.WithInstructor {
		selectFragmentBuilder.WriteString(`,`)
		selectFragmentBuilder.WriteString(instr.InstructorWithUserCols())
	}

	selectFragmentBuilder.WriteString(` from reservations r 
		left join users uu on r.use_user_acc = true and uu.id = r.user_id `)

	if args.WithInstructor {
		selectFragmentBuilder.WriteString(`
			inner join instructors i on i.id = r.instructor_id
			inner join users u on u.id = i.user_id`)
	}

	selectFragmentBuilder.WriteString(` where 1=1 `)

	whereArgs := make([]interface{}, 0, 4)
	pageArgs := make([]interface{}, 0, 4)
	nextParamIndex := 1

	// if args.RefID != nil {
	// 	whereFragmentBuilder.WriteString(fmt.Sprintf(` and r.ref_id = $%d `, nextParamIndex))
	// 	nextParamIndex++
	// 	whereArgs = append(whereArgs, *args.RefID)
	// }

	if args.UserID != nil {
		whereFragmentBuilder.WriteString(fmt.Sprintf(` and r.user_id = $%d `, nextParamIndex))
		nextParamIndex++
		whereArgs = append(whereArgs, *args.UserID)
	}

	if args.ID != nil {
		whereFragmentBuilder.WriteString(fmt.Sprintf(` and r.id = $%d `, nextParamIndex))
		nextParamIndex++
		whereArgs = append(whereArgs, *args.ID)
	}

	if args.AccessToken != nil {
		whereFragmentBuilder.WriteString(fmt.Sprintf(` and r.access_token = $%d `, nextParamIndex))
		nextParamIndex++
		whereArgs = append(whereArgs, *args.AccessToken)
	}

	if args.IDs != nil {
		whereFragmentBuilder.WriteString(" and ")
		// whereFragmentBuilder.WriteString((fmt.Sprintf(" and r.id = any($%d) ", nextParamIndex)))
		// nextParamIndex++
		// whereArgs = append(whereArgs, pq.StringArray(args.IDs))
		helpers.PgUuidMatch(&whereFragmentBuilder, args.IDs, nextParamIndex, &whereArgs, "r.id")
	}

	if args.InstrUserID != nil {
		if !args.WithInstructor {
			return nil, fmt.Errorf("DalReadReservations: cant have instrUserID without WithInstructor")
		}
		whereFragmentBuilder.WriteString(fmt.Sprintf(` and u.id = $%d `, nextParamIndex))
		nextParamIndex++
		whereArgs = append(whereArgs, *args.InstrUserID)
	}

	if args.TrainingID != nil {
		whereFragmentBuilder.WriteString(fmt.Sprintf(` and r.training_id = $%d `, nextParamIndex))
		nextParamIndex++
		whereArgs = append(
			whereArgs,
			*args.TrainingID)
	}

	if args.IsActive != nil {
		whereFragmentBuilder.WriteString(fmt.Sprintf(" and r.is_active = $%d ", nextParamIndex))
		nextParamIndex++
		whereArgs = append(
			whereArgs,
			*args.IsActive)
	}

	if args.IsConfirmed != nil {
		whereFragmentBuilder.WriteString(fmt.Sprintf(" and r.is_confirmed = $%d ", nextParamIndex))
		nextParamIndex++
		whereArgs = append(
			whereArgs,
			*args.IsConfirmed)
	}

	if args.InstructorID != nil {
		whereFragmentBuilder.WriteString(fmt.Sprintf(" and r.instructor_id = $%d ", nextParamIndex))
		nextParamIndex++
		whereArgs = append(
			whereArgs,
			*args.InstructorID)
	}

	//filter
	if args.DateRange.IsNotZero() {
		whereFragmentBuilder.WriteString(fmt.Sprintf(` and r.date_start <= $%d `, nextParamIndex))
		nextParamIndex++
		whereFragmentBuilder.WriteString(fmt.Sprintf(` and r.date_end >= $%d `, nextParamIndex))
		whereArgs = append(
			whereArgs,
			args.DateRange.End.In(time.UTC),
			args.DateRange.Start.In(time.UTC))
		nextParamIndex++
	}

	// pagination
	if args.Pagination != nil {
		paginationFragment = args.Pagination.ToSqlQuery(&nextParamIndex)
		pageArgs = args.Pagination.ToSqlQueryParams()
	}

	whereFragment := whereFragmentBuilder.String()

	dbres, err := ctx.Dal.Db.Query(
		selectFragmentBuilder.String()+whereFragment+paginationFragment,
		append(whereArgs, pageArgs...)...)
	if err != nil {
		return nil, err
	}

	defer dbres.Close()

	var tmp DDLRsvWithInstr
	var tmpInstr instr.InstructorWithUser
	rsv := make([]DDLRsvWithInstr, 0, 10)
	var uid *uuid.UUID
	var userInfo user.PubUserInfo
	var userContactData user.ContactData
	var d dc.Dc
	scanArgs := []interface{}{
		&tmp.ID,
		&tmp.InstructorID,

		&tmp.TrainingID,
		&tmp.Training,
		&tmp.Occ,
		&tmp.Groups,

		&tmp.UserID,
		&tmp.UserInfo.UserData,

		&tmp.DateStart,
		&tmp.DateEnd,

		&tmp.IsConfirmed,
		&tmp.InstructorDecision,
		&tmp.State,
		&tmp.SmTimeout,
		&tmp.CreatedOn,

		&tmp.LinkID,
		&tmp.LinkUrl,

		&tmp.SmRetries,
		&tmp.IsActive,

		//

		&tmp.OrderID,

		&tmp.SmCache,

		&tmp.ProcessingFee,
		&tmp.SplitIncomeFee,
		&tmp.SplitPayout,
		&tmp.RefundAmount,

		&tmp.AccessToken,
		&tmp.QrConfirmed,

		&tmp.UserContactData,
		&d,
		&tmp.DcRollback,
		&tmp.UseUserAcc,

		&uid,
		&userInfo.UserData,
		&userInfo.AvatarRelpath,
		&userContactData,
	}
	if args.WithInstructor {
		scanArgs = append(scanArgs, tmpInstr.ScanFields()...)
	}
	for dbres.Next() {
		if err := dbres.Scan(scanArgs...); err != nil {
			return nil, err
		}
		if args.WithInstructor {
			tmp.Instructor = new(instr.InstructorWithUser)
			tmpInstr.PubInstructorInfo.PostprocessAfterDbScan(ctx.Instr)
			tmpInstr.UserInfo.PostprocessAfterDbScan(ctx.User)
			*tmp.Instructor = tmpInstr
		}

		// if we managed to join user table then replace user data
		// with joined ones
		if uid != nil && tmp.UseUserAcc {
			userInfo.PostprocessAfterDbScan(ctx.User)
			tmp.UserInfo = userInfo
			tmp.UserContactData = userContactData
		}
		if d.ID != uuid.Nil {
			tmp.Dc = new(dc.Dc)
			*tmp.Dc = d
		} else {
			tmp.Dc = nil
		}
		rsv = append(rsv, tmp)
	}

	var res ReadDDLRsvResponse
	res.Rsv = rsv

	// response pagination info
	if args.Pagination != nil {
		row := ctx.Dal.Db.QueryRow(
			"select count(1) from reservations r where 1=1 "+whereFragment, whereArgs...)
		var count int
		if err := row.Scan(&count); err != nil {
			return nil, err
		}
		res.Pagination.PaginationRequest = *args.Pagination
		if res.Pagination.Size != 0 && count > 0 {
			res.Pagination.NPages = (count / res.Pagination.Size) + 1
		}
	}

	return &res, nil
}

func (ctx *Ctx) ReadRsvWithInstr(args ReadRsvsArgs) (*RsvWithInstrPagination, error) {
	res, err := ctx.ReadDDLRsvs(args)
	if err != nil {
		return nil, err
	}
	return res.ToRsvWithInstrPagination(), nil
}

func (ctx *Ctx) ReadSingleRsvInfo(args ReadRsvsArgs) (*RsvWithInstr, error) {
	res, err := ctx.ReadRsvWithInstr(args)
	if err != nil {
		return nil, err
	}
	r := res.Rsv
	if len(r) != 1 {
		return nil, sql.ErrNoRows
	}
	return &r[0], nil
}

// note that this will always merge instructor.
// if you dont want to do it, then use ReadSingleRsv
func (ctx *Ctx) ReadRsvByID(
	id uuid.UUID,
) (*DDLRsvWithInstr, error) {
	return ctx.ReadSingleRsv(ReadRsvsArgs{
		ID:             &id,
		WithInstructor: true,
	})
}

func (ctx *Ctx) ReadSingleRsv(
	args ReadRsvsArgs,
) (*DDLRsvWithInstr, error) {
	res, err := ctx.ReadDDLRsvs(args)
	if err != nil {
		return nil, err
	}
	r := res.Rsv
	if len(r) != 1 {
		return nil, sql.ErrNoRows
	}
	return &r[0], nil
}

func (ctx *Ctx) CreateRsv(r *DDLRsv, tx *sql.Tx) error {

	const q = `insert into reservations (
		id,
		instructor_id,

		training_id,
		training,
		occ,
		groups,

		user_id,
		user_data,
		user_contact_data,

		date_start,
		date_end,

		is_confirmed,
		instructor_decision,
		state,
		sm_timeout,
		created_on,

		order_id,

		link_id,
		link_url,

		sm_retries,
		sm_cache,

		processing_fee,
		split_income_fee,
		split_payout,
		refund_amount,

		is_active,
		access_token,
		qr_confirmed,

		dc,
		dc_rollback
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
		$17,
		$18,
		$19,
		$20,
		$21,
		$22,
		$23,
		$24,
		$25,
		$26,
		$27,
		$28,
		$29,
		$30 )`
	args := []interface{}{
		r.ID,
		r.InstructorID,

		r.TrainingID,
		r.Training,
		&r.Occ,
		&r.Groups,

		r.UserID,
		&r.UserInfo.UserData,
		&r.UserContactData,

		r.DateStart,
		r.DateEnd,

		r.IsConfirmed,
		r.InstructorDecision,
		r.State,
		r.SmTimeout,
		r.CreatedOn,

		r.OrderID,
		r.LinkID,
		r.LinkUrl,

		r.SmRetries,
		&r.SmCache,

		r.ProcessingFee,
		r.SplitIncomeFee,
		r.SplitPayout,
		r.RefundAmount,

		r.IsActive,
		r.AccessToken,
		r.QrConfirmed,

		r.Dc,
		r.DcRollback,
	}
	//var res sql.Result
	var err error
	if tx != nil {
		_, err = tx.Exec(q, args...)
	} else {
		_, err = ctx.Dal.Db.Exec(q, args...)
	}
	return err
	//return helpers.SqlEnsureOneRowWasAffected(res)
}

func (ctx *Ctx) UpdateRsvSmTimeout(id uuid.UUID, timeout time.Time) error {
	res, err := ctx.Dal.Db.Exec(
		`update reservations set sm_timeout = $1 where id = $2`,
		timeout, id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) UpdateRsvSmDateStart(id uuid.UUID, start, end time.Time) error {
	res, err := ctx.Dal.Db.Exec(
		`update reservations set date_start = $1, date_end = $2 where id = $3`,
		start, end, id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalConfirmQrCode(rsvID uuid.UUID, instrID uuid.UUID) error {
	const q = `
		update reservations 
		set qr_confirmed = true 
		where id = $1 
			and qr_confirmed = false 
			and is_confirmed = true
			and instructor_id = $2`
	res, err := ctx.Dal.Db.Exec(q, rsvID, instrID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}
