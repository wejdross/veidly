package main

import (
	"flag"
	"log"
	"os"
	"sport/adyen"
	"sport/api"
	"sport/charts"
	"sport/chat"
	"sport/chat_integrator"
	"sport/config"
	"sport/dal"
	"sport/dbg"
	"sport/dc"
	"sport/helpers"
	"sport/instr"
	"sport/invoicing"
	"sport/notify"
	"sport/review"
	"sport/rsv"
	"sport/rsv_qr"
	"sport/schedule"
	"sport/search"
	"sport/static"
	"sport/sub"
	"sport/sub_qr"
	"sport/train"
	"sport/user"
)

func main() {

	dbname := "sportdb"

	dbg.LogInfo("Starting api...")

	var ver string
	var keyfile string
	var confpath string
	var rmkeyfile bool
	var resetdb bool
	var earlyexit bool
	flag.StringVar(&ver, "ver", "local", "api version: local, dev, prod")
	flag.StringVar(&keyfile, "keyfile", "", "decryption key file path")
	flag.BoolVar(&rmkeyfile, "rmkeyfile", false, "remove key file after bootstrap")
	flag.StringVar(&confpath, "config", "", "config file path")
	flag.BoolVar(&resetdb, "resetdb", false, "completely redeploy db")
	flag.BoolVar(&earlyexit, "earlyexit", false, "exit right after prep is done")
	flag.Parse()

	if confpath == "" {
		panic("config cannot be empty")
	}

	// if err := api.ValidateApiConfigVersion(confPath, isProd); err != nil {
	// 	log.Fatal(err)
	// }

	configCtx := config.NewCtx(confpath, ver, keyfile)

	if rmkeyfile {
		if err := os.Remove(keyfile); err != nil {
			panic("failed to remove keyfile: " + err.Error())
		}
	}

	apiCtx := api.NewApi(configCtx)
	dalCtx := dal.NewDal(configCtx, dbname)

	if resetdb {
		dal.DeployDb(configCtx, false, dbname, true)

		// verify that ddl is up to date
	} else if err := dal.ValidateNoPendingChanges(
		configCtx, true, dbname, false, ver != "local"); err != nil {

		if helpers.PgIsDbNotExists(err) {
			dal.DeployDb(configCtx, false, dbname, false)
		} else {
			log.Fatal(err)
		}
	}

	if earlyexit {
		return
	}

	/*
		setting up email client used for notifications
	*/
	noReplyEmailConfig := notify.NewNoReplySmtpConfig(configCtx)
	noReplyCtx := notify.NewReusableSmtpCtx(noReplyEmailConfig, 8)

	/*
		static files serving
	*/
	staticCtx := static.NewCtx(apiCtx)

	/*
		creating, managing, remvoing usert accounts
	*/
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, noReplyCtx)
	// that will delete expired user requests
	go userCtx.RunRetentionLoop()

	/*
		adyen api client, common adyen webhook handler
	*/
	adyenCtx := adyen.NewMockupCtx(apiCtx)
	adyenCtx.Dal = dalCtx
	/*
		creating, removing, updating instructor accounts
	*/
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, adyenCtx)

	trainCtx := train.NewCtx(
		apiCtx, dalCtx, userCtx, staticCtx, instrCtx)
	reviewCtx := review.NewCtx(apiCtx, dalCtx, trainCtx, userCtx)
	// that will delete expired reviews
	go reviewCtx.RunAgent()

	invoicingCtx := invoicing.NewCtx(apiCtx, dalCtx, instrCtx)

	subCtx := sub.NewCtx(apiCtx, dalCtx, userCtx, instrCtx,
		adyenCtx, noReplyCtx, invoicingCtx)
	// this is subscription payments sm agent
	go subCtx.RunAgent()
	_ = sub_qr.NewCtx(apiCtx, dalCtx, subCtx)

	/*
		discount codes
	*/
	dcCtx := dc.NewCtx(apiCtx, dalCtx, instrCtx)
	rsvCtx := rsv.NewCtx(
		apiCtx, dalCtx,
		userCtx, instrCtx,
		trainCtx, adyenCtx,
		reviewCtx, noReplyCtx,
		dcCtx, subCtx,
		invoicingCtx)
	_ = rsv_qr.NewCtx(apiCtx, dalCtx, rsvCtx)
	// this is rsv payments sm agent
	go rsvCtx.RunAgent()

	scheduleCtx := schedule.NewCtx(
		apiCtx, instrCtx, trainCtx, subCtx, rsvCtx)

	searchCtx := search.NewCtx(
		apiCtx, dalCtx, userCtx, instrCtx, trainCtx, rsvCtx, scheduleCtx)

	c := make(chan struct{})
	go searchCtx.RunCacheGenerationAgent(c)

	_ = charts.NewCtx(apiCtx, dalCtx, instrCtx)

	chatCtx := chat.NewCtx(apiCtx, userCtx.LangCtx, noReplyCtx, "", false)

	chatIntegratorCtx := chat_integrator.NewCtx(apiCtx, chatCtx, userCtx)
	_ = chatIntegratorCtx

	// waiting for cache to be generated
	<-c

	err := apiCtx.Run()
	if err != nil {
		log.Fatal(err)
	}
}
