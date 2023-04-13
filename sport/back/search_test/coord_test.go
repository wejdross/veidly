package search_test

import (
	"sport/api"
	"sport/helpers"
	"sport/search"
	"sport/train"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCoord(t *testing.T) {

	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	iuid, err := trainCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)

	}

	s := time.Now()
	trs := []train.CreateTrainingRequest{
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
	}

	trs[0].Training.Title = "coordTest0"
	trs[0].Training.LocationLat = 20
	trs[0].Training.LocationLng = 60

	trs[1].Training.Title = "coordTest1"
	trs[1].Training.LocationLat = 21
	trs[1].Training.LocationLng = 60.5

	trs[2].Training.Title = "coordTest2"
	trs[2].Training.LocationLat = 20.3
	trs[2].Training.LocationLng = 60.5

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

	// alter training one (lng out of range)
	trs[1].Training.LocationLng = 87

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:  []string{"pl"},
					Lat:    20.5,
					Lng:    60,
					DistKm: 150,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 0", "coordTest0", "coordTest1", "coordTest2"),
		}, {
			RequestMethod: "PATCH",
			RequestUrl:    "/api/training",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			RequestReader: helpers.JsonMustSerializeReader(train.UpdateTrainingRequest{
				TrainingRequest: trs[1].Training,
				ID:              tids[1],
			}),
			ExpectedStatusCode: 204,
		},
	})

	q, err := searchCtx.GetPgChanges()
	if err != nil {
		t.Fatal(err)
	}
	if err := searchCtx.UpdateSearchCache(q); err != nil {
		t.Fatal(err)
	}

	// reset training one (lng in range - also duplicate)
	trs[1].Training.LocationLng = 60.5

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:  []string{"pl"},
					Lat:    20.5,
					Lng:    60,
					DistKm: 150,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 1", "coordTest0", "coordTest2"),
		}, {
			RequestMethod: "PATCH",
			RequestUrl:    "/api/training",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			RequestReader: helpers.JsonMustSerializeReader(train.UpdateTrainingRequest{
				TrainingRequest: trs[1].Training,
				ID:              tids[1],
			}),
			ExpectedStatusCode: 204,
		},
	})

	q, err = searchCtx.GetPgChanges()
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
					Langs:  []string{"pl"},
					Lat:    20.5,
					Lng:    60,
					DistKm: 150,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 2", "coordTest0", "coordTest1", "coordTest2"),
		},
	})

	// delete 2 trainings
	if err := searchCtx.Train.DalDeleteTraining(tids[1], iuid); err != nil {
		t.Fatal(err)
	}
	if err := searchCtx.Train.DalDeleteTraining(tids[0], iuid); err != nil {
		t.Fatal(err)
	}

	q, err = searchCtx.GetPgChanges()
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
					Langs:  []string{"pl"},
					Lat:    20.5,
					Lng:    60,
					DistKm: 150,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 3", "coordTest2"),
		},
	})

	// delete last training
	if err := searchCtx.Train.DalDeleteTraining(tids[2], iuid); err != nil {
		t.Fatal(err)
	}

	q, err = searchCtx.GetPgChanges()
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
					Langs:  []string{"pl"},
					Lat:    20.5,
					Lng:    60,
					DistKm: 150,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 4"),
		},
	})
}
