package instr_test

import (
	"encoding/json"
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"sport/user"
	"testing"
)

func TestInfo(t *testing.T) {
	token, err := instrCtx.User.ApiCreateAndLoginUser(&user.UserRequest{
		Email:    helpers.CRNG_stringPanic(10) + "@test.test",
		Password: "!asdasdSADASDAS56456$%^$%^$%",
		UserData: user.UserData{
			Language: "pl",
			Country:  "PL",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	uid, err := instrCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}
	_, err = instrCtx.CreateTestInstructor(uid, nil)
	if err != nil {
		t.Fatal(err)
	}
	m := "GET"
	u := "/api/instructor/info"
	instrCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var r instr.InfoResponse
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				exp := instr.CS_NoName |
					instr.CS_NoAboutMe
				if r.ConfiguredState != exp {
					return fmt.Errorf("invalid configured state, expected: %v, got: %v",
						instr.StateToStr(exp),
						instr.StateToStr(r.ConfiguredState))
				}

				return nil
			},
		},
		//
		{
			RequestMethod:      "PATCH",
			RequestUrl:         "/api/user/data",
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 204,
			RequestReader: helpers.JsonMustSerializeReader(user.UserData{
				Name:     "foo",
				Language: "pl",
				Country:  "PL",
				AboutMe:  "about me",
			}),
		}, {
			RequestMethod:  "PATCH",
			RequestUrl:     "/api/instructor/payout",
			RequestHeaders: user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.CardInfo{
				CardNumber:  "123",
				ExpiryMonth: "123",
				ExpiryYear:  "2030",
				HolderName:  "kz",
				Cvc:         "123",
			}),
			ExpectedStatusCode: 204,
		},
	})

}
