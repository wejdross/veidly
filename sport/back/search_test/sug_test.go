package search_test

import (
	"encoding/json"
	"sport/api"
	"sport/helpers"
	"sport/search"
	"sport/train"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSug(t *testing.T) {

	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	// iuid, err := trainCtx.Api.AuthorizeUserFromToken(token)
	// if err != nil {
	// 	t.Fatal(err)

	// }

	s := time.Now()
	trs := []train.CreateTrainingRequest{
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
	}

	trs[0].Training.Title = "sugTest0"
	trs[0].Training.LocationLat = 70
	trs[0].Training.LocationLng = 10

	trs[1].Training.Title = "sugTest1"
	trs[1].Training.LocationLat = 70.5
	trs[1].Training.LocationLng = 10

	trs[2].Training.Title = "sugTest2"
	trs[2].Training.LocationLat = 71
	trs[2].Training.LocationLng = 10

	trs[3].Training.Title = "sugTest3"
	trs[3].Training.LocationLat = 70
	trs[3].Training.LocationLng = 12

	tids := make([]uuid.UUID, len(trs))
	for i := range trs {
		_tid, err := trainCtx.ApiCreateTraining(token, &trs[i])
		if err != nil {
			t.Fatal(err)
		}
		tids[i] = _tid
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
					Langs:      []string{"pl"},
					Lat:        70,
					Lng:        10,
					DistKm:     80,
					SugEnabled: true,
					SugDistKm:  150,
				},
				Query: "sugTest",
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var r search.SearchResultsWithMetadata
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				if err := matchResults(r.Data, "search 0", "sugTest0", "sugTest1"); err != nil {
					return err
				}
				return matchResults(r.SugData, "search 0 sug", "sugTest2", "sugTest3")
			},
		},
	})
}
