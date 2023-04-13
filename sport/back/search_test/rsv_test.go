package search_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"sport/adyen_sm"
// 	"sport/api"
// 	"sport/helpers"
// 	"sport/rsv"
// 	"sport/search"
// 	"sport/train"
// 	"testing"
// 	"time"
// )

// func RsvCountSchedValidator(
// 	title string, s time.Time, c int, trainings ...string,
// ) func(b []byte, i interface{}) error {
// 	return func(b []byte, i interface{}) error {
// 		var r search.SearchResultsWithMetadata
// 		if err := json.Unmarshal(b, &r); err != nil {
// 			return err
// 		}
// 		d := r.Data[0]
// 		if d.Schedule[0].Start != s {
// 			return fmt.Errorf("%s: unexpected schedule", title)
// 		}
// 		if d.Schedule[0].Count != c {
// 			return fmt.Errorf("%s: unexpected Count", title)
// 		}
// 		return mustMatchTrainings(title, trainings...)(b, i)
// 	}
// }

// func TestRsvCache(t *testing.T) {

// 	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	s := helpers.NowMin().Add(time.Hour * 48)
// 	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Minute))

// 	tr.Training.Title = "rsvTest"
// 	tr.Training.LocationCountry = "SL"
// 	tr.Training.Capacity = 100

// 	tid, err := trainCtx.ApiCreateTraining(token, &tr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := rsv.NewTestCreateReservationRequest(tid, s)
// 	r1, err := rsvCtx.ApiCreateAndQueryRsv(&token, &rr, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := rsvCtx.AdyenSm.MoveStateToCapture(
// 		rsvCtx.RsvResponseToSmPassPtr(r1),
// 		adyen_sm.ManualSrc, nil); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := searchCtx.RegenerateCache(); err != nil {
// 		t.Fatal(err)
// 	}

// 	searchCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod: "POST",
// 			RequestUrl:    searchPath,
// 			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
// 				SearchRequestOpts: search.SearchRequestOpts{
// 					Langs:     []string{"pl"},
// 					Country:   "SL",
// 					DateStart: s,
// 					DateEnd:   s.Add(time.Hour),
// 				},
// 			}),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal:    RsvCountSchedValidator("search 0", s, 1, tr.Training.Title),
// 		},
// 	})

// 	// l := searchCtx.NewListener()
// 	// defer func() {
// 	// 	if err := l.UnlistenAll(); err != nil {
// 	// 		fmt.Fprintln(os.Stderr, "couldnt unlisten")
// 	// 	}
// 	// }()

// 	rr2 := rsv.NewTestCreateReservationRequest(tid, s)
// 	r2, err := rsvCtx.ApiCreateAndQueryRsv(&token, &rr2, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := rsvCtx.AdyenSm.MoveStateToCapture(
// 		rsvCtx.RsvResponseToSmPassPtr(r2),
// 		adyen_sm.ManualSrc, nil); err != nil {
// 		t.Fatal(err)
// 	}

// 	searchCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod: "POST",
// 			RequestUrl:    searchPath,
// 			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
// 				SearchRequestOpts: search.SearchRequestOpts{
// 					Langs:     []string{"pl"},
// 					Country:   "SL",
// 					DateStart: s,
// 					DateEnd:   s.Add(time.Hour),
// 				},
// 			}),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal:    RsvCountSchedValidator("search 1", s, 1, tr.Training.Title),
// 		},
// 	})

// 	// q, err := GetListenerQueue(l)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	q, err := searchCtx.GetPgChanges()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := searchCtx.UpdateSearchCache(q); err != nil {
// 		t.Fatal(err)
// 	}

// 	searchCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod: "POST",
// 			RequestUrl:    searchPath,
// 			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
// 				SearchRequestOpts: search.SearchRequestOpts{
// 					Langs:     []string{"pl"},
// 					Country:   "SL",
// 					DateStart: s,
// 					DateEnd:   s.Add(time.Hour),
// 				},
// 			}),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal:    RsvCountSchedValidator("search 2", s, 2, tr.Training.Title),
// 		},
// 	})

// 	r2, err = rsvCtx.ReadRsvByID(r2.ID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := rsvCtx.AdyenSm.MoveStateToRefund(
// 		rsvCtx.RsvResponseToSmPassPtr(r2), adyen_sm.ManualSrc); err != nil {
// 		t.Fatal(err)
// 	}

// 	// q, err = GetListenerQueue(l)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	q, err = searchCtx.GetPgChanges()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := searchCtx.UpdateSearchCache(q); err != nil {
// 		t.Fatal(err)
// 	}

// 	searchCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod: "POST",
// 			RequestUrl:    searchPath,
// 			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
// 				SearchRequestOpts: search.SearchRequestOpts{
// 					Langs:     []string{"pl"},
// 					Country:   "SL",
// 					DateStart: s,
// 					DateEnd:   s.Add(time.Hour),
// 				},
// 			}),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal:    RsvCountSchedValidator("search 3", s, 1, tr.Training.Title),
// 		},
// 	})
// }
