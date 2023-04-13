package train_test

import (
	"encoding/json"
	"sport/api"
	"sport/helpers"
	"sport/train"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestOcc(t *testing.T) {
	token, iid, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now().In(time.UTC).Round(time.Minute)
	tr := train.NewTestCreateTrainingRequest(start, start.Add(time.Hour))
	tid, err := trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}

	// ensure that training is as expected

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "GET",
			RequestUrl:    "/api/training",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []*train.TrainingWithJoins
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				return AssertGetTrainingResponse(res, []*train.CreateTrainingRequest{&tr})
			},
			ExpectedStatusCode: 200,
		},
	})

	// add another occ
	tr.Occurrences = append(tr.Occurrences, train.CreateOccRequest{
		OccRequest: train.NewTestOccRequest(start.Add(time.Hour), start.Add(time.Hour*3)),
	})

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "PUT",
			RequestUrl:         "/api/training/occ",
			ExpectedStatusCode: 401,
		}, {
			RequestMethod: "PUT",
			RequestUrl:    "/api/training/occ",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			ExpectedStatusCode: 400,
		}, {
			RequestMethod: "PUT",
			RequestUrl:    "/api/training/occ",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			RequestReader: helpers.JsonMustSerializeReader(&train.PutTrainingOccsRequest{
				Occurrences: tr.Occurrences,
				TrainingID:  uuid.New(),
			}),
			ExpectedStatusCode: 404,
		}, {
			RequestMethod: "PUT",
			RequestUrl:    "/api/training/occ",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			RequestReader: helpers.JsonMustSerializeReader(&train.PutTrainingOccsRequest{
				Occurrences: tr.Occurrences,
				TrainingID:  tid,
			}),
			ExpectedStatusCode: 204,
		}, {
			RequestMethod: "GET",
			RequestUrl:    "/api/training",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []*train.TrainingWithJoins
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				return AssertGetTrainingResponse(res, []*train.CreateTrainingRequest{&tr})
			},
			ExpectedStatusCode: 200,
		}, {
			RequestMethod: "PUT",
			RequestUrl:    "/api/training/occ",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			RequestReader: helpers.JsonMustSerializeReader(&train.PutTrainingOccsRequest{
				Occurrences: tr.Occurrences,
				TrainingID:  tid,
			}),
			ExpectedStatusCode: 204,
		}, {
			RequestMethod: "GET",
			RequestUrl:    "/api/training",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []*train.TrainingWithJoins
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				return AssertGetTrainingResponse(res, []*train.CreateTrainingRequest{&tr})
			},
			ExpectedStatusCode: 200,
		},
	})

	// add secondary occs
	tr.Occurrences[0].DateEnd = tr.Occurrences[0].DateStart.Add(time.Minute * 864)
	tr.Occurrences[0].SecondaryOccs = []train.SecondaryOccRequest{
		train.NewTest2OccRequest(65, 120),
		train.NewTest2OccRequest(0, 65),
		train.NewTest2OccRequest(130, 200),
		train.NewTest2OccRequest(400, 864),
	}
	tr.Occurrences[1].DateEnd = tr.Occurrences[1].DateStart.Add(time.Minute * 150)
	tr.Occurrences[1].SecondaryOccs = []train.SecondaryOccRequest{
		train.NewTest2OccRequest(0, 120),
		train.NewTest2OccRequest(130, 150),
	}
	tr.Occurrences = append(tr.Occurrences, train.CreateOccRequest{
		OccRequest: train.NewTestOccRequest(start.Add(time.Hour*20), start.Add(time.Hour*30)),
	})

	start2 := time.Now().In(time.UTC).Round(time.Hour)
	tr2 := train.NewTestCreateTrainingRequest(start2, start2.Add(time.Hour))
	_, err = trainCtx.ApiCreateTraining(token, &tr2)
	if err != nil {
		t.Fatal(err)
	}

	tr3 := train.NewTestCreateTrainingRequestNoOcc()
	_, err = trainCtx.ApiCreateTraining(token, &tr3)
	if err != nil {
		t.Fatal(err)
	}

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "PUT",
			RequestUrl:    "/api/training/occ",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			RequestReader: helpers.JsonMustSerializeReader(&train.PutTrainingOccsRequest{
				Occurrences: tr.Occurrences,
				TrainingID:  tid,
			}),
			ExpectedStatusCode: 204,
		}, {
			RequestMethod: "GET",
			RequestUrl:    "/api/training",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []*train.TrainingWithJoins
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				return AssertGetTrainingResponse(res, []*train.CreateTrainingRequest{&tr, &tr2, &tr3})
			},
			ExpectedStatusCode: 200,
		},
	})

	_ = iid
}
