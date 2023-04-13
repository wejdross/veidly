package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sport/adyen"
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/dc"
	"sport/instr"
	"sport/rsv"
	"sport/static"
	"sport/sub"
	"sport/train"
	"sport/user"
	"time"
)

func main() {
	fmt.Println("will reset sportdb_sg, ok?")
	var x string
	fmt.Scanln(&x)
	rand.Seed(time.Now().UTC().UnixNano())

	testdb := "sportdb_sg"
	apiCtx := api.NewApi(config.NewLocalCtx())

	// completely reset target database
	dal.DeployDb(apiCtx.Config, false, testdb, true)
	dalCtx := dal.NewDal(apiCtx.Config, testdb)
	staticCtx := static.NewCtx(apiCtx)
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	adyenCtx := adyen.NewMockupCtx(apiCtx)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, adyenCtx)
	dcCtx := dc.NewCtx(apiCtx, dalCtx, instrCtx)
	trainCtx := train.NewCtx(apiCtx, dalCtx, userCtx, staticCtx, instrCtx)
	subCtx := sub.NewCtx(apiCtx, dalCtx, userCtx, instrCtx,
		adyenCtx, nil, nil)
	rsvCtx := rsv.NewCtx(
		apiCtx, dalCtx, userCtx,
		instrCtx, trainCtx, adyenCtx, nil,
		nil, dcCtx, subCtx, nil)

	maxth := runtime.NumCPU()
	runtime.GOMAXPROCS(maxth)

	if maxth > 20 {
		maxth = 20
	}

	instrs, err := InstrGen(instrCtx, maxth, 100000)
	if err != nil {
		log.Fatal(err)
	}

	// why tf parallel training insert is slower than users
	ts, err := TrainGen(trainCtx, 4, 1000000, instrs)

	if err != nil {
		log.Fatal(err)
	}

	if err := RsvGen(rsvCtx, 4, 100000, ts); err != nil {
		log.Fatal(err)
	}

}
