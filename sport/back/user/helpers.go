package user

import (
	"fmt"
	"io/ioutil"
	"path"
	"sport/api"
	"sport/helpers"
	"time"

	"github.com/google/uuid"
)

func GetAuthHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
}

func (ctx *Ctx) ApiUploadAvatarFromPath(t, p string) error {
	ch := GetAuthHeader(t)

	fc, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	f, ct, err := api.CreateMultipartForm("image", path.Base(p), fc)
	if err != nil {
		return err
	}
	ch["Content-Type"] = ct

	cases := []api.TestCase{
		{
			RequestMethod:      "PUT",
			RequestUrl:         "/api/user/avatar",
			RequestHeaders:     ch,
			RequestReader:      &f,
			ExpectedStatusCode: 204,
		},
	}

	return ctx.Api.TestAssertCasesErr(cases)
}

func (ctx *Ctx) ApiLoginUser(ur *UserRequest) (string, error) {
	var token string

	testCases := []api.TestCase{

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

	if err := ctx.Api.TestAssertCasesErr(testCases); err != nil {
		return "", err
	}

	return token, nil
}

func ApiTestUserRequest() UserRequest {
	return UserRequest{
		UserData: UserData{
			Language: LangEn,
			Country:  "US",
		},
		Password: "!aDB&!@#$#$dsfsdfjk98ASDASASD",
		Email:    helpers.CRNG_stringPanic(12),
	}
}

func (ctx *Ctx) ApiCreateAndLoginUser(ur *UserRequest) (string, error) {
	if ur == nil {
		_ur := ApiTestUserRequest()
		ur = &_ur
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

	if err := ctx.Api.TestAssertCasesErr(testCases); err != nil {
		return "", err
	}

	return token, nil
}

func (ctx *Ctx) CreateAndLoginUserWithID(id uuid.UUID) (bool, string, error) {

	ur := UserRequest{
		Email: "foobar@gmail.com",
		UserData: UserData{
			Name:     "foobar",
			Country:  "PL",
			Language: "pl",
		},
	}
	u, err := ur.NewUser(false)
	if err != nil {
		return false, "", err
	}
	u.ID = id
	createdUser := false
	if err := ctx.DalCreateUser(u); err != nil {
		if !helpers.PgIsUqViolation(err) {
			return false, "", err
		}
	} else {
		createdUser = true
	}
	tp := new(api.JwtPayload)
	tp.Exp = time.Now().In(time.UTC).Add(time.Duration(time.Hour * 48))
	tp.UserID = u.ID
	token, err := ctx.Api.Jwt.GenToken(tp)
	return createdUser, token, err
}
