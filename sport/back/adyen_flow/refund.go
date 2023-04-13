package main

import (
	"fmt"
	"sport/adyen_sm"
	"sport/rsv"
)

/*
	link_expresss
		-> wait_capture
		-> capture
		-> wait_refund
		-> refund
*/

func RunRefundFlow(ctx *Ctx, instrToken string) error {

	var returnInstr *string = nil
	if instrToken == "" {
		returnInstr = &instrToken
	}

	r, err := ctx.prepRsvAndWaitForAuth(
		PrepOpts{
			ExpressRsv:       true,
			ReturnInstrToken: returnInstr,
			InstrToken:       instrToken,
		})
	if err != nil {
		return err
	}

	if err := ctx.WaitForAuth(r.ID); err != nil {
		return err
	}

	fmt.Println("payment has been authorized. waiting for capture notify")

	if err := ctx.WaitForCapture(r.ID); err != nil {
		return err
	}

	fmt.Println("payment has been captured. Refunding rsv")

	if err = ctx.Rsv.ApiRefundRsv(r.ID, instrToken, 0, true); err != nil {
		return err
	}

	if r, err = ctx.Rsv.RefreshRsvAndValidate(
		r.ID,
		rsv.RsvValidationOpts{
			State: adyen_sm.WaitRefund,
		},
	); err != nil {
		return err
	}

	fmt.Println("waiting for refund notify")

	if err = ctx.WaitForRefund(r.ID); err != nil {
		return err
	}

	fmt.Println("rsv has been refunded")

	return nil
}
