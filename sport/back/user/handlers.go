package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"sport/api"
	"sport/helpers"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PatchPasswordRequest struct {
	OldPassword string
	NewPassword string
}

func ValidatePassword(passw string) error {
	if len(passw) < 8 {
		return fmt.Errorf("password len < 8")
	}
	if len(passw) > UserPassMaxLength {
		return fmt.Errorf("too long password")
	}
	flags := 0
	for _, c := range passw {
		switch {
		case unicode.IsUpper(c):
			flags |= 1
		case unicode.IsLower(c):
			flags |= 2
		case unicode.IsNumber(c):
			flags |= 4
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			flags |= 8
		}
	}
	if flags != 15 {
		return fmt.Errorf("no required char")
	}
	return nil
}

func (pr *PatchPasswordRequest) Validate() error {
	if pr.OldPassword == "" {
		return fmt.Errorf("Validate PatchPasswordRequest: invalid old password")
	}
	if pr.NewPassword == "" || len(pr.NewPassword) > UserPassMaxLength {
		return fmt.Errorf("Validate PatchPasswordRequest: invalid new password")
	}
	return ValidatePassword(pr.NewPassword)
}

/*

PATCH /api/user/password

DESC
	update user password

REQUEST
	- Authorization: Bearer <token> in header
	- application/json body of type PatchPasswordRequest

RETURNS
	- status code:
		- 204 on success
		- error status code i.e. 401, 400
*/
func (ctx *Ctx) PatchPasswordHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var req PatchPasswordRequest
		var status = 401
		var newHash []byte
		var jb []byte
		var u *User
		var userID = g.MustGet("UserID").(uuid.UUID)

		if jb, err = ioutil.ReadAll(g.Request.Body); err != nil {
			goto end
		}

		if err = json.Unmarshal(jb, &req); err != nil {
			goto end
		}

		if err = req.Validate(); err != nil {
			goto end
		}

		if u, err = ctx.DalReadUser(userID, KeyTypeID, true); err != nil {
			goto end
		}

		if err = api.CmpHash(req.OldPassword, u.Passh); err != nil {
			goto end
		}

		// im setting default err status to 500 now since everything is validated
		// and errors are 95% not users fault
		status = 500

		if newHash, err = api.GetHash(req.NewPassword); err != nil {
			goto end
		}

		if err = ctx.DalUpdatePassh(userID, newHash); err != nil {
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		_ = g.AbortWithError(status, err)
	}
}

/*
	GET /api/user/password/validate
*/
func (ctx *Ctx) ValidatePassword() gin.HandlerFunc {
	return func(g *gin.Context) {
		p := g.Query("password")
		if err := ValidatePassword(p); err != nil {
			g.AbortWithStatus(400)
		} else {
			g.AbortWithStatus(204)
		}
	}
}

/*
DELETE /api/user

DESC
	delete user from the system

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code: 204 on success
*/
func (ctx *Ctx) DeleteHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var status = 400
		userID := g.MustGet("UserID").(uuid.UUID)

		if err = ctx.DalDeleteUser(userID); err != nil {
			goto end
		}

		g.AbortWithStatus(204)

		return

	end:
		_ = g.AbortWithError(status, err)
	}
}

/*

GET /api/user

DESC
	get user details

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- application/json of type User.
		note that some fields may be sanitized for security reasons

EXAMPLE
	curl -H "Authorization: Bearer 123" \
		-X GET http://127.0.0.1:1580/api/user
*/
func (ctx *Ctx) GetHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var status = 500
		var u *User

		if u, err = ctx.DalReadUser(
			g.MustGet("UserID").(uuid.UUID), KeyTypeID,
			true,
		); err != nil {
			if helpers.IsENF(err) {
				status = 404
			}
			goto end
		}

		g.AbortWithStatusJSON(200, u.PrivUserInfo)

		return

	end:
		_ = g.AbortWithError(status, err)
	}
}

type LoginRequest struct {
	Email    string
	Password string
}

/*

POST /api/user/login

DESC
	generate authorization token

REQUEST
	- application/json body of type LoginRequest

RETURNS
	- text/plain jwt token with status code 200 on success

EXAMPLE
	curl -H "Content-Type: application/json" \
		 -d '{"email":"marian@marian.pl","password":"password"}' \
		 -X POST http://127.0.0.1:1580/api/user/login
*/
func (ctx *Ctx) LoginHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		var ur LoginRequest
		var status = 401
		var jb []byte
		var u *User
		var tp *api.JwtPayload
		var token string

		if jb, err = ioutil.ReadAll(g.Request.Body); err != nil {
			goto end
		}

		if err = json.Unmarshal(jb, &ur); err != nil {
			goto end
		}

		if u, err = ctx.DalReadUser(ur.Email, KeyTypeEmail, true); err != nil {
			goto end
		}

		if err = api.CmpHash(ur.Password, u.Passh); err != nil {
			goto end
		}

		tp = new(api.JwtPayload)
		tp.Exp = time.Now().In(time.UTC).Add(time.Duration(time.Hour * 48))
		tp.UserID = u.ID

		if token, err = ctx.Api.Jwt.GenToken(tp); err != nil {
			goto end
		}

		if _, err = g.Writer.Write([]byte(token)); err != nil {
			status = 500
			goto end
		}
		g.AbortWithStatus(200)
		return

	end:
		_ = g.AbortWithError(status, err)
	}
}

/*

PATCH /api/user/data

DESC
	update user fields.
	merge will be made in accordance to
		https://datatracker.ietf.org/doc/html/rfc7396

REQUEST
	- Authorization: Bearer <token> in header
	- application/json body of type UserData
		note that only non-empty fields will be updated

RETURNS
	- status code:
		- 204 on success
		- error status code i.e. 401, 400
*/
func (ctx *Ctx) PatchUserDataHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var status = 400
		var u UserData
		userID := g.MustGet("UserID").(uuid.UUID)
		var jb []byte

		if jb, err = ioutil.ReadAll(g.Request.Body); err != nil {
			goto end
		}

		// validate that json request conforms to UserData
		// this wont be used in db update but i just want to verify that request
		// can indeed be parsed to userData,
		// so situations like Name: [] or Name = null
		//      dont emerge later and cause 500
		if err = json.Unmarshal(jb, &u); err != nil {
			goto end
		}

		if err = u.Validate(ctx, true); err != nil {
			status = 400
			goto end
		}

		if err = ctx.DalUpdateUserDataJson(userID, jb); err != nil {
			if helpers.IsENF(err) {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		_ = g.AbortWithError(status, err)
	}
}

/*

PATCH /api/user/contact

DESC
	update user fields.
	merge will be made in accordance to
		https://datatracker.ietf.org/doc/html/rfc7396

REQUEST
	- Authorization: Bearer <token> in header
	- application/json body of type ContactData
		note that only non-empty fields will be updated

RETURNS
	- status code:
		- 204 on success
		- error status code i.e. 401, 400
*/
func (ctx *Ctx) PatchUserContactDataHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var status = 400
		var u ContactData
		userID := g.MustGet("UserID").(uuid.UUID)
		var jb []byte

		if jb, err = ioutil.ReadAll(g.Request.Body); err != nil {
			goto end
		}

		// validate that json request conforms to UserData
		// this wont be used in db update but i just want to verify that request
		// can indeed be parsed to userData,
		// so situations like Name: [] or Name = null
		//      dont emerge later and cause 500
		if err = json.Unmarshal(jb, &u); err != nil {
			goto end
		}

		if err = u.Validate(); err != nil {
			status = 400
			goto end
		}

		if err = ctx.DalUpdateUserContactDataJson(userID, jb); err != nil {
			if helpers.IsENF(err) {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		_ = g.AbortWithError(status, err)
	}
}

/*

GET /api/user/stat

DESC
	validates token

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 204 if token if valid
		- 401 if token is not valid

*/
func (ctx *Ctx) Stat() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.AbortWithStatus(204)
	}
}

/*
desc to be done <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

PUT /api/user/photo

DESC
	upload photo endpoint

REQUEST
	- Authorization: Bearer <token> in header
	- multipart/form-data with file named "file" <- thats gonna be avatar photo

RETURNS
	- status code:
		- 204 if token if valid
		- 401 if token is not valid

*/
func (ctx *Ctx) PutUserAvatar() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		var status = 400
		userID := g.MustGet("UserID").(uuid.UUID)
		var fpath string
		var u *User
		var tx *sql.Tx
		var h *multipart.FileHeader

		u, err = ctx.DalReadUser(userID, KeyTypeID, true)
		if err != nil {
			status = 404
			goto end
		}

		if h, err = g.FormFile("image"); err != nil {
			status = 400
			goto end
		}

		if fpath, err = ctx.Static.ValidateImgAndGetRelpath(h, "user"); err != nil {
			status = 400
			goto end
		}

		if tx, err = ctx.Dal.Db.Begin(); err != nil {
			status = 500
			goto end
		}

		if err = ctx.DalUpdateUserAvatar(tx, userID, fpath); err != nil {
			if err == sql.ErrNoRows {
				status = 404
			} else {
				status = 500
			}
			tx.Rollback()
			goto end
		}

		if err = ctx.Static.UpsertImage(h, fpath, u.AvatarRelpath); err != nil {
			status = 500
			tx.Rollback()
			goto end
		}

		if err = tx.Commit(); err != nil {
			status = 500
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		if e := helpers.HttpErr(err); e != nil {
			e.WriteAndAbort(g)
		} else {
			_ = g.AbortWithError(status, err)
		}
	}
}

/*
	DELETE /api/user/photo
*/
func (ctx *Ctx) DeleteUserAvatar() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)

		u, err := ctx.DalReadUser(userID, KeyTypeID, true)
		if err != nil {
			g.AbortWithError(404, err)
			return
		}
		tx, err := ctx.Dal.Db.Begin()
		if err != nil {
			g.AbortWithError(500, err)
			return
		}
		if u.AvatarRelpath == "" {
			g.AbortWithError(404, fmt.Errorf("avatar does not exist"))
			return
		}
		if err := ctx.DalUpdateUserAvatar(tx, userID, ""); err != nil {
			g.AbortWithError(500, err)
			return
		}
		if err := ctx.Static.DeleteImg(u.AvatarRelpath); err != nil {
			g.AbortWithError(500, err)
			tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			g.AbortWithError(500, err)
			return
		}
		g.AbortWithStatus(204)
	}
}
