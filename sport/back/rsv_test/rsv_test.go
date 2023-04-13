package rsv_test

// import (
// 	"sport/adyen_sm"
// 	"sport/api"
// 	"sport/rsv"
// 	"sport/train"
// 	"sport/user"
// 	"testing"
// 	"time"
// )

// func TestCanDeleteInstr(t *testing.T) {

// 	var pass = "!@#!@312312#@!#sadasASDSA"

// 	token, err := rsvCtx.User.ApiCreateAndLoginUser(&user.UserRequest{
// 		Email:    "foo@foo.foo",
// 		Password: pass,
// 		UserData: user.UserData{
// 			Country:  "PL",
// 			Language: "pl",
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if err = rsvCtx.Instr.ApiCreateInstructor(token, nil); err != nil {
// 		t.Fatal(err)
// 	}

// 	rsvCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod: "GET",
// 			RequestUrl:    "/api/instructor/can_delete",
// 			RequestHeaders: map[string]string{
// 				"Authorization": "Bearer " + token,
// 			},
// 			ExpectedStatusCode: 204,
// 		},
// 	})

// 	start := time.Now()
// 	start = time.Date(
// 		start.Year(),
// 		start.Month(),
// 		start.Day(),
// 		start.Hour(),
// 		0, 0, 0, time.UTC)
// 	start = start.Add(time.Hour * 24 * 2)

// 	trainingID, err := rsvCtx.Train.ApiCreateTraining(token, &train.CreateTrainingRequest{
// 		Training: train.TrainingRequest{
// 			Title:    "foo",
// 			Price:    1000,
// 			Currency: "PLN",
// 			Capacity: 1,
// 		},
// 		Occurrences: []train.CreateOccRequest{
// 			{
// 				OccRequest: train.OccRequest{
// 					DateStart: start,
// 					DateEnd:   start.Add(time.Hour),
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	r, err := rsvCtx.ApiCreateAndQueryRsv(&token, &rsv.ApiReservationRequest{
// 		TrainingID: trainingID,
// 		Occurrence: start,
// 		UserData: user.UserData{
// 			Country:  "PL",
// 			Language: "pl",
// 			Name:     "vv",
// 		},
// 		NoRedirect: true,
// 	}, 200)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rsvCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod: "GET",
// 			RequestUrl:    "/api/instructor/can_delete",
// 			RequestHeaders: map[string]string{
// 				"Authorization": "Bearer " + token,
// 			},
// 			ExpectedStatusCode: 409,
// 		},
// 		{
// 			RequestMethod:      "DELETE",
// 			RequestUrl:         "/api/instructor",
// 			ExpectedStatusCode: 401,
// 		},
// 		{
// 			RequestMethod: "DELETE",
// 			RequestUrl:    "/api/instructor",
// 			RequestHeaders: map[string]string{
// 				"Authorization": "Bearer " + token,
// 			},
// 			ExpectedStatusCode: 409,
// 		},
// 	})

// 	// make rsv inactive but expiring link
// 	if err := rsvCtx.AdyenSm.MoveStateToLinkExpire(
// 		rsvCtx.RsvResponseToSmPassPtr(r), adyen_sm.ManualSrc,
// 	); err != nil {
// 		t.Fatal(err)
// 	}

// 	rsvCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod: "GET",
// 			RequestUrl:    "/api/instructor/can_delete",
// 			RequestHeaders: map[string]string{
// 				"Authorization": "Bearer " + token,
// 			},
// 			ExpectedStatusCode: 204,
// 		},
// 		{
// 			RequestMethod:      "DELETE",
// 			RequestUrl:         "/api/instructor",
// 			ExpectedStatusCode: 401,
// 		},
// 		{
// 			RequestMethod: "DELETE",
// 			RequestUrl:    "/api/instructor",
// 			RequestHeaders: map[string]string{
// 				"Authorization": "Bearer " + token,
// 			},
// 			ExpectedStatusCode: 204,
// 		},
// 	})
// }
