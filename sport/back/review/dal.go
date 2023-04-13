package review

import (
	"database/sql"
	"fmt"
	"sport/helpers"
	"sport/user"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (ctx *Ctx) DalCreateReviewToken(o *Review) error {
	q := `insert into reviews (
		id,
		access_token,
		created_on,

		training_id,
		rsv_id,
		user_id,
		email,
		user_data,

		mark,
		review
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
		$10
	)`
	_, err := ctx.Dal.Db.Exec(q,
		o.ID,
		o.AccessToken,
		o.CreatedOn,

		o.TrainingID,
		o.RsvID,
		o.UserID,
		o.Email,
		&o.UserInfo.UserData,

		o.Mark,
		o.Review)
	return err
}

type ReadReviewsOpts struct {
	TrainingID *uuid.UUID
}

func (ctx *Ctx) DalReadPubReviews(opts ReadReviewsOpts, tx *sql.Tx) ([]PubReview, error) {
	args := make([]interface{}, 0, 2)
	q := strings.Builder{}
	q.WriteString(`select 
		r.created_on,

		r.user_data,

		r.mark,
		r.review,

		u.id,
		coalesce(u.user_data, '{}'::jsonb),
		coalesce(u.avatar_relpath, '')

	from reviews r
		left join users u on u.id = r.user_id
	where access_token = '' `)

	if opts.TrainingID != nil {
		q.WriteString(" and r.training_id = $1 ")
		args = append(args, *opts.TrainingID)
	}

	var err error
	var rows *sql.Rows

	if tx != nil {
		rows, err = tx.Query(q.String(), args...)
	} else {
		rows, err = ctx.Dal.Db.Query(q.String(), args...)
	}

	if err != nil {
		return nil, err
	}

	var res = make([]PubReview, 0, 2)

	var joinedUserID *uuid.UUID
	var joinedUserInfo user.PubUserInfo
	for i := 0; rows.Next(); i++ {
		res = append(res, PubReview{})
		joinedUserID = nil
		if err := rows.Scan(
			&res[i].CreatedOn,

			&res[i].UserInfo.UserData,

			&res[i].Mark,
			&res[i].Review,

			&joinedUserID,
			&joinedUserInfo.UserData,
			&joinedUserInfo.AvatarRelpath,
		); err != nil {
			return nil, err
		}

		if joinedUserID != nil {
			res[i].UserInfo = joinedUserInfo
		}
	}

	return res, nil
}

type ReadSingleRevKeyType int

type SingleRevKeyUserRsvKey struct {
	UserID uuid.UUID
	RsvID  uuid.UUID
}

const (
	SingleRevKeyID          ReadSingleRevKeyType = 1
	SingleRevKeyAccessToken                      = 2
	SingleRevKeyRsvID                            = 3
	SingleRevKeyUserRsv                          = 4
	SingleRevKeyTrainingID                       = 5
)

func (ctx *Ctx) DalReadSingleReview(
	key interface{},
	keyType ReadSingleRevKeyType,
	t ReviewType,
	tx *sql.Tx,
) (*Review, error) {

	q := strings.Builder{}
	q.WriteString(`select 
		r.id,
		r.access_token,
		r.created_on,

		r.training_id,
		r.rsv_id,
		r.user_id,
		r.email,
		r.user_data,

		r.mark,
		r.review,

		u.id,
		coalesce(u.user_data, '{}'::jsonb),
		coalesce(u.avatar_relpath, ''),
		coalesce(u.contact_data->>'email', '')

	from reviews r
		left join users u on u.id = r.user_id
	where `)

	args := make([]interface{}, 0, 2)

	switch keyType {
	case SingleRevKeyID:
		q.WriteString(" id = $1")
		args = append(args, key)
		break
	case SingleRevKeyAccessToken:
		accessToken := key.(string)
		if accessToken == "" {
			return nil, fmt.Errorf("no access token provided")
		}
		q.WriteString(" access_token = $1")
		args = append(args, key)
		break
	case SingleRevKeyRsvID:
		q.WriteString(" rsv_id = $1")
		args = append(args, key)
		break
	case SingleRevKeyUserRsv:
		k := key.(SingleRevKeyUserRsvKey)
		q.WriteString(" rsv_id = $1 and user_id = $2 ")
		args = append(args, k.RsvID, k.UserID)
		break
	default:
		return nil, fmt.Errorf("invalid key type")
	}

	switch t {
	case AnyReviewType:
		break
	case ContentReviewType:
		q.WriteString(" and access_token = '' ")
		break
	case TokenReviewType:
		q.WriteString(" and access_token != '' ")
		break
	default:
		return nil, fmt.Errorf("invalid review type")
	}

	var row *sql.Row

	if tx != nil {
		row = tx.QueryRow(q.String(), args...)
	} else {
		row = ctx.Dal.Db.QueryRow(q.String(), args...)
	}

	res := new(Review)

	var uid *uuid.UUID
	var ui user.PubUserInfo
	var cd user.ContactData

	if err := row.Scan(
		&res.ID,
		&res.AccessToken,
		&res.CreatedOn,

		&res.TrainingID,
		&res.RsvID,
		&res.UserID,
		&res.Email,
		&res.UserInfo.UserData,

		&res.Mark,
		&res.Review,

		&uid,
		&ui.UserData,
		&ui.AvatarRelpath,
		&cd.Email,
	); err != nil {
		return nil, err
	}

	if uid != nil {
		res.UserInfo = ui
		res.Email = cd.Email
		res.UserInfo.PostprocessAfterDbScan(ctx.User)
	}

	return res, nil
}

func (ctx *Ctx) DalUpdateReview(req *UpdateReviewRequest) error {

	if req.AccessToken == "" {
		return fmt.Errorf("no access token provided")
	}

	tx, err := ctx.Dal.Db.Begin()
	if err != nil {
		return err
	}

	o, err := ctx.DalReadSingleReview(
		req.AccessToken,
		SingleRevKeyAccessToken,
		TokenReviewType,
		tx)
	if err != nil {
		return err
	}

	uq := `update reviews set 
		mark = $1,
		review = $2,
		access_token = ''
	where access_token != '' and access_token = $3`
	res, err := ctx.Dal.Db.Exec(uq, req.Mark, req.Review, req.AccessToken)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := helpers.PgMustBeOneRow(res); err != nil {
		tx.Rollback()
		return err
	}

	if req.Mark > 0 && req.Mark <= MaxMark && o.TrainingID != nil {
		if err = ctx.Train.DalUpdateTrainingAvgMark(*o.TrainingID, req.Mark, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (ctx *Ctx) DalDeleteReview(userID uuid.UUID, req *DeleteReviewRequest) error {
	if userID == uuid.Nil {
		return fmt.Errorf("invalid userID")
	}

	tx, err := ctx.Dal.Db.Begin()
	if err != nil {
		return err
	}

	o, err := ctx.DalReadSingleReview(
		req.ID,
		SingleRevKeyID,
		AnyReviewType,
		tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	q := `delete from reviews where id = $1 and user_id = $2`
	res, err := tx.Exec(q, req.ID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = helpers.PgMustBeOneRow(res)
	if err != nil {
		tx.Rollback()
		return err
	}

	if o.AccessToken == "" && o.Mark > 0 && o.TrainingID != nil {
		if err = ctx.Train.DalUpdateTrainingAvgMark(*o.TrainingID, -o.Mark, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (ctx *Ctx) DalDeleteExpiredReviewTokens(exp time.Duration) error {
	q := `delete from reviews where access_token != '' and created_on < $1`
	_, err := ctx.Dal.Db.Exec(
		q,
		time.Now().In(time.UTC).Add(-time.Duration(ctx.Config.ReviewExp)))
	return err
}

// resets created on for rsv to specified time.
// meant to be used in testing only
func (ctx *Ctx) DalTestResetCreatedOn(accessToken string, t time.Time) error {
	q := `update reviews set created_on = $1 where access_token = $2`
	res, err := ctx.Dal.Db.Exec(q, t, accessToken)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}
