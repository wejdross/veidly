package rsv_test

// import (
// 	"sport/adyen_sm"
// 	"sport/api"
// 	"sport/helpers"
// 	"sport/rsv"
// 	"sport/train"
// 	"sport/user"
// 	"testing"
// 	"time"
// )

// func TestContact(t *testing.T) {
// 	itoken, instrID, err := rsvCtx.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	token, err := rsvCtx.User.ApiCreateAndLoginUser(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	m := "GET"
// 	u := "/api/rsv/instr/contact"

// 	rsvCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         u,
// 			ExpectedStatusCode: 400,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         u + "?instructor_id=" + instrID.String(),
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         u + "?instructor_id=" + instrID.String(),
// 			RequestHeaders:     user.GetAuthHeader(token),
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:  "PATCH",
// 			RequestUrl:     "/api/user/contact",
// 			RequestHeaders: user.GetAuthHeader(itoken),
// 			RequestReader: helpers.JsonMustSerializeReader(user.ContactData{
// 				Share: true,
// 			}),
// 			ExpectedStatusCode: 204,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         u + "?instructor_id=" + instrID.String(),
// 			RequestHeaders:     user.GetAuthHeader(token),
// 			ExpectedStatusCode: 200,
// 		}, {
// 			RequestMethod:  "PATCH",
// 			RequestUrl:     "/api/user/contact",
// 			RequestHeaders: user.GetAuthHeader(itoken),
// 			RequestReader: helpers.JsonMustSerializeReader(user.ContactData{
// 				Share: false,
// 			}),
// 			ExpectedStatusCode: 204,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         u + "?instructor_id=" + instrID.String(),
// 			RequestHeaders:     user.GetAuthHeader(token),
// 			ExpectedStatusCode: 401,
// 		},
// 	})

// 	{
// 		start := rsv.NowMin().Add(time.Hour)
// 		end := start.Add(time.Hour)
// 		ct := train.NewTestCreateTrainingRequest(start, end)
// 		trainingID, err := rsvCtx.Train.ApiCreateTraining(itoken, &ct)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		tr := rsv.NewTestCreateReservationRequest(trainingID, start)
// 		r, err := rsvCtx.ApiCreateAndQueryRsv(&token, &tr, 0)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if err = rsvCtx.AdyenSm.MoveStateToCapture(
// 			rsvCtx.RsvResponseToSmPassPtr(r),
// 			adyen_sm.ManualSrc, nil,
// 		); err != nil {
// 			t.Fatal(err)
// 		}

// 		rsvCtx.Api.TestAssertCases(t, []api.TestCase{
// 			{
// 				RequestMethod:      m,
// 				RequestUrl:         u + "?instructor_id=" + instrID.String(),
// 				RequestHeaders:     user.GetAuthHeader(token),
// 				ExpectedStatusCode: 200,
// 			},
// 		})
// 	}

// 	{
// 		start := rsv.NowMin().Add(time.Hour)
// 		end := start.Add(time.Hour)
// 		ct := train.NewTestCreateTrainingRequest(start, end)
// 		trainingID, err := rsvCtx.Train.ApiCreateTraining(itoken, &ct)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		tr := rsv.NewTestCreateReservationRequest(trainingID, start)
// 		r, err := rsvCtx.ApiCreateAndQueryRsv(&token, &tr, 0)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if err = rsvCtx.AdyenSm.MoveStateToCapture(
// 			rsvCtx.RsvResponseToSmPassPtr(r),
// 			adyen_sm.ManualSrc, nil,
// 		); err != nil {
// 			t.Fatal(err)
// 		}

// 		rsvCtx.Api.TestAssertCases(t, []api.TestCase{
// 			{
// 				RequestMethod:      m,
// 				RequestUrl:         u + "?instructor_id=" + instrID.String() + "&access_token=" + r.AccessToken.String(),
// 				ExpectedStatusCode: 200,
// 			},
// 		})
// 	}
// }
