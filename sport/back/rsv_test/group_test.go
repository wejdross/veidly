package rsv_test

// import (
// 	"sport/adyen_sm"
// 	"sport/rsv"
// 	"sport/train"
// 	"testing"
// 	"time"
// )

// func TestGroupPeople(t *testing.T) {
// 	token, _, err := rsvCtx.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	uid, err := rsvCtx.Api.AuthorizeUserFromToken(token)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	s := time.Now().Add(time.Hour).In(time.UTC)
// 	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Hour))
// 	tr.Training.Capacity = 100 // i dont want to worry about it
// 	tr.Training.AllowExpress = true
// 	tid, err := rsvCtx.Train.ApiCreateTraining(token, &tr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// create groups
// 	grA := train.NewTestGroupRequest()
// 	grA.Name = "group A"
// 	grA.MaxPeople = 2
// 	grA.MaxTrainings = 1
// 	grAid, err := rsvCtx.Train.ApiCreateGroup(token, &grA)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// grB := train.NewTestGroupRequest()
// 	// grB.Name = "group B"
// 	// grB.MaxPeople = 1
// 	// grB.MaxTrainings = 1
// 	// grBid, err := rsvCtx.Train.ApiCreateGroup(token, &grB)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	// create bindings

// 	if err := rsvCtx.Train.DalAddTvgBinding(tid, grAid, uid); err != nil {
// 		t.Fatal(err)
// 	}
// 	// if err := rsvCtx.Train.DalAddTvgBinding(tid, grBid, uid); err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	// rsv

// 	rr := rsv.NewTestCreateReservationRequest(tid, s)
// 	if r, err := rsvCtx.ApiCreateAndQueryRsv(&token, &rr, 0); err != nil {
// 		t.Fatal(err)
// 	} else {
// 		if err := rsvCtx.AdyenSm.MoveStateToCapture(
// 			rsvCtx.RsvResponseToSmPassPtr(r), adyen_sm.ManualSrc, nil,
// 		); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	if r, err := rsvCtx.ApiCreateAndQueryRsv(nil, &rr, 0); err != nil {
// 		t.Fatal(err)
// 	} else {
// 		if err := rsvCtx.AdyenSm.MoveStateToCapture(
// 			rsvCtx.RsvResponseToSmPassPtr(r), adyen_sm.ManualSrc, nil,
// 		); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	if _, _, err := rsvCtx.ApiCreateRsv(nil, &rr, 409); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestGroupTrainings(t *testing.T) {
// 	token, _, err := rsvCtx.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	uid, err := rsvCtx.Api.AuthorizeUserFromToken(token)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// create trainings

// 	s1 := time.Now().In(time.UTC).Add(time.Hour)
// 	tr1 := train.NewTestCreateTrainingRequest(s1, s1.Add(time.Hour))
// 	tr1.Training.Capacity = 100 // i dont want to worry about it
// 	tr1.Training.AllowExpress = true
// 	tr1.Training.Title = "training 1"
// 	tr1id, err := rsvCtx.Train.ApiCreateTraining(token, &tr1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// this one doesnt overlap
// 	s2 := time.Now().In(time.UTC).Add(time.Hour * 2)
// 	tr2 := train.NewTestCreateTrainingRequest(s2, s2.Add(time.Hour))
// 	tr2.Training.Capacity = 100
// 	tr2.Training.AllowExpress = true
// 	tr2.Training.Title = "training 2"
// 	tr2id, err := rsvCtx.Train.ApiCreateTraining(token, &tr2)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// this one does
// 	s3 := time.Now().In(time.UTC).Add(time.Hour)
// 	tr3 := train.NewTestCreateTrainingRequest(s3, s3.Add(time.Minute*30))
// 	tr3.Training.Title = "training 3"
// 	tr3.Training.Capacity = 100
// 	tr3.Training.AllowExpress = true
// 	tr3id, err := rsvCtx.Train.ApiCreateTraining(token, &tr3)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// create groups
// 	grA := train.NewTestGroupRequest()
// 	grA.Name = "group A"
// 	grA.MaxPeople = 2
// 	grA.MaxTrainings = 1
// 	grAid, err := rsvCtx.Train.ApiCreateGroup(token, &grA)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// create bindings

// 	if err := rsvCtx.Train.DalAddTvgBinding(tr1id, grAid, uid); err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := rsvCtx.Train.DalAddTvgBinding(tr2id, grAid, uid); err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := rsvCtx.Train.DalAddTvgBinding(tr3id, grAid, uid); err != nil {
// 		t.Fatal(err)
// 	}

// 	// rsv

// 	rr1 := rsv.NewTestCreateReservationRequest(tr1id, s1)
// 	if r, err := rsvCtx.ApiCreateAndQueryRsv(nil, &rr1, 0); err != nil {
// 		t.Fatal(err)
// 	} else {
// 		if err := rsvCtx.AdyenSm.MoveStateToCapture(
// 			rsvCtx.RsvResponseToSmPassPtr(r), adyen_sm.ManualSrc, nil,
// 		); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// 	// this doesnt overlap so it should be ok
// 	rr2 := rsv.NewTestCreateReservationRequest(tr2id, s2)
// 	if r, err := rsvCtx.ApiCreateAndQueryRsv(nil, &rr2, 0); err != nil {
// 		t.Fatal(err)
// 	} else {
// 		if err := rsvCtx.AdyenSm.MoveStateToCapture(
// 			rsvCtx.RsvResponseToSmPassPtr(r), adyen_sm.ManualSrc, nil,
// 		); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// 	// this overlaps and group has max_trainings set to 1 so this should not be ok
// 	rr3 := rsv.NewTestCreateReservationRequest(tr3id, s3)
// 	if _, _, err := rsvCtx.ApiCreateRsv(nil, &rr3, 409); err != nil {
// 		t.Fatal(err)
// 	}
// }
