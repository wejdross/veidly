package main

import (
	"fmt"
	"log"
	"sport/adyen"
	"sport/adyen_sm"
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/dc"
	"sport/instr"
	"sport/invoicing"
	"sport/review"
	"sport/rsv"
	"sport/static"
	"sport/sub"
	"sport/train"
	"sport/user"

	"github.com/google/uuid"
)

func main() {
	dbname := "sportdb"
	apiCtx := api.NewApi(config.NewLocalCtx())
	dalCtx := dal.NewDal(apiCtx.Config, dbname)
	staticCtx := static.NewCtx(apiCtx)
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, nil)
	dcCtx := dc.NewCtx(apiCtx, dalCtx, instrCtx)
	trainCtx := train.NewCtx(apiCtx, dalCtx, userCtx, staticCtx, instrCtx)
	reviewCtx := review.NewCtx(apiCtx, dalCtx, trainCtx, userCtx)
	adyenCtx := adyen.NewMockupCtx(apiCtx)
	invoiceCtx := invoicing.NewCtx(apiCtx, dalCtx, instrCtx)
	subCtx := sub.NewCtx(apiCtx, dalCtx, userCtx, instrCtx, adyenCtx, nil, nil)
	rsvCtx := rsv.NewCtx(
		apiCtx, dalCtx,
		userCtx, instrCtx,
		trainCtx, adyenCtx,
		reviewCtx, nil, dcCtx, subCtx, invoiceCtx)

	opts := []adyen_sm.State{
		adyen_sm.Payout,
	}

	isToken := false

	var t string
	fmt.Println("token/id")
	fmt.Print("> ")
	fmt.Scanln(&t)
	switch t {
	case "token":
		isToken = true
	case "id":
		break
	default:
		log.Fatal(fmt.Errorf("invalid type"))
	}

	if isToken {
		fmt.Println("enter rsv token")
	} else {
		fmt.Println("enter rsv id")
	}

	var rsvIDs string
	fmt.Print("> ")
	fmt.Scanln(&rsvIDs)
	key, err := uuid.Parse(rsvIDs)
	if err != nil {
		log.Fatal(err)
	}

	var r *rsv.DDLRsvWithInstr

	if isToken {
		r, err = rsvCtx.ReadSingleRsv(rsv.ReadRsvsArgs{AccessToken: &key, WithInstructor: true})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		r, err = rsvCtx.ReadSingleRsv(rsv.ReadRsvsArgs{ID: &key, WithInstructor: true})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("select option")
	for k := range opts {
		fmt.Printf("\t%s\n", opts[k])
	}
	var optStr string
	fmt.Print("> ")
	fmt.Scanln(&optStr)
	var opt = adyen_sm.State(optStr)

	switch opt {
	case adyen_sm.Payout:
		p := rsvCtx.RsvResponseToSmPass(r)
		err = rsvCtx.AdyenSm.MoveStateToPayout(&p, adyen_sm.ManualSrc)
	default:
		fmt.Println("unsupported state")
	}

	if err != nil {
		log.Fatal(err)
	}
}
