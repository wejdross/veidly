package sub

import (
	"database/sql"
	"fmt"
	"sport/helpers"
	"sport/instr"
	"sport/user"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (ctx *Ctx) DalCreateSub(s *Sub) error {
	q := `insert into subs (
		id,

		sub_model,
		sub_model_id,
		instructor_id,
		instr_user_id,
		user_id,
		ref_id,

		instructor_decision,
		state,

		sm_cache,
		sm_retries,
		sm_timeout,

		is_confirmed,
		is_active,
		order_id,

		created_on,

		link_id,
		link_url,

		date_start,
		date_end,

		remaining_entries
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
		$21
	)`
	_, err := ctx.Dal.Db.Exec(q,
		s.ID,

		&s.SubModel,
		s.SubModelID,
		s.InstructorID,
		s.InstrUserID,
		s.UserID,
		s.RefID,

		s.InstructorDecision,
		s.State,

		&s.SmCache,
		s.SmRetries,
		s.SmTimeout,

		s.IsConfirmed,
		s.IsActive,
		s.OrderID,

		s.CreatedOn,

		s.LinkID,
		s.LinkUrl,

		s.DateStart,
		s.DateEnd,

		s.RemainingEntries)
	return err
}

type ReadSubRequest struct {
	ID           *uuid.UUID
	RefID        *uuid.UUID
	InstructorID *uuid.UUID

	UserID      *uuid.UUID
	InstrUserID *uuid.UUID

	DateRange helpers.DateRange

	IDs []string
}

var EmptyReadSubRequest ReadSubRequest

type SubWithJoins struct {
	Sub
	Instructor  instr.PubInstructorWithUser
	UserInfo    user.PubUserInfo
	ContactData user.ContactData
}

func (r *ReadSubRequest) WriteSubWhereFragment(qb *strings.Builder, args *[]interface{}) {
	nextArg := 1

	if r.ID != nil {
		qb.WriteString(fmt.Sprintf(" s.id = $%d and ", nextArg))
		nextArg++
		*args = append(*args, *r.ID)
	}

	if r.IDs != nil {
		qb.WriteString((fmt.Sprintf(" s.id = any($%d) and ", nextArg)))
		nextArg++
		*args = append(*args, pq.StringArray(r.IDs))
	}

	if r.InstructorID != nil {
		qb.WriteString(fmt.Sprintf(" s.instructor_id=$%d and ", nextArg))
		nextArg++
		*args = append(*args, *r.InstructorID)
	}

	if r.RefID != nil {
		qb.WriteString(fmt.Sprintf(" s.ref_id = $%d and ", nextArg))
		nextArg++
		*args = append(*args, *r.RefID)
	}

	if r.UserID != nil {
		qb.WriteString(fmt.Sprintf(" s.user_id = $%d and ", nextArg))
		nextArg++
		*args = append(*args, *r.UserID)
	}

	if r.InstrUserID != nil {
		qb.WriteString(fmt.Sprintf(" s.instr_user_id = $%d and ", nextArg))
		nextArg++
		*args = append(*args, *r.InstrUserID)
	}

	if r.DateRange.IsNotZero() {
		qb.WriteString(fmt.Sprintf(` S.date_start <= $%d and `, nextArg))
		nextArg++
		qb.WriteString(fmt.Sprintf(` S.date_end >= $%d and `, nextArg))
		nextArg++
		*args = append(
			*args,
			r.DateRange.End.In(time.UTC),
			r.DateRange.Start.In(time.UTC))
	}
}

func (ctx *Ctx) DalReadSubWithJoins(r ReadSubRequest) ([]SubWithJoins, error) {
	qb := strings.Builder{}
	qb.WriteString(fmt.Sprintf(`
			select 
				%s,

				%s,
				%s,
				%s
				
			from subs s 
			left join users uu on uu.id = s.user_id
			left join instructors i on i.id = s.instructor_id
			left join users iu on iu.id = s.instr_user_id
			where `,
		SubColumns(),
		user.GetPubUserInfoSelectColsFmt("uu"),
		instr.PubInstructorWithUserCols("iu"),
		user.GetContactDataSelectCols("uu")))

	args := make([]interface{}, 0, 3)

	r.WriteSubWhereFragment(&qb, &args)

	qb.WriteString(" 1=1 ")

	len := 2
	if r.ID != nil || r.RefID != nil {
		len = 1
	}

	var s SubWithJoins
	var res = make([]SubWithJoins, 0, len)

	rdr, err := ctx.Dal.Db.Query(qb.String(), args...)
	if err != nil {
		return nil, err
	}

	defer rdr.Close()

	scanArgs := s.ScanFields()

	scanArgs = append(scanArgs, s.UserInfo.ScanFields()...)
	scanArgs = append(scanArgs, s.Instructor.ScanFields()...)
	scanArgs = append(scanArgs, s.ContactData.ScanFields()...)

	for rdr.Next() {
		if err := rdr.Scan(scanArgs...); err != nil {
			return nil, err
		}
		s.Instructor.PubInstructorInfo.PostprocessAfterDbScan(ctx.Instr)
		s.UserInfo.PostprocessAfterDbScan(ctx.User)
		res = append(res, s)
	}

	return res, nil
}

func (ctx *Ctx) DalReadSubs(r ReadSubRequest) ([]Sub, error) {

	qb := strings.Builder{}
	qb.WriteString(fmt.Sprintf(`
			select %s
			from subs s where `,
		SubColumns()))

	args := make([]interface{}, 0, 3)

	r.WriteSubWhereFragment(&qb, &args)

	qb.WriteString(" 1=1 ")

	len := 2
	if r.ID != nil || r.RefID != nil {
		len = 1
	}

	var s Sub
	var res = make([]Sub, 0, len)

	rdr, err := ctx.Dal.Db.Query(qb.String(), args...)
	if err != nil {
		return nil, err
	}

	defer rdr.Close()

	scanArgs := s.ScanFields()

	for rdr.Next() {
		if err := rdr.Scan(scanArgs...); err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, nil
}

func (ctx *Ctx) DalReadSingleSubWithJoins(r ReadSubRequest) (*SubWithJoins, error) {
	res, err := ctx.DalReadSubWithJoins(r)
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, sql.ErrNoRows
	}
	return &res[0], nil
}

func (ctx *Ctx) UpdateSubSmTimeout(id uuid.UUID, timeout time.Time) error {
	res, err := ctx.Dal.Db.Exec(
		`update subs set sm_timeout = $1 where id = $2`,
		timeout, id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalReadSingleSub(r ReadSubRequest) (*Sub, error) {
	res, err := ctx.DalReadSubs(r)
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, sql.ErrNoRows
	}
	return &res[0], nil
}

func (ctx *Ctx) DalConfirmSub(subID, instrUserID uuid.UUID) error {
	res, err := ctx.Dal.Db.Exec(
		`update subs set 
			remaining_entries = remaining_entries - 1
		where 
			remaining_entries > 0 and 
			id = $1 and 
			instr_user_id = $2 and
			is_confirmed = true`,
		subID, instrUserID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}
