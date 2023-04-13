package rsv_test

import (
	"sport/rsv"
	"sport/train"
	"sport/user"
	"testing"
	"time"
)

func TestDisabledRsv(t *testing.T) {
	it, err := rsvCtx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}
	iuid, err := rsvCtx.Api.AuthorizeUserFromToken(it)
	if err != nil {
		t.Fatal(err)
	}
	iid, err := rsvCtx.Instr.CreateTestInstructor(iuid, nil)
	if err != nil {
		t.Fatal(err)
	}
	start := rsv.NowMin().Add(time.Hour)
	end := start.Add(time.Hour)
	tr := train.TrainingRequest{
		Title:         "vesuvius",
		Capacity:      1,
		Currency:      "PLN",
		Price:         10 * 100,
		ManualConfirm: false,
		AllowExpress:  true,
		Disabled:      true,
	}
	trainingID, err := rsvCtx.Train.ApiCreateTraining(it, &train.CreateTrainingRequest{
		Training: tr,
		Occurrences: []train.CreateOccRequest{
			{
				OccRequest: train.OccRequest{
					DateStart:  start,
					DateEnd:    end,
					RepeatDays: 1,
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	req := rsv.ApiReservationRequest{
		TrainingID: trainingID,
		Occurrence: start,
		UserData: user.UserData{
			Name:     "foo",
			Language: "pl",
		},
	}

	if _, _, err = rsvCtx.ApiCreateRsv(nil, &req, 409); err != nil {
		t.Fatal(err)
	}

	tr.Disabled = false

	if err := rsvCtx.Train.DalUpdateTraining(
		&train.UpdateTrainingRequest{
			ID:              trainingID,
			TrainingRequest: tr,
		}, iid, iuid); err != nil {
		t.Fatal(err)
	}

	if _, _, err = rsvCtx.ApiCreateRsv(nil, &req, 0); err != nil {
		t.Fatal(err)
	}

}
