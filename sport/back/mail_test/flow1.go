package main

import (
	"fmt"
	"sport/adyen_sm"
	"sport/rsv"
	"sport/user"
	"time"
)

// LINK -> HOLD -> WAIT_CAPTURE -> CAPTURE
//	-> WAIT_REFUND -> REFUND -> DISPUTE
func Flow1(rsvCtx *rsv.Ctx) error {
	email := "support@veidly.com"
	instrUR := &user.UserRequest{
		Email:    email,
		Password: "123",
		UserData: user.UserData{
			Language: "en",
			Name:     "English instructor",
		},
	}
	traineeUserData := user.UserData{
		Name:     "Polish user",
		Language: "pl",
	}
	traineeContactData := user.ContactData{
		Email: email,
	}
	var v string
	var it string
	r, err := rsvCtx.PrepRsvTesting(rsv.RsvTestOpts{
		AltInstructorUserRq:  instrUR,
		UserIsInstr:          true,
		UsersMayExist:        true,
		RsvStartDateSinceNow: time.Hour * 24 * 3,
		AltRsvRequest: &rsv.ApiReservationRequest{
			UserData:    traineeUserData,
			ContactData: traineeContactData,
		},
		AnonRsvUser:      true,
		ReturnInstrToken: &it,
	})
	if err != nil {
		return err
	}

	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveAuth(r)); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.Hold,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	fmt.Println("Hold (positive payment confirmation from adyen)")
	fmt.Scanln(&v)

	if err := rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.WaitCapture,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	fmt.Println("WaitCapture (performed capture request)")
	fmt.Scanln(&v)

	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveCapture(r)); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.Capture,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	fmt.Println("Capture (got capture confirmation)")
	fmt.Scanln(&v)

	if err = rsvCtx.ApiRefundRsv(r.ID, it, 0, true); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.WaitRefund,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	fmt.Println("WaitRefund (performed refund request)")
	fmt.Scanln(&v)

	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Refund(r, true)); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.Refund,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	fmt.Println("Refund (got refund confirmation)")
	fmt.Scanln(&v)

	if err := rsvCtx.ApiDisputeRsv(r.ID, it, 0, true); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State: adyen_sm.Dispute,
	}); err != nil {
		return err
	}

	fmt.Println("Dispute (simulated dispute call)")
	return nil
}
