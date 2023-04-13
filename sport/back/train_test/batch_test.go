package train_test

// func TestBatchOcc(t *testing.T) {
// 	tk, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tr := train.NewTestCreateTrainingRequestNoOcc()
// 	id, err := trainCtx.ApiCreateTraining(tk, &tr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	s := time.Now().Round(time.Minute).In(time.UTC)
// 	var req = train.PostTrainingOccBatchRequest{
// 		TrainingID: id,
// 		Occs: []train.OccRequest{
// 			train.NewTestOccRequest(s, s.Add(time.Hour)),
// 			train.NewTestOccRequest(s.Add(time.Hour), s.Add(2*time.Hour)),
// 			train.NewTestOccRequest(s.Add(2*time.Hour), s.Add(3*time.Hour)),
// 		},
// 	}

// 	m := "POST"
// 	u := "/api/training/occ/batch"

// 	trainCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         u,
// 			RequestReader:      nil,
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         u,
// 			RequestHeaders:     user.GetAuthHeader(tk),
// 			RequestReader:      helpers.JsonMustSerializeReader(req),
// 			ExpectedStatusCode: 204,
// 		},
// 	})
// 	u = "/api/training"
// 	m = "GET"
// 	tr.Occurrences = req.Occs
// 	trainCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         u,
// 			RequestHeaders:     user.GetAuthHeader(tk),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res []*train.TrainingRes
// 				if err := json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				req := []*train.CreateTrainingRequest{
// 					&tr,
// 				}
// 				return AssertGetTrainingResponse(res, req)
// 			},
// 		},
// 	})
// }
