package main

import (
	"fmt"
	"sport/adyen_sm"
	"sport/helpers"
	"sport/rsv"
	"sport/search"
	"sport/train"
	"time"
)

func updateCache(ctx *search.Ctx) error {
	q, err := ctx.GetPgChanges()
	if err != nil {
		return err
	}
	if err := ctx.UpdateSearchCache(q); err != nil {
		return err
	}
	fmt.Println("--- UNLOCKED TIME ---")
	fmt.Println(helpers.JsonMustSerializeFormatStr(ctx.Cache.UnlockedUpdatePerf))
	fmt.Println("--- LOCKED TIME ---")
	fmt.Println(helpers.JsonMustSerializeFormatStr(ctx.Cache.LockedUpdatePerf))
	return nil
}

func Diff(ctx *search.Ctx) error {

	fmt.Println("--- UPDATE PERF BEFORE ---")
	updateCache(ctx)

	token, _, err := ctx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		return err
	}
	s := helpers.NowMin().Add(time.Hour * 48)
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Minute))
	tr.Training.Title = fmt.Sprintf("clusterfuck")
	tid, err := ctx.Train.ApiCreateTraining(token, &tr)
	if err != nil {
		return err
	}
	fmt.Printf("created training with id = %v\n", tid)

	rr := rsv.NewTestCreateReservationRequest(tid, s)
	r1, err := ctx.Rsv.ApiCreateAndQueryRsv(&token, &rr, 0)
	if err != nil {
		return err
	}
	if err := ctx.Rsv.AdyenSm.MoveStateToCapture(
		ctx.Rsv.RsvResponseToSmPassPtr(r1),
		adyen_sm.ManualSrc, nil); err != nil {
		return err
	}
	fmt.Printf("created rsv with id = %v\n", r1.ID)

	fmt.Println("--- UPDATE PERF AFTER ---")
	updateCache(ctx)

	return nil
}
