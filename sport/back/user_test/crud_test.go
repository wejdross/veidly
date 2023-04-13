package user_test

import (
	"encoding/json"
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/user"
	"testing"
)

func TestCrud(t *testing.T) {
	pass := "!aDB&!@#$#$dsfsdfjk98ASDASASD"
	ur := &user.UserRequest{
		UserData: user.UserData{
			Language: user.LangEn,
			Country:  "US",
		},
		Password: pass,
		Email:    helpers.CRNG_stringPanic(12),
	}

	var token string

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

	testCases = []api.TestCase{
		{
			RequestMethod: "GET",
			RequestUrl:    "/api/user",
			RequestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res user.User
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				if err := helpers.AssertJsonErr("user data does not match", res.UserData, ur.UserData); err != nil {
					return err
				}
				if err := helpers.AssertErr("Emails do not match", res.Email, ur.Email); err != nil {
					return err
				}
				return nil
			},
		},
	}

	userCtx.Api.TestAssertCases(t, testCases)

	ur.Language = user.LangPl

	testCases = []api.TestCase{
		{
			RequestMethod: "PATCH",
			RequestUrl:    "/api/user/data",
			RequestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			RequestReader:      helpers.JsonMustSerializeReader(ur),
			ExpectedStatusCode: 204,
		}, {
			RequestMethod: "GET",
			RequestUrl:    "/api/user",
			RequestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var u user.User
				if err := json.Unmarshal(b, &u); err != nil {
					return err
				}
				if err := helpers.AssertJsonErr("user data does not match", u.UserData, ur.UserData); err != nil {
					return err
				}
				if err := helpers.AssertErr("Emails do not match", u.Email, ur.Email); err != nil {
					return err
				}
				return nil
			},
		},
		{
			RequestMethod:      "DELETE",
			RequestUrl:         "/api/user",
			ExpectedStatusCode: 401,
		},
		{
			RequestMethod: "DELETE",
			RequestUrl:    "/api/user",
			RequestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			ExpectedStatusCode: 204,
		},
		{
			RequestMethod: "GET",
			RequestUrl:    "/api/user",
			RequestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			ExpectedStatusCode: 404,
		},
	}

	userCtx.Api.TestAssertCases(t, testCases)
}
