package helpers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// to be used in dynamic sql queries
var ZeroTimeString = "'0001-01-01 00:00:00'::timestamp without time zone"
var ZeroUuidString = fmt.Sprintf("'%s'", uuid.UUID{})

func PgIsUqViolation(err error) bool {
	if e, ok := err.(*pq.Error); ok && e != nil {
		if e.Code == "23505" {
			return true
		}
	}
	return false
}

func PgIsFkViolation(err error) bool {
	if e, ok := err.(*pq.Error); ok && e != nil {
		if e.Code == "23503" {
			return true
		}
	}
	return false
}

func PgIsDbNotExists(err error) bool {
	if e, ok := err.(*pq.Error); ok && e != nil {
		if e.Code == "3D000" {
			return true
		}
	}
	return false
}

func PgIsConcurrentUpdate(err error) bool {
	if e, ok := err.(*pq.Error); ok && e != nil {
		if e.Code == "40001" {
			return true
		}
	}
	return false
}

func PgIsInvalidDatetimeFormat(err error) bool {
	if e, ok := err.(*pq.Error); ok && e != nil {
		if e.Code == "22007" {
			return true
		}
	}
	return false
}

func PgAddUpdateField(query *string, name string, currentArgCount *int) {
	if *currentArgCount > 0 {
		*query += fmt.Sprintf(",%s=$%d", name, *currentArgCount+1)
	} else {
		*query += fmt.Sprintf("%s=$%d", name, *currentArgCount+1)
	}
	*currentArgCount++
}

func PgMustBeOneRow(res sql.Result) error {
	if q, err := res.RowsAffected(); err != nil {
		return err
	} else {
		if q < 1 {
			return NewElementNotFoundErr("Update target")
		} else if q > 1 {
			return fmt.Errorf("invalid number of rows affected by update: %d", q)
		}
	}
	return nil
}

func PgUuidMatch(qb *strings.Builder, ids []string, next int, params *[]interface{}, col string) {
	if len(ids) == 0 {
		return
	}
	qb.WriteString(" ( ")
	for i := range ids {
		if i > 0 {
			qb.WriteString("or")
		}
		qb.WriteString(fmt.Sprintf(" %s = $%d::uuid ", col, next))
		next++
		(*params) = append(*params, ids[i])
	}
	qb.WriteString(") ")
}
