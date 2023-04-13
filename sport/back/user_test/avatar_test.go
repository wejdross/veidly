package user_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sport/api"
	"sport/helpers"
	"sport/user"
	"testing"
)

func TestAvatar(t *testing.T) {
	token, err := userCtx.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}

	uid, err := userCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}

	url := "/api/user/avatar"
	method := "PUT"

	{
		ch := user.GetAuthHeader(token)

		fc := []byte(helpers.CRNG_stringPanic(100))

		ext := ".png"

		f, ct, err := api.CreateMultipartForm("image", "image"+ext, fc)
		if err != nil {
			t.Fatal(err)
		}
		ch["Content-Type"] = ct

		cases := []api.TestCase{
			{
				RequestMethod:      method,
				RequestUrl:         url,
				ExpectedStatusCode: 401,
			},
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestHeaders:     ch,
				ExpectedStatusCode: 400,
			},
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestHeaders:     ch,
				RequestReader:      &f,
				ExpectedStatusCode: 400,
			},
		}

		userCtx.Api.TestAssertCases(t, cases)
	}

	{
		ch := user.GetAuthHeader(token)

		fc, err := ioutil.ReadFile(helpers.TestFilePath)
		if err != nil {
			t.Fatal(err)
		}

		ext := ".jpg"

		f, ct, err := api.CreateMultipartForm("image", "image"+ext, fc)
		if err != nil {
			t.Fatal(err)
		}
		ch["Content-Type"] = ct

		cases := []api.TestCase{
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestHeaders:     ch,
				RequestReader:      &f,
				ExpectedStatusCode: 204,
			},
		}

		userCtx.Api.TestAssertCases(t, cases)
	}

	u, err := userCtx.DalReadUser(uid, user.KeyTypeID, true)
	if err != nil {
		t.Fatal(err)
	}

	ppath := path.Join(userCtx.Static.Config.Basepath, u.AvatarRelpath)

	_, err = os.Stat(ppath)
	if err != nil {
		t.Fatal(err)
	}

	cases := []api.TestCase{
		{
			RequestMethod:      "GET",
			RequestUrl:         u.AvatarUrl,
			ExpectedStatusCode: 200,
		},
	}

	userCtx.Api.TestAssertCases(t, cases)

	{
		ch := user.GetAuthHeader(token)

		fc, err := ioutil.ReadFile(helpers.TestFilePath)
		if err != nil {
			t.Fatal(err)
		}

		ext := ".jpg"

		f, ct, err := api.CreateMultipartForm("image", "image"+ext, fc)
		if err != nil {
			t.Fatal(err)
		}
		ch["Content-Type"] = ct

		cases := []api.TestCase{
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestHeaders:     ch,
				RequestReader:      &f,
				ExpectedStatusCode: 204,
			},
		}

		userCtx.Api.TestAssertCases(t, cases)
	}

	u, err = userCtx.DalReadUser(uid, user.KeyTypeID, true)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(ppath)
	if err == nil || !os.IsNotExist(err) {
		t.Fatal(fmt.Println("avatar should be gone by now"))
	}

	ppath = path.Join(userCtx.Static.Config.Basepath, u.AvatarRelpath)
	_, err = os.Stat(ppath)
	if err != nil {
		t.Fatal(err)
	}

	method = "DELETE"

	{
		cases := []api.TestCase{
			{
				RequestMethod:      method,
				RequestUrl:         url,
				ExpectedStatusCode: 401,
			},
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestHeaders:     user.GetAuthHeader(token),
				ExpectedStatusCode: 204,
			},
		}

		userCtx.Api.TestAssertCases(t, cases)
	}

	_, err = os.Stat(ppath)
	if err == nil || !os.IsNotExist(err) {
		t.Fatal(fmt.Println("avatar should be gone by now"))
	}

	u, err = userCtx.DalReadUser(uid, user.KeyTypeID, true)
	if err != nil {
		t.Fatal(err)
	}

	if u.AvatarRelpath != "" {
		t.Fatal("invalid relpath")
	}

}
