package user_test

import (
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/user"
	"testing"
)

func TestPatchPassword(t *testing.T) {
	// create user
	var token string
	ur := &user.UserRequest{
		UserData: user.UserData{
			Name:     helpers.CRNG_stringPanic(user.UserNameMaxLength),
			Language: user.LangPl,
			Country:  "PL",
		},
		Password: "!aDB&!@#$#$dsfsdfjk98ASDASASD",
		Email:    helpers.CRNG_stringPanic(12),
	}

	{
		testCases := []api.TestCase{
			{
				RequestMethod:      "POST",
				RequestUrl:         "/api/user/register",
				RequestReader:      helpers.JsonMustSerializeReader(ur),
				ExpectedStatusCode: 200,
			},
			{
				RequestMethod:      "POST",
				RequestUrl:         "/api/user/login",
				ExpectedStatusCode: 200,
				RequestReader:      helpers.JsonMustSerializeReader(ur),
				ExpectedBodyVal: func(b []byte, i interface{}) error {
					token = string(b)
					return nil
				},
			},
		}

		userCtx.Api.TestAssertCases(t, testCases)
	}

	// patch password

	pp := user.PatchPasswordRequest{
		OldPassword: ur.Password,
		NewPassword: "!aDB&!@#$#$11111dsfsdfjk98ASDASASD",
	}

	lr := user.LoginRequest{
		Email:    ur.Email,
		Password: pp.NewPassword,
	}

	olr := user.LoginRequest{
		Email: pp.OldPassword,
	}

	testCases := []api.TestCase{
		// no auth must result in 401
		{
			RequestMethod:      "PATCH",
			RequestUrl:         "/api/user/password",
			RequestReader:      helpers.JsonMustSerializeReader(pp),
			ExpectedStatusCode: 401,
		},
		// correct request must result in 204
		{
			RequestMethod: "PATCH",
			RequestUrl:    "/api/user/password",
			RequestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			RequestReader:      helpers.JsonMustSerializeReader(pp),
			ExpectedStatusCode: 204,
		},
		// providing invalid old password must result in 401
		{
			RequestMethod: "PATCH",
			RequestUrl:    "/api/user/password",
			RequestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			RequestReader:      helpers.JsonMustSerializeReader(pp),
			ExpectedStatusCode: 401,
		},
		// try to login with new password
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/user/login",
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(lr),
		},
		// old password shouldnt work any more
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/user/login",
			ExpectedStatusCode: 401,
			RequestReader:      helpers.JsonMustSerializeReader(olr),
		},
	}

	userCtx.Api.TestAssertCases(t, testCases)
}
