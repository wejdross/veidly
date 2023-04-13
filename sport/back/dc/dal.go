package dc

import (
	"context"
	"database/sql"
	"fmt"
	"sport/helpers"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (ctx *Ctx) DalUpdateDcUsesTx(dcID uuid.UUID, c int, tx *sql.Tx) error {
	if c == 0 {
		return fmt.Errorf("invalid value")
	}
	w := ""
	if c > 0 {
		// if adding value cannot go above quantity
		w = " quantity >= (redeemed_quantity + $1)"
	} else {
		// if subtracting value cant go below 0
		w = " redeemed_quantity >= -1 * $1"
	}
	// if subtracting dc uses
	const q = `
		update dc set 
		redeemed_quantity = redeemed_quantity + $1 
		where id = $2 and `
	res, err := tx.Exec(q+w, c, dcID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalUpdateDcUses(dcID uuid.UUID, c int) error {
	tx, err := ctx.Dal.Db.Begin()
	if err != nil {
		return err
	}
	if err = ctx.DalUpdateDcUsesTx(dcID, c, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (ctx *Ctx) DalCreateDcBinding(trainingID, dcID, userID uuid.UUID) error {
	const q = `
		insert into dc_v_train (
			training_id, dc_id
		)
		select t.id, dc.id
		from instructors i
			inner join trainings t on t.instructor_id = i.id
			inner join dc on dc.instr_id = i.id
		where i.user_id = $3 and t.id = $1 and dc.id = $2`
	res, err := ctx.Dal.Db.Exec(q, trainingID, dcID, userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalDeleteDcBinding(trainingID, dcID, userID uuid.UUID) error {
	const q = `
		delete from dc_v_train 
		where 
			dc_id = $2
			and training_id = (
				select t.id 
				from trainings t, instructors i 
				where t.instructor_id = i.id and i.user_id = $3 and t.id = $1 
			)`
	res, err := ctx.Dal.Db.Exec(q, trainingID, dcID, userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalGetDcCount(iid uuid.UUID, tx *sql.Tx) (int, error) {
	const q = "select count(1) from dc where instructor_id = $1"
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(q, iid)
	} else {
		row = ctx.Dal.Db.QueryRow(q, iid)
	}
	var c int
	err := row.Scan(&c)
	return c, err
}

func (ctx *Ctx) DalCreateDc(tx *sql.Tx, dc *Dc) error {
	const q = `
		insert into dc (
				id, 
				instr_id, 
				redeemed_quantity, 
				quantity, 
				valid_start, 
				valid_end, 
				name,
				discount
		) values (
				$1, $2, $3, $4, $5, $6, $7, $8)`
	args := []interface{}{
		dc.ID,
		dc.InstrID,
		dc.RedeemedQuantity,
		dc.Quantity,
		time.Time(dc.ValidStart),
		time.Time(dc.ValidEnd),
		dc.Name,
		dc.Discount,
	}
	var err error
	if tx == nil {
		_, err = ctx.Dal.Db.Exec(q, args...)
	} else {
		_, err = tx.Exec(q, args...)
	}
	return err
}

func (ctx *Ctx) DalValidateAndCreateDc(dc *Dc) error {

	const maxRpt = 100
	for i := 0; i < maxRpt; i++ {

		tx, err := ctx.Dal.Db.BeginTx(context.Background(), &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			return err
		}

		c, err := ctx.DalGetDcCount(dc.InstrID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		if c >= ctx.Config.MaxCodesPerInstr {
			tx.Rollback()
			return helpers.NewHttpError(409, "", fmt.Errorf("too many objects"))
		}

		if err = ctx.DalCreateDc(tx, dc); err != nil {
			tx.Rollback()
			if helpers.PgIsConcurrentUpdate(err) {
				continue
			} else {
				return err
			}
		}

		if err = tx.Commit(); err != nil {
			tx.Rollback()
			if helpers.PgIsConcurrentUpdate(err) {
				continue
			} else {
				return err
			}
		}

		return nil
	}

	return helpers.NewHttpError(503, "", fmt.Errorf("too many requests"))
}

type ReadDcRequest struct {
	ID         *uuid.UUID
	InstrID    *uuid.UUID
	Name       string
	TrainingID *uuid.UUID
}

func (ctx *Ctx) DalReadDc(req ReadDcRequest, tx *sql.Tx) ([]Dc, error) {
	qb := strings.Builder{}
	args := []interface{}{}
	next := 1
	jq := ""

	if req.TrainingID != nil {
		jq = `inner join dc_v_train dvt on dvt.dc_id = dc.id and dvt.training_id = $1 `
		next++
		args = append(args, *req.TrainingID)
	}

	qb.WriteString(fmt.Sprintf(`
	select 
		%s
	from dc 
	%s 
	where 1 = 1 `, DcSelectColumns(), jq))

	if req.InstrID != nil {
		qb.WriteString(fmt.Sprintf(` and instr_id = $%d `, next))
		next++
		args = append(args, *req.InstrID)
	}

	if req.ID != nil {
		qb.WriteString(fmt.Sprintf(" and id = $%d ", next))
		next++
		args = append(args, *req.ID)
	}

	if req.Name != "" {
		qb.WriteString(fmt.Sprintf(" and name = $%d ", next))
		next++
		args = append(args, req.Name)
	}

	q := qb.String()
	var dbres *sql.Rows
	var err error
	if tx == nil {
		dbres, err = ctx.Dal.Db.Query(q, args...)
	} else {
		dbres, err = tx.Query(q, args...)
	}

	if err != nil {
		return nil, err
	}
	defer dbres.Close()

	res := make([]Dc, 0, 2)
	var tmp Dc
	scanArgs := tmp.ScanFields()

	for dbres.Next() {
		if err := dbres.Scan(scanArgs...); err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}

	return res, nil
}

func (ctx *Ctx) DalUpdateDc(dc *UpdateDcRequest, userID uuid.UUID) error {
	const q = `
	update dc set 
		quantity = $2, 
		valid_start = $3, 
		valid_end = $4, 
		name = $5,
		discount = $6
	where id = $1 and instr_id = (select id from instructors where user_id = $7)`
	res, err := ctx.Dal.Db.Exec(q,
		dc.ID,
		dc.Quantity,
		time.Time(dc.ValidStart),
		time.Time(dc.ValidEnd),
		dc.Name,
		dc.Discount,
		userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalDeleteDc(id, userID uuid.UUID) error {
	const q = `
		delete from dc 
		where id = $1 
			and instr_id = (select id from instructors where user_id = $2)`
	res, err := ctx.Dal.Db.Exec(q, id, userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}
