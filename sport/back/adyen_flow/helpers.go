package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"os/exec"
	"sport/adyen_sm"
	"sport/user"
	"time"

	"github.com/google/uuid"
)

const browser = "firefox"

func WaitFor(fn func() (bool, error), delMs int) error {
	for {
		c, err := fn()
		if err != nil {
			return err
		}
		if c {
			return nil
		}
		if delMs > 0 {
			time.Sleep(time.Millisecond * time.Duration(delMs))
		} else {
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (ctx *Ctx) WaitForCardInfo(uid uuid.UUID, req *user.UserRequest) error {

	iid, err := ctx.Instr.DalReadInstructorID(uid)
	if err != nil {
		return err
	}
	_, err = ctx.Instr.DalReadCardInfo(iid)
	if err == nil {
		return nil
	}

	if err != sql.ErrNoRows {
		return err
	}

	url := "http://127.0.0.1:3000/login?return_url=/payments&email=" + url.PathEscape(req.Email)
	fmt.Printf("login now to your account @ %s and insert card info\n", url)
	fmt.Printf("username: %s\npassword: %s\n", req.Email, req.Password)

	err = exec.Command(browser, url).Start()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("waiting for card data")

	return WaitFor(func() (bool, error) {
		_, err := ctx.Instr.DalReadCardInfo(iid)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, nil
			} else {
				return false, err
			}
		} else {
			return true, nil
		}
	}, 0)
}

func (ctx *Ctx) WaitForSimpleTransition(
	rsvID uuid.UUID, from, to adyen_sm.State,
) error {
	return WaitFor(func() (bool, error) {
		r, err := ctx.Rsv.ReadRsvByID(rsvID)
		if err != nil {
			return false, err
		}
		switch r.State {
		case from:
			return false, nil
		case to:
			return true, nil
		default:
			return false, fmt.Errorf("unexpected rsv state encountered: %v", r.State)
		}
	}, 0)
}

func (ctx *Ctx) WaitForCapture(rsvID uuid.UUID) error {
	return ctx.WaitForSimpleTransition(
		rsvID,
		adyen_sm.WaitCapture,
		adyen_sm.Capture)
}

func (ctx *Ctx) WaitForCancelOrRefund(rsvID uuid.UUID) error {
	return ctx.WaitForSimpleTransition(
		rsvID,
		adyen_sm.WaitCancelOrRefund,
		adyen_sm.CancelOrRefund)
}

func (ctx *Ctx) WaitForRefund(rsvID uuid.UUID) error {
	return ctx.WaitForSimpleTransition(
		rsvID,
		adyen_sm.WaitRefund,
		adyen_sm.Refund)
}

func (ctx *Ctx) WaitForAuth(rsvID uuid.UUID) error {

	r, err := ctx.Rsv.ReadRsvByID(rsvID)
	if err != nil {
		return err
	}

	if r.State != adyen_sm.Link && r.State != adyen_sm.LinkExpress {
		if r.State == adyen_sm.Hold || r.State == adyen_sm.WaitCapture {
			return nil
		} else {
			return fmt.Errorf("invalid rsv state: %v", r.State)
		}
	}

	fmt.Printf("opening %s in the browser. go there and pay\n", r.LinkUrl)

	if err = exec.Command(browser, r.LinkUrl).Start(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Waiting for payment...")

	return WaitFor(func() (bool, error) {
		r, err = ctx.Rsv.ReadRsvByID(r.ID)
		if err != nil {
			return false, err
		}
		switch r.State {
		case adyen_sm.Link:
			fallthrough
		case adyen_sm.LinkExpress:
			return false, nil
		case adyen_sm.Hold:
			fallthrough
		case adyen_sm.WaitCapture:
			return true, nil
		default:
			return false, fmt.Errorf("unexpected rsv state encountered: %v", r.State)
		}
	}, 0)
}

func (ctx *Ctx) WaitForPayout(rsvID uuid.UUID) error {
	return ctx.WaitForSimpleTransition(
		rsvID,
		adyen_sm.WaitPayout,
		adyen_sm.Payout)
}
