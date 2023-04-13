package main

import (
	"fmt"
	"sport/adyen_sm"
	"sport/rsv"
)

/*
	link
		-> hold
		-> wait_capture
		-> capture
		-> wait_payout
		-> payout
*/

func RunNormalFlow(ctx *Ctx, instrToken string) error {

	r, err := ctx.prepRsvAndWaitForAuth(PrepOpts{
		InstrToken: instrToken,
	})
	if err != nil {
		return err
	}

	if err := ctx.WaitForAuth(r.ID); err != nil {
		return err
	}

	fmt.Println("payment has been authorized.")

	if r, err = ctx.Rsv.RefreshRsvAndValidate(
		r.ID,
		rsv.RsvValidationOpts{
			State:        adyen_sm.Hold,
			ResetTimeout: true,
		},
	); err != nil {
		return err
	}

	fmt.Println("moving to capture")

	if err := ctx.Rsv.AdyenSm.AgentDoOne(ctx.Rsv.RsvResponseToSmPassPtr(r)); err != nil {
		return err
	}

	if r, err = ctx.Rsv.RefreshRsvAndValidate(
		r.ID,
		rsv.RsvValidationOpts{
			State: adyen_sm.WaitCapture,
		},
	); err != nil {
		return err
	}

	fmt.Println("waiting for capture notification")

	if err := ctx.WaitForCapture(r.ID); err != nil {
		return err
	}

	fmt.Println("payment has been captured.")

	if r, err = ctx.Rsv.RefreshRsvAndValidate(
		r.ID,
		rsv.RsvValidationOpts{
			Rsv:          r,
			ResetTimeout: true,
		}); err != nil {
		return err
	}

	if err := ctx.Rsv.AdyenSm.AgentDoOne(ctx.Rsv.RsvResponseToSmPassPtr(r)); err != nil {
		return err
	}

	if r, err = ctx.Rsv.RefreshRsvAndValidate(
		r.ID,
		rsv.RsvValidationOpts{
			State: adyen_sm.WaitPayout,
		},
	); err != nil {
		return err
	}

	fmt.Println("waiting for payout notify")

	if err = ctx.WaitForPayout(r.ID); err != nil {
		return err
	}

	fmt.Println("Payout has been sent. Normal flow is done")

	return nil
}
