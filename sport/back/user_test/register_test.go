package user_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"sport/api"
	"sport/helpers"
	"sport/user"
	"strconv"
	"strings"
	"testing"
)

/*
	test register without 2fa or oauth enabled
*/
func TestSimpleRegister(t *testing.T) {
	method := "POST"
	href := "/api/user/register"
	var cases []api.TestCase

	/*providing invalid json must result in 401 */
	cases = append(cases, api.TestCase{
		RequestMethod:      method,
		RequestUrl:         href,
		RequestReader:      strings.NewReader(helpers.CRNG_stringPanic(256)),
		ExpectedStatusCode: 401,
	})

	ur := &user.UserRequest{
		UserData: user.UserData{
			Name:     helpers.CRNG_stringPanic(user.UserNameMaxLength),
			Language: "en",
			Country:  "US",
		},
		Password: "!aDB&!@#$#$dsfsdfjk98ASDASASD",
		Email:    helpers.CRNG_stringPanic(user.UserNameMaxLength),
	}

	/* providing correct data must result in 200 */
	cases = append(cases, api.TestCase{
		RequestMethod:      method,
		RequestUrl:         href,
		RequestReader:      helpers.JsonMustSerializeReader(ur),
		ExpectedStatusCode: 200,
	})

	tooLongUr := &user.UserRequest{
		UserData: user.UserData{
			Name:     helpers.CRNG_stringPanic(user.UserNameMaxLength + 1),
			Language: "en",
			Country:  "US",
		},
		Password: "!aDB&!@#$#$dsfsdfjk98ASDASASD",
		Email:    helpers.CRNG_stringPanic(user.UserNameMaxLength),
	}

	/* providing too long username must result in 401 */
	cases = append(cases, api.TestCase{
		RequestMethod:      method,
		RequestUrl:         href,
		RequestReader:      helpers.JsonMustSerializeReader(tooLongUr),
		ExpectedStatusCode: 400,
	})

	emptyPassUr := &user.UserRequest{
		UserData: user.UserData{
			Name:     helpers.CRNG_stringPanic(user.UserNameMaxLength),
			Language: "en",
			Country:  "US",
		},
		Password: "",
		Email:    helpers.CRNG_stringPanic(user.UserNameMaxLength),
	}

	/* empty password in request must result in 401 */
	cases = append(cases, api.TestCase{
		RequestMethod:      method,
		RequestUrl:         href,
		RequestReader:      helpers.JsonMustSerializeReader(emptyPassUr),
		ExpectedStatusCode: 400,
	})

	userCtx.Api.TestAssertCases(t, cases)
}

func mfaRegisterTest() error {
	method := "POST"
	href := "/api/user/register"
	var cases []api.TestCase

	ur := &user.UserRequest{
		UserData: user.UserData{
			Name:     helpers.CRNG_stringPanic(user.UserNameMaxLength),
			Language: "en",
			Country:  "US",
		},
		Password: "!aDB&!@#$#$dsfsdfjk98ASDASASD",
		Email:    helpers.CRNG_stringPanic(user.UserNameMaxLength),
	}

	returnUrl := helpers.CRNG_stringPanic(12)
	token := ""

	cases = append(cases, api.TestCase{
		RequestMethod:      method,
		RequestUrl:         href + "?mfa_return_url=" + returnUrl,
		RequestReader:      helpers.JsonMustSerializeReader(&ur),
		ExpectedStatusCode: 200,
		ExpectedBodyVal: func(b []byte, i interface{}) error {
			var res user.MockupRegisterResponse
			if err := json.Unmarshal(b, &res); err != nil {
				return err
			}
			hr := res.HrefFromEmail
			res.HrefFromEmail = ""
			if err := helpers.AssertJsonErr(
				"mfa register: invalid HandlerRegisterResponse",
				res,
				user.MockupRegisterResponse{
					HrefFromEmail: "",
					HandlerRegisterResponse: user.HandlerRegisterResponse{
						MFA:        true,
						TTLSeconds: userCtx.Config.RetentionLoopRecordTTLseconds,
					},
				}); err != nil {
				return err
			}
			u, err := url.Parse(hr)
			if err != nil {
				return err
			}
			if err := helpers.AssertErr(
				"mfa register: invalid path in return url",
				u.Path,
				path.Join(api.ApiBaseHref, user.ConfirmRegisterPath),
			); err != nil {
				return err
			}
			q := u.Query()
			if err := helpers.AssertErr(
				"mfa register: invalid mfa_return_url",
				q.Get("return_url"),
				returnUrl,
			); err != nil {
				return err
			}
			token = q.Get("token")
			if token == "" {
				return fmt.Errorf("mfa register: empty token")
			}
			return nil
		},
	})

	// login mustnt work
	cases = append(cases, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/user/login",
		RequestReader: helpers.JsonMustSerializeReader(user.LoginRequest{
			Email:    ur.Email,
			Password: ur.Password,
		}),
		ExpectedStatusCode: 401,
	})

	cases = append(cases, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/user/login",
		RequestReader: helpers.JsonMustSerializeReader(user.LoginRequest{
			Email:    ur.Email,
			Password: ur.Password,
		}),
		ExpectedStatusCode: 401,
	})

	// resend register email must work
	cases = append(cases, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/user/register/resend",
		RequestReader: helpers.JsonMustSerializeReader(user.ResendRegisterEmailRequest{
			Email: ur.Email,
		}),
		ExpectedStatusCode: 200,
	})

	if err := userCtx.Api.TestAssertCasesErr(cases); err != nil {
		return err
	}

	cases = []api.TestCase{}

	// random confirm results in 401
	cases = append(cases, api.TestCase{
		RequestMethod:      "GET",
		RequestUrl:         "/api" + user.ConfirmRegisterPath + "?token=" + helpers.CRNG_stringPanic(128),
		ExpectedStatusCode: 307,
		ExpectedHeaders: map[string][]string{
			"Location": {fmt.Sprintf(
				userCtx.Config.RegisterRedirectURL,
				"",
				strconv.Itoa(401),
				"",
			)},
		},
	})

	// confirm email must work
	cases = append(cases, api.TestCase{
		RequestMethod:      "GET",
		RequestUrl:         "/api" + user.ConfirmRegisterPath + "?token=" + token,
		ExpectedStatusCode: 307,
		ExpectedHeaders: map[string][]string{
			"Location": {fmt.Sprintf(
				userCtx.Config.RegisterRedirectURL,
				url.QueryEscape(ur.Email),
				strconv.Itoa(200),
				"",
			)},
		},
	})

	// confirming on enabled user must result in forbidden
	cases = append(cases, api.TestCase{
		RequestMethod:      "GET",
		RequestUrl:         "/api" + user.ConfirmRegisterPath + "?token=" + token,
		ExpectedStatusCode: 307,
		ExpectedHeaders: map[string][]string{
			"Location": {fmt.Sprintf(
				userCtx.Config.RegisterRedirectURL,
				"",
				strconv.Itoa(401),
				"",
			)},
		},
	})

	// login must work now
	cases = append(cases, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/user/login",
		RequestReader: helpers.JsonMustSerializeReader(user.LoginRequest{
			Email:    ur.Email,
			Password: ur.Password,
		}),
		ExpectedStatusCode: 200,
	})

	return userCtx.Api.TestAssertCasesErr(cases)
}

/*
	test register without 2fa or oauth enabled
*/
func Test2FARegister(t *testing.T) {

	var e = userCtx.Email
	if err := userCtx.EnableMockup2FA(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := userCtx.DisableMockup2FA(e)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// simple test still has to work
	TestSimpleRegister(t)

	if err := mfaRegisterTest(); err != nil {
		t.Fatal(err)
	}
}
