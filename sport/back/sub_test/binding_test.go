package sub_test

import (
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/sub"
	"sport/train"
	"sport/user"
	"testing"
)

func TestBindings(t *testing.T) {
	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	var smb sub.SubModelBinding
	req := sub.NewTestSubModelRequest()
	smb.SubModelID, err = subCtx.ApiCreateSubModel(token, &req)
	if err != nil {
		t.Fatal(err)
	}
	treq := train.NewTestCreateTrainingRequestNoOcc()
	smb.TrainingID, err = trainCtx.ApiCreateTraining(token, &treq)
	if err != nil {
		t.Fatal(err)
	}

	token2, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	var smb2 sub.SubModelBinding
	req2 := sub.NewTestSubModelRequest()
	smb2.SubModelID, err = subCtx.ApiCreateSubModel(token2, &req2)
	if err != nil {
		t.Fatal(err)
	}
	treq2 := train.NewTestCreateTrainingRequestNoOcc()
	smb2.TrainingID, err = trainCtx.ApiCreateTraining(token2, &treq2)
	if err != nil {
		t.Fatal(err)
	}

	var false1, false2 sub.SubModelBinding
	false1 = sub.SubModelBinding{
		SubModelID: smb.SubModelID,
		TrainingID: smb2.TrainingID,
	}
	false2 = sub.SubModelBinding{
		SubModelID: smb2.SubModelID,
		TrainingID: smb.TrainingID,
	}

	// POST
	m := "POST"
	u := "/api/sub/model/binding"
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
			ExpectedStatusCode: 404,
			RequestReader:      helpers.JsonMustSerializeReader(smb),
			RequestHeaders:     user.GetAuthHeader(token2),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader:      helpers.JsonMustSerializeReader(smb2),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader:      helpers.JsonMustSerializeReader(false1),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader:      helpers.JsonMustSerializeReader(false2),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader:      helpers.JsonMustSerializeReader(false1),
			RequestHeaders:     user.GetAuthHeader(token2),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader:      helpers.JsonMustSerializeReader(false2),
			RequestHeaders:     user.GetAuthHeader(token2),
		},
		// correct requests
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(smb),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(smb2),
			RequestHeaders:     user.GetAuthHeader(token2),
		},
		// validate
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/training?sm_id=" + smb.SubModelID.String(),
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var x []*train.TrainingWithJoins
				helpers.JsonMustDeserialize(b, &x)
				if len(x) != 1 {
					return fmt.Errorf("invalid trainings 1")
				}
				if x[0].Training.Title != treq.Training.Title {
					return fmt.Errorf("invalid training 1")
				}
				if len(x[0].Sms) != 1 {
					return fmt.Errorf("invalid Sm's 1")
				}
				g := x[0].Sms[0]
				if err := helpers.AssertJsonErr("assert train 1 sm", &g.SubModelEditable, &req.SubModelEditable); err != nil {
					return err
				}
				pi := subCtx.NewSubPricing(req.Price)
				if err := helpers.AssertJsonErr(
					"assert train 1 pricing", &g.PricingInfo, &pi,
				); err != nil {
					return err
				}
				return nil
			},
		}, {
			RequestMethod:      "GET",
			RequestUrl:         "/api/training?sm_id=" + smb2.SubModelID.String(),
			RequestHeaders:     user.GetAuthHeader(token2),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var x []*train.TrainingWithJoins
				helpers.JsonMustDeserialize(b, &x)
				if len(x) != 1 {
					return fmt.Errorf("invalid trainings 2")
				}
				if x[0].Training.Title != treq2.Training.Title {
					return fmt.Errorf("invalid training 2")
				}
				if len(x[0].Sms) != 1 {
					return fmt.Errorf("invalid Sm's 2")
				}
				g := x[0].Sms[0]
				if err := helpers.AssertJsonErr(
					"assert train 2 sm", &g.SubModelEditable, &req2.SubModelEditable); err != nil {
					return err
				}
				pi := subCtx.NewSubPricing(req2.Price)
				if err := helpers.AssertJsonErr(
					"assert train 2 pricing", &g.PricingInfo, &pi,
				); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
