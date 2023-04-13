package instr_test

import (
	"encoding/json"
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"sport/user"
	"strings"
	"testing"
	"time"
)

func TestCRUD(t *testing.T) {

	// create user
	var token string

	var pass = "!aDB&!@#$#$dsfsdfjk98ASDASASD"

	{
		ur := user.UserRequest{
			Email:    helpers.CRNG_stringPanic(12),
			Password: pass,
			UserData: user.UserData{
				Language: "en",
				Country:  "US",
			},
		}
		cases := []api.TestCase{
			{
				RequestMethod:      "POST",
				RequestUrl:         "/api/user/register",
				RequestReader:      helpers.JsonMustSerializeReader(&ur),
				ExpectedStatusCode: 200,
			},
			{
				RequestMethod:      "POST",
				RequestUrl:         "/api/user/login",
				RequestReader:      helpers.JsonMustSerializeReader(&ur),
				ExpectedStatusCode: 200,
				ExpectedBodyVal: func(b []byte, i interface{}) error {
					token = string(b)
					return nil
				},
			},
		}
		instrCtx.Api.TestAssertCases(t, cases)
	}

	// userID, err := global.Api.AuthorizeUserFromToken(token)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	thisYear := time.Now().In(time.UTC).Year()

	correctInstructorRequest := instr.InstructorRequest{
		Tags: []string{
			"abc",
			"def",
			"ghi",
		},
		YearExp: thisYear,
		ProfileSections: instr.ProfileSectionArr{
			instr.ProfileSection{
				Title:   "1",
				Content: "11",
			},
			instr.ProfileSection{
				Title:   "2",
				Content: "22",
			},
		},
		InvoiceLines: []string{
			"1",
			"2",
			"3",
			"4",
			"5",
		},
	}

	// CREATE
	{
		url := "/api/instructor"
		method := "POST"

		invalidRequest := instr.InstructorRequest{
			Tags: []string{
				"a",
				"b",
				"c",
				"d",
				"e",
				"f",
			},
			ProfileSections: instr.ProfileSectionArr{},
		}

		cases := []api.TestCase{
			// no auth
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      nil,
				ExpectedStatusCode: 401,
			},
			// correct request
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      helpers.JsonMustSerializeReader(&correctInstructorRequest),
				ExpectedStatusCode: 204,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
			},
			// fuzzing
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				RequestReader:      strings.NewReader(helpers.CRNG_stringPanic(128)),
				ExpectedStatusCode: 400,
			},
			// already exists
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      helpers.JsonMustSerializeReader(&correctInstructorRequest),
				ExpectedStatusCode: 409,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
			},
			// invalid request
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      helpers.JsonMustSerializeReader(&invalidRequest),
				ExpectedStatusCode: 400,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
			},
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestReader: helpers.JsonMustSerializeReader(&instr.InstructorRequest{
					YearExp: 30,
				}),
				ExpectedStatusCode: 400,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
			},
		}

		instrCtx.Api.TestAssertCases(t, cases)
	}

	// ENSURE CORRECT TAGS
	// {
	// 	tags, err := DalReadTags(InstructorUserID, userID)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if err := helpers.AssertErr("invalid tags length after create instructor", len(tags), len(correctInstructorRequest.Tags)); err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	for i := range tags {
	// 		if err := helpers.AssertErr("invalid tag after create instructor", tags[i].TagPlain, correctInstructorRequest.Tags[i]); err != nil {
	// 			t.Fatal(err)
	// 		}
	// 	}
	// }

	// READ
	{
		url := "/api/instructor"
		method := "GET"

		cases := []api.TestCase{
			// no auth
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      nil,
				ExpectedStatusCode: 404,
			},
			// correct
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedBodyVal: func(b []byte, i interface{}) error {
					var res instr.PubInstructorInfo
					if err := json.Unmarshal(b, &res); err != nil {
						return err
					}
					return helpers.AssertJsonErr(
						"instructor Get",
						correctInstructorRequest,
						res.InstructorRequest)
				},
				ExpectedStatusCode: 200,
			},
		}

		instrCtx.Api.TestAssertCases(t, cases)
	}

	patchRequest := instr.InstructorRequest{
		Tags: []string{
			"abc",
			"def",
			"ghi",
		},
		YearExp:         2011,
		ProfileSections: instr.ProfileSectionArr{},
		InvoiceLines: []string{
			"1",
			"2",
			"3",
		},
	}

	// UPDATE
	{
		url := "/api/instructor"
		method := "PATCH"

		cases := []api.TestCase{
			// no auth
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      nil,
				ExpectedStatusCode: 401,
			},
			// correct patch
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				RequestReader:      helpers.JsonMustSerializeReader(patchRequest),
				ExpectedStatusCode: 204,
			},
			// correct
			{
				RequestMethod: "GET",
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedBodyVal: func(b []byte, i interface{}) error {
					var res instr.PubInstructorInfo
					if err := json.Unmarshal(b, &res); err != nil {
						return err
					}
					return helpers.AssertJsonErr(
						"instructor Patch",
						patchRequest,
						res.InstructorRequest)
				},
				ExpectedStatusCode: 200,
			},
		}

		instrCtx.Api.TestAssertCases(t, cases)
	}

	// ENSURE CORRECT TAGS
	// {
	// 	tags, err := DalReadTags(InstructorUserID, userID)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if err := helpers.AssertErr("invalid tags length after patch instructor", len(tags), len(patchRequest.Tags)); err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	for i := range tags {
	// 		if err := helpers.AssertErr("invalid tag after patch instructor", tags[i].TagPlain, patchRequest.Tags[i]); err != nil {
	// 			t.Fatal(err)
	// 		}
	// 	}
	// }

	// DELETE
	{
		url := "/api/instructor"
		method := "DELETE"

		cases := []api.TestCase{
			// no auth
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      nil,
				ExpectedStatusCode: 401,
			},
			{
				RequestMethod:      "GET",
				RequestUrl:         "/api/instructor/can_delete",
				ExpectedStatusCode: 401,
			},
			// can delete
			{
				RequestMethod: "GET",
				RequestUrl:    "/api/instructor/can_delete",
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedStatusCode: 204,
			},
			// correct delete
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedStatusCode: 204,
			},
			// get
			{
				RequestMethod: "GET",
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedStatusCode: 404,
			},
			{
				RequestMethod: "GET",
				RequestUrl:    "/api/instructor/can_delete",
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedStatusCode: 404,
			},
		}

		instrCtx.Api.TestAssertCases(t, cases)
	}

	// ENSURE CORRECT TAGS
	// {
	// 	tags, err := DalReadTags(InstructorUserID, userID)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if err := helpers.AssertErr("invalid tags length after delete instructor", len(tags), 0); err != nil {
	// 		t.Fatal(err)
	// 	}
	// }

}
