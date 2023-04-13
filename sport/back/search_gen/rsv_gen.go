package main

import (
	"math/rand"
	"sport/adyen_sm"
	"sport/helpers"
	"sport/rsv"
)

func RsvGen(rsvCtx *rsv.Ctx, maxTh int, rsvNo int, ts []TrainingIDWithOccs) error {

	err := helpers.Sem(maxTh, rsvNo, []func(int) error{
		func(i int) error {
			t := ts[rand.Int()%len(ts)]
			rr := rsv.NewTestCreateReservationRequest(t.TID, t.Occs[0].DateStart)
			rr.UserData.Language = "pl"
			rr.UserData.Name = n[rand.Int()%len(n)] + ln[rand.Int()%len(ln)]
			r, err := rsvCtx.ApiCreateAndQueryRsv(nil, &rr, 0)
			if err != nil {
				return err
			}
			return rsvCtx.AdyenSm.MoveStateToCapture(
				rsvCtx.RsvResponseToSmPassPtr(r),
				adyen_sm.ManualSrc, nil)
		},
	})

	return err
}
