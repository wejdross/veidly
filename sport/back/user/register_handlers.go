package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/mail"
	"net/url"
	"path"
	"sport/api"
	"sport/helpers"
	"strconv"

	"github.com/gin-gonic/gin"
)

const ConfirmRegisterPath = "/user/register/confirm"

/*
GET /api/{confirmRegisterPath}

DESC
	finish 2fa registration process

REQUEST
	- query string:
		- token: REQUIRED token obtained from registration email

RETURNS
	- status code 307 with redirect to UI login page.
		- query string
			- email - on success user email
			- state - status information about registration

*/

func (ctx *Ctx) ConfirmRegisterHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var status = 401
		var u *User
		var email string

		token := g.Query("token")
		if token == "" {
			err = fmt.Errorf("token was empty")
			goto end
		}

		if !ctx.Config.Use2Fa {
			status = 501
			err = fmt.Errorf("user requested confirm register without 2fa being enabled on server")
			goto end
		}

		status = 500

		if u, err = ctx.DalReadUser(token, KeyType2FAToken, false); err != nil {
			if helpers.IsENF(err) {
				status = 401
			}
			goto end
		}

		u.Enabled = true

		if err = ctx.DalConfirmUser(u.ID); err != nil {
			if helpers.IsENF(err) {
				status = 403
			}
			goto end
		}

		status = 200

	end:
		if u != nil {
			email = u.Email
		}
		g.Redirect(307, fmt.Sprintf(
			ctx.Config.RegisterRedirectURL,
			url.QueryEscape(email),
			strconv.Itoa(status),
			url.QueryEscape(g.Query("return_url")),
		))
		//_ = ctx.AbortWithError(status, err)
	}
}

type MockupRegisterResponse struct {
	HrefFromEmail string
	HandlerRegisterResponse
}

type HandlerRegisterResponse struct {
	// this field informs user whether 2fa is enabled and
	// one needs to perform more steps (like email validation)
	// to finish register
	MFA bool `json:"mfa"`
	// this field informs how much time user has to activate account (via 2fa)
	// before it will be removed
	// it is set to -1 if 2fa is disabled
	TTLSeconds int `json:"ttl_seconds"`
}

/*
POST /api/user/register

REQUEST
	- application/json body of type UserRequest with user data

RETURNS:
	- application/json body of type HandlerRegisterResponse

EXAMPLE USAGE:

	- curl -H "Content-Type: application/json"  \
			-d '{"email":"marian@marian.pl","password":"password"}' \
			-X POST http://127.0.0.1:1580/api/user/register
*/
func (ctx *Ctx) RegisterHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var status = 401
		var ur UserRequest
		var u *User
		var jb []byte
		var confirmUrl string
		var html string

		if jb, err = ioutil.ReadAll(g.Request.Body); err != nil {
			goto end
		}

		if err = json.Unmarshal(jb, &ur); err != nil {
			goto end
		}

		if err = ur.Validate(ctx, true); err != nil {
			status = 400
			goto end
		}

		if u, err = ur.NewUser(ctx.Config.Use2Fa); err != nil {
			goto end
		}

		if err = ctx.DalCreateUser(u); err != nil {
			if helpers.PgIsUqViolation(err) {
				if u, err = ctx.DalReadUser(u.Email, KeyTypeEmail, false); err != nil {
					goto end
				}
				if u.Enabled {
					err = fmt.Errorf("Requested register on enabled user")
					goto end
				}
			} else {
				goto end
			}
		}

		if ctx.Config.Use2Fa {
			href := path.Join(
				api.ApiBaseHref,
				ConfirmRegisterPath) + "?token=" + u.MFAToken +
				"&return_url=" + url.QueryEscape(g.Query("mfa_return_url"))
			confirmUrl = ctx.Config.PublicURL + href

			if html, err = ctx.GetRegisterHtml(ur.Language, &RegisterData{
				Name: ur.Name,
				Url:  confirmUrl,
			}); err != nil {
				status = 500
				goto end
			}

			if err = ctx.Email.SendHtmlMail(
				mail.Address{Name: ur.Name, Address: u.Email},
				Locale[ur.Language][RegisterTitle],
				html,
			); err != nil {
				goto end
			}
			if ctx.IsMockup2Fa {
				g.AbortWithStatusJSON(200, MockupRegisterResponse{
					href, HandlerRegisterResponse{
						true, ctx.Config.RetentionLoopRecordTTLseconds,
					}})
			} else {
				g.AbortWithStatusJSON(200, HandlerRegisterResponse{
					true, ctx.Config.RetentionLoopRecordTTLseconds})
			}
		} else {
			g.AbortWithStatusJSON(200, HandlerRegisterResponse{false, -1})
		}

		return

	end:
		_ = g.AbortWithError(status, err)
	}
}

type ResendRegisterEmailRequest struct {
	Email string
}

func (r *ResendRegisterEmailRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("Validate ResendRegisterEmailRequest: invalid email")
	}
	return nil
}

/*

POST /api/user/register/resend

REQUEST
	- application/json body of type ResendRegisterEmailRequest with user data

RETURNS:
	- status code

*/
func (ctx *Ctx) ResendRegisterEmailHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		var status = 401
		var ur ResendRegisterEmailRequest
		var u *User
		var confirmUrl string
		var html string

		if !ctx.Config.Use2Fa {
			status = 501
			err = fmt.Errorf("user requested resend register email" +
				" without 2fa being enabled on server")
			goto end
		}

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &ur, ur.Validate,
		); err != nil {
			status = 400
			goto end
		}

		if u, err = ctx.DalReadUser(ur.Email, KeyTypeEmail, false); err != nil {
			if err == sql.ErrNoRows {
				status = 404
			}
			goto end
		}

		if u.Enabled {
			err = fmt.Errorf("Cant resend registration email" +
				" since user is already enabled")
			status = 409
			goto end
		}

		confirmUrl = ctx.Config.PublicURL + path.Join(api.ApiBaseHref,
			ConfirmRegisterPath) + "?token=" + u.MFAToken +
			"&return_url=" + url.QueryEscape(g.Query("mfa_return_url"))

		{
			if html, err = ctx.GetRegisterHtml(u.Language, &RegisterData{
				Name: u.Name,
				Url:  confirmUrl,
			}); err != nil {
				status = 500
				goto end
			}

			if err = ctx.Email.SendHtmlMail(
				mail.Address{Name: u.Name, Address: u.Email},
				Locale[u.Language][RegisterTitle],
				html,
			); err != nil {
				goto end
			}
		}

		g.AbortWithStatusJSON(200, HandlerRegisterResponse{
			true, ctx.Config.RetentionLoopRecordTTLseconds})

		return

	end:
		_ = g.AbortWithError(status, err)
	}
}
