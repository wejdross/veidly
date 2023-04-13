package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sport/adyen"
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/dc"
	"sport/helpers"
	"sport/instr"
	"sport/rsv"
	"sport/static"
	"sport/sub"
	"sport/train"
	"sport/user"

	"github.com/kzaag/dp/cmn"
)

type FlowHandler func(*Ctx, string) error

func main() {
	dbname := "sportdb_test"
	apiCtx := api.NewApi(config.NewLocalCtx())
	if err := dal.ValidateNoPendingChanges(apiCtx.Config, true, dbname, false, false); err != nil {
		if helpers.PgIsDbNotExists(err) {
			dal.DeployDb(apiCtx.Config, false, dbname, false)
		} else {
			log.Fatal(err)
		}
	}

	dalCtx := dal.NewDal(apiCtx.Config, dbname)

	staticCtx := static.NewCtx(apiCtx)
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	adyenCtx := adyen.NewCtx(apiCtx, nil)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, adyenCtx)
	dcCtx := dc.NewCtx(apiCtx, dalCtx, instrCtx)
	trainCtx := train.NewCtx(apiCtx, dalCtx, userCtx, staticCtx, instrCtx)
	subCtx := sub.NewCtx(apiCtx, dalCtx, userCtx, instrCtx, adyenCtx, nil, nil)
	rsvCtx := rsv.NewCtx(
		apiCtx, dalCtx, userCtx,
		instrCtx, trainCtx, adyenCtx,
		nil, nil, dcCtx, subCtx, nil)

	ctx := &Ctx{
		Rsv:   rsvCtx,
		User:  userCtx,
		Instr: instrCtx,
		Api:   apiCtx,
		Train: trainCtx,
	}

	/*
		im not running daemon here since i will be calling iterations by myself.
		more control that way
	*/
	//go instructor.RunRsvDaemon()

	go func() {
		err := apiCtx.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// validate that UI is running
	c, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldnt connect to UI, err was: %v\n", err)
		os.Exit(1)
	}
	if err = c.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf(
		"%sNote that following tests require user interaction and WONT mock adyen communication.%s\n",
		cmn.ForeYellow,
		cmn.AttrOff)

	opts := map[string]FlowHandler{
		"express":          RunExpressFlow,
		"normal":           RunNormalFlow,
		"cancel_or_refund": RunCancelOrRefundFlow,
		"refund":           RunRefundFlow,
		//
		"all": RunAllFlow,
	}

	fmt.Println("Which flow do you want to test? Your options are:")
	for k := range opts {
		fmt.Printf("\t%s\n", k)
	}

	var res string
	fmt.Print("> ")
	if _, err := fmt.Scanln(&res); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if h := opts[res]; h != nil {
		if err = h(ctx, ""); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "unkown option: %s\n", res)
		os.Exit(1)
	}

}
