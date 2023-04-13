package train_test

import (
	"encoding/json"
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"sport/train"
	"sport/user"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func AssertGetTrainingResponse(res []*train.TrainingWithJoins, req []*train.CreateTrainingRequest) error {
	if err := helpers.AssertErr(
		"instructor training read len", len(res), len(req),
	); err != nil {
		return err
	}
	for j := 0; j < len(req); j++ {
		var r *train.TrainingWithJoins
		foundRes := false
		for jj := 0; jj < len(res); jj++ {
			if err := helpers.AssertJsonErr(
				"",
				res[jj].Training.TrainingRequest,
				req[j].Training,
			); err == nil {
				r = res[jj]
				foundRes = true
				break
			}
		}
		if !foundRes {
			return fmt.Errorf(
				"Couldnt find match for training: \n%s\nIn response:\n%s",
				helpers.JsonMustSerializeFormatStr(req[j]),
				helpers.JsonMustSerializeFormatStr(res))
		}
		if err := helpers.AssertErr(
			"occurrences length",
			len(req[j].Occurrences), len(r.Occurrences),
		); err != nil {
			return err
		}
		for x := 0; x < len(req[j].Occurrences); x++ {
			var foundOcc *train.OccWithSecondary
			for xx := 0; xx < len(r.Occurrences); xx++ {
				if err := helpers.AssertJsonErr(
					"",
					r.Occurrences[xx].OccRequest,
					req[j].Occurrences[x].OccRequest,
				); err == nil {
					foundOcc = &r.Occurrences[xx]
					break
				}
			}
			if foundOcc == nil {
				return fmt.Errorf(
					"Couldnt find match for occurrency: \n%s\nIn response:\n%s",
					helpers.JsonMustSerializeFormatStr(req[j].Occurrences[x]),
					helpers.JsonMustSerializeFormatStr(r.Occurrences))
			}
			if err := helpers.AssertErr(
				"2nd occs length",
				len(req[j].Occurrences[x].SecondaryOccs),
				len(foundOcc.SecondaryOccs),
			); err != nil {
				return err
			}
			for y := range req[j].Occurrences[x].SecondaryOccs {
				found2Occ := false
				for yy := range foundOcc.SecondaryOccs {
					if err := helpers.AssertJsonErr(
						"",
						&req[j].Occurrences[x].SecondaryOccs[y],
						&foundOcc.SecondaryOccs[yy].SecondaryOccRequest,
					); err == nil {
						found2Occ = true
						break
					}
				}
				if !found2Occ {
					return fmt.Errorf(
						"Couldnt find match for 2occurrency: \n%s\nIn response:\n%s",
						helpers.JsonMustSerializeFormatStr(req[j].Occurrences[x].SecondaryOccs[y]),
						helpers.JsonMustSerializeFormatStr(foundOcc.SecondaryOccs))
				}
			}
		}
	}
	return nil
}

func TestTrainingCRUD(t *testing.T) {

	// create user

	var token string

	{
		ur := user.UserRequest{
			Email:    helpers.CRNG_stringPanic(12),
			Password: "!asdasdSADASDAS56456$%^$%^$%",
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
		trainCtx.Api.TestAssertCases(t, cases)
	}

	// create instructor

	instructorRequest := instr.InstructorRequest{}

	{
		cases := []api.TestCase{
			{
				RequestMethod:      "POST",
				RequestUrl:         "/api/instructor",
				RequestReader:      helpers.JsonMustSerializeReader(&instructorRequest),
				ExpectedStatusCode: 204,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
			},
		}
		trainCtx.Api.TestAssertCases(t, cases)
	}

	correctRequest := train.CreateTrainingRequest{
		Training: train.TrainingRequest{
			Title:           "1",
			Capacity:        1,
			Price:           1000,
			Currency:        "PLN",
			LocationCountry: "PL",
			AllowExpress:    true,
			ManualConfirm:   false,
			RequiredGear: []string{
				"foo",
				"bar",
			},
			RecommendedGear: []string{
				"foobar",
			},
			InstructorGear: []string{
				"foooo",
				"barrr",
			},
			MinAge:                   1,
			MaxAge:                   40,
			TrainingSupportsDisabled: true,
			PlaceSupportsDisabled:    true,
		},
		Occurrences: []train.CreateOccRequest{
			{
				OccRequest: train.OccRequest{
					DateStart:  time.Date(2020, 02, 01, 8, 0, 0, 0, time.UTC),
					DateEnd:    time.Date(2020, 02, 01, 10, 0, 0, 0, time.UTC),
					RepeatDays: 7, // weekly
					Remarks:    helpers.CRNG_stringPanic(128),
					Color:      "ffffff",
				},
			},
		},
	}
	correctRequest2 := train.CreateTrainingRequest{
		Training: train.TrainingRequest{
			Title:         "2",
			Capacity:      12,
			Price:         1546,
			Currency:      "PLN",
			AllowExpress:  false,
			ManualConfirm: true,
		},
		Occurrences: []train.CreateOccRequest{
			{
				OccRequest: train.OccRequest{
					DateStart:  time.Date(2020, 02, 01, 8, 0, 0, 0, time.UTC),
					DateEnd:    time.Date(2020, 02, 01, 10, 0, 0, 0, time.UTC),
					RepeatDays: 7, // weekly
				},
			},
			{
				OccRequest: train.OccRequest{
					DateStart:  time.Date(2020, 02, 02, 13, 0, 0, 0, time.UTC),
					DateEnd:    time.Date(2020, 02, 02, 15, 0, 0, 0, time.UTC),
					RepeatDays: 7, // weekly
				},
			},
		},
	}

	emptyRequest := train.CreateTrainingRequest{
		Training: train.TrainingRequest{
			Title:    "3",
			Capacity: 12,
			Price:    999,
			Currency: "PLN",
		},
	}

	// CREATE test
	{
		url := "/api/training"
		method := "POST"

		invalidRequest := train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title: "",
			},
			Occurrences: correctRequest.Occurrences,
			ReturnID:    false,
		}

		// invalidRequest2 := CreateTrainingRequest{
		// 	Training:    correctRequest.Training,
		// 	Occurrences: []OccRequest{},
		// }

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
				RequestReader:      helpers.JsonMustSerializeReader(&correctRequest),
				ExpectedStatusCode: 204,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
			},
			// correct request 2
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      helpers.JsonMustSerializeReader(&correctRequest2),
				ExpectedStatusCode: 204,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
			},
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      helpers.JsonMustSerializeReader(&emptyRequest),
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
			// {
			// 	RequestMethod:      method,
			// 	RequestUrl:         url,
			// 	RequestReader:      helpers.JsonMustSerializeReader(&correctRequest),
			// 	ExpectedStatusCode: 409,
			// 	RequestHeaders: map[string]string{
			// 		"Authorization": "Bearer " + token,
			// 	},
			// },
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
			// invalid request2
			// {
			// 	RequestMethod:      method,
			// 	RequestUrl:         url,
			// 	RequestReader:      helpers.JsonMustSerializeReader(&invalidRequest2),
			// 	ExpectedStatusCode: 400,
			// 	RequestHeaders: map[string]string{
			// 		"Authorization": "Bearer " + token,
			// 	},
			// },
		}

		trainCtx.Api.TestAssertCases(t, cases)
	}

	var instructorTrainingID uuid.UUID
	var instructorOccIDs [2]uuid.UUID

	// READ
	{
		url := "/api/training"
		method := "GET"

		cases := []api.TestCase{
			// no auth
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      nil,
				ExpectedStatusCode: 403,
			},
			// correct
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedBodyVal: func(b []byte, i interface{}) error {
					var res []*train.TrainingWithJoins
					if err := json.Unmarshal(b, &res); err != nil {
						return err
					}
					req := []*train.CreateTrainingRequest{
						&correctRequest,
						&correctRequest2,
						&emptyRequest,
					}
					for i := range res {
						if res[i].Training.Title == "1" {
							instructorTrainingID = res[i].Training.ID
							instructorOccIDs[0] = res[i].Occurrences[0].ID
							break
						}
					}
					if instructorTrainingID == uuid.Nil {
						return fmt.Errorf("training 1 not found")
					}
					return AssertGetTrainingResponse(res, req)
				},
				ExpectedStatusCode: 200,
			},
			{
				RequestMethod: method,
				RequestUrl:    url + "?q=0",
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedBodyVal: func(b []byte, i interface{}) error {
					var res []train.TrainingWithJoins
					if err := json.Unmarshal(b, &res); err != nil {
						return err
					}
					if len(res) != 0 {
						return fmt.Errorf("get training by query, expected 0 results")
					}
					return nil
				},
				ExpectedStatusCode: 200,
			},
			{
				RequestMethod: method,
				RequestUrl:    url + "?q=2",
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedBodyVal: func(b []byte, i interface{}) error {
					var res []*train.TrainingWithJoins
					if err := json.Unmarshal(b, &res); err != nil {
						return err
					}
					req := []*train.CreateTrainingRequest{
						&correctRequest2,
					}
					return AssertGetTrainingResponse(res, req)
				},
				ExpectedStatusCode: 200,
			},
		}

		trainCtx.Api.TestAssertCases(t, cases)
	}

	patchRequest := train.UpdateTrainingRequest{
		TrainingRequest: train.TrainingRequest{
			Title:                    helpers.CRNG_stringPanic(20),
			Capacity:                 10,
			Price:                    2000,
			Currency:                 "PLN",
			TrainingSupportsDisabled: false,
		},
		ID: instructorTrainingID,
	}

	// UPDATE
	{
		url := "/api/training"
		method := "PATCH"

		invalidPatchRequest := train.UpdateTrainingRequest{
			TrainingRequest: patchRequest.TrainingRequest,
			ID:              uuid.New(),
		}

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
			// not found
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				RequestReader:      helpers.JsonMustSerializeReader(invalidPatchRequest),
				ExpectedStatusCode: 404,
			},
			// correct value
			{
				RequestMethod: "GET",
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedBodyVal: func(b []byte, i interface{}) error {
					var res []*train.TrainingWithJoins
					if err := json.Unmarshal(b, &res); err != nil {
						return err
					}
					req := []*train.CreateTrainingRequest{
						{
							Training:    patchRequest.TrainingRequest,
							Occurrences: correctRequest.Occurrences,
							ReturnID:    false,
						},
						&correctRequest2,
						&emptyRequest,
					}
					return AssertGetTrainingResponse(res, req)
				},
				ExpectedStatusCode: 200,
			},
		}

		trainCtx.Api.TestAssertCases(t, cases)

		correctRequest.Training = patchRequest.TrainingRequest
	}

	// occ
	// {
	// 	url := "/api/training/occ"

	// 	addr := train.PostTrainingOccRequest{
	// 		TrainingID: instructorTrainingID,
	// 		OccRequest: train.OccRequest{
	// 			DateStart:  time.Date(2020, 03, 12, 13, 30, 0, 0, time.UTC),
	// 			DateEnd:    time.Date(2020, 03, 12, 20, 0, 0, 0, time.UTC),
	// 			RepeatDays: 2,
	// 			Remarks:    "123",
	// 		},
	// 	}

	// 	cases := []api.TestCase{
	// 		// no auth
	// 		{
	// 			RequestMethod:      "POST",
	// 			RequestUrl:         url,
	// 			RequestReader:      nil,
	// 			ExpectedStatusCode: 401,
	// 		},
	// 		// correct patch
	// 		{
	// 			RequestMethod: "POST",
	// 			RequestUrl:    url,
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			RequestReader:      helpers.JsonMustSerializeReader(addr),
	// 			ExpectedStatusCode: 204,
	// 		},
	// 		{
	// 			RequestMethod: "GET",
	// 			RequestUrl:    "/api/training",
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
	// 				var res []*train.TrainingRes
	// 				if err := json.Unmarshal(b, &res); err != nil {
	// 					return err
	// 				}
	// 				req := []*train.CreateTrainingRequest{
	// 					{
	// 						Training: correctRequest.Training,
	// 						Occurrences: []train.CreateOccRequest{
	// 							correctRequest.Occurrences[0],
	// 							{
	// 								OccRequest: addr.OccRequest,
	// 							},
	// 						},
	// 						ReturnID: false,
	// 					},
	// 					&correctRequest2,
	// 					&emptyRequest,
	// 				}
	// 				if err := AssertGetTrainingResponse(res, req); err != nil {
	// 					return err
	// 				}
	// 				for j := 0; j < len(res); j++ {
	// 					for i := 0; i < len(res[j].Occurrences); i++ {
	// 						if res[j].Occurrences[i].Remarks == "123" {
	// 							instructorOccIDs[1] = res[j].Occurrences[i].ID
	// 							break
	// 						}
	// 					}
	// 				}
	// 				if instructorOccIDs[1] == uuid.Nil {
	// 					return fmt.Errorf("couldnt find occ")
	// 				}
	// 				return nil
	// 			},
	// 			ExpectedStatusCode: 200,
	// 		},
	// 		// no auth
	// 		{
	// 			RequestMethod:      "DELETE",
	// 			RequestUrl:         url,
	// 			RequestReader:      nil,
	// 			ExpectedStatusCode: 401,
	// 		},
	// 		{
	// 			RequestMethod: "DELETE",
	// 			RequestUrl:    url,
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			RequestReader: helpers.JsonMustSerializeReader(train.DeleteTrainingOccRequest{
	// 				TrainingID: uuid.New(),
	// 				ObjectKey:  train.ObjectKey{ID: instructorOccIDs[0]},
	// 			}),
	// 			ExpectedStatusCode: 404,
	// 		},
	// 		{
	// 			RequestMethod: "DELETE",
	// 			RequestUrl:    url,
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			RequestReader: helpers.JsonMustSerializeReader(train.DeleteTrainingOccRequest{
	// 				TrainingID: instructorTrainingID,
	// 				ObjectKey:  train.ObjectKey{ID: uuid.New()},
	// 			}),
	// 			ExpectedStatusCode: 404,
	// 		},
	// 		{
	// 			RequestMethod: "DELETE",
	// 			RequestUrl:    url,
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			RequestReader: helpers.JsonMustSerializeReader(train.DeleteTrainingOccRequest{
	// 				TrainingID: instructorTrainingID,
	// 				ObjectKey:  train.ObjectKey{ID: instructorOccIDs[0]},
	// 			}),
	// 			ExpectedStatusCode: 204,
	// 		},
	// 		{
	// 			RequestMethod: "GET",
	// 			RequestUrl:    "/api/training",
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
	// 				var res []*train.TrainingRes
	// 				if err := json.Unmarshal(b, &res); err != nil {
	// 					return err
	// 				}
	// 				req := []*train.CreateTrainingRequest{
	// 					{
	// 						Training: correctRequest.Training,
	// 						Occurrences: []train.CreateOccRequest{
	// 							{
	// 								OccRequest: addr.OccRequest,
	// 							},
	// 						},
	// 						ReturnID: false,
	// 					},
	// 					&correctRequest2,
	// 					&emptyRequest,
	// 				}
	// 				return AssertGetTrainingResponse(res, req)
	// 			},
	// 			ExpectedStatusCode: 200,
	// 		},
	// 		{
	// 			RequestMethod: "DELETE",
	// 			RequestUrl:    url,
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			RequestReader: helpers.JsonMustSerializeReader(train.DeleteTrainingOccRequest{
	// 				TrainingID: instructorTrainingID,
	// 				ObjectKey:  train.ObjectKey{ID: instructorOccIDs[0]},
	// 			}),
	// 			ExpectedStatusCode: 404,
	// 		},
	// 	}
	// 	trainCtx.Api.TestAssertCases(t, cases)

	// 	patchReq := train.Occ{
	// 		OccRequest: train.OccRequest{
	// 			DateStart:  time.Date(2023, 05, 12, 13, 30, 0, 0, time.UTC),
	// 			DateEnd:    time.Date(2023, 05, 12, 20, 0, 0, 0, time.UTC),
	// 			RepeatDays: 2,
	// 		},
	// 	}

	// 	cases = []api.TestCase{
	// 		// no auth
	// 		{
	// 			RequestMethod:      "PATCH",
	// 			RequestUrl:         url,
	// 			RequestReader:      nil,
	// 			ExpectedStatusCode: 401,
	// 		},
	// 		{
	// 			RequestMethod: "PATCH",
	// 			RequestUrl:    url,
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			RequestReader:      helpers.JsonMustSerializeReader(patchReq),
	// 			ExpectedStatusCode: 400,
	// 		},
	// 	}
	// 	trainCtx.Api.TestAssertCases(t, cases)

	// 	patchReq.ID = instructorOccIDs[1]
	// 	patchReq.TrainingID = instructorTrainingID

	// 	cases = []api.TestCase{
	// 		{
	// 			RequestMethod: "PATCH",
	// 			RequestUrl:    url,
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			RequestReader:      helpers.JsonMustSerializeReader(patchReq),
	// 			ExpectedStatusCode: 204,
	// 		},
	// 		{
	// 			RequestMethod: "GET",
	// 			RequestUrl:    "/api/training",
	// 			RequestHeaders: map[string]string{
	// 				"Authorization": "Bearer " + token,
	// 			},
	// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
	// 				var res []*train.TrainingRes
	// 				if err := json.Unmarshal(b, &res); err != nil {
	// 					return err
	// 				}
	// 				req := []*train.CreateTrainingRequest{
	// 					{
	// 						Training: correctRequest.Training,
	// 						Occurrences: []train.CreateOccRequest{
	// 							{
	// 								OccRequest: patchReq.OccRequest,
	// 							},
	// 						},
	// 						ReturnID: false,
	// 					},
	// 					&correctRequest2,
	// 					&emptyRequest,
	// 				}
	// 				return AssertGetTrainingResponse(res, req)
	// 			},
	// 			ExpectedStatusCode: 200,
	// 		},
	// 	}
	// 	trainCtx.Api.TestAssertCases(t, cases)
	// }

	// DELETE
	{
		url := "/api/training"
		method := "DELETE"

		deleteRequest := train.ObjectKey{
			ID: instructorTrainingID,
		}

		invalidDeleteRequest := train.ObjectKey{
			ID: uuid.New(),
		}

		cases := []api.TestCase{
			// no auth
			{
				RequestMethod:      method,
				RequestUrl:         url,
				RequestReader:      nil,
				ExpectedStatusCode: 401,
			},
			// empty delete
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedStatusCode: 400,
			},
			// invalid delete
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				RequestReader:      helpers.JsonMustSerializeReader(invalidDeleteRequest),
				ExpectedStatusCode: 404,
			},
			// correct delete
			{
				RequestMethod: method,
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				RequestReader:      helpers.JsonMustSerializeReader(deleteRequest),
				ExpectedStatusCode: 204,
			},
			// get
			{
				RequestMethod: "GET",
				RequestUrl:    url,
				RequestHeaders: map[string]string{
					"Authorization": "Bearer " + token,
				},
				ExpectedBodyVal: func(b []byte, i interface{}) error {
					var res []*train.TrainingWithJoins
					if err := json.Unmarshal(b, &res); err != nil {
						return err
					}
					req := []*train.CreateTrainingRequest{
						&correctRequest2,
						&emptyRequest,
					}
					return AssertGetTrainingResponse(res, req)
				},
				ExpectedStatusCode: 200,
			},
		}

		trainCtx.Api.TestAssertCases(t, cases)

	}
}

// func TestTrainingsSemaphore(t *testing.T) {

// 	token, err := user.CreateAndLoginUser(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	userID, err := global.Api.AuthorizeUserFromToken(token)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	iid, err := CreateTestInstructor(userID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	numTests := runtime.NumCPU() * 5
// 	numSem := runtime.NumCPU()
// 	apiReq := make([]api.TestCase, numTests)
// 	for i := 0; i < numTests; i++ {
// 		apiReq[i] = api.TestCase{
// 			RequestMethod: "POST",
// 			RequestUrl:    "/api/training",
// 			RequestReader: helpers.JsonMustSerializeReader(CreateTrainingRequest{
// 				Training: TrainingRequest{
// 					Title:           helpers.CRNG_stringPanic(20),
// 					Capacity:        1,
// 					Price:           1,
// 					Currency:        "PLN",
// 					LocationCountry: "PL",
// 				},
// 				// this doesnt matter right now
// 				Occurrences: []OccRequest{
// 					{
// 						DateStart:  time.Date(2020, 02, 01, 8, 0, 0, 0, time.UTC),
// 						DateEnd:    time.Date(2020, 02, 01, 10, 0, 0, 0, time.UTC),
// 						RepeatDays: 7, // weekly
// 						Remarks:    helpers.CRNG_stringPanic(128),
// 						Color:      "ffffff",
// 					},
// 					{
// 						DateStart:  time.Date(2020, 02, 01, 11, 0, 0, 0, time.UTC),
// 						DateEnd:    time.Date(2020, 02, 01, 12, 0, 0, 0, time.UTC),
// 						RepeatDays: 3,
// 						Remarks:    helpers.CRNG_stringPanic(128),
// 						Color:      "ffffff",
// 					},
// 				},
// 			}),
// 			ExpectedStatusCode: 204,
// 			RequestHeaders: map[string]string{
// 				"Authorization": "Bearer " + token,
// 			},
// 		}
// 	}

// 	global.Api.TestAssertCasesSemaphore(t, apiReq, numSem)

// 	res, err := DalReadTrainings(DalReadTrainingsRequest{
// 		InstructorID: &iid,
// 	}, true)

// 	numOcc := 0

// 	for i := 0; i < len(res); i++ {
// 		numOcc += len(res[i].Occurrences)
// 	}

// 	if err := helpers.AssertErr("invalid number of trainings", numTests, len(res)); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := helpers.AssertErr("invalid number of occs", numTests*2, numOcc); err != nil {
// 		t.Fatal(err)
// 	}

// }
