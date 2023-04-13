package instr

import (
	"database/sql"
	"sport/helpers"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (ctx *Ctx) DalCreateVacation(v *Vacation) error {
	q := `insert into instr_vacations (
		id,
		instructor_id,
		date_start,
		date_end
	) values (
		$1,
		$2,
		$3,
		$4
	)`
	res, err := ctx.Dal.Db.Exec(
		q, v.ID, v.InstructorID, v.DateStart, v.DateEnd)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalUpdateVacation(iid uuid.UUID, vr *UpdateVacationRequest) error {
	q := `
		update instr_vacations
			set date_start = $1,
				date_end = $2
		where id = $3 and instructor_id = $4
	`
	res, err := ctx.Dal.Db.Exec(q, vr.DateStart, vr.DateEnd, vr.ID, iid)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalDeleteVacation(iid uuid.UUID, dr *DeleteVacationRequest) error {
	q := `delete from instr_vacations where id = $1 and instructor_id = $2`
	res, err := ctx.Dal.Db.Exec(q, dr.ID, iid)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalReadVacations(iid uuid.UUID, laterThan time.Time) ([]VacationInfo, error) {
	q := `select 
		id,
		date_start,
		date_end
	from instr_vacations
	where instructor_id = $1 
		and date_end >= $2
	order by date_start desc`
	dbres, err := ctx.Dal.Db.Query(q, iid, laterThan)
	if err != nil {
		return nil, err
	}
	defer dbres.Close()
	res := make([]VacationInfo, 0, 5)
	var tmp VacationInfo
	for dbres.Next() {
		if err = dbres.Scan(&tmp.ID, &tmp.DateStart, &tmp.DateEnd); err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (ctx *Ctx) DalIsInstructorOnVacation(iid uuid.UUID, start, end time.Time) (bool, error) {
	q := `select 1 
		from instr_vacations
		where instructor_id = $1
			and date_end > $2 and date_start < $3 
		limit 1`
	dbres := ctx.Dal.Db.QueryRow(q, iid, start, end)
	var foo int
	if err := dbres.Scan(&foo); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type VacationsGroupedByInstr map[uuid.UUID][]VacationInfo

func (ctx *Ctx) DalCreateVacationIndex(iids []string) (VacationsGroupedByInstr, error) {
	if len(iids) == 0 {
		return VacationsGroupedByInstr{}, nil
	}
	q := `select 
		id,
		instructor_id,
		date_start,
		date_end
	from instr_vacations
	where instructor_id = any($1::uuid[]) and date_end >= $2`
	dbres, err := ctx.Dal.Db.Query(q, pq.StringArray(iids), time.Now().In(time.UTC))
	if err != nil {
		return nil, err
	}
	defer dbres.Close()
	res := make(VacationsGroupedByInstr)
	var tmp VacationInfo
	var instrID uuid.UUID
	for dbres.Next() {
		if err = dbres.Scan(&tmp.ID, &instrID, &tmp.DateStart, &tmp.DateEnd); err != nil {
			return nil, err
		}
		if res[instrID] == nil {
			res[instrID] = make([]VacationInfo, 0, 5)
		}
		res[instrID] = append(res[instrID], tmp)
	}
	return res, nil
}
