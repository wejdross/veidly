package sub_test

// import (
// 	"sport/adyen_sm"
// 	"sport/helpers"
// 	"sport/sub"
// 	"testing"
// 	"time"
// )

// /*
// 	link_express
// 		-> wait_capture
// 		-> error
// */
// func testWaitCaptureTimeouts() error {

// 	r, err := subCtx.PrepSubTesting(sub.EmptySubTestOpts)
// 	if err != nil {
// 		return helpers.WrapErr("prep", err)
// 	}

// 	if err = sub.EnsureSubIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return helpers.WrapErr("sub.EnsureSubIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State:        adyen_sm.WaitCapture,
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return helpers.WrapErr("WaitCapture", err)
// 	}

// 	for i := 1; i <= subCtx.AdyenSm.Config.SmMaxRetries; i++ {
// 		if err := subCtx.AdyenSm.AgentDoOne(subCtx.SubToSmPassPtr(r)); err != nil {
// 			return helpers.WrapErr("WaitCapture timeout iter", err)
// 		}
// 		if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 			State:        adyen_sm.WaitCapture,
// 			ResetTimeout: true,
// 			SmRetries:    i,
// 		}); err != nil {
// 			return helpers.WrapErr("WaitCapture", err)
// 		}
// 	}

// 	if err := subCtx.AdyenSm.AgentDoOne(subCtx.SubToSmPassPtr(r)); err != nil {
// 		return helpers.WrapErr("WaitCapture timeout iter", err)
// 	}
// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.Error,
// 	}); err != nil {
// 		return helpers.WrapErr("WaitCapture", err)
// 	}

// 	return nil
// }

// func TestWaitCaptureTimeouts(t *testing.T) {
// 	if err := testWaitCaptureTimeouts(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// /*
// 	link_express
// 		-> { wait_capture -> hold } * limit
// 		-> wait_cancel_or_refund
// 		-> cancel_or_refund
// */
// func testFailCapture() error {

// 	r, err := subCtx.PrepSubTesting(sub.EmptySubTestOpts)
// 	if err != nil {
// 		return helpers.WrapErr("prep", err)
// 	}

// 	if err = sub.EnsureSubIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return helpers.WrapErr("sub.EnsureSubIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_PositiveAuth", err)
// 	}

// 	for i := 1; i <= subCtx.AdyenSm.Config.SmMaxRetries; i++ {
// 		if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 			State: adyen_sm.WaitCapture,
// 		}); err != nil {
// 			return helpers.WrapErr("WaitCapture", err)
// 		}

// 		// simulate positive capture adyen.Notification
// 		if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_Capture(r, false)); err != nil {
// 			return helpers.WrapErr("sub.AdyenNtf_Capture", err)
// 		}

// 		if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 			State:        adyen_sm.Hold,
// 			Timeout:      time.Now().Add(time.Duration(subCtx.AdyenSm.Config.SmRetryTimeout)),
// 			ResetTimeout: true,
// 		}); err != nil {
// 			return helpers.WrapErr("Hold", err)
// 		}

// 		// daemon should perform payout request
// 		if err = subCtx.AdyenSm.AgentDoOne(subCtx.SubToSmPassPtr(r)); err != nil {
// 			return helpers.WrapErr("daemon Hold->WaitCapture", err)
// 		}
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return helpers.WrapErr("WaitCapture", err)
// 	}

// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_Capture(r, false)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_Capture", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.WaitCancelOrRefund,
// 	}); err != nil {
// 		return helpers.WrapErr("WaitCancelOrRefund", err)
// 	}

// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_CancelOrRefund(r, true)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_CancelOrRefund", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.CancelOrRefund,
// 	}); err != nil {
// 		return helpers.WrapErr("WaitCancelOrRefund", err)
// 	}

// 	return nil
// }

// func TestFailCapture(t *testing.T) {
// 	if err := testFailCapture(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// /*
// 	link_express
// 		-> wait_capture
// 		-> capture
// 		-> dispute
// */
// func testDispute() error {
// 	userToken := ""

// 	r, err := subCtx.PrepSubTesting(sub.SubTestOpts{
// 		ReturnUserToken: &userToken,
// 	})
// 	if err != nil {
// 		return helpers.WrapErr("prep", err)
// 	}

// 	if err = sub.EnsureSubIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return helpers.WrapErr("sub.EnsureSubIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return helpers.WrapErr("WaitCapture", err)
// 	}

// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_Capture(r, true)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_Capture", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.Capture,
// 	}); err != nil {
// 		return helpers.WrapErr("Capture", err)
// 	}

// 	if err = subCtx.ApiRefundSub(r.ID, userToken, 404, false); err != nil {
// 		return helpers.WrapErr("apiRefundRequest", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.Capture,
// 	}); err != nil {
// 		return helpers.WrapErr("Capture", err)
// 	}

// 	if err = subCtx.ApiDisputeSub(r.ID, userToken, 0, false); err != nil {
// 		return helpers.WrapErr("apiDisputeSub", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.Dispute,
// 	}); err != nil {
// 		return helpers.WrapErr("Dispute", err)
// 	}

// 	return nil
// }

// func TestDispute(t *testing.T) {
// 	if err := testDispute(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// /*
// 	link_express
// 		-> wait_capture
// 		-> capture
// 		-> wait_payout
// 		-> payout
// 		-> dispute
// */
// func testPayoutDispute() error {

// 	userToken := ""

// 	r, err := subCtx.PrepSubTesting(sub.SubTestOpts{
// 		ReturnUserToken: &userToken,
// 	})
// 	if err != nil {
// 		return helpers.WrapErr("prep", err)
// 	}

// 	if err = sub.EnsureSubIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return helpers.WrapErr("sub.EnsureSubIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return helpers.WrapErr("WaitCapture", err)
// 	}

// 	// simulate positive capture adyen.Notification
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositiveCapture(r)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_PositiveCapture", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State:        adyen_sm.Capture,
// 		Timeout:      time.Now().In(time.UTC).Add(time.Duration(subCtx.Config.PayoutDelay)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return helpers.WrapErr("Capture", err)
// 	}

// 	// daemon should perform payout request
// 	if err = subCtx.AdyenSm.AgentDoOne(subCtx.SubToSmPassPtr(r)); err != nil {
// 		return helpers.WrapErr("daemon Capture->WaitPayout", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State:   adyen_sm.WaitPayout,
// 		Timeout: time.Now().In(time.UTC).Add(time.Duration(subCtx.AdyenSm.Config.InstantPayoutTimeout)),
// 	}); err != nil {
// 		return helpers.WrapErr("WaitPayout", err)
// 	}

// 	// simulate positive payout adyen.Notification
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositivePayout(r)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_PositivePayout", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.Payout,
// 	}); err != nil {
// 		return helpers.WrapErr("StatePayout", err)
// 	}

// 	if err = subCtx.ApiDisputeSub(r.ID, userToken, 0, false); err != nil {
// 		return helpers.WrapErr("apiDisputeRequest", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.Dispute,
// 	}); err != nil {
// 		return helpers.WrapErr("Dispute", err)
// 	}

// 	return nil
// }

// func TestPayoutDispute(t *testing.T) {
// 	if err := testPayoutDispute(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// /*
// 	link_express
// 		-> wait_capture
// 		-> { capture -> wait_payout } * limit
// 		-> error
// */
// func testFailPayout() error {
// 	userToken := ""

// 	r, err := subCtx.PrepSubTesting(sub.SubTestOpts{
// 		ReturnUserToken: &userToken,
// 	})
// 	if err != nil {
// 		return helpers.WrapErr("prep", err)
// 	}

// 	if err = sub.EnsureSubIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return helpers.WrapErr("sub.EnsureSubIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return helpers.WrapErr("WaitCapture", err)
// 	}

// 	// simulate positive capture adyen.Notification
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositiveCapture(r)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_PositiveCapture", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State:        adyen_sm.Capture,
// 		Timeout:      time.Now().In(time.UTC).Add(time.Duration(subCtx.Config.PayoutDelay)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return helpers.WrapErr("Capture", err)
// 	}

// 	for i := 1; i <= subCtx.AdyenSm.Config.SmMaxRetries; i++ {

// 		if err = subCtx.AdyenSm.AgentDoOne(subCtx.SubToSmPassPtr(r)); err != nil {
// 			return helpers.WrapErr("daemon Capture->WaitPayout", err)
// 		}

// 		if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 			State:   adyen_sm.WaitPayout,
// 			Timeout: time.Now().In(time.UTC).Add(time.Duration(subCtx.AdyenSm.Config.InstantPayoutTimeout)),
// 		}); err != nil {
// 			return helpers.WrapErr("WaitPayout", err)
// 		}

// 		// simulate positive payout adyen.Notification
// 		if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_Payout(r, false)); err != nil {
// 			return helpers.WrapErr("sub.AdyenNtf_Payout", err)
// 		}

// 		if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 			State:        adyen_sm.Capture,
// 			Timeout:      time.Now().Add(time.Duration(subCtx.AdyenSm.Config.SmRetryTimeout)),
// 			ResetTimeout: true,
// 		}); err != nil {
// 			return helpers.WrapErr("Capture", err)
// 		}
// 	}

// 	if err = subCtx.AdyenSm.AgentDoOne(subCtx.SubToSmPassPtr(r)); err != nil {
// 		return helpers.WrapErr("daemon out of loop Capture->WaitPayout", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State:   adyen_sm.WaitPayout,
// 		Timeout: time.Now().In(time.UTC).Add(time.Duration(subCtx.AdyenSm.Config.InstantPayoutTimeout)),
// 	}); err != nil {
// 		return helpers.WrapErr("WaitPayout", err)
// 	}

// 	// simulate positive payout adyen.Notification
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_Payout(r, false)); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_Payout", err)
// 	}

// 	if r, err = subCtx.RefreshSubAndValidate(r.ID, sub.SubValidationOpts{
// 		State: adyen_sm.Error,
// 	}); err != nil {
// 		return helpers.WrapErr("Error", err)
// 	}

// 	return nil
// }

// func TestFailPayout(t *testing.T) {
// 	if err := testFailPayout(); err != nil {
// 		t.Fatal(err)
// 	}
// }
