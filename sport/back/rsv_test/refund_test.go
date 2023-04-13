package rsv_test

// import (
// 	"sport/adyen_sm"
// 	"sport/helpers"
// 	"sport/rsv"
// 	"sport/user"
// 	"testing"
// 	"time"
// )

// func TestInstrRefundInFlow(t *testing.T) {
// 	ir := &user.UserRequest{
// 		Email:    helpers.CRNG_stringPanic(10) + "@foo.bar",
// 		Password: "sadasASDAS21312#!@#!@",
// 		UserData: user.UserData{
// 			Name:     "Velvet Velour",
// 			Language: "pl",
// 			Country:  "PL",
// 		},
// 	}

// 	prevPenaltyPercent := rsvCtx.AdyenSm.Config.InstrShotPenaltyPercent
// 	prevPenaltyNoMoreThan := rsvCtx.AdyenSm.Config.InstrShotNoMoreThanPercent
// 	svcFee := rsvCtx.Config.ServiceFee
// 	fis := rsvCtx.AdyenSm.Config.FreeInstrShots
// 	rsvCtx.AdyenSm.Config.InstrShotPenaltyPercent = 20
// 	rsvCtx.AdyenSm.Config.InstrShotNoMoreThanPercent = 50
// 	rsvCtx.Config.ServiceFee = 20
// 	rsvCtx.AdyenSm.Config.FreeInstrShots = 0
// 	// remove free instructor shots
// 	defer func() {
// 		rsvCtx.AdyenSm.Config.InstrShotPenaltyPercent = prevPenaltyPercent
// 		rsvCtx.AdyenSm.Config.InstrShotNoMoreThanPercent = prevPenaltyNoMoreThan
// 		rsvCtx.Config.ServiceFee = svcFee
// 		rsvCtx.AdyenSm.Config.FreeInstrShots = fis
// 	}()

// 	const rsvPrice = 100 * 100
// 	const expectedPenalty = 2000
// 	// move rsv to refund
// 	{
// 		instrToken := ""
// 		r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 			AltInstructorUserRq:  ir,
// 			RsvStartDateSinceNow: time.Hour,
// 			ReturnInstrToken:     &instrToken,
// 			RsvPrice:             rsvPrice,
// 		})
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if err = rsv.EnsureRsvIsInState(r, adyen_sm.LinkExpress); err != nil {
// 			t.Fatal(err)
// 		}

// 		// simulate positive auth from adyen
// 		if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 			t.Fatal(err)
// 		}

// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State: adyen_sm.WaitCapture,
// 		}); err != nil {
// 			t.Fatal(err)
// 		}

// 		if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Capture(r, true)); err != nil {
// 			t.Fatal(err)
// 		}

// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State: adyen_sm.Capture,
// 		}); err != nil {
// 			t.Fatal(err)
// 		}

// 		if err = rsvCtx.ApiRefundRsv(r.ID, instrToken, 0, true); err != nil {
// 			t.Fatal(err)
// 		}

// 		c, err := rsvCtx.Instr.DalReadNextPenalty(*r.InstructorID, nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if c != expectedPenalty {
// 			t.Fatalf("before payout expected penalty: %d, instead got: %d", expectedPenalty, c)
// 		}

// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State: adyen_sm.WaitRefund,
// 		}); err != nil {
// 			t.Fatal(err)
// 		}

// 		if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Refund(r, true)); err != nil {
// 			t.Fatal(err)
// 		}

// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State: adyen_sm.Refund,
// 		}); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	// create another rsv and move to payout

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 		AltInstructorUserRq:  ir,
// 		RsvPrice:             1000,
// 		UsersMayExist:        true,
// 	})

// 	if err := executeExpressFlowForRsv(r); err != nil {
// 		t.Fatal(err)
// 	}

// 	// payout value for price 1000 is 800 (i set svc fee to 20%)
// 	// 50% of that is 400 and this should be extracted
// 	// remaining penalty should be 2000 - 400 = 1600
// 	c, err := rsvCtx.Instr.DalReadNextPenalty(*r.InstructorID, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if c != 1600 {
// 		t.Fatalf("after payout expected penalty: %d, instead got: %d", 1600, c)
// 	}

// 	r, err = rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 		AltInstructorUserRq:  ir,
// 		RsvPrice:             100000,
// 		UsersMayExist:        true,
// 	})

// 	if err := executeExpressFlowForRsv(r); err != nil {
// 		t.Fatal(err)
// 	}

// 	c, err = rsvCtx.Instr.DalReadNextPenalty(*r.InstructorID, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if c != 0 {
// 		t.Fatalf("after payout expected penalty: %d, instead got: %d", 0, c)
// 	}
// }
