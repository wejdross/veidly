package train

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"sport/api"
	"sport/helpers"
	"sport/user"
	"time"

	"github.com/google/uuid"
)

func NewTestCreateTrainingRequestNoOcc() CreateTrainingRequest {
	return CreateTrainingRequest{
		Training: TrainingRequest{
			Title:         helpers.CRNG_stringPanic(12),
			Capacity:      rand.Int()%10 + 1,
			Currency:      "PLN",
			Price:         10 * 100,
			ManualConfirm: false,
			AllowExpress:  true,
		},
	}
}

func NewTestCreateTrainingRequest(start, end time.Time) CreateTrainingRequest {
	return CreateTrainingRequest{
		Training: TrainingRequest{
			Title:         helpers.CRNG_stringPanic(12),
			Capacity:      1,
			Currency:      "PLN",
			Price:         10 * 100,
			ManualConfirm: false,
			AllowExpress:  true,
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
			MinAge: 1,
			MaxAge: 40,
		},
		Occurrences: []CreateOccRequest{
			{
				OccRequest: OccRequest{
					DateStart:  start,
					DateEnd:    end,
					RepeatDays: 7,
				},
			},
		},
	}
}

func NewTestOccRequest(start, end time.Time) OccRequest {
	return OccRequest{
		DateStart:  start,
		DateEnd:    end,
		RepeatDays: 1,
	}
}

func NewTest2OccRequest(ostart, oend int) SecondaryOccRequest {
	return SecondaryOccRequest{
		OffsetStart: ostart,
		OffsetEnd:   oend,
		Color:       "",
		Remarks:     helpers.CRNG_stringPanic(12),
	}
}

func (ctx *Ctx) ApiCreateTraining(
	token string, cr *CreateTrainingRequest,
) (uuid.UUID, error) {
	if !cr.ReturnID {
		cr.ReturnID = true
	}
	var ret uuid.UUID
	tc := []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/training",
			RequestReader:      helpers.JsonMustSerializeReader(cr),
			ExpectedStatusCode: 200,
			RequestHeaders:     api.JwtHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				s := string(b)
				if s == "" {
					return fmt.Errorf("empty id returned")
				}
				var err error
				ret, err = uuid.Parse(s)
				return err
			},
		},
	}
	err := ctx.Api.TestAssertCasesErr(tc)
	return ret, err
}

func (ctx *Ctx) ApiUploadImgFromPath(t, p string, trainingID uuid.UUID) error {

	fc, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	r, ct, err := api.CreateMultipartFormWithValues(
		"image", path.Base(p), fc,
		[]api.AdditionalFormField{
			{
				Key: "training_id",
				Val: trainingID.String(),
			},
			{
				Key: "main",
				Val: "1",
			},
		})
	h := user.GetAuthHeader(t)
	h["Content-Type"] = ct

	cases := []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/training/img",
			RequestHeaders:     h,
			RequestReader:      &r,
			ExpectedStatusCode: 204,
		},
	}

	return ctx.Api.TestAssertCasesErr(cases)
}

func NewTestGroupRequest() GroupRequest {
	return GroupRequest{
		Name:         helpers.CRNG_stringPanic(12),
		MaxPeople:    1,
		MaxTrainings: 1,
	}
}

func (ctx *Ctx) ApiCreateGroup(token string, tr *GroupRequest) (uuid.UUID, error) {
	if tr == nil {
		_tr := NewTestGroupRequest()
		tr = &_tr
	}
	var id uuid.UUID
	err := ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         "/api/training/group",
		ExpectedStatusCode: 200,
		RequestReader:      helpers.JsonMustSerializeReader(tr),
		RequestHeaders:     user.GetAuthHeader(token),
		ExpectedBodyVal: func(b []byte, i interface{}) error {
			var res struct{ ID uuid.UUID }
			if err := json.Unmarshal(b, &res); err != nil {
				return err
			}
			if res.ID == uuid.Nil {
				return fmt.Errorf("invalid group id")
			}
			id = res.ID
			return nil
		},
	})
	return id, err
}
