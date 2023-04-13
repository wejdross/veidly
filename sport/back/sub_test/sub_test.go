package sub_test

import (
	"encoding/json"
	"fmt"
	"sport/adyen_sm"
	"sport/api"
	"sport/helpers"
	"sport/sub"
	"sport/user"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSub(t *testing.T) {

	instrToken, iid, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	userToken, err := trainCtx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}

	instrUserID, err := trainCtx.Api.AuthorizeUserFromToken(instrToken)
	if err != nil {
		t.Fatal(err)
	}

	userID, err := trainCtx.Api.AuthorizeUserFromToken(userToken)
	if err != nil {
		t.Fatal(err)
	}

	modelReq := sub.NewTestSubModelRequest()
	subModelID, err := subCtx.ApiCreateSubModel(instrToken, &modelReq)
	if err != nil {
		t.Fatal(err)
	}

	m := "POST"
	u := "/api/sub"
	var subID uuid.UUID
	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(instrToken),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestReader:      helpers.JsonMustSerializeReader(sub.ApiSubRequest{}),
			RequestHeaders:     user.GetAuthHeader(instrToken),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(sub.ApiSubRequest{
				SubModelID: uuid.New(),
			}),
			RequestHeaders: user.GetAuthHeader(instrToken),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader: helpers.JsonMustSerializeReader(sub.ApiSubRequest{
				SubModelID: subModelID,
			}),
			RequestHeaders: user.GetAuthHeader(userToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res sub.ApiSubResponse
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				if res.ID == uuid.Nil {
					return fmt.Errorf("invalid ID")
				}
				subID = res.ID
				if res.Url == "" {
					return fmt.Errorf("invalid Url")
				}
				return nil
			},
		},
	})

	s, err := subCtx.DalReadSingleSubWithJoins(sub.ReadSubRequest{
		ID: &subID,
	})
	if err != nil {
		t.Fatal(err)
	}

	helpers.Assert(t, "sub ID", s.ID, subID)
	helpers.AssertJson(t,
		"sub SubModel",
		s.SubModel.SubModelEditable, modelReq.SubModelEditable)
	helpers.Assert(t, "sub SubModelID", *s.SubModelID, subModelID)
	helpers.Assert(t, "sub InstructorID", s.InstructorID, iid)
	helpers.Assert(t, "sub InstrUserID", s.InstrUserID, instrUserID)
	helpers.Assert(t, "sub UserID", s.UserID, userID)
	if s.RefID == uuid.Nil {
		t.Fatal("invalid RefID")
	}
	helpers.Assert(t, "sub State", s.State, adyen_sm.Hold)
	n := time.Now().In(time.UTC)
	helpers.EnsuretTimeInRange(s.DateStart, n)
	//helpers.EnsuretTimeInRange(s.DateEnd, n.Add(time.Duration(modelReq.Durations[0])*24*time.Hour))
	// sm_timeout
	// sm_cache
	// created_on

	m = "GET"
	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(userToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []sub.Sub
				helpers.JsonMustDeserialize(b, &res)
				if len(res) != 1 {
					return fmt.Errorf("1 invalid res")
				}
				return nil
			},
		}, {
			RequestMethod:      m,
			RequestUrl:         u + "/user",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(userToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []sub.Sub
				helpers.JsonMustDeserialize(b, &res)
				if len(res) != 1 {
					return fmt.Errorf("2 invalid res")
				}
				return nil
			},
		}, {
			RequestMethod:      m,
			RequestUrl:         u + "/instructor",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(instrToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []sub.Sub
				helpers.JsonMustDeserialize(b, &res)
				if len(res) != 1 {
					return fmt.Errorf("3 invalid res")
				}
				return nil
			},
		}, {
			RequestMethod:      m,
			RequestUrl:         u + "/instructor",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(userToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []sub.Sub
				helpers.JsonMustDeserialize(b, &res)
				if len(res) != 0 {
					return fmt.Errorf("4 invalid res")
				}
				return nil
			},
		}, {
			RequestMethod:      m,
			RequestUrl:         u + "/user",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(instrToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []sub.Sub
				helpers.JsonMustDeserialize(b, &res)
				if len(res) != 0 {
					return fmt.Errorf("5 invalid res")
				}
				return nil
			},
		},
	})
}
