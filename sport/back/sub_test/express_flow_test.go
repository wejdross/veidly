package sub_test

// import (
// 	"sport/adyen_sm"
// 	"sport/helpers"
// 	h "sport/helpers"
// 	"sport/sub"
// 	"testing"
// 	"time"
// )

// func testExpressFlow() error {

// 	s, err := subCtx.PrepSubTesting(sub.EmptySubTestOpts)
// 	if err != nil {
// 		return h.WrapErr("prep", err)
// 	}

// 	if err = sub.EnsureSubIsInState(s, adyen_sm.LinkExpress); err != nil {
// 		return h.WrapErr("sub.EnsureRsvIsInState(LinkExpress)", err)
// 	}

// 	// simulate positive auth from adyen
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositiveAuth(s)); err != nil {
// 		return h.WrapErr("sub.AdyenNtf_PositiveAuth", err)
// 	}

// 	if s, err = subCtx.RefreshSubAndValidate(s.ID, sub.SubValidationOpts{
// 		State: adyen_sm.WaitCapture,
// 	}); err != nil {
// 		return h.WrapErr("WaitCapture", err)
// 	}

// 	// simulate positive capture adyen.Notification
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositiveCapture(s)); err != nil {
// 		return h.WrapErr("sub.AdyenNtf_PositiveCapture", err)
// 	}

// 	if s, err = subCtx.RefreshSubAndValidate(s.ID, sub.SubValidationOpts{
// 		State:        adyen_sm.Capture,
// 		Timeout:      helpers.NowMin().Add(time.Duration(subCtx.Config.PayoutDelay)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return h.WrapErr("Capture", err)
// 	}

// 	// daemon should perform payout request
// 	if err = subCtx.AdyenSm.AgentDoOne(subCtx.SubToSmPassPtr(s)); err != nil {
// 		return h.WrapErr("daemon Capture->WaitPayout", err)
// 	}

// 	if s, err = subCtx.RefreshSubAndValidate(s.ID, sub.SubValidationOpts{
// 		State:   adyen_sm.WaitPayout,
// 		Timeout: time.Now().In(time.UTC).Add(time.Duration(subCtx.AdyenSm.Config.InstantPayoutTimeout)),
// 	}); err != nil {
// 		return h.WrapErr("WaitPayout", err)
// 	}

// 	// simulate positive payout adyen.Notification
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_PositivePayout(s)); err != nil {
// 		return h.WrapErr("sub.AdyenNtf_PositivePayout", err)
// 	}

// 	if s, err = subCtx.RefreshSubAndValidate(s.ID, sub.SubValidationOpts{
// 		State: adyen_sm.Payout,
// 	}); err != nil {
// 		return h.WrapErr("StatePayout", err)
// 	}

// 	return nil
// }

// func TestExpressFlow(t *testing.T) {
// 	if err := testExpressFlow(); err != nil {
// 		t.Fatal(err)
// 	}
// }
