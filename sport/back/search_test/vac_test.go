package search_test

import (
	"encoding/json"
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"sport/search"
	"sport/train"
	"sport/user"
	"testing"
	"time"

	"github.com/google/uuid"
)

func validateVacations(b []byte, i interface{}, isAvailable bool, sn, title string) error {
	var r search.SearchResultsWithMetadata
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	t := sn
	if err := mustMatchTrainings(t, title)(b, i); err != nil {
		return err
	}
	s := r.Data[0].Schedule
	if len(s) != 1 {
		return fmt.Errorf("%s: invalid schedule", t)
	}
	if s[0].IsAvailable != isAvailable {
		return fmt.Errorf("%s: invalid sched availability", t)
	}
	return nil
}

func TestVacations(t *testing.T) {

	token, err := searchCtx.User.ApiCreateAndLoginUser(nil)
	uid, err := searchCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}
	_, err = searchCtx.Instr.CreateTestInstructor(uid, nil)
	if err != nil {
		t.Fatal(err)
	}

	country := "AX"
	// ensure training is available
	s := time.Now().Add(time.Hour*24 + 1)
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Minute))
	tr.Training.Title = "vactest1"
	tr.Occurrences[0].RepeatDays = 0
	title := tr.Training.Title
	tr.Training.LocationCountry = country
	_, err = trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}

	if err := searchCtx.RegenerateCache(); err != nil {
		t.Fatal(err)
	}

	sr := search.ApiSearchRequest{
		SearchRequestOpts: search.SearchRequestOpts{
			Langs:     []string{"pl"},
			Country:   country,
			DateStart: s,
			DateEnd:   s.Add(time.Hour * 24 * 30),
		},
	}

	var vid uuid.UUID

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         searchPath,
			RequestReader:      helpers.JsonMustSerializeReader(sr),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return validateVacations(b, i, true, "search0", title)
			},
		}, {
			RequestMethod:      "POST",
			RequestUrl:         "/api/instructor/vacation",
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 204,
			RequestReader: helpers.JsonMustSerializeReader(instr.VacationRequest{
				DateStart: s,
				DateEnd:   s.Add(time.Hour),
			}),
		},
		// add vacation
		{
			RequestMethod:      "POST",
			RequestUrl:         searchPath,
			RequestReader:      helpers.JsonMustSerializeReader(sr),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return validateVacations(b, i, true, "search1", title)
			},
		}, {
			RequestMethod:      "GET",
			RequestUrl:         "/api/instructor/vacation",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var got []instr.VacationInfo
				if err := json.Unmarshal(b, &got); err != nil {
					return err
				}
				if len(got) != 1 {
					return fmt.Errorf("invalid vacation")
				}
				vid = got[0].ID
				return nil
			},
		},
	})

	q, err := searchCtx.GetPgChanges()
	if err != nil {
		t.Fatal(err)
	}
	if err := searchCtx.UpdateSearchCache(q); err != nil {
		t.Fatal(err)
	}

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         searchPath,
			RequestReader:      helpers.JsonMustSerializeReader(sr),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return validateVacations(b, i, false, "search2", title)
			},
		},
	})

	// modify vacation
	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "PATCH",
			RequestUrl:         "/api/instructor/vacation",
			ExpectedStatusCode: 204,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.UpdateVacationRequest{
				ID: vid,
				VacationRequest: instr.VacationRequest{
					// ensure that this vacation starts AFTER training
					DateStart: s.Add(time.Hour * 1),
					DateEnd:   s.Add(time.Hour * 2),
				},
			}),
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
			RequestMethod:      "POST",
			RequestUrl:         searchPath,
			RequestReader:      helpers.JsonMustSerializeReader(sr),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return validateVacations(b, i, true, "search3", title)
			},
		},
		// restore
		{
			RequestMethod:      "PATCH",
			RequestUrl:         "/api/instructor/vacation",
			ExpectedStatusCode: 204,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.UpdateVacationRequest{
				ID: vid,
				VacationRequest: instr.VacationRequest{
					// ensure that this vacation starts AFTER training
					DateStart: s,
					DateEnd:   s.Add(time.Hour * 2),
				},
			}),
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
			RequestMethod:      "POST",
			RequestUrl:         searchPath,
			RequestReader:      helpers.JsonMustSerializeReader(sr),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return validateVacations(b, i, false, "search4", title)
			},
		}, {
			RequestMethod:      "DELETE",
			RequestUrl:         "/api/instructor/vacation",
			ExpectedStatusCode: 204,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.DeleteVacationRequest{
				ID: vid,
			}),
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
			RequestMethod:      "POST",
			RequestUrl:         searchPath,
			RequestReader:      helpers.JsonMustSerializeReader(sr),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return validateVacations(b, i, true, "search5", title)
			},
		},
	})
}
