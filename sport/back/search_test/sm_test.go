package search_test

import (
	"encoding/json"
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/search"
	"sport/sub"
	"sport/train"
	"sport/user"
	"strings"
	"testing"
	"time"
)

func SmSchedValidator(
	title string, smNames []string, trainings ...string,
) func(b []byte, i interface{}) error {
	return func(b []byte, i interface{}) error {
		var r search.SearchResultsWithMetadata
		if err := json.Unmarshal(b, &r); err != nil {
			return err
		}
		d := r.Data[0]
		foundSmNames := map[string]bool{}
		for i := range smNames {
			foundSmNames[smNames[i]] = false
		}
		for i := range d.Sms {
			n := d.Sms[i].Name
			if c, f := foundSmNames[n]; f {
				if c {
					return fmt.Errorf("%s: found duplicate sm", title)
				} else {
					foundSmNames[n] = true
				}
			} else {
				return fmt.Errorf("%s: found unexpected item: %s", title, n)
			}
		}

		var notFoundSmsB strings.Builder
		for x := range foundSmNames {
			if !foundSmNames[x] {
				notFoundSmsB.WriteString(x)
				notFoundSmsB.WriteString(" ")
			}
		}
		if notFoundSmsB.Len() > 0 {
			notFoundSms := notFoundSmsB.String()
			return fmt.Errorf("%s: didnt find specified sub models: %s", title, notFoundSms)
		}
		return mustMatchTrainings(title, trainings...)(b, i)
	}
}

func TestSubCache(t *testing.T) {

	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	country := "US"

	s := helpers.NowMin().Add(time.Hour * 48)
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Minute))

	tr.Training.Title = "rsvTest"
	tr.Training.LocationCountry = country
	tr.Training.Capacity = 100

	tid, err := trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}

	smReq := sub.NewTestSubModelRequest()
	smReq.Name = "sched1"
	smReq.AllTrainingsByDef = false
	smID, err := subCtx.ApiCreateSubModel(token, &smReq)
	if err != nil {
		t.Fatal(err)
	}

	if err := subCtx.ApiCreateSmBinding(token, &sub.SubModelBinding{
		SubModelID: smID,
		TrainingID: tid,
	}); err != nil {
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
			ExpectedBodyVal: SmSchedValidator("search 0", []string{
				"sched1",
			}, tr.Training.Title),
		},
	})

	// l := searchCtx.NewListener()
	// defer func() {
	// 	if err := l.UnlistenAll(); err != nil {
	// 		fmt.Fprintln(os.Stderr, "couldnt unlisten")
	// 	}
	// }()

	if err := subCtx.ApiRmSmBinding(token, &sub.SubModelBinding{
		SubModelID: smID,
		TrainingID: tid,
	}); err != nil {
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
			ExpectedBodyVal: SmSchedValidator("search 1", []string{
				"sched1",
			}, tr.Training.Title),
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
					Langs:     []string{"pl"},
					Country:   country,
					DateStart: s,
					DateEnd:   s.Add(time.Hour),
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    SmSchedValidator("search 2", []string{}, tr.Training.Title),
		},
	})

	if err := subCtx.ApiCreateSmBinding(token, &sub.SubModelBinding{
		SubModelID: smID,
		TrainingID: tid,
	}); err != nil {
		t.Fatal(err)
	}

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
					Country: country,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: SmSchedValidator("search 3", []string{
				"sched1",
			}, tr.Training.Title),
		},
	})

	smReq.Name = "sched1_updated"

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "PATCH",
			RequestUrl:         "/api/sub/model",
			ExpectedStatusCode: 204,
			RequestReader: helpers.JsonMustSerializeReader(sub.UpdateSubModelRequest{
				ID:              smID,
				SubModelRequest: smReq,
			}),
			RequestHeaders: user.GetAuthHeader(token),
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
					Country: country,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: SmSchedValidator("search 4", []string{
				"sched1_updated",
			}, tr.Training.Title),
		},
	})
}
