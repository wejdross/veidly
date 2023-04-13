package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/mail"
	"sport/api"
	"sport/helpers"

	"github.com/gin-gonic/gin"
)

/*

GET /api/user/password/forgot

DESC
	used to initialize password recovery procedure

REQUEST
	- query string:
		- lang: OPTIONAL parameter used to determine
			which language should be used during procedure
		- email: REQUIRED user email

RETURNS
	- 204 status code. Procedure will be continued via email

*/
func (ctx *Ctx) ForgotPasswordHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var status = 401
		var token, email string
		var resetPassUrl string
		var lang string
		var html string

		if !ctx.Config.Use2Fa {
			status = 501
			err = fmt.Errorf("user requested password reset without 2fa being enabled on server")
			goto end
		}

		lang = ctx.LangCtx.ApiLangOrDefault(g.Query("lang"))

		if email = g.Query("email"); email == "" {
			err = fmt.Errorf("email was not provided")
			goto end
		}

		token = helpers.GetUniqueToken()

		if err = ctx.DalUpdateUserPassToken(email, token); err != nil {
			if !helpers.IsENF(err) {
				status = 500
			}
			goto end
		}
		resetPassUrl = fmt.Sprintf(ctx.Config.UIPassResetUrl, token)

		html, err = ctx.GetForgotPassHtml(lang, &ForgotPassData{
			Name: "",
			Url:  resetPassUrl,
		})
		if err != nil {
			status = 500
			goto end
		}

		if err = ctx.Email.SendHtmlMail(
			mail.Address{Address: email},
			Locale[lang][PassRecoveryTitle],
			html,
		); err != nil {
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		_ = g.AbortWithError(status, err)
	}
}

type ResetPasswordRequest struct {
	// new password
	Password string
	// confirmation token obtained from email
	Token string
}

func (r *ResetPasswordRequest) Validate() error {
	if r.Token == "" {
		return fmt.Errorf("ResetPasswordRequest validate: invalid token")
	}
	if r.Password == "" {
		return fmt.Errorf("ResetPasswordRequest validate: invalid password")
	}
	return ValidatePassword(r.Password)
}

/*

POST /api/user/password/reset

DESC
	finishes password reset procedure

REQUEST
	- application/json body of type ResetPasswordRequest

RETURNS
	- status code:
		- 204 on success
		- error status code i.e. 401, 400
*/
func (ctx *Ctx) ResetPasswordHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var status = 401
		var ur ResetPasswordRequest
		var u *User
		var passh []byte
		var jb []byte

		if jb, err = ioutil.ReadAll(g.Request.Body); err != nil {
			goto end
		}

		if err = json.Unmarshal(jb, &ur); err != nil {
			goto end
		}

		if err = ur.Validate(); err != nil {
			goto end
		}

		if u, err = ctx.DalReadUser(ur.Token, KeyTypePassToken, true); err != nil {
			goto end
		}

		status = 500

		if passh, err = api.GetHash(ur.Password); err != nil {
			goto end
		}

		if err = ctx.DalUpdatePassh(u.ID, passh); err != nil {
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		_ = g.AbortWithError(status, err)
	}
}
