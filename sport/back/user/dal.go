package user

import (
	"database/sql"
	"fmt"
	"sport/dal"
	"sport/helpers"
	"time"

	"github.com/google/uuid"
)

func (ctx *Ctx) DalConfirmUser(id uuid.UUID) error {
	r, err := ctx.Dal.Db.Exec(`update users set 
			enabled = true, 
			mfa_token = '' 
		where id = $1 and enabled = false`, id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(r)
}

func (ctx *Ctx) DalUpdatePassh(id uuid.UUID, passh []byte) error {
	r, err := ctx.Dal.Db.Exec(`update users set 
			passh = $1, 
			forgot_password_token = '' 
		where id = $2`,
		passh,
		id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(r)
}

func (ctx *Ctx) DalUpdateUserAvatar(tx *sql.Tx, id uuid.UUID, relpath string) error {
	const q = "update users set avatar_relpath = $1 where id = $2"
	var args = []interface{}{relpath, id}
	var r sql.Result
	var err error
	if tx != nil {
		r, err = tx.Exec(q, args...)
	} else {
		r, err = ctx.Dal.Db.Exec(q, args...)
	}
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(r)
}

func (ctx *Ctx) DalUpdateUserDataJson(
	id uuid.UUID,
	j []byte,
) error {
	q := `update users 
		set user_data = user_data || $1 
	where enabled = true and id = $2`
	r, err := ctx.Dal.Db.Exec(q, j, id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(r)
}

func (ctx *Ctx) DalUpdateUserContactDataJson(
	id uuid.UUID,
	j []byte,
) error {
	q := `update users 
		set contact_data = contact_data || $1 
	where enabled = true and id = $2`
	r, err := ctx.Dal.Db.Exec(q, j, id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(r)
}

func (ctx *Ctx) DalUpdateUserData(
	id uuid.UUID,
	u *UserData,
) error {
	q := `update users 
		set user_data = $1 
	where enabled = true and id = $2`
	r, err := ctx.Dal.Db.Exec(q, u, id)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(r)
}

func (ctx *Ctx) DalCreateUser(u *User) error {
	_, err := ctx.Dal.Db.Exec(`insert into users (
			id,
			email,
			passh,
			access_failed, 
			enabled,
			created_on,
			mfa_token,
			forgot_password_token,
			oauth_provider,
			oauth_id,
			user_data,
			avatar_relpath,
			contact_data
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
			$13
		)`,
		u.ID,
		u.Email,
		u.Passh,
		u.AccessFailed,
		u.Enabled,
		u.CreatedOn,
		u.MFAToken,
		u.ForgotPassToken,
		u.OauthProvider,
		u.OauthID,
		&u.UserData,
		u.AvatarRelpath,
		&u.ContactData,
	)
	return err
}

type GetUserKeyType int

const (
	KeyTypeEmail         GetUserKeyType = 1
	KeyTypeID            GetUserKeyType = 2
	KeyType2FAToken      GetUserKeyType = 3
	KeyTypePassToken     GetUserKeyType = 4
	KeyTypeOauthUserInfo GetUserKeyType = 5
)

func (ctx *Ctx) DalUpdateUserPassToken(email, token string) error {
	q := `update users set forgot_password_token = $1 where enabled = true and email = $2 and oauth_provider = ''`
	res, err := ctx.Dal.Db.Exec(q, token, email)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalReadUser(
	key interface{}, keyType GetUserKeyType,
	onlyEnabled bool,
) (*User, error) {
	d := ctx.Dal
	q := `select
			id,
			email,
			passh,
			access_failed,
			enabled,
			created_on,
			mfa_token,
			forgot_password_token,
			oauth_provider,
			oauth_id,
			user_data,
			avatar_relpath,
			contact_data
		from users where `
	if onlyEnabled {
		q += " enabled = true and "
	}
	var i *sql.Rows
	var err error
	switch keyType {
	case KeyTypeEmail:
		i, err = d.Db.Query(
			q+" email != '' and email = $1", key)
	case KeyTypeID:
		i, err = d.Db.Query(
			q+" id = $1", key)
	case KeyType2FAToken:
		i, err = d.Db.Query(
			q+" mfa_token != '' and mfa_token = $1", key)
	case KeyTypePassToken:
		i, err = d.Db.Query(
			q+" forgot_password_token != '' and forgot_password_token = $1", key)
	case KeyTypeOauthUserInfo:
		ui, ok := key.(*OauthUserInfo)
		if !ok {
			return nil, fmt.Errorf("unsupported key value")
		}
		i, err = d.Db.Query(
			q+" email != '' and email = $1 and "+
				"oauth_id != '' and oauth_id = $2"+
				" and oauth_provider != '' and oauth_provider = $3",
			ui.Email,
			ui.OauthID,
			ui.Provider)
	default:
		return nil, fmt.Errorf("unsupported search key type")
	}
	if err != nil {
		return nil, err
	}
	var ret User
	defer i.Close()
	if !i.Next() {
		return nil, helpers.NewElementNotFoundErr(key)
	}
	if err = i.Scan(
		&ret.ID,
		&ret.Email,
		&ret.Passh,
		&ret.AccessFailed,
		&ret.Enabled,
		&ret.CreatedOn,
		&ret.MFAToken,
		&ret.ForgotPassToken,
		&ret.OauthProvider,
		&ret.OauthID,
		&ret.UserData,
		&ret.AvatarRelpath,
		&ret.ContactData,
	); err != nil {
		return nil, err
	}
	ret.PostprocessAfterDbScan(ctx)
	return &ret, err
}

func (ctx *Ctx) DalReadUserLangOrDefault(
	key interface{}, keyType GetUserKeyType,
) (string, error) {
	u, err := ctx.DalReadUser(key, keyType, true)
	if err != nil {
		return "", err
	}
	if u.Language == "" {
		return LangDefault, nil
	}
	return u.Language, nil
}

func (ctx *Ctx) DalDeleteInactiveUsers(allowedTTL time.Duration) error {
	_, err := ctx.Dal.Db.Exec(
		`delete from users where enabled = false and created_on < $1`,
		time.Now().In(time.UTC).Add(-allowedTTL))
	return err
}

func (ctx *Ctx) DalDeleteUser(id uuid.UUID) error {
	_, err := ctx.Dal.Db.Exec("delete from users where id = $1", id)
	return err
}

/*
	convert [disabled] user whose registration was already started via 2fa
	to oauth
*/
func DalMigrateUserToOauth(d *dal.Ctx, id uuid.UUID, ui *OauthUserInfo) error {
	_, err := d.Db.Exec(
		`update users set 
			passh = '', 
			oauth_id = $1, 
			oauth_provider = $2,
			enabled = true
		where 
			enabled = false and 
			passh != '' and 
			id = $3 and 
			email = $4`,
		ui.OauthID,
		ui.Provider,
		id,
		ui.Email,
	)
	return err
}
