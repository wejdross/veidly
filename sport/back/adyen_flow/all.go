package main

import (
	"fmt"

	"github.com/kzaag/dp/cmn"
)

type FlowRes struct {
	Name  string
	Error error
}

func RunAllFlow(ctx *Ctx, instrToken string) error {

	var err error

	if instrToken, err = ctx.createInstructorWithPayoutInfo(NoPrepOpts); err != nil {
		return err
	}

	type flowItem struct {
		Name    string
		Handler FlowHandler
	}

	todo := []flowItem{
		{
			Name:    "express",
			Handler: RunExpressFlow,
		}, {
			Name:    "normal",
			Handler: RunNormalFlow,
		}, {
			Name:    "cancel_or_refund",
			Handler: RunCancelOrRefundFlow,
		}, {
			Name:    "refund",
			Handler: RunRefundFlow,
		},
	}

	resChan := make(chan FlowRes, len(todo))

	for i := range todo {
		go func(i int) {
			resChan <- FlowRes{
				Name:  todo[i].Name,
				Error: todo[i].Handler(ctx, instrToken),
			}
		}(i)
	}

	for range todo {
		r := <-resChan
		if r.Error == nil {
			fmt.Printf("%s%s completed%s\n",
				cmn.ForeGreen, r.Name, cmn.AttrOff)
		} else {
			fmt.Printf("%s%s completed with error: %v%s\n",
				cmn.ForeRed, r.Name, r.Error, cmn.AttrOff)
		}
	}

	return nil
}
