package main

import (
	"sport/adyen_sm"
	"sport/rsv"
	"sport/user"
	"time"
)

// failed payout
func Flow3(rsvCtx *rsv.Ctx) error {
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

	if err := rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.WaitCapture,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_PositiveCapture(r)); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.Capture,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	if err := rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.WaitPayout,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Payout(r, false)); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.Capture,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	if err := rsvCtx.AdyenSm.AgentDoOne(rsvCtx.RsvResponseToSmPassPtr(r)); err != nil {
		return err
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State:        adyen_sm.WaitPayout,
		ResetTimeout: true,
	}); err != nil {
		return err
	}

	if err := rsvCtx.ApiSimulateAdyenWHCall(rsv.AdyenNtf_Payout(r, false)); err != nil {
		return err
	}

	return nil
}
