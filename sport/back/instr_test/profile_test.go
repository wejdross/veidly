package instr_test

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"sport/user"
	"testing"
)

func reqFromFile(token string, fc []byte) (io.Reader, map[string]string, error) {
	ch := user.GetAuthHeader(token)
	ext := ".png"
	f, ct, err := api.CreateMultipartForm("image", "image"+ext, fc)
	if err != nil {
		return nil, nil, err
	}
	ch["Content-Type"] = ct
	return &f, ch, nil
}

func TestProfileImg(t *testing.T) {
	token, iid, err := instrCtx.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	url := "/api/instructor/profile/img"
	method := "POST"

	fc := []byte(helpers.CRNG_stringPanic(100))
	f, ch, err := reqFromFile(token, fc)
	if err != nil {
		t.Fatal(err)
	}

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
			RequestReader:      f,
			ExpectedStatusCode: 400,
		},
	}

	instrCtx.Api.TestAssertCases(t, cases)

	fc, err = ioutil.ReadFile(helpers.TestFilePath)
	if err != nil {
		t.Fatal(err)
	}

	f, ch, err = reqFromFile(token, fc)
	if err != nil {
		t.Fatal(err)
	}

	cases = []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         url + "?primary=1",
			RequestHeaders:     ch,
			RequestReader:      f,
			ExpectedStatusCode: 204,
		},
	}
	instrCtx.Api.TestAssertCases(t, cases)

	i, err := instrCtx.DalReadInstructor(iid, instr.InstructorID)
	if err != nil {
		t.Fatal(err)
	}

	if i.BgImgPath == "" {
		t.Fatal("no primary img")
	}

	ppath := path.Join(instrCtx.Static.Config.Basepath, i.BgImgPath)
	if _, err = os.Stat(ppath); err != nil {
		t.Fatal(err)
	}

	f, ch, err = reqFromFile(token, fc)
	if err != nil {
		t.Fatal(err)
	}

	cases = []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         url + "?primary=1",
			RequestHeaders:     ch,
			RequestReader:      f,
			ExpectedStatusCode: 204,
		},
	}
	instrCtx.Api.TestAssertCases(t, cases)

	path.Join(instrCtx.Static.Config.Basepath, i.BgImgPath)
	if _, err = os.Stat(ppath); err == nil {
		t.Fatal("1: instr profile img should be gone by now")
	}

	i, err = instrCtx.DalReadInstructor(iid, instr.InstructorID)
	if err != nil {
		t.Fatal(err)
	}

	ppath = path.Join(instrCtx.Static.Config.Basepath, i.BgImgPath)
	if _, err = os.Stat(ppath); err != nil {
		t.Fatal(err)
	}

	cases = []api.TestCase{
		{
			RequestMethod:  "DELETE",
			RequestUrl:     url + "?primary=1",
			RequestHeaders: user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.DeleteProfileImgRequest{
				Path: i.BgImgPath,
			}),
			ExpectedStatusCode: 204,
		},
	}
	instrCtx.Api.TestAssertCases(t, cases)

	ppath = path.Join(instrCtx.Static.Config.Basepath, i.BgImgPath)
	if _, err = os.Stat(ppath); err == nil {
		t.Fatal("2: instr profile img should be gone by now")
	}

	i, err = instrCtx.DalReadInstructor(iid, instr.InstructorID)
	if err != nil {
		t.Fatal(err)
	}

	if i.BgImgPath != "" {
		t.Fatal("3: instr profile img should be gone by now")
	}
}

func TestProfileExtraImg(t *testing.T) {

	token, iid, err := instrCtx.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	url := "/api/instructor/profile/img"
	method := "POST"

	fc, err := ioutil.ReadFile(helpers.TestFilePath)
	if err != nil {
		t.Fatal(err)
	}

	c := instr.MaxExtraProfileImgs

	for i := 0; i < c; i++ {
		f, ch, err := reqFromFile(token, fc)
		if err != nil {
			t.Fatal(err)
		}
		cases := []api.TestCase{
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestHeaders:     ch,
				RequestReader:      f,
				ExpectedStatusCode: 204,
			},
		}
		instrCtx.Api.TestAssertCases(t, cases)
	}

	instructor, err := instrCtx.DalReadInstructor(iid, instr.InstructorID)
	if err != nil {
		t.Fatal(err)
	}

	if len(instructor.ExtraImgPaths) != c {
		t.Fatal("1: invalid paths")
	}

	for _, ep := range instructor.ExtraImgPaths {
		ppath := path.Join(instrCtx.Static.Config.Basepath, ep)
		if _, err = os.Stat(ppath); err != nil {
			t.Fatal(err)
		}
	}

	f, ch, err := reqFromFile(token, fc)
	if err != nil {
		t.Fatal(err)
	}
	cases := []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         url,
			RequestHeaders:     ch,
			RequestReader:      f,
			ExpectedStatusCode: 413,
		},
	}
	instrCtx.Api.TestAssertCases(t, cases)

	for i := 0; i < c; i++ {
		cases := []api.TestCase{
			{
				RequestMethod:  "DELETE",
				RequestUrl:     url,
				RequestHeaders: user.GetAuthHeader(token),
				RequestReader: helpers.JsonMustSerializeReader(instr.DeleteProfileImgRequest{
					Path: instructor.ExtraImgPaths[i],
				}),
				ExpectedStatusCode: 204,
			},
		}
		instrCtx.Api.TestAssertCases(t, cases)
	}

	oldpaths := instructor.ExtraImgPaths

	instructor, err = instrCtx.DalReadInstructor(iid, instr.InstructorID)
	if err != nil {
		t.Fatal(err)
	}

	if len(instructor.ExtraImgPaths) != 0 {
		t.Fatal("2: invalid paths")
	}

	for i := range oldpaths {
		ppath := path.Join(instrCtx.Static.Config.Basepath, oldpaths[i])
		if _, err = os.Stat(ppath); err == nil {
			t.Fatal("img should be gone now")
		}
	}
}
