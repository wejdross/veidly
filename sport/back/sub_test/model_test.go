package sub_test

import (
	"encoding/json"
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/sub"
	"sport/train"
	"sport/user"
	"testing"

	"github.com/google/uuid"
)

func AssertGetSubModelsRes(
	rawRes []byte,
	expected []sub.SubModelRequest,
	creatorsUserID uuid.UUID,
	ids map[string]uuid.UUID) error {
	//
	var res []sub.SubModel
	if err := json.Unmarshal(rawRes, &res); err != nil {
		return err
	}
	if len(res) != len(expected) {
		return fmt.Errorf("AssertGetSubModelsRes: invalid length")
	}
	for i := range expected {
		found := false
		for j := range res {
			if err := helpers.AssertJsonErr(
				"", res[j].SubModelEditable, expected[i].SubModelEditable); err != nil {
				continue
			}
			pi := subCtx.NewSubPricing(expected[i].Price)
			if err := helpers.AssertJsonErr("", res[j].PricingInfo, pi); err != nil {
				return fmt.Errorf("invalid pricing info")
			}
			if res[j].InstrUserID != creatorsUserID {
				return fmt.Errorf("invalid userID")
			}
			if res[j].ID == uuid.Nil {
				return fmt.Errorf("invalid ID")
			}
			if ids != nil {
				ids[res[j].Name] = res[j].ID
			}
			found = true
			break
		}
		if !found {
			return fmt.Errorf("unable to find item: \n%s\n in response: \n%s\n",
				helpers.JsonMustSerializeFormatStr(expected[i]),
				helpers.JsonMustSerializeFormatStr(res))
		}
	}
	return nil
}

func TestTrainingGroupCrud(t *testing.T) {
	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	userID, err := trainCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}

	// POST
	m := "POST"
	u := "/api/sub/model"

	sms := []sub.SubModelRequest{
		{
			SubModelEditable: sub.SubModelEditable{
				Name:           "foo",
				MaxEntrances:   12,
				Duration:       100,
				Currency:       "PLN",
				MaxActive:      -1,
				IsFreeEntrance: true,
			},
			Price: 10000,
		},
		{
			SubModelEditable: sub.SubModelEditable{
				Name:           "bar",
				MaxEntrances:   3,
				Duration:       31,
				Currency:       "PLN",
				MaxActive:      -1,
				IsFreeEntrance: false,
			},
			Price: 5000,
		},
	}

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestReader:      helpers.JsonMustSerializeReader(sub.SubModelRequest{}),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(sms[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(sms[1]),
			RequestHeaders:     user.GetAuthHeader(token),
		},
	})

	m = "GET"
	ids := make(map[string]uuid.UUID)

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return AssertGetSubModelsRes(b, sms, userID, ids)
			},
		},
	})

	m = "PATCH"

	sms = []sub.SubModelRequest{
		{
			SubModelEditable: sub.SubModelEditable{
				Name:           "foo",
				MaxEntrances:   102,
				Duration:       80,
				Currency:       "PLN",
				MaxActive:      -1,
				IsFreeEntrance: true,
			},
			Price: 100000,
		},
		{
			SubModelEditable: sub.SubModelEditable{
				Name:           "bar",
				MaxEntrances:   102,
				Duration:       67,
				Currency:       "PLN",
				MaxActive:      12,
				IsFreeEntrance: false,
			},
			Price: 100000,
		},
	}

	updateReqs := make([]sub.UpdateSubModelRequest, len(sms))
	for i := range sms {
		id := ids[sms[i].Name]
		if id == uuid.Nil {
			t.Fatal("invalid id")
		}
		updateReqs[i] = sub.UpdateSubModelRequest{
			ID:              ids[sms[i].Name],
			SubModelRequest: sms[i],
		}
	}

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestReader:      helpers.JsonMustSerializeReader(sub.SubModelRequest{}),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestReader:      helpers.JsonMustSerializeReader(sms[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(sub.UpdateSubModelRequest{
				ID:              uuid.New(),
				SubModelRequest: sms[0],
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(updateReqs[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(updateReqs[1]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      "GET",
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return AssertGetSubModelsRes(b, sms, userID, ids)
			},
		},
	})

	m = "DELETE"
	sms = []sub.SubModelRequest{
		// first one 'foo' is removed
		sms[1],
	}

	delReq := train.DeleteGroupRequest{
		ID: ids["foo"],
	}

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestReader:      helpers.JsonMustSerializeReader(sms[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(sub.DeleteSubModelRequest{
				ID: uuid.New(),
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(delReq),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      "GET",
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return AssertGetSubModelsRes(b, sms, userID, ids)
			},
		},
	})
}
