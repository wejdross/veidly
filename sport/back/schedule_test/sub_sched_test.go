package schedule_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"sport/adyen_sm"
// 	"sport/api"
// 	"sport/rsv"
// 	"sport/schedule"
// 	"sport/sub"
// 	"sport/train"
// 	"testing"
// 	"time"
// )

// func TestSubSched(t *testing.T) {

// 	iToken, _, err := rsvCtx.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	uToken, err := rsvCtx.User.ApiCreateAndLoginUser(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	tDateStart := rsv.NowMin().Add(time.Hour * 24 * 7) // 7 days
// 	tDateEnd := tDateStart.Add(time.Hour)
// 	tReq := train.NewTestCreateTrainingRequest(tDateStart, tDateEnd)
// 	tid, err := rsvCtx.Train.ApiCreateTraining(iToken, &tReq)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// valid
// 	sm1req := sub.NewTestSubModelRequest()
// 	sm1req.Duration = 8

// 	// not valid
// 	sm2req := sub.NewTestSubModelRequest()
// 	sm2req.Duration = 3

// 	// no training
// 	sm3req := sub.NewTestSubModelRequest()
// 	sm1req.Duration = 8

// 	// create test sms
// 	sm1id, err := rsvCtx.Sub.ApiCreateSubModel(iToken, &sm1req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sm2id, err := rsvCtx.Sub.ApiCreateSubModel(iToken, &sm2req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sm3id, err := rsvCtx.Sub.ApiCreateSubModel(iToken, &sm3req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// create sm bindings

// 	if err := rsvCtx.Sub.ApiCreateSmBinding(iToken, &sub.SubModelBinding{
// 		SubModelID: sm1id,
// 		TrainingID: tid,
// 	}); err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := rsvCtx.Sub.ApiCreateSmBinding(iToken, &sub.SubModelBinding{
// 		SubModelID: sm2id,
// 		TrainingID: tid,
// 	}); err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := rsvCtx.Sub.ApiCreateSmBinding(iToken, &sub.SubModelBinding{
// 		SubModelID: sm3id,
// 		TrainingID: tid,
// 	}); err != nil {
// 		t.Fatal(err)
// 	}

// 	// create subs

// 	sub1, err := rsvCtx.Sub.ApiCreateAndQuerySub(uToken, sm1id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sub2, err := rsvCtx.Sub.ApiCreateAndQuerySub(uToken, sm2id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sub3, err := rsvCtx.Sub.ApiCreateAndQuerySub(uToken, sm3id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// confirm subs

// 	p := rsvCtx.Sub.SubToSmPass(sub1)
// 	if err := rsvCtx.Sub.AdyenSm.MoveStateToCapture(&p, adyen_sm.ManualSrc, nil); err != nil {
// 		t.Fatal(err)
// 	}
// 	p = rsvCtx.Sub.SubToSmPass(sub2)
// 	if err := rsvCtx.Sub.AdyenSm.MoveStateToCapture(&p, adyen_sm.ManualSrc, nil); err != nil {
// 		t.Fatal(err)
// 	}
// 	p = rsvCtx.Sub.SubToSmPass(sub3)
// 	if err := rsvCtx.Sub.AdyenSm.MoveStateToCapture(&p, adyen_sm.ManualSrc, nil); err != nil {
// 		t.Fatal(err)
// 	}

// 	rsvCtx.Api.TestAssertCaseErr(&api.TestCase{
// 		RequestMethod: "GET",
// 		RequestUrl: fmt.Sprintf("/api/schedule?start=%d&end=%d&flags=1",
// 			tDateStart.Unix(), tDateStart.Unix()),
// 		ExpectedStatusCode: 200,
// 		RequestHeaders: map[string]string{
// 			"Authorization": "Bearer " + iToken,
// 		},
// 		ExpectedBodyVal: func(b []byte, i interface{}) error {
// 			var res []schedule.TrainingSchedule
// 			if err = json.Unmarshal(b, &res); err != nil {
// 				return err
// 			}
// 			if len(res[0].Schedule[0].Subs) != 2 {
// 				return fmt.Errorf("invalid subs")
// 			}
// 			for i := range res[0].Schedule[0].Subs {
// 				s := res[0].Schedule[0].Subs[i]
// 				if s.SubModel.Name != sub2.SubModel.Name && s.SubModel.Name != sub1.SubModel.Name {
// 					return fmt.Errorf("invalid subs")
// 				}
// 			}
// 			return nil
// 		},
// 	})

// }
