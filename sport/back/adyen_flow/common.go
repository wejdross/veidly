package main

import (
	"fmt"
	"sport/helpers"
	"sport/rsv"
	"sport/train"
	"sport/user"
	"time"

	"github.com/google/uuid"
)

type PrepOpts struct {
	ExpressRsv        bool
	ReturnClientToken *string
	ReturnInstrToken  *string
	// if provided wont create new instructor account
	InstrToken string
}

var NoPrepOpts = PrepOpts{}

func (ctx *Ctx) createClient(opts PrepOpts) (string, error) {
	_, clientToken, err := ctx.User.CreateAndLoginUserWithID(
		uuid.MustParse("eb0c325c-aac4-4891-86c3-4a101c46425c"))
	if err != nil {
		return "", err
	}

	if opts.ReturnClientToken != nil {
		*opts.ReturnClientToken = clientToken
	}

	return clientToken, nil
}

func (ctx *Ctx) createInstructorWithPayoutInfo(opts PrepOpts) (string, error) {

	if opts.InstrToken != "" {
		return opts.InstrToken, nil
	}

	// instructor
	req := user.UserRequest{
		Email:    helpers.CRNG_stringPanic(4) + "@foobar.com",
		Password: "foobar123",
		UserData: user.UserData{
			Name:     "Jeanette Voerman",
			Language: "pl",
			Country:  "PL",
		},
	}

	instrToken, err := ctx.User.ApiCreateAndLoginUser(&req)
	if err != nil {
		return "", err
	}

	if opts.ReturnInstrToken != nil {
		*opts.ReturnInstrToken = instrToken
	}

	err = ctx.Instr.ApiCreateInstructor(instrToken, nil)
	if err != nil {
		return "", err
	}

	uid, err := ctx.Api.AuthorizeUserFromToken(instrToken)
	if err != nil {
		return "", err
	}

	if err = ctx.WaitForCardInfo(uid, &req); err != nil {
		return "", err
	}

	fmt.Println("got card info")

	return instrToken, nil
}

/*
when this function is done you have:
1. client account
2. instructor account with provided payout data
3. training
4. rsv either in state:
	- link - if no opts were provided
	- link_express if opts.ExpressRsv was set to true
*/
func (ctx *Ctx) prepRsvAndWaitForAuth(opts PrepOpts) (*rsv.DDLRsvWithInstr, error) {

	var err error

	// instructor

	instrToken, err := ctx.createInstructorWithPayoutInfo(opts)
	if err != nil {
		return nil, err
	}

	// client

	clientToken, err := ctx.createClient(opts)
	if err != nil {
		return nil, err
	}

	// create training & rsv

	var delay time.Duration = 0
	if opts.ExpressRsv {
		delay = time.Hour * 2
	} else {
		delay = time.Hour * 48
	}

	s := rsv.NowMin().Add(delay)

	trainingID, err := ctx.Train.ApiCreateTraining(
		instrToken,
		&train.CreateTrainingRequest{
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						DateStart: s,
						DateEnd:   s.Add(time.Hour),
					},
				},
			},
			Training: train.TrainingRequest{
				Title:        "123",
				Capacity:     1,
				Price:        10000,
				Currency:     "PLN",
				AllowExpress: opts.ExpressRsv,
			},
		})
	if err != nil {
		return nil, err
	}

	rsv, err := ctx.Rsv.ApiCreateAndQueryRsv(
		&clientToken,
		&rsv.ApiReservationRequest{
			TrainingID: trainingID,
			Occurrence: s,
			UserData: user.UserData{
				Name:     "foo",
				Country:  "PL",
				Language: "pl",
			},
		}, 0)
	if err != nil {
		return nil, err
	}

	// done

	return rsv, nil
}
