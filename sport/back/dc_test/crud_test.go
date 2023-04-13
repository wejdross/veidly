package dc_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"sport/api"
	"sport/dc"
	"sport/helpers"
	"sport/train"
	"sport/user"
	"testing"
	"time"

	"github.com/google/uuid"
)

func AssertGetRes(
	rawRes []byte,
	expected []dc.DcRequest,
	creatorsIid uuid.UUID,
	ids map[string]uuid.UUID) error {
	//
	var res []dc.Dc
	if err := json.Unmarshal(rawRes, &res); err != nil {
		return err
	}
	if len(res) != len(expected) {
		return fmt.Errorf("AssertGetRes: invalid length")
	}
	for i := range expected {
		found := false
		for j := range res {
			if err := helpers.AssertJsonErr("", res[j].DcRequest, expected[i]); err != nil {
				continue
			}
			if res[j].InstrID != creatorsIid {
				return fmt.Errorf("invalid InstrID")
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

func TestDcCrud(t *testing.T) {
	token, iid, err := instrCtx.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	// POST
	m := "POST"
	u := "/api/dc"

	s := time.Now().In(time.UTC).Round(time.Second)

	dcs := []dc.DcRequest{
		{
			Name:       "foo",
			Quantity:   1,
			ValidStart: dc.Time(s),
			ValidEnd:   dc.Time(s.Add(time.Hour)),
			Discount:   30,
		},
		{
			Name:       "bar",
			Quantity:   2,
			ValidStart: dc.Time(s),
			ValidEnd:   dc.Time(s.Add(time.Hour)),
			Discount:   30,
		},
		{
			Name:       "foobar",
			Quantity:   30,
			ValidStart: dc.Time(s),
			ValidEnd:   dc.Time(s.Add(time.Hour)),
			Discount:   32,
		},
	}

	dcCtx.Api.TestAssertCases(t, []api.TestCase{
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
			RequestReader: helpers.JsonMustSerializeReader(dc.DcRequest{
				Name: "",
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(dcs[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(dcs[1]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(dcs[2]),
			RequestHeaders:     user.GetAuthHeader(token),
		},
	})

	m = "GET"
	ids := make(map[string]uuid.UUID)

	dcCtx.Api.TestAssertCases(t, []api.TestCase{
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
				return AssertGetRes(b, dcs, iid, ids)
			},
		},
	})

	m = "PATCH"

	s = time.Now().In(time.UTC).Round(time.Hour)
	dcs = []dc.DcRequest{
		{
			Name:       "foo",
			Quantity:   10,
			ValidStart: dc.Time(s),
			ValidEnd:   dc.Time(s.Add(time.Hour)),
			Discount:   44,
		},
		{
			Name:       "bar",
			Quantity:   20,
			ValidStart: dc.Time(s),
			ValidEnd:   dc.Time(s.Add(time.Hour)),
			Discount:   22,
		},
		{
			Name:       "foobar",
			Quantity:   30,
			ValidStart: dc.Time(s),
			ValidEnd:   dc.Time(s.Add(time.Hour)),
			Discount:   32,
		},
	}

	newGrpsUpdateReqs := make([]dc.UpdateDcRequest, len(dcs))
	for i := range dcs {
		id := ids[dcs[i].Name]
		if id == uuid.Nil {
			t.Fatal("invalid id")
		}
		newGrpsUpdateReqs[i] = dc.UpdateDcRequest{
			ID:        ids[dcs[i].Name],
			DcRequest: dcs[i],
		}
	}

	dcCtx.Api.TestAssertCases(t, []api.TestCase{
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
			RequestReader: helpers.JsonMustSerializeReader(dc.DcRequest{
				Name: "",
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestReader:      helpers.JsonMustSerializeReader(dcs[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(dc.UpdateDcRequest{
				ID:        uuid.New(),
				DcRequest: newGrpsUpdateReqs[0].DcRequest,
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
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(newGrpsUpdateReqs[2]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      "GET",
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return AssertGetRes(b, dcs, iid, ids)
			},
		},
	})

	m = "DELETE"
	dcs = []dc.DcRequest{
		// first one 'foo' is removed
		dcs[1],
		dcs[2],
	}

	newGrpsDeleteReqs := make([]dc.DeleteDcRequest, len(dcs))
	newGrpsDeleteReqs[0] = dc.DeleteDcRequest{
		ID: ids["foo"],
	}

	dcCtx.Api.TestAssertCases(t, []api.TestCase{
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
			RequestReader:      helpers.JsonMustSerializeReader(dcs[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(dc.DeleteDcRequest{
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
				return AssertGetRes(b, dcs, iid, ids)
			},
		},
	})

	var d, d2 *dc.Dc

	// existing group
	dcCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "GET",
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []dc.Dc
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				if len(res) != 2 {
					panic(1)
				}
				for i := range res {
					if res[i].Name == "bar" {
						d = &res[i]
					}
					if res[i].Name == "foobar" {
						d2 = &res[i]
					}
				}
				if d2 == nil || d == nil {
					return errors.New("invalid grps for bindings")
				}
				return nil
			},
		},
	})

	// bindings
	m = "POST"
	u = "/api/dc/binding"

	tr := train.NewTestCreateTrainingRequestNoOcc()
	tid, err := trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}

	//
	t2, _, err := instrCtx.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	tri2 := train.NewTestCreateTrainingRequestNoOcc()
	tid_i2, err := trainCtx.ApiCreateTraining(t2, &tri2)
	if err != nil {
		t.Fatal(err)
	}

	dcCtx.Api.TestAssertCases(t, []api.TestCase{
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
			RequestReader: helpers.JsonMustSerializeReader(dc.DcBindingRequest{
				TrainingID: uuid.Nil,
				DcID:       d.ID,
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
			RequestReader: helpers.JsonMustSerializeReader(dc.DcBindingRequest{
				TrainingID: tid,
				DcID:       uuid.Nil,
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(dc.DcBindingRequest{
				TrainingID: uuid.New(),
				DcID:       d.ID,
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(dc.DcBindingRequest{
				TrainingID: tid,
				DcID:       uuid.New(),
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 404,
			RequestReader: helpers.JsonMustSerializeReader(dc.DcBindingRequest{
				TrainingID: tid_i2,
				DcID:       d.ID,
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader: helpers.JsonMustSerializeReader(dc.DcBindingRequest{
				TrainingID: tid,
				DcID:       d.ID,
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader: helpers.JsonMustSerializeReader(dc.DcBindingRequest{
				TrainingID: tid,
				DcID:       d2.ID,
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      "GET",
			RequestUrl:         "/api/training",
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var x []*train.TrainingWithJoins
				helpers.JsonMustDeserialize(b, &x)
				if len(x) != 1 {
					return fmt.Errorf("invalid trainings 1")
				}
				return nil
			},
		}, {
			RequestMethod:      "GET",
			RequestUrl:         "/api/training?dc_id=" + d.ID.String(),
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var x []*train.TrainingWithJoins
				helpers.JsonMustDeserialize(b, &x)
				if len(x) != 1 {
					return fmt.Errorf("invalid trainings 2")
				}
				if x[0].Training.Title != tr.Training.Title {
					return fmt.Errorf("invalid training")
				}
				if len(x[0].Dcs) != 1 {
					return fmt.Errorf("invalid Dcs")
				}
				g := x[0].Dcs[0]
				if err := helpers.AssertJsonErr("assert train dc", &g, &d); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
