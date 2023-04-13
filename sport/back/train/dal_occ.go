package train

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// // warning: no validation, make sure to validate that user can access this occ
// // before calling this function
func (ctx *Ctx) DalCreateOccs(occ []Occ, tx *sql.Tx) error {

	if len(occ) == 0 {
		return nil
	}

	qb := strings.Builder{}
	qb.WriteString(`insert into occurrences (
		id,
		training_id,
		date_start,
		date_end,
		repeat_days,
		color,
		remarks
	) values `)

	next := 1
	args := make([]interface{}, len(occ)*7)

	for i := range occ {
		if i > 0 {
			qb.WriteString(",")
		}
		qb.WriteString(fmt.Sprintf(
			`($%d, $%d, $%d, $%d, $%d, $%d, $%d)`,
			next, next+1, next+2, next+3, next+4, next+5, next+6))
		// -1 since next starts with 1
		args[next-1] = occ[i].ID
		args[next] = occ[i].TrainingID
		args[next+1] = occ[i].DateStart
		args[next+2] = occ[i].DateEnd
		args[next+3] = occ[i].RepeatDays
		args[next+4] = occ[i].Color
		args[next+5] = occ[i].Remarks
		next += 7
	}

	var err error
	if tx == nil {
		_, err = ctx.Dal.Db.Exec(qb.String(), args...)
	} else {
		_, err = tx.Exec(qb.String(), args...)
	}
	return err
}

func (ctx *Ctx) DalCreate2Occs(occ []SecondaryOcc, tx *sql.Tx) error {
	if len(occ) == 0 {
		return nil
	}

	qb := strings.Builder{}
	qb.WriteString(`insert into secondary_occs (
		id,
		occ_id,
		training_id,
		offset_start,
		offset_end,
		color,
		remarks
	) values `)

	next := 1
	args := make([]interface{}, len(occ)*7)

	for i := range occ {
		if i > 0 {
			qb.WriteString(",")
		}
		qb.WriteString(fmt.Sprintf(
			`($%d, $%d, $%d, $%d, $%d, $%d, $%d)`,
			next, next+1, next+2, next+3, next+4, next+5, next+6))
		args[next-1] = occ[i].ID
		args[next] = occ[i].OccID
		args[next+1] = occ[i].TrainingID
		args[next+2] = occ[i].OffsetStart
		args[next+3] = occ[i].OffsetEnd
		args[next+4] = occ[i].Color
		args[next+5] = occ[i].Remarks
		next += 7
	}

	var err error
	if tx == nil {
		_, err = ctx.Dal.Db.Exec(qb.String(), args...)
	} else {
		_, err = tx.Exec(qb.String(), args...)
	}
	return err
}

func (ctx *Ctx) DalDeleteTrainingOccs(trainingID uuid.UUID, tx *sql.Tx) error {
	const q = `delete from occurrences where training_id = $1`
	args := []interface{}{trainingID}
	var err error
	if tx == nil {
		_, err = ctx.Dal.Db.Exec(q, args...)
	} else {
		_, err = tx.Exec(q, args...)
	}
	return err
}
