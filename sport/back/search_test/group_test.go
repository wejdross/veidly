package search_test

import (
	"encoding/json"
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/search"
	"sport/train"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {

	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	uid, err := trainCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}

	country := "CA"

	s := helpers.NowMin().Add(time.Hour * 48)
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Minute))

	tr.Training.Title = "rsvTest"
	tr.Training.LocationCountry = country
	tr.Training.Capacity = 100

	tid, err := trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}

	grReq := train.NewTestGroupRequest()
	gid, err := trainCtx.ApiCreateGroup(token, &grReq)
	if err != nil {
		t.Fatal(err)
	}

	if err := trainCtx.DalAddTvgBinding(tid, gid, uid); err != nil {
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
				title := "search 0"

				var r search.SearchResultsWithMetadata
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}

				if len(r.Data) != 1 {
					return fmt.Errorf(title + ": invalid search response 1")
				}

				d := r.Data[0]

				if len(d.Groups) != 1 {
					return fmt.Errorf(title + ": invalid search response 2")
				}

				if err := helpers.AssertJsonErr(title+" grp", d.Groups[0].GroupRequest, grReq); err != nil {
					return err
				}

				return mustMatchTrainings(title, tr.Training.Title)(b, i)
			},
		},
	})
}
