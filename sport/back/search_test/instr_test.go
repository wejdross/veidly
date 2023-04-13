package search_test

import (
	"encoding/json"
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"sport/search"
	"sport/train"
	"testing"
	"time"
)

func TestInstr(t *testing.T) {

	token, err := searchCtx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}
	uid, err := searchCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}
	ir := &instr.InstructorRequest{
		Tags:            []string{"instrTag1", "instrTag2"},
		YearExp:         1997,
		ProfileSections: instr.ProfileSectionArr{},
	}
	_, err = searchCtx.Instr.CreateTestInstructor(uid, ir)
	if err != nil {
		t.Fatal(err)
	}

	country := "AU"
	s := time.Now()
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Minute))
	tr.Training.Title = "instrTest1"
	tr.Training.LocationCountry = country
	_, err = trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}

	if err := searchCtx.RegenerateCache(); err != nil {
		t.Fatal(err)
	}

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"pl"},
					Country: country,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				if err := mustMatchTrainings("search 0", "instrTest1")(b, i); err != nil {
					return err
				}
				var r search.SearchResultsWithMetadata
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				if err := helpers.AssertJsonErr(
					"search 0 - instr",
					r.Data[0].Instructor.InstructorRequest, ir); err != nil {
					return err
				}
				return nil
			},
		},
	})

	ir.Tags = []string{"instrTag3", "instrTag4"}
	ir.YearExp = 2012

	// l := searchCtx.NewListener()
	// defer func() {
	// 	if err := l.UnlistenAll(); err != nil {
	// 		fmt.Fprintln(os.Stderr, "couldnt unlisten")
	// 	}
	// }()

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "PATCH",
			RequestUrl:    "/api/instructor",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			RequestReader:      helpers.JsonMustSerializeReader(ir),
			ExpectedStatusCode: 204,
		}, {
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"pl"},
					Country: country,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				if err := mustMatchTrainings("search 1", "instrTest1")(b, i); err != nil {
					return err
				}
				var r search.SearchResultsWithMetadata
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				if err := helpers.AssertJsonErr(
					"search 1 - instr",
					r.Data[0].Instructor.InstructorRequest, ir); err == nil {
					return err
				}
				return nil
			},
		},
	})

	// q, err := GetListenerQueue(l)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	q, err := searchCtx.GetPgChanges()
	if err != nil {
		t.Fatal(err)
	}
	if err := searchCtx.UpdateSearchCache(q); err != nil {
		t.Fatal(err)
	}

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"pl"},
					Country: country,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				if err := mustMatchTrainings("search 2", "instrTest1")(b, i); err != nil {
					return err
				}
				var r search.SearchResultsWithMetadata
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				if err := helpers.AssertJsonErr(
					"search 2 - instr",
					r.Data[0].Instructor.InstructorRequest, ir); err != nil {
					return err
				}
				return nil
			},
		},
	})

}
