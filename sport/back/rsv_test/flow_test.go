package rsv_test

// import (
// 	"fmt"
// 	"sport/adyen_sm"
// 	"sport/rsv"
// 	"testing"
// 	"time"

// 	"github.com/google/uuid"
// )

// /*
// 	following are automatic versions of ../adyen_test/*_flow.go

// 	those tests will be run with all of the other tests,
// 	this will mock adyen responses / webhook requests
// 	instead of using real ones
// */

// func getAndEnsureRsvIsInState(
// 	id uuid.UUID,
// 	state adyen_sm.State,
// ) error {
// 	rr, err := getRsvByID(id)
// 	if err != nil {
// 		return err
// 	}
// 	return rsv.EnsureRsvIsInState(rr, state)
// }

// func getRsvByID(id uuid.UUID) (*rsv.DDLRsvWithInstr, error) {
// 	return rsvCtx.ReadRsvByID(id)
// }

// func wrapErr(prefix string, err error) error {
// 	return fmt.Errorf("%s: %s", prefix, err)
// }

// func executeExpressFlowForRsv(r *rsv.DDLRsvWithInstr) error {
// 	var err error

// 	if err = rsv.EnsureRsvIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return wrapErr("rsv.EnsureRsvIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return wrapErr("WaitCapture", err)
// 	}

// 	// simulate positive capture adyen.Notification
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveCapture(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveCapture", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:        adyen_sm.Capture,
// 		Timeout:      r.DateEnd.Add(time.Duration(rsvCtx.Config.PayoutDelay)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Capture", err)
// 	}

// 	// daemon should perform payout request
// 	if err = rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 		return wrapErr("daemon Capture->WaitPayout", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:   adyen_sm.WaitPayout,
// 		Timeout: time.Now().In(time.UTC).Add(time.Duration(rsvCtx.AdyenSm.Config.InstantPayoutTimeout)),
// 	}); err != nil {
// 		return wrapErr("WaitPayout", err)
// 	}

// 	// simulate positive payout adyen.Notification
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositivePayout(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositivePayout", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Payout,
// 	}); err != nil {
// 		return wrapErr("StatePayout", err)
// 	}

// 	return nil
// }

// /*
// 	link_express -> wait_capture -> capture -> wait_payout -> payout
// */
// func testExpressFlow() error {

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	return executeExpressFlowForRsv(r)
// }

// /*
// 	link -> hold -> wait_capture -> capture -> wait_payout -> payout
// */
// func testNormalFlow() error {
// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour*24 + time.Minute,
// 		DisableLinkExpress:   true,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if err = rsv.EnsureRsvIsInState(r, adyen_sm.Link); err != nil {
// 		return wrapErr("Link", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:        adyen_sm.Hold,
// 		Timeout:      time.Now().In(time.UTC).Add(time.Hour * 24),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Hold", err)
// 	}

// 	// daemon should perform capture request
// 	if err = rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 		return wrapErr("Hold->WaitCapture", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return wrapErr("WaitCapture", err)
// 	}

// 	// simulate positive capture adyen.Notification
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveCapture(r)); err != nil {
// 		return err
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:        adyen_sm.Capture,
// 		Timeout:      r.DateEnd.Add(time.Duration(rsvCtx.Config.PayoutDelay)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Capture", err)
// 	}

// 	// daemon should perform payout request
// 	if err = rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 		return wrapErr("Capture->WaitPayout", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:   adyen_sm.WaitPayout,
// 		Timeout: time.Now().In(time.UTC).Add(time.Duration(time.Duration(rsvCtx.AdyenSm.Config.InstantPayoutTimeout))),
// 	}); err != nil {
// 		return wrapErr("WaitPayout", err)
// 	}

// 	// simulate positive payout adyen.Notification
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositivePayout(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositivePayout", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Payout,
// 	}); err != nil {
// 		return wrapErr("StatePayout", err)
// 	}

// 	return nil
// }

// /*
// 	link -> link_expire
// */
// func testNormalLinkExpireFlow() error {
// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour*24 + time.Minute,
// 		DisableLinkExpress:   true,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		Rsv:          r,
// 		State:        adyen_sm.Link,
// 		Timeout:      time.Now().In(time.UTC).Add(time.Duration(rsvCtx.Config.LinkExpire)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Link", err)
// 	}

// 	if err = rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 		return wrapErr("Link->LinkExpire", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.LinkExpire,
// 	}); err != nil {
// 		return wrapErr("LinkExpire", err)
// 	}

// 	return nil
// }

// /*
// 	link -> link_expire
// */
// func testForceLinkExpireFlow() error {
// 	var utoken string
// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour*24 + time.Minute,
// 		DisableLinkExpress:   true,
// 		ReturnUserToken:      &utoken,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		Rsv:     r,
// 		State:   adyen_sm.Link,
// 		Timeout: time.Now().In(time.UTC).Add(time.Duration(rsvCtx.Config.LinkExpire)),
// 	}); err != nil {
// 		return wrapErr("Link", err)
// 	}

// 	if err = rsvCtx.ApiExpireRsv(r.ID, utoken, 0); err != nil {
// 		return wrapErr("Link->LinkExpire", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.LinkExpire,
// 	}); err != nil {
// 		return wrapErr("LinkExpire", err)
// 	}

// 	return nil
// }

// /*
// 	link_express -> link_expire

// 	will also test that loop is created if receiving false auth notifications
// */
// func testExpressLinkExpireFlow() error {
// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		Rsv:     r,
// 		State:   adyen_sm.LinkExpress,
// 		Timeout: time.Now().Add(time.Duration(rsvCtx.Config.LinkExpressExpire)),
// 	}); err != nil {
// 		return wrapErr("LinkExpress", err)
// 	}

// 	// simulate false auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Auth(r, rsv.AdyenNtf_AuthOpts{
// 		Success: false,
// 	})); err != nil {
// 		return wrapErr("rsv.AdyenNtf_Auth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:        adyen_sm.LinkExpress,
// 		Timeout:      time.Now().Add(time.Duration(rsvCtx.Config.LinkExpressExpire)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("LinkExpress", err)
// 	}

// 	if err = rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 		return wrapErr("LinkExpress->LinkExpire", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.LinkExpire,
// 	}); err != nil {
// 		return wrapErr("LinkExpire", err)
// 	}

// 	return nil
// }

// /*
// 	link
// 		-> hold
// 		-> wait_cancel_or_refund
// 		-> cancel_or_refund
// */
// func testSuccessCancelFlow() error {

// 	var userToken string
// 	var instrToken string

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour*24 + time.Minute,
// 		DisableLinkExpress:   true,
// 		ReturnUserToken:      &userToken,
// 		ReturnInstrToken:     &instrToken,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		Rsv:          r,
// 		State:        adyen_sm.Link,
// 		Timeout:      time.Now().In(time.UTC).Add(time.Duration(rsvCtx.Config.LinkExpire)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Link", err)
// 	}

// 	// negative cancel req. scenarios
// 	if err = rsvCtx.ApiCancelRsv(r.ID, "", 401); err != nil {
// 		return wrapErr("apiUserCancelRequest", err)
// 	}
// 	if err = rsvCtx.ApiCancelRsv(r.ID, userToken, 412); err != nil {
// 		return wrapErr("apiUserCancelRequest", err)
// 	}
// 	if err = rsvCtx.ApiCancelRsv(r.ID, instrToken, 404); err != nil {
// 		return wrapErr("apiUserCancelRequest", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:        adyen_sm.Hold,
// 		Timeout:      time.Now().In(time.UTC).Add(time.Hour * 24),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Hold", err)
// 	}

// 	if err = rsvCtx.ApiCancelRsv(r.ID, instrToken, 404); err != nil {
// 		return wrapErr("apiUserCancelRequest", err)
// 	}
// 	if err = rsvCtx.ApiCancelRsv(r.ID, userToken, 204); err != nil {
// 		return wrapErr("apiUserCancelRequest", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCancelOrRefund,
// 	}); err != nil {
// 		return wrapErr("WaitCancelOrRefund", err)
// 	}

// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_CancelOrRefund(r, true)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_CancelOrRefund", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.CancelOrRefund,
// 	}); err != nil {
// 		return wrapErr("CancelOrRefund", err)
// 	}

// 	return nil
// }

// /*
// 	link
// 		-> hold
// 		-> { wait_cancel_or_refund -> retry_cancel_or_refund } * limit
// 		-> error
// */
// func testFailCancelFlow() error {

// 	x := rsvCtx.AdyenSm.Config.SmMaxRetries
// 	rsvCtx.AdyenSm.Config.SmMaxRetries = 3
// 	defer func() {
// 		rsvCtx.AdyenSm.Config.SmMaxRetries = x
// 	}()

// 	var userToken string
// 	var instrToken string

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour*24 + time.Minute,
// 		DisableLinkExpress:   true,
// 		ReturnUserToken:      &userToken,
// 		ReturnInstrToken:     &instrToken,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		Rsv:          r,
// 		State:        adyen_sm.Link,
// 		Timeout:      time.Now().In(time.UTC).Add(time.Duration(rsvCtx.Config.LinkExpire)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Link", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:        adyen_sm.Hold,
// 		Timeout:      time.Now().In(time.UTC).Add(time.Hour * 24),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Hold", err)
// 	}

// 	if err = rsvCtx.ApiCancelRsv(r.ID, userToken, 204); err != nil {
// 		return wrapErr("apiUserCancelRequest", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCancelOrRefund,
// 	}); err != nil {
// 		return wrapErr("WaitCancelOrRefund", err)
// 	}

// 	//

// 	for i := 1; i <= rsvCtx.AdyenSm.Config.SmMaxRetries; i++ {

// 		if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_CancelOrRefund(r, false)); err != nil {
// 			return wrapErr("rsv.AdyenNtf_CancelOrRefund", err)
// 		}

// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State:               adyen_sm.RetryCancelOrRefund,
// 			Timeout:             time.Now().Add(time.Duration(rsvCtx.AdyenSm.Config.SmRetryTimeout)),
// 			ExpectedCancelCount: i,
// 			ResetTimeout:        true,
// 		}); err != nil {
// 			return wrapErr("RetryCancelOrRefund", err)
// 		}

// 		if err := rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 			return wrapErr(
// 				fmt.Sprintf("loop: %d retry_cancel_or_refund->wait_cancel_or_refund", i),
// 				err)
// 		}

// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State:               adyen_sm.WaitCancelOrRefund,
// 			ExpectedCancelCount: i,
// 		}); err != nil {
// 			return wrapErr("WaitCancelOrRefund", err)
// 		}

// 	}

// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_CancelOrRefund(r, false)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_CancelOrRefund", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:               adyen_sm.Error,
// 		ExpectedCancelCount: 3,
// 	}); err != nil {
// 		return wrapErr("Error", err)
// 	}

// 	return nil
// }

// /*
// 	link_express
// 		-> wait_capture
// 		-> error
// */
// func testWaitCaptureTimeouts() error {

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if err = rsv.EnsureRsvIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return wrapErr("rsv.EnsureRsvIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:        adyen_sm.WaitCapture,
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("WaitCapture", err)
// 	}

// 	for i := 1; i <= rsvCtx.AdyenSm.Config.SmMaxRetries; i++ {
// 		if err := rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 			return wrapErr("WaitCapture timeout iter", err)
// 		}
// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State:        adyen_sm.WaitCapture,
// 			ResetTimeout: true,
// 			SmRetries:    i,
// 		}); err != nil {
// 			return wrapErr("WaitCapture", err)
// 		}
// 	}

// 	if err := rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 		return wrapErr("WaitCapture timeout iter", err)
// 	}
// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Error,
// 	}); err != nil {
// 		return wrapErr("WaitCapture", err)
// 	}

// 	return nil
// }

// /*
// 	link_express
// 		-> { wait_capture -> hold } * limit
// 		-> wait_cancel_or_refund
// 		-> cancel_or_refund
// */
// func testFailCapture() error {

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if err = rsv.EnsureRsvIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return wrapErr("rsv.EnsureRsvIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	for i := 1; i <= rsvCtx.AdyenSm.Config.SmMaxRetries; i++ {
// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State: adyen_sm.WaitCapture,
// 		}); err != nil {
// 			return wrapErr("WaitCapture", err)
// 		}

// 		// simulate positive capture adyen.Notification
// 		if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Capture(r, false)); err != nil {
// 			return wrapErr("rsv.AdyenNtf_Capture", err)
// 		}

// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State:        adyen_sm.Hold,
// 			Timeout:      time.Now().Add(time.Duration(rsvCtx.AdyenSm.Config.SmRetryTimeout)),
// 			ResetTimeout: true,
// 		}); err != nil {
// 			return wrapErr("Hold", err)
// 		}

// 		// daemon should perform payout request
// 		if err = rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 			return wrapErr("daemon Hold->WaitCapture", err)
// 		}
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return wrapErr("WaitCapture", err)
// 	}

// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Capture(r, false)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_Capture", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCancelOrRefund,
// 	}); err != nil {
// 		return wrapErr("WaitCancelOrRefund", err)
// 	}

// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_CancelOrRefund(r, true)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_CancelOrRefund", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.CancelOrRefund,
// 	}); err != nil {
// 		return wrapErr("WaitCancelOrRefund", err)
// 	}

// 	return nil
// }

// /*
// 	link_express
// 		-> wait_capture
// 		-> capture
// 		-> wait_refund
// 		-> refund
// */
// func testUserCaptureRefundBeforeT() error {

// 	userToken := ""

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 		ReturnUserToken:      &userToken,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if err = rsv.EnsureRsvIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return wrapErr("rsv.EnsureRsvIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return wrapErr("WaitCapture", err)
// 	}

// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Capture(r, true)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_Capture", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Capture,
// 	}); err != nil {
// 		return wrapErr("Capture", err)
// 	}

// 	if err = rsvCtx.ApiRefundRsv(r.ID, userToken, 0, false); err != nil {
// 		return wrapErr("apiRefundRequest", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitRefund,
// 	}); err != nil {
// 		return wrapErr("WaitRefund", err)
// 	}

// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Refund(r, true)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_Refund", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Refund,
// 	}); err != nil {
// 		return wrapErr("Refund", err)
// 	}

// 	return nil
// }

// /*
// 	link_express
// 		-> wait_capture
// 		-> capture
// 		-> dispute
// */
// func testDisputeAfterTrain() error {
// 	userToken := ""

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 		ReturnUserToken:      &userToken,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if err = rsv.EnsureRsvIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return wrapErr("rsv.EnsureRsvIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return wrapErr("WaitCapture", err)
// 	}

// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Capture(r, true)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_Capture", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:          adyen_sm.Capture,
// 		ResetDateStart: time.Now().In(time.UTC).Add(-time.Minute),
// 	}); err != nil {
// 		return wrapErr("Capture", err)
// 	}

// 	if err = rsvCtx.ApiRefundRsv(r.ID, userToken, 412, false); err != nil {
// 		return wrapErr("apiRefundRequest", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Capture,
// 	}); err != nil {
// 		return wrapErr("Capture", err)
// 	}

// 	if err = rsvCtx.ApiDisputeRsv(r.ID, userToken, 0, false); err != nil {
// 		return wrapErr("apiDisputeRsv", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Dispute,
// 	}); err != nil {
// 		return wrapErr("Dispute", err)
// 	}

// 	return nil
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

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 		ReturnUserToken:      &userToken,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if err = rsv.EnsureRsvIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return wrapErr("rsv.EnsureRsvIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return wrapErr("WaitCapture", err)
// 	}

// 	// simulate positive capture adyen.Notification
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveCapture(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveCapture", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:        adyen_sm.Capture,
// 		Timeout:      r.DateEnd.Add(time.Duration(rsvCtx.Config.PayoutDelay)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Capture", err)
// 	}

// 	// daemon should perform payout request
// 	if err = rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 		return wrapErr("daemon Capture->WaitPayout", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:   adyen_sm.WaitPayout,
// 		Timeout: time.Now().In(time.UTC).Add(time.Duration(rsvCtx.AdyenSm.Config.InstantPayoutTimeout)),
// 	}); err != nil {
// 		return wrapErr("WaitPayout", err)
// 	}

// 	// simulate positive payout adyen.Notification
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositivePayout(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositivePayout", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Payout,
// 	}); err != nil {
// 		return wrapErr("StatePayout", err)
// 	}

// 	if err = rsvCtx.ApiDisputeRsv(r.ID, userToken, 0, false); err != nil {
// 		return wrapErr("apiDisputeRequest", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Dispute,
// 	}); err != nil {
// 		return wrapErr("Dispute", err)
// 	}

// 	return nil
// }

// /*
// 	link_express
// 		-> wait_capture
// 		-> { capture -> wait_payout } * limit
// 		-> error
// */
// func testFailPayout() error {
// 	userToken := ""

// 	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
// 		RsvStartDateSinceNow: time.Hour,
// 		ReturnUserToken:      &userToken,
// 	})
// 	if err != nil {
// 		return wrapErr("prep", err)
// 	}

// 	if err = rsv.EnsureRsvIsInState(r, adyen_sm.LinkExpress); err != nil {
// 		return wrapErr("rsv.EnsureRsvIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveAuth", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return wrapErr("WaitCapture", err)
// 	}

// 	// simulate positive capture adyen.Notification
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveCapture(r)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_PositiveCapture", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:        adyen_sm.Capture,
// 		Timeout:      r.DateEnd.Add(time.Duration(rsvCtx.Config.PayoutDelay)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return wrapErr("Capture", err)
// 	}

// 	for i := 1; i <= rsvCtx.AdyenSm.Config.SmMaxRetries; i++ {

// 		if err = rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 			return wrapErr("daemon Capture->WaitPayout", err)
// 		}

// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State:   adyen_sm.WaitPayout,
// 			Timeout: time.Now().In(time.UTC).Add(time.Duration(rsvCtx.AdyenSm.Config.InstantPayoutTimeout)),
// 		}); err != nil {
// 			return wrapErr("WaitPayout", err)
// 		}

// 		// simulate positive payout adyen.Notification
// 		if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Payout(r, false)); err != nil {
// 			return wrapErr("rsv.AdyenNtf_Payout", err)
// 		}

// 		if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 			State:        adyen_sm.Capture,
// 			Timeout:      time.Now().Add(time.Duration(rsvCtx.AdyenSm.Config.SmRetryTimeout)),
// 			ResetTimeout: true,
// 		}); err != nil {
// 			return wrapErr("Capture", err)
// 		}
// 	}

// 	if err = rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
// 		return wrapErr("daemon out of loop Capture->WaitPayout", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State:   adyen_sm.WaitPayout,
// 		Timeout: time.Now().In(time.UTC).Add(time.Duration(rsvCtx.AdyenSm.Config.InstantPayoutTimeout)),
// 	}); err != nil {
// 		return wrapErr("WaitPayout", err)
// 	}

// 	// simulate positive payout adyen.Notification
// 	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Payout(r, false)); err != nil {
// 		return wrapErr("rsv.AdyenNtf_Payout", err)
// 	}

// 	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
// 		State: adyen_sm.Error,
// 	}); err != nil {
// 		return wrapErr("Error", err)
// 	}

// 	return nil
// }

// //

// func TestExpressFlow(t *testing.T) {
// 	if err := testExpressFlow(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestNormalFlow(t *testing.T) {
// 	if err := testNormalFlow(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestNormalLinkExpireFlow(t *testing.T) {
// 	if err := testNormalLinkExpireFlow(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestExpressLinkExpireFlow(t *testing.T) {
// 	if err := testExpressLinkExpireFlow(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestForceLinkExpireFlow(t *testing.T) {
// 	if err := testForceLinkExpireFlow(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestSuccessCancelFlow(t *testing.T) {
// 	if err := testSuccessCancelFlow(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestFailCancelFlow(t *testing.T) {
// 	if err := testFailCancelFlow(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestWaitCaptureTimeouts(t *testing.T) {
// 	if err := testWaitCaptureTimeouts(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestFailCapture(t *testing.T) {
// 	if err := testFailCapture(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestUserCaptureRefundBeforeTrain(t *testing.T) {
// 	if err := testUserCaptureRefundBeforeT(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestDisputeAfterTrain(t *testing.T) {
// 	if err := testDisputeAfterTrain(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestPayoutDispute(t *testing.T) {
// 	if err := testPayoutDispute(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestFailPayout(t *testing.T) {
// 	if err := testFailPayout(); err != nil {
// 		t.Fatal(err)
// 	}
// }
