package dal

import "database/sql"

// will perform QueryRow
// db can be either *sql.Tx or *dal.Dal
func QueryRowOnIface(db interface{}, q string, args ...interface{}) *sql.Row {
	var res *sql.Row
	switch db.(type) {
	case *Ctx:
		res = db.(*Ctx).Db.QueryRow(q, args...)
		break
	case *sql.Tx:
		res = db.(*sql.Tx).QueryRow(q, args...)
		break
	default:
		panic("invalid argument provided to QueryRowOnIface")
	}
	return res
}

// will perform Exec
// db can be either *sql.Tx or *dal.Dal
func ExecOnIface(db interface{}, q string, args ...interface{}) (sql.Result, error) {
	var res sql.Result
	var err error
	switch db.(type) {
	case *Ctx:
		res, err = db.(*Ctx).Db.Exec(q, args...)
		break
	case *sql.Tx:
		res, err = db.(*sql.Tx).Exec(q, args...)
		break
	default:
		panic("invalid argument provided to ExecOnIface")
	}
	return res, err
}

// will perform QueryRow
// db can be either *sql.Tx or *dal.Dal
func QueryOnIface(db interface{}, q string, args ...interface{}) (*sql.Rows, error) {
	var res *sql.Rows
	var err error
	switch db.(type) {
	case *Ctx:
		res, err = db.(*Ctx).Db.Query(q, args...)
		break
	case *sql.Tx:
		res, err = db.(*sql.Tx).Query(q, args...)
		break
	default:
		panic("invalid argument provided to QueryOnIface")
	}
	return res, err
}
