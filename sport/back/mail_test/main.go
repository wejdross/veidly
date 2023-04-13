package main

import (
	"fmt"
	"log"
	"os"
	"sport/adyen"
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/dc"
	"sport/instr"
	"sport/notify"
	"sport/rsv"
	"sport/static"
	"sport/sub"
	"sport/train"
	"sport/user"
	"strconv"
)

func main() {
	apiCtx := api.NewApi(config.NewLocalCtx())
	dal.DeployDb(apiCtx.Config, false, "sportdb_test", false)
	dalCtx := dal.NewDal(apiCtx.Config, "sportdb_test")
	staticCtx := static.NewCtx(apiCtx)
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, nil)
	dcCtx := dc.NewCtx(apiCtx, dalCtx, instrCtx)
	trainCtx := train.NewCtx(apiCtx, dalCtx, userCtx, staticCtx, instrCtx)
	noReplyEmailConfig := notify.NewNoReplySmtpConfig(apiCtx.Config)
	noReplyCtx := notify.NewReusableSmtpCtx(noReplyEmailConfig, 8)
	adyenCtx := adyen.NewMockupCtx(apiCtx)
	subCtx := sub.NewCtx(apiCtx, dalCtx, userCtx, instrCtx, adyenCtx, nil, nil)
	rsvCtx := rsv.NewCtx(apiCtx, dalCtx, userCtx, instrCtx, trainCtx,
		adyenCtx, nil, noReplyCtx, dcCtx, subCtx, nil)

	type Flow struct {
		Exec func(*rsv.Ctx) error
		Desc string
	}

	flows := []Flow{
		{
			Exec: Flow1,
			Desc: "LINK-> HOLD-> WAIT_CAPTURE-> CAPTURE-> WAIT_REFUND-> REFUND-> DISPUTE",
		},
		{
			Exec: Flow2,
			Desc: "LINK -> HOLD -> WAIT_CAPTURE -> HOLD -> WAIT_CAPTURE ...",
		},
		{
			Exec: Flow3,
			Desc: "LINK -> HOLD -> WAIT_CAPTURE -> CAPTURE -> WAIT_PAYOUT -> CAPTURE -> WAIT_PAYOUT ...",
		},
	}

	for {
		fmt.Println("select flow to run:")

		for i := range flows {
			fmt.Printf("%d: %s\n", i, flows[i].Desc)
		}

		var v string
		fmt.Print("> ")
		fmt.Scanln(&v)
		vi, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		if vi < 0 || vi > (len(flows)-1) {
			log.Fatal("flow doesnt exist")
		}

		if err := flows[vi].Exec(rsvCtx); err != nil {
			fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
		}
	}
}
