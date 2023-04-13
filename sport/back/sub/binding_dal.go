package sub

import (
	"sport/helpers"

	"github.com/google/uuid"
)

func (ctx *Ctx) DalCreateSubModelBinding(smb *SubModelBinding, userID uuid.UUID) error {
	const q = `insert into sub_model_bindings (
		sub_model_id,
		training_id
	)
	select sm.id, t.id
	from sub_models sm
	inner join trainings t on t.instructor_id = sm.instructor_id
	where sm.instr_user_id = $3 and t.id = $2 and sm.id = $1`
	res, err := ctx.Dal.Db.Exec(q, smb.SubModelID, smb.TrainingID, userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalDeleteSubModelBinding(smb *SubModelBinding, userID uuid.UUID) error {
	const q = `delete from sub_model_bindings 
		where sub_model_id = $1 and training_id = $2 and sub_model_id = (
			select sm.id from sub_models sm where sm.instr_user_id = $3 and sm.id = $1
		)`
	res, err := ctx.Dal.Db.Exec(q, smb.SubModelID, smb.TrainingID, userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}
