package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sport/api"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

/*
GET /api/user/oauth/url/:provider

DESC
	- generate redirect url to Oauth provider

REQUEST
	- url
		- provider, required can be one of supported providers:
			- google
	- query string
		- mode, optional : can be:
			- 'redirect' to generate response 307 redirect on success
			- 'raw' to return url (to avoid issues related to cors)
		  by default 'redirect' is assumed

RETURNS:
	- temporary redirect to oauth provider on success
	- status code

*/
func (h *Ctx) OauthUrlHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var status = 401
		var oauthProvider = ctx.Param("provider")
		var config *oauth2.Config
		var retUrl string

		if !h.Config.UseOauth {
			status = 501
			err = fmt.Errorf("User requested Oauth redirect without Oauth being enabled on server")
			goto end
		}

		switch oauthProvider {
		case OauthProviderGoogle:
			config = &h.Oauth.GoogleConfig
		default:
			// oauth provider should be validated by now
			status = 501
			err = fmt.Errorf("User requested Oauth provider \"%s\" which is not supported", oauthProvider)
			goto end
		}

		retUrl = OauthGetProviderAuthUrl(config, ctx.Query("return_url"))

		switch ctx.Query("mode") {
		case "raw":
			ctx.String(200, retUrl)
		default:
			ctx.Redirect(307, retUrl)
		}

		return

	end:
		_ = ctx.AbortWithError(status, err)
	}
}

type OauthCodeRequest struct {
	Code string
}

/*
GET /api/user/oauth/code/:provider

DESC:
	login or register user using Oauth

REQUEST
	- query string
		- provider can be one of supported providers:
			- google
	- body: OauthCodeRequest
RETURNS:
	- 200 with text/plain jwt authorization token
	- error status code

*/
func (ctx *Ctx) OauthCodeHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		var status = 401
		var provider string
		var req OauthCodeRequest
		var oauthToken *oauth2.Token
		var ui *OauthUserInfo
		var u *User
		var jb []byte
		var tp *api.JwtPayload
		var token string

		if !ctx.Config.UseOauth {
			status = 501
			err = fmt.Errorf("User requested Oauth login " +
				"without Oauth being enabled on the server")
			goto end
		}

		if jb, err = ioutil.ReadAll(g.Request.Body); err != nil {
			goto end
		}

		if err = json.Unmarshal(jb, &req); err != nil {
			goto end
		}

		if req.Code == "" {
			err = fmt.Errorf(
				"Could not find code in oauth response. Reason: \"%s\"",
				g.Request.FormValue("error_reason"))
			goto end
		}

		provider = g.Param("provider")

		switch provider {
		case OauthProviderGoogle:
			if oauthToken, err = ctx.Oauth.GoogleConfig.Exchange(
				oauth2.NoContext, req.Code,
			); err != nil {
				goto end
			}
			if ui, err = OauthGoogleGetUserInfo(oauthToken); err != nil {
				goto end
			}
		default:
			err = fmt.Errorf("Invalid provider: %s", provider)
			status = 400
			goto end
		}

		if err = ui.Validate(); err != nil {
			goto end
		}

		//UserCheck:

		// if u, err = ctx.DalReadUser(ui, KeyTypeOauthUserInfo, false); err != nil {
		// 	if helpers.IsENF(err) {
		// 		u = ui.NewUser()
		// 		if err = u.UserData.Validate(ctx.Config); err != nil {
		// 			goto end
		// 		}
		// 		if err = ctx.DalCreateUser(u); err != nil {
		// 			goto end
		// 		}
		// 	} else {
		// 		goto end
		// 	}
		// }

		u, err = ctx.DalReadUser(ui.Email, KeyTypeEmail, false)
		if err != nil {
			if err == sql.ErrNoRows {
				u = ui.NewUser()
				if err = u.UserData.Validate(ctx, false); err != nil {
					goto end
				}
				if err = ctx.DalCreateUser(u); err != nil {
					goto end
				}
			} else {
				goto end
			}
		}

		if u.OauthProvider != ui.Provider {
			status = 401
			err = fmt.Errorf("provider does not match")
			goto end
		}

		if u.OauthID != ui.OauthID {
			status = 401
			err = fmt.Errorf("oauth id does not match")
			goto end
		}

		if !u.Enabled {
			status = 401
			err = fmt.Errorf("user is disabled")
			goto end
		}

		tp = new(api.JwtPayload)
		tp.Exp = time.Now().In(time.UTC).Add(time.Duration(time.Hour * 480))
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
