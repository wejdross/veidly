package user_test

import (
	"sport/api"
	"sport/helpers"
	"sport/user"
	"testing"
)

func TestLang(t *testing.T) {
	{
		ur := &user.UserRequest{
			UserData: user.UserData{
				Name:     helpers.CRNG_stringPanic(user.UserNameMaxLength),
				Language: user.LangPl,
				Country:  "PL",
			},
			Password: "!aDB&!@#$#$dsfsdfjk98ASDASASD",
			Email:    helpers.CRNG_stringPanic(12),
		}

		testCases := []api.TestCase{
			{
				RequestMethod:      "POST",
				RequestUrl:         "/api/user/register",
				RequestReader:      helpers.JsonMustSerializeReader(ur),
				ExpectedStatusCode: 200,
			},
		}

		userCtx.Api.TestAssertCases(t, testCases)
	}

	{
		ur := &user.UserRequest{
			UserData: user.UserData{
				Name:     helpers.CRNG_stringPanic(user.UserNameMaxLength),
				Language: helpers.CRNG_stringPanic(12),
			},
			Password: "!aDB&!@#$#$dsfsdfjk98ASDASASD",
			Email:    helpers.CRNG_stringPanic(12),
		}

		testCases := []api.TestCase{
			{
				RequestMethod:      "POST",
				RequestUrl:         "/api/user/register",
				RequestReader:      helpers.JsonMustSerializeReader(ur),
				ExpectedStatusCode: 200,
			},
		}

		userCtx.Api.TestAssertCases(t, testCases)
	}
}
