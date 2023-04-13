package user_test

import (
	"encoding/json"
	"sport/api"
	"sport/helpers"
	"sport/user"
	"testing"
)

func TestContact(t *testing.T) {
	ur := user.ApiTestUserRequest()
	token, err := userCtx.ApiCreateAndLoginUser(&ur)
	if err != nil {
		t.Fatal(err)
	}

	m := "PATCH"
	u := "/api/user/contact"

	const email = "email@email.email"
	const email2 = "email2@email.email"
	const phone = "123 123 123 1323s"

	userCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		},
		// ensure that initial contact data is empty
		{
			RequestMethod:  "GET",
			RequestUrl:     "/api/user",
			RequestHeaders: user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res user.PrivUserInfo
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				expected := user.ContactData{
					Email: ur.Email,
				}
				return helpers.AssertJsonErr("init contact data", res.ContactData, expected)
			},
		},
		{
			RequestMethod:      m,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestUrl:         u,
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod:      m,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestUrl:         u,
			RequestReader:      helpers.JsonMustSerializeReader(user.ContactData{}),
			ExpectedStatusCode: 204,
		},
		{
			RequestMethod:  m,
			RequestHeaders: user.GetAuthHeader(token),
			RequestUrl:     u,
			RequestReader: helpers.JsonMustSerializeReader(user.ContactData{
				Email: "sadasdas",
			}),
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod:  m,
			RequestHeaders: user.GetAuthHeader(token),
			RequestUrl:     u,
			RequestReader: helpers.JsonMustSerializeReader(user.ContactData{
				Email: email,
			}),
			ExpectedStatusCode: 204,
		},
		{
			RequestMethod:  "GET",
			RequestUrl:     "/api/user",
			RequestHeaders: user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res user.PrivUserInfo
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				expected := user.ContactData{
					Email: email,
				}
				return helpers.AssertJsonErr("init contact data", res.ContactData, expected)
			},
		},
		{
			RequestMethod:  m,
			RequestHeaders: user.GetAuthHeader(token),
			RequestUrl:     u,
			RequestReader: helpers.JsonMustSerializeReader(map[string]string{
				"Phone": phone,
			}),
			ExpectedStatusCode: 204,
		},
		{
			RequestMethod:  "GET",
			RequestUrl:     "/api/user",
			RequestHeaders: user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res user.PrivUserInfo
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				expected := user.ContactData{
					Email: email,
					Phone: phone,
				}
				return helpers.AssertJsonErr("init contact data", res.ContactData, expected)
			},
		},
		{
			RequestMethod:  m,
			RequestHeaders: user.GetAuthHeader(token),
			RequestUrl:     u,
			RequestReader: helpers.JsonMustSerializeReader(map[string]string{
				"Email": email2,
			}),
			ExpectedStatusCode: 204,
		},
		{
			RequestMethod:  "GET",
			RequestUrl:     "/api/user",
			RequestHeaders: user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res user.PrivUserInfo
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				expected := user.ContactData{
					Email: email2,
					Phone: phone,
				}
				return helpers.AssertJsonErr("init contact data", res.ContactData, expected)
			},
		},
	})

}
