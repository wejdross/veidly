package sub_test

// import (
// 	"sport/adyen_sm"
// 	"sport/helpers"
// 	"sport/sub"
// 	"testing"
// 	"time"
// )

// /*
// 	link -> link_expire
// */
// func testForceLinkExpireFlow() error {
// 	var utoken string
// 	s, err := subCtx.PrepSubTesting(sub.SubTestOpts{
// 		ReturnUserToken: &utoken,
// 	})
// 	if err != nil {
// 		return helpers.WrapErr("prep", err)
// 	}

// 	if s, err = subCtx.RefreshSubAndValidate(s.ID, sub.SubValidationOpts{
// 		Sub:     s,
// 		State:   adyen_sm.LinkExpress,
// 		Timeout: time.Now().In(time.UTC).Add(time.Duration(subCtx.Config.LinkExpire)),
// 	}); err != nil {
// 		return helpers.WrapErr("LinkExpress", err)
// 	}

// 	if err = subCtx.ApiExpireSub(s.ID, utoken, 0); err != nil {
// 		return helpers.WrapErr("LinkExpress->LinkExpire", err)
// 	}

// 	if s, err = subCtx.RefreshSubAndValidate(s.ID, sub.SubValidationOpts{
// 		State: adyen_sm.LinkExpire,
// 	}); err != nil {
// 		return helpers.WrapErr("LinkExpire", err)
// 	}

// 	return nil
// }

// /*
// 	link_express -> link_expire

// 	will also test that loop is created if receiving false auth notifications
// */
// func testExpressLinkExpireFlow() error {
// 	s, err := subCtx.PrepSubTesting(sub.EmptySubTestOpts)
// 	if err != nil {
// 		return helpers.WrapErr("prep", err)
// 	}

// 	if s, err = subCtx.RefreshSubAndValidate(s.ID, sub.SubValidationOpts{
// 		Sub:     s,
// 		State:   adyen_sm.LinkExpress,
// 		Timeout: time.Now().Add(time.Duration(subCtx.Config.LinkExpire)),
// 	}); err != nil {
// 		return helpers.WrapErr("LinkExpress", err)
// 	}

// 	// simulate false auth from adyen
// 	if err := subCtx.ApiSimulateAdyenWHCall(sub.AdyenNtf_Auth(s, sub.AdyenNtf_AuthOpts{
// 		Success: false,
// 	})); err != nil {
// 		return helpers.WrapErr("sub.AdyenNtf_Auth", err)
// 	}

// 	if s, err = subCtx.RefreshSubAndValidate(s.ID, sub.SubValidationOpts{
// 		State:        adyen_sm.LinkExpress,
// 		Timeout:      time.Now().Add(time.Duration(subCtx.Config.LinkExpire)),
// 		ResetTimeout: true,
// 	}); err != nil {
// 		return helpers.WrapErr("LinkExpress", err)
// 	}

// 	if err = subCtx.AdyenSm.AgentDoOne(subCtx.SubToSmPassPtr(s)); err != nil {
// 		return helpers.WrapErr("LinkExpress->LinkExpire", err)
// 	}

// 	if s, err = subCtx.RefreshSubAndValidate(s.ID, sub.SubValidationOpts{
// 		State: adyen_sm.LinkExpire,
// 	}); err != nil {
// 		return helpers.WrapErr("LinkExpire", err)
// 	}

// 	return nil
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
