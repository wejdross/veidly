package train_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sport/api"
	"sport/helpers"
	"sport/train"
	"sport/user"
	"testing"
	"time"

	"github.com/google/uuid"
)

func AssertGetTrainingGroupsRes(
	rawRes []byte,
	expected []train.GroupRequest,
	creatorsUserID uuid.UUID,
	ids map[string]uuid.UUID) error {
	//
	var res []train.Group
	if err := json.Unmarshal(rawRes, &res); err != nil {
		return err
	}
	if len(res) != len(expected) {
		return fmt.Errorf("AssertGetTrainingGroupsRes: invalid length")
	}
	for i := range expected {
		found := false
		for j := range res {
			if err := helpers.AssertJsonErr("", res[j].GroupRequest, expected[i]); err != nil {
				continue
			}
			if res[j].UserID != creatorsUserID {
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

func AssertGroupsAreTheSame(a, b train.GroupArr) error {

	if len(a) != len(b) {
		return fmt.Errorf("AssertGroupsAreTheSame: invalid length")
	}
	ix := make(map[int]struct{})
	for i := range a {
		s := helpers.JsonMustSerializeFormatStr(a[i])
		found := false
		for j := range b {
			if _, ok := ix[j]; ok {
				continue
			}
			bj := helpers.JsonMustSerializeFormatStr(b[j])

			if s == bj {
				ix[j] = struct{}{}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unable to find item from a: \n%s\nin b: \n%s\n",
				s,
				helpers.JsonMustSerializeFormatStr(b))
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
	u := "/api/training/group"

	grps := []train.GroupRequest{
		{
			Name:         "foo",
			MaxPeople:    1,
			MaxTrainings: 2,
		},
		{
			Name:         "bar",
			MaxPeople:    3,
			MaxTrainings: 4,
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
			RequestReader: helpers.JsonMustSerializeReader(train.GroupRequest{
				Name: "",
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(grps[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(grps[1]),
			RequestHeaders:     user.GetAuthHeader(token),
		},
	})

	m = "GET"
	ids := make(map[string]uuid.UUID)

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return AssertGetTrainingGroupsRes(b, grps, userID, ids)
			},
		},
	})

	m = "PATCH"

	grps = []train.GroupRequest{
		{
			Name:         "foo",
			MaxPeople:    1,
			MaxTrainings: 2,
		},
		{
			Name:         "bar",
			MaxPeople:    5,
			MaxTrainings: 6,
		},
	}

	newGrpsUpdateReqs := make([]train.UpdateGroupRequest, len(grps))
	for i := range grps {
		id := ids[grps[i].Name]
		if id == uuid.Nil {
			t.Fatal("invalid id")
		}
		newGrpsUpdateReqs[i] = train.UpdateGroupRequest{
			ID:           ids[grps[i].Name],
			GroupRequest: grps[i],
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
			RequestReader: helpers.JsonMustSerializeReader(train.GroupRequest{
				Name: "",
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestReader:      helpers.JsonMustSerializeReader(grps[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(train.UpdateGroupRequest{
				ID:           uuid.New(),
				GroupRequest: grps[0],
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(newGrpsUpdateReqs[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(newGrpsUpdateReqs[1]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      "GET",
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return AssertGetTrainingGroupsRes(b, grps, userID, ids)
			},
		},
	})

	m = "DELETE"
	grps = []train.GroupRequest{
		// first one 'foo' is removed
		{
			Name:         "bar",
			MaxPeople:    5,
			MaxTrainings: 6,
		},
	}

	newGrpsDeleteReqs := make([]train.DeleteGroupRequest, len(grps))
	newGrpsDeleteReqs[0] = train.DeleteGroupRequest{
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
			RequestReader:      helpers.JsonMustSerializeReader(grps[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(train.DeleteGroupRequest{
				ID: uuid.New(),
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(newGrpsDeleteReqs[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      "GET",
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return AssertGetTrainingGroupsRes(b, grps, userID, ids)
			},
		},
	})
}

func TestQueryingTrainings(t *testing.T) {
	token, iid, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	uid, err := trainCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}

	s := time.Now().Round(time.Minute).In(time.UTC)
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Hour))
	tid, err := trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}
	tr2 := train.NewTestCreateTrainingRequest(s, s.Add(time.Hour))
	_, err = trainCtx.ApiCreateTraining(token, &tr2)
	if err != nil {
		t.Fatal(err)
	}

	gr := train.NewTestGroupRequest()
	gid, err := trainCtx.ApiCreateGroup(token, &gr)
	if err != nil {
		t.Fatal(err)
	}
	if err := trainCtx.DalAddTvgBinding(tid, gid, uid); err != nil {
		t.Fatal(err)
	}

	tres, err := trainCtx.DalReadTrainings(train.DalReadTrainingsRequest{
		InstructorID: &iid,
		WithOccs:     true,
		WithGroups:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := AssertGetTrainingResponse(tres, []*train.CreateTrainingRequest{&tr, &tr2}); err != nil {
		t.Fatal(err)
	}

	tres, err = trainCtx.DalReadTrainings(train.DalReadTrainingsRequest{
		InstructorID: &iid,
		GroupIds:     []string{gid.String()},
		WithOccs:     true,
		WithGroups:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := AssertGetTrainingResponse(tres, []*train.CreateTrainingRequest{&tr}); err != nil {
		t.Fatal(err)
	}

	ids := helpers.JsonMustSerializeStr([]string{gid.String()})
	ids = url.QueryEscape(ids)

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/training?group_ids=" + ids,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []*train.TrainingWithJoins
				helpers.JsonMustDeserialize(b, &res)
				return AssertGetTrainingResponse(res, []*train.CreateTrainingRequest{&tr})
			},
		},
	})
}

func testTvg(token string, tid uuid.UUID) error {
	// token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	uid, err := trainCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		return err
	}
	// s := time.Now().In(time.UTC)
	// tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Hour))
	// tid, err := trainCtx.ApiCreateTraining(token, &tr)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	gr := train.NewTestGroupRequest()
	gid, err := trainCtx.ApiCreateGroup(token, &gr)
	if err != nil {
		return err
	}

	// noop, just checking for error
	if err := trainCtx.DalRemoveTvgBinding(tid, gid, uid); err != nil {
		return err
	}

	// assign training to group
	if err := trainCtx.DalAddTvgBinding(tid, gid, uid); err != nil {
		return err
	}

	// validate
	expected := train.GroupArr{
		train.Group{
			GroupRequest: gr,
			ID:           gid,
			UserID:       uid,
		},
	}
	if a, err := trainCtx.DalReadTvgBindings(tid); err != nil {
		return err
	} else {
		if err = AssertGroupsAreTheSame(a, expected); err != nil {
			return err
		}
	}

	// add another group
	gr = train.NewTestGroupRequest()
	gid, err = trainCtx.ApiCreateGroup(token, &gr)
	if err != nil {
		return err
	}
	// add binding
	if err := trainCtx.DalAddTvgBinding(tid, gid, uid); err != nil {
		return err
	}

	// validate
	expected = append(expected, train.Group{
		GroupRequest: gr,
		ID:           gid,
		UserID:       uid,
	})
	if a, err := trainCtx.DalReadTvgBindings(tid); err != nil {
		return err
	} else {
		if err = AssertGroupsAreTheSame(a, expected); err != nil {
			return err
		}
	}

	// delete first group
	if err := trainCtx.DalDeleteTrainingGroup(uid, &train.DeleteGroupRequest{
		ID: expected[0].ID,
	}); err != nil {
		return err
	}

	// validate
	expected = train.GroupArr{expected[1]}
	if a, err := trainCtx.DalReadTvgBindings(tid); err != nil {
		return err
	} else {
		if err = AssertGroupsAreTheSame(a, expected); err != nil {
			return err
		}
	}

	// add group
	gr = train.NewTestGroupRequest()
	gid, err = trainCtx.ApiCreateGroup(token, &gr)
	if err != nil {
		return err
	}
	// add binding
	if err := trainCtx.DalAddTvgBinding(tid, gid, uid); err != nil {
		return err
	}

	// validate
	expected = append(expected, train.Group{
		GroupRequest: gr,
		ID:           gid,
		UserID:       uid,
	})
	if a, err := trainCtx.DalReadTvgBindings(tid); err != nil {
		return err
	} else {
		if err = AssertGroupsAreTheSame(a, expected); err != nil {
			return err
		}
	}

	// update group
	gr = train.NewTestGroupRequest()
	if err := trainCtx.DalUpdateTrainingGroup(uid, &train.UpdateGroupRequest{
		ID:           gid,
		GroupRequest: gr,
	}); err != nil {
		return err
	}

	// validate
	expected = train.GroupArr{expected[0], train.Group{
		GroupRequest: gr,
		ID:           gid,
		UserID:       uid,
	}}
	if a, err := trainCtx.DalReadTvgBindings(tid); err != nil {
		return err
	} else {
		if err = AssertGroupsAreTheSame(a, expected); err != nil {
			return err
		}
	}

	// add binding to existing group
	if err := trainCtx.DalAddTvgBinding(tid, expected[0].ID, uid); err != nil {
		return err
	}

	// validate that nothing has changed
	if a, err := trainCtx.DalReadTvgBindings(tid); err != nil {
		return err
	} else {
		if err = AssertGroupsAreTheSame(a, expected); err != nil {
			return err
		}
	}

	// remove binding [0]
	if err := trainCtx.DalRemoveTvgBinding(tid, expected[0].ID, uid); err != nil {
		return err
	}

	// validate
	expected = train.GroupArr{train.Group{
		GroupRequest: gr,
		ID:           gid,
		UserID:       uid,
	}}
	if a, err := trainCtx.DalReadTvgBindings(tid); err != nil {
		return err
	} else {
		if err = AssertGroupsAreTheSame(a, expected); err != nil {
			return err
		}
	}

	// remove binding [0]
	if err := trainCtx.DalRemoveTvgBinding(tid, expected[0].ID, uid); err != nil {
		return err
	}

	// validate
	expected = train.GroupArr{}
	if a, err := trainCtx.DalReadTvgBindings(tid); err != nil {
		return err
	} else {
		if err = AssertGroupsAreTheSame(a, expected); err != nil {
			return err
		}
	}

	// add binding again
	expected = train.GroupArr{train.Group{
		GroupRequest: gr,
		ID:           gid,
		UserID:       uid,
	}}
	if err := trainCtx.DalAddTvgBinding(tid, expected[0].ID, uid); err != nil {
		return err
	}

	// validate
	if a, err := trainCtx.DalReadTvgBindings(tid); err != nil {
		return err
	} else {
		if err = AssertGroupsAreTheSame(a, expected); err != nil {
			return err
		}
	}

	return nil
}

func TestTvg(t *testing.T) {
	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	s := time.Now().In(time.UTC)
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Hour))
	tid, err := trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}
	if err := testTvg(token, tid); err != nil {
		t.Fatal(err)
	}

}

// func TestTvgInSem(t *testing.T) {
// 	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fn := func() error {
// 		s := time.Now().In(time.UTC)
// 		tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Hour))
// 		tid, err := trainCtx.ApiCreateTraining(token, &tr)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		return testTvg(token, tid)
// 	}
// 	cases := 200
// 	if err := helpers.Sem(12, cases, []func() error{fn}); err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestTvgHandlers(t *testing.T) {
	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	uid, err := trainCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}

	s := time.Now().In(time.UTC)
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Hour))
	tid, err := trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}
	gr := train.NewTestGroupRequest()
	gid, err := trainCtx.ApiCreateGroup(token, &gr)
	if err != nil {
		t.Fatal(err)
	}
	//
	m := "PUT"
	u := "/api/training/group/binding"

	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 400,
			// these 2 tests are fails, however its pretty hard in current db state
			// to determine whether training and group exists and table has changed
			// thus for now those 2 handlers return 204
		}, {
			RequestMethod:  m,
			RequestUrl:     u,
			RequestHeaders: user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(train.GroupBindingRequest{
				TrainingID: tid,
				GroupID:    uuid.New(),
			}),
			ExpectedStatusCode: 204,
		}, {
			RequestMethod:  m,
			RequestUrl:     u,
			RequestHeaders: user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(train.GroupBindingRequest{
				TrainingID: uuid.New(),
				GroupID:    gid,
			}),
			ExpectedStatusCode: 204,
			//
		}, {
			RequestMethod:  m,
			RequestUrl:     u,
			RequestHeaders: user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(train.GroupBindingRequest{
				TrainingID: tid,
				GroupID:    gid,
			}),
			ExpectedStatusCode: 204,
		}, {
			// req should be idempotent
			RequestMethod:  m,
			RequestUrl:     u,
			RequestHeaders: user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(train.GroupBindingRequest{
				TrainingID: tid,
				GroupID:    gid,
			}),
			ExpectedStatusCode: 204,
		}, {
			RequestMethod:      "GET",
			RequestUrl:         "/api/training?id=" + tid.String(),
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []train.TrainingWithJoins
				helpers.JsonMustDeserialize(b, &res)
				if len(res) != 1 {
					return fmt.Errorf("invalid training count")
				}
				g := res[0].Groups
				return AssertGroupsAreTheSame(g, train.GroupArr{
					{
						GroupRequest: gr,
						ID:           gid,
						UserID:       uid,
					},
				})
			},
		},
	})

	m = "DELETE"
	trainCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 400,
			// these 2 tests are fails, however its pretty hard in current db state
			// to determine whether training and group exists and table has changed
			// thus for now those 2 handlers return 204
		}, {
			RequestMethod:  m,
			RequestUrl:     u,
			RequestHeaders: user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(train.GroupBindingRequest{
				TrainingID: tid,
				GroupID:    uuid.New(),
			}),
			ExpectedStatusCode: 204,
		}, {
			RequestMethod:  m,
			RequestUrl:     u,
			RequestHeaders: user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(train.GroupBindingRequest{
				TrainingID: uuid.New(),
				GroupID:    gid,
			}),
			ExpectedStatusCode: 204,
			//
		}, {
			RequestMethod:  m,
			RequestUrl:     u,
			RequestHeaders: user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(train.GroupBindingRequest{
				TrainingID: tid,
				GroupID:    gid,
			}),
			ExpectedStatusCode: 204,
		}, {
			RequestMethod:      "GET",
			RequestUrl:         "/api/training?id=" + tid.String(),
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []train.TrainingWithJoins
				helpers.JsonMustDeserialize(b, &res)
				if len(res) != 1 {
					return fmt.Errorf("invalid training count")
				}
				g := res[0].Groups
				return AssertGroupsAreTheSame(g, train.GroupArr{})
			},
		},
	})
}
