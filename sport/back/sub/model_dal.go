package sub

import (
	"database/sql"
	"fmt"
	"sport/helpers"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (ctx *Ctx) DalCreateSubModel(sm *SubModel) error {
	const q = `insert into sub_models (
		id,
		instructor_id,
		instr_user_id,
		name,
		max_entrances,
		duration,
		price,
		processing_fee,
		payout_value,
		refund_value,
		currency,
		max_active,
		is_free_entrance,
		all_trainings_by_def
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
	)`
	_, err := ctx.Dal.Db.Exec(q,
		sm.ID,
		sm.InstructorID,
		sm.InstrUserID,
		sm.Name,
		sm.MaxEntrances,
		sm.Duration,
		sm.PricingInfo.Price,
		sm.PricingInfo.ProcessingFee,
		sm.PricingInfo.PayoutValue,
		sm.PricingInfo.RefundValue,
		sm.Currency,
		sm.MaxActive,
		sm.IsFreeEntrance,
		sm.AllTrainingsByDef)
	return err
}

type ReadSubModelRequest struct {
	InstrUserID  *uuid.UUID
	ID           *uuid.UUID
	InstructorID *uuid.UUID
	IDs          []string
}

func (ctx *Ctx) DalReadSubModels(r ReadSubModelRequest) ([]SubModel, error) {
	const q = `select %s from sub_models sm`

	qb := strings.Builder{}
	qb.WriteString(fmt.Sprintf(q, SmSelectColumns()))
	args := make([]interface{}, 0, 2)

	if r.InstrUserID != nil {
		qb.WriteString(" where sm.instr_user_id = $1 ")
		args = append(args, *r.InstrUserID)
	}

	if r.ID != nil {
		next := len(args) + 1
		if next == 1 {
			qb.WriteString(" where ")
		} else {
			qb.WriteString(" and ")
		}
		qb.WriteString(fmt.Sprintf(" sm.id = $%d ", next))
		args = append(args, *r.ID)
	}

	if r.IDs != nil {
		next := len(args) + 1
		if next == 1 {
			qb.WriteString(" where ")
		} else {
			qb.WriteString(" and ")
		}
		qb.WriteString(fmt.Sprintf(" sm.id = any($%d) ", next))
		//qb.WriteString(fmt.Sprintf(" sm.id @> $%d ", next))
		args = append(args, pq.StringArray(r.IDs))
	}

	if r.InstructorID != nil {
		next := len(args) + 1
		if next == 1 {
			qb.WriteString(" where ")
		} else {
			qb.WriteString(" and ")
		}
		qb.WriteString(fmt.Sprintf(" instructor_id = $%d ", next))
		args = append(args, *r.InstructorID)
	}

	dbres, err := ctx.Dal.Db.Query(qb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer dbres.Close()
	res := make([]SubModel, 0, 4)
	var tmp SubModel
	var sf = tmp.ScanFields()
	for dbres.Next() {
		if err := dbres.Scan(sf...); err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (ctx *Ctx) DalReadSingleSubModel(r ReadSubModelRequest) (*SubModel, error) {
	res, err := ctx.DalReadSubModels(r)
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, sql.ErrNoRows
	}
	return &res[0], nil
}

func (ctx *Ctx) DalUpdateSubModel(
	sm *UpdateSubModelRequest, userID uuid.UUID,
) error {
	pi := ctx.NewSubPricing(sm.Price)
	const q = `update sub_models set
		max_entrances = $1,
		duration = $2,
		price = $3,
		processing_fee = $4,
		payout_value = $5,
		refund_value = $6,
		currency = $7,
		max_active = $8,
		is_free_entrance = $9,
		name = $10,
		all_trainings_by_def = $11
		where id = $12 and instr_user_id = $13`
	res, err := ctx.Dal.Db.Exec(q,
		sm.MaxEntrances,
		sm.Duration,
		pi.Price,
		pi.ProcessingFee,
		pi.PayoutValue,
		pi.RefundValue,
		sm.Currency,
		sm.MaxActive,
		sm.IsFreeEntrance,
		sm.Name,
		sm.AllTrainingsByDef,
		sm.ID,
		userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalDeleteSubModel(smID uuid.UUID, userID uuid.UUID) error {
	const q = "delete from sub_models where id = $1 and instr_user_id = $2"
	res, err := ctx.Dal.Db.Exec(q, smID, userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}
