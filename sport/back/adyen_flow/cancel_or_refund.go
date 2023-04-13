package main

import (
	"fmt"
	"sport/adyen_sm"
	"sport/rsv"
)

/*
	link
		-> hold
		-> wait_cancel_or_refund
		-> cancel_or_refund
*/

func RunCancelOrRefundFlow(ctx *Ctx, instrToken string) error {

	var userToken string

	r, err := ctx.prepRsvAndWaitForAuth(
		PrepOpts{
			ReturnClientToken: &userToken,
			InstrToken:        instrToken,
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
			State: adyen_sm.Hold,
		},
	); err != nil {
		return err
	}

	fmt.Println("cancelling rsv")

	if err = ctx.Rsv.ApiCancelRsv(r.ID, userToken, 0); err != nil {
		return err
	}

	fmt.Println("waiting for cancel notify")

	if err = ctx.WaitForCancelOrRefund(r.ID); err != nil {
		return err
	}

	if r, err = ctx.Rsv.RefreshRsvAndValidate(
		r.ID,
		rsv.RsvValidationOpts{
			State: adyen_sm.CancelOrRefund,
		},
	); err != nil {
		return err
	}

	fmt.Println("rsv has been cancelled")

	return nil
}
