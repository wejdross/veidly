package search_test

import (
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/search"
	"sport/train"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTrainingCache(t *testing.T) {
	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	// create training
	s := time.Now()
	trs := []train.CreateTrainingRequest{
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
	}

	for i := range trs {
		trs[i].Training.Title = fmt.Sprintf("cacheTest%d", i)
		trs[i].Training.LocationCountry = "FR"
	}

	//ts := make(search.TestListenerMap)
	tids := make([]uuid.UUID, len(trs))

	for i := range trs {
		_tid, err := trainCtx.ApiCreateTraining(token, &trs[i])
		if err != nil {
			t.Fatal(err)
		}
		// ts[search.TestListenerMapKey{
		// 	ID:      _tid,
		// 	Channel: search.Trainings,
		// }] = 1
		// ts[search.TestListenerMapKey{
		// 	ID:      _tid,
		// 	Channel: search.Occurrences,
		// }] = 1
		tids[i] = _tid
	}

	if err := searchCtx.RegenerateCache(); err != nil {
		t.Fatal(err)
	}

	trs[1].Training.Title = "modified training"

	// l := searchCtx.NewListener()
	// defer func() {
	// 	if err := l.UnlistenAll(); err != nil {
	// 		fmt.Fprintln(os.Stderr, "couldnt unlisten")
	// 	}
	// }()

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{

					Langs:   []string{"pl"},
					Country: "FR",
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 0", "cacheTest0", "cacheTest1"),
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
		}, {
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{

					Langs:   []string{"pl"},
					Country: "FR",
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 1", "cacheTest0", "cacheTest1"),
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
					Country: "FR",
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 2", "cacheTest0", "modified training"),
		}, {
			RequestMethod: "DELETE",
			RequestUrl:    "/api/training",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			RequestReader: helpers.JsonMustSerializeReader(train.ObjectKey{
				ID: tids[0],
			}),
			ExpectedStatusCode: 204,
		},
	})

	// q, err = GetListenerQueue(l)
	// if err != nil {
	// 	t.Fatal(err)
	// }
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

					Langs:   []string{"pl"},
					Country: "FR",
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 3", "modified training"),
		},
	})

	trs[1].Training.Title = "recreated training"
	_tid, err := trainCtx.ApiCreateTraining(token, &trs[1])
	if err != nil {
		t.Fatal(err)
	}
	tids[1] = _tid

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{

					Langs:   []string{"pl"},
					Country: "FR",
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 4", "modified training"),
		},
	})

	// q, err = GetListenerQueue(l)
	// if err != nil {
	// 	t.Fatal(err)
	// }

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

					Langs:   []string{"pl"},
					Country: "FR",
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 5", "modified training", "recreated training"),
		},
	})
}
