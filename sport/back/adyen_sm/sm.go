package adyen_sm

import (
	"database/sql"
	"fmt"
	"sport/adyen"
	"sport/helpers"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (ctx *Ctx) UpdateInstrDecision(
	p *Pass,
	decision InstructorDecision,
) error {
	res, err := ctx.dal.Db.Exec(fmt.Sprintf(`
			update %s set
				%s = $1
			where %s = $2 and %s = $3`,
		ctx.Model.TableName,
		ctx.Model.Columns.InstructorDecision,
		ctx.Model.Columns.ID,
		ctx.Model.Columns.InstructorID,
	),
		decision, p.Data.ID, p.Data.InstructorID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) UpdateState(
	objID uuid.UUID, state, nextState State, eargs map[string]interface{}, tx *sql.Tx,
) error {
	var err error

	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(
		`update %s set %s = $1, %s = 0`,
		ctx.Model.TableName,
		ctx.Model.Columns.State,
		ctx.Model.Columns.SmRetries)) // always set retries to 0 on successful state change
	args := []interface{}{nextState}
	next := 2
	for k := range eargs {
		q.WriteString(fmt.Sprintf(" , %s = $%d", k, next))
		args = append(args, eargs[k])
		next++
	}
	q.WriteString(fmt.Sprintf(` where %s = $%d and %s = $%d`,
		ctx.Model.Columns.ID, next, ctx.Model.Columns.State, next+1))
	args = append(args, objID, state)
	next += 2
	var res sql.Result
	if tx == nil {
		res, err = ctx.dal.Db.Exec(q.String(), args...)
	} else {
		res, err = tx.Exec(q.String(), args...)
	}

	if err != nil {
		return err
	}

	err = helpers.PgMustBeOneRow(res)

	return err
}

func (ctx *Ctx) UpdateRetries(
	retries int,
	id uuid.UUID, state State,
) error {
	to := time.Now().In(time.UTC).Add(time.Duration(ctx.Config.SmRetryTimeout))
	res, err := ctx.dal.Db.Exec(fmt.Sprintf(`
			update %s set 
				%s = $1, 
				%s = $2 
			where %s = $3 and %s = $4`,
		ctx.Model.TableName,
		ctx.Model.Columns.SmRetries,
		ctx.Model.Columns.SmTimeout,
		ctx.Model.Columns.ID,
		ctx.Model.Columns.State,
	),
		retries, to, id, state)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) MoveStateToError(obj *Pass, sd StateChangeEventSource) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        obj.Data.ID,
		PreviousState: string(obj.Data.State),
		NextState:     string(Error),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	// Notify admin about error situation
	// before i also was setting:
	//		"is_active":    false,
	// but i think thats incorrect. If rsv enters error state we must manually resolve it
	// and it still should be considered active
	err := ctx.UpdateState(
		obj.Data.ID, obj.Data.State, Error,
		map[string]interface{}{
			ctx.Model.Columns.IsConfirmed: false,
		},
		nil,
	)

	if err != nil {
		sce.Success = false
		sce.Error = err.Error()
	}

	ctx.LogStateChangeEvent(&sce, nil)

	return err
}

func (ctx *Ctx) MoveStateToWaitCancelOrRefund(
	p *Pass,
	sd StateChangeEventSource,
	args map[string]interface{},
) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(WaitCancelOrRefund),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	tx, err := ctx.dal.Db.Begin()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	if args == nil {
		args = make(map[string]interface{})
	}
	args[ctx.Model.Columns.IsConfirmed] = false
	args[ctx.Model.Columns.SmTimeout] = time.Now().In(time.UTC).Add(time.Hour * 24)

	if err = ctx.UpdateState(
		p.Data.ID, p.Data.State, WaitCancelOrRefund, args, tx,
	); err != nil {
		sce.LogWithError(ctx, err, nil)
		tx.Rollback()
		return err
	}

	if ctx.adyen.Mockup {
		err = tx.Commit()
		if err != nil {
			sce.LogWithError(ctx, err, nil)
			return err
		}
		sce.Log(ctx, nil)
		return nil
	}

	req := adyen.RefundTransactionRequest{
		RequestID: uuid.NewString(),
		Refunds: []adyen.RefundTransactionItem{
			{
				OrderID:     "",
				SessionID:   p.Data.SmType + "_" + p.Data.RefID,
				Amount:      p.Data.RefundValue,
				Description: "",
			},
		},
	}

	_, err = ctx.adyen.RefundTransaction(&req)
	if err != nil {

		e := tx.Rollback()
		if e != nil {
			sce.LogWithError(
				ctx,
				fmt.Errorf("couldnt rollback transaction, err was: %v", e), nil)
		}

		sce.LogWithError(ctx, err, nil)

		// there cant be any retrying during this state transition

		// if r.SmRetries >= global.Request.SmMaxRetries {
		// 	return MoveStateToError(r, source)
		// }

		//return DalUpdateReservationRetries(r)

		return ctx.MoveStateToError(p, sd)
	}

	err = tx.Commit()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	sce.Log(ctx, nil)

	return nil
}

func (ctx *Ctx) MoveStateToWaitPayout(
	p *Pass, sd StateChangeEventSource,
) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(WaitPayout),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	tx, err := ctx.dal.Db.Begin()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	args := map[string]interface{}{
		ctx.Model.Columns.SmTimeout: time.Now().In(time.UTC).Add(time.Duration(ctx.Config.InstantPayoutTimeout)),
	}

	if ctx.BeforeWaitPayout != nil {
		if c, err := ctx.BeforeWaitPayout(p, args, tx, sd); err != nil {
			tx.Rollback()
			sce.Success = false
			sce.Error = err.Error()
			ctx.LogStateChangeEvent(&sce, nil)
			return err
		} else if !c {
			tx.Rollback()
			sce.Success = false
			sce.Error = "interrupted by callback"
			ctx.LogStateChangeEvent(&sce, nil)
			return nil
		}
	}

	// var res *adyen.PayoutResponse
	// var req adyen.PayoutRequest
	pv := p.Data.PayoutValue
	var cut int

	if p.Data.ShopperReference == "" {
		err = fmt.Errorf("cannot perform payout without instructor having specified ShopperReference")
		goto ERR
	}

	if cut, err = ctx.instr.DalGetNextPenaltySerial(
		p.Data.InstructorID, pv,
		ctx.Config.InstrShotNoMoreThanPercent); err != nil {
		goto ERR
	}

	if cut > 0 {
		pv -= cut
	}

	// req = adyen.PayoutRequest{
	// 	ShopperReference: p.Data.ShopperReference,
	// 	Amount: adyen.Amount{
	// 		Currency: p.Data.Currency,
	// 		Value:    int64(pv),
	// 	},
	// }

	// if ctx.adyen.Mockup {
	// 	// res = &adyen.PayoutResponse{
	// 	// 	PspReference: "test psp res",
	// 	// }
	// } else {

	// 	// req.Reference = p.Data.SmType + "_" + p.Data.RefID
	// 	// req.MerchantAccount = ctx.adyen.Config.MerchAcc
	// 	// req.SelectedRecurringDetailReference = "LATEST"
	// 	// req.Recurring = new(adyen.Recurring)
	// 	// req.Recurring.Contract = "PAYOUT"
	// 	// req.ShopperInteraction = "ContAuth"
	// 	// res, err = ctx.adyen.AdyenPayout(&req)
	// 	// if err != nil {
	// 	// 	goto ERR
	// 	// }
	// }

	err = ctx.UpdateState(p.Data.ID, p.Data.State, WaitPayout, args, tx)

ERR:
	if err != nil {

		if cut > 0 {
			if err := ctx.instr.DalUpdatePenalty(p.Data.InstructorID, cut, nil); err != nil {
				sce.LogWithError(ctx, fmt.Errorf("couldnt rollback penalty transaction, err was: %v", err), nil)
			}
		}

		e := tx.Rollback()
		if e != nil {
			sce.LogWithError(ctx, fmt.Errorf("couldnt rollback transaction, err was: %v", e), nil)
		}

		sce.LogWithError(ctx, err, nil)

		if p.Data.SmRetries >= ctx.Config.SmMaxRetries {
			return ctx.MoveStateToError(p, sd)
		}

		return ctx.UpdateRetries(p.Data.SmRetries+1, p.Data.ID, p.Data.State)
	}

	err = tx.Commit()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	// if res != nil {
	// 	j, err := json.Marshal(res)
	// 	if err == nil {
	// 		sce.Source.Data = string(j)
	// 	}
	// }

	sce.Log(ctx, nil)
	return nil
}

func (ctx *Ctx) MoveStateToHold(
	p *Pass, sd StateChangeEventSource, args map[string]interface{}) error {

	now := time.Now().In(time.UTC)

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(Hold),
		Timestamp:     now,
		Success:       true,
		Source:        sd,
	}

	tx, err := ctx.dal.Db.Begin()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	if ctx.BeforeHold != nil {
		if c, err := ctx.BeforeHold(p, args, tx, sd); err != nil {
			tx.Rollback()
			sce.LogWithError(ctx, err, nil)
			return err
		} else if !c {
			tx.Rollback()
			sce.Success = false
			sce.Error = "interrupted by callback"
			ctx.LogStateChangeEvent(&sce, nil)
			return nil
		}
	}

	err = ctx.UpdateState(p.Data.ID, p.Data.State, Hold, args, tx)
	if err != nil {
		tx.Rollback()
		sce.LogWithError(ctx, err, nil)
		return err
	}

	sce.Log(ctx, nil)

	// mailing: make sure its the first time we hit hold
	if p.Data.SmCache == nil || p.Data.SmCache[WaitCapture].Retries == 0 {
		if ctx.EmailUserAboutHold != nil {
			ctx.EmailUserAboutHold(p)
		}
	}

	return tx.Commit()
}

func (ctx *Ctx) MoveStateToCancelOrRefund(
	p *Pass,
	sd StateChangeEventSource,
	args map[string]interface{}) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(CancelOrRefund),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	tx, err := ctx.dal.Db.Begin()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	if args == nil {
		args = make(map[string]interface{})
	}

	args[ctx.Model.Columns.IsActive] = false
	args[ctx.Model.Columns.IsConfirmed] = false

	if ctx.BeforeCancelOrRefund != nil {
		if c, err := ctx.BeforeCancelOrRefund(p, args, tx, sd); err != nil {
			tx.Rollback()
			sce.LogWithError(ctx, err, nil)
			return err
		} else if !c {
			tx.Rollback()
			sce.Success = false
			sce.Error = "interrupted by callback"
			ctx.LogStateChangeEvent(&sce, nil)
			return nil
		}
	}

	err = ctx.UpdateState(p.Data.ID, p.Data.State, CancelOrRefund, args, tx)
	if err != nil {
		tx.Rollback()
		sce.LogWithError(ctx, err, nil)
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		sce.LogWithError(ctx, err, nil)
		return err
	}

	sce.Log(ctx, nil)

	if ctx.EmailInstrAboutCancel != nil {
		ctx.EmailInstrAboutCancel(p)
	}
	if ctx.EmailUserAboutCancel != nil {
		ctx.EmailUserAboutCancel(p)
	}

	return nil
}

func (ctx *Ctx) MoveStateToRetryCancelOrRefund(
	p *Pass, sd StateChangeEventSource, args map[string]interface{}) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(RetryCancelOrRefund),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	if p.Data.IsConfirmed {
		if args == nil {
			args = make(map[string]interface{})
		}
		args[ctx.Model.Columns.IsConfirmed] = false
	}

	err := ctx.UpdateState(p.Data.ID, p.Data.State, RetryCancelOrRefund, args, nil)
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}
	sce.Log(ctx, nil)
	return nil
}

func (ctx *Ctx) MoveStateToPayout(p *Pass, sd StateChangeEventSource) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(Payout),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	err := ctx.UpdateState(
		p.Data.ID, p.Data.State, Payout, map[string]interface{}{
			ctx.Model.Columns.IsActive:    false,
			ctx.Model.Columns.IsConfirmed: true,
		}, nil,
	)
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	sce.Log(ctx, nil)

	if ctx.AfterPayout != nil {
		ctx.AfterPayout(p)
	}

	if ctx.invoicing != nil {
		if err := ctx.invoicing.CreateInvoice(
			p.Data.TotalPrice-p.Data.PayoutValue,
			p.Data.InstructorID,
			p.Data.ID,
			p.Data.SmType,
		); err != nil {
			sce.LogWithError(ctx, err, nil)
		}
	}

	return nil
}

// WARN: there is another way to get rsv into CAPTURE state
// by using ConfirmRsvPayment from dal.go
// bear that in mind when modifying this function
// i know
// i dont like it either - its gonna be refactored when refactoring happens
func (ctx *Ctx) MoveStateToCapture(
	p *Pass,
	sd StateChangeEventSource,
	orderID int,
	args map[string]interface{}) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(Capture),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	if args == nil {
		args = map[string]interface{}{}
	}

	args[ctx.Model.Columns.OrderID] = orderID

	tx, err := ctx.dal.Db.Begin()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	// if its first time we hit capture, set is confirmed
	if p.Data.SmCache == nil || p.Data.SmCache[WaitPayout].Retries == 0 {
		args[ctx.Model.Columns.IsConfirmed] = true
	}

	// possible use case:
	// obj is no longer available, refund payment
	// TODO: we may create some dispute-like state
	// which would give client way to change reservation without having to pay again
	// and we wouldnt have to refund
	// win - win ?
	if ctx.BeforeCapture != nil {
		if c, err := ctx.BeforeCapture(p, args, tx, sd); err != nil {
			tx.Rollback()
			sce.LogWithError(ctx, err, nil)
			return err
		} else if !c {
			tx.Rollback()
			sce.Success = false
			sce.Error = "interrupted by callback"
			ctx.LogStateChangeEvent(&sce, nil)
			return nil
		}
	}

	err = ctx.UpdateState(p.Data.ID, p.Data.State, Capture, args, tx)
	if err != nil {
		tx.Rollback()
		sce.LogWithError(ctx, err, nil)
		return err
	}

	if ctx.adyen.Mockup {
		err = tx.Commit()
		if err != nil {
			sce.LogWithError(ctx, err, nil)
			return err
		}
		sce.Log(ctx, nil)
		return nil
	}

	req := adyen.VerifyTransactionRequest{
		SessionID: p.Data.SmType + "_" + p.Data.ID.String(),
		Amount:    p.Data.TotalPrice,
		Currency:  p.Data.Currency,
		OrderID:   orderID,
	}

	_, err = ctx.adyen.VerifyTransaction(&req)
	if err != nil {

		e := tx.Rollback()
		if e != nil {
			sce.LogWithError(
				ctx,
				fmt.Errorf("couldnt rollback transaction, err was: %v", e), nil)
		}

		sce.LogWithError(ctx, err, nil)

		return ctx.MoveStateToError(p, sd)
	}

	sce.Log(ctx, nil)

	if p.Data.SmCache == nil || p.Data.SmCache[WaitPayout].Retries == 0 {
		if ctx.EmailUserAboutCapture != nil {
			ctx.EmailUserAboutCapture(p)
		}
	}

	return tx.Commit()
}

func (ctx *Ctx) MoveStateToWaitRefund(
	p *Pass, sd StateChangeEventSource,
	onErrorReturn bool,
	fullRefund bool,
) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(WaitRefund),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	tx, err := ctx.dal.Db.Begin()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	args := map[string]interface{}{
		ctx.Model.Columns.SmTimeout: time.Now().In(time.UTC).Add(time.Duration(ctx.Config.RefundTimeout)),
	}

	if ctx.BeforeWaitRefund != nil {
		if c, err := ctx.BeforeWaitRefund(p, args, tx, sd); err != nil {
			tx.Rollback()
			sce.Success = false
			sce.Error = err.Error()
			ctx.LogStateChangeEvent(&sce, nil)
			return err
		} else if !c {
			tx.Rollback()
			sce.Success = false
			sce.Error = "interrupted by callback"
			ctx.LogStateChangeEvent(&sce, nil)
			return nil
		}
	}

	if err = ctx.UpdateState(p.Data.ID, p.Data.State, WaitRefund, args, tx); err != nil {
		sce.LogWithError(ctx, err, nil)
		tx.Rollback()
		return err
	}

	if ctx.adyen.Mockup {
		err = tx.Commit()
		if err != nil {
			sce.LogWithError(ctx, err, nil)
			return err
		}
		sce.Log(ctx, nil)
		return nil
	}

	v := 0
	if fullRefund {
		v = p.Data.TotalPrice
	} else {
		v = p.Data.RefundValue
	}

	req := adyen.RefundTransactionRequest{
		RequestID: uuid.NewString(),
		Refunds: []adyen.RefundTransactionItem{
			{
				OrderID:     p.Data.OrderID,
				SessionID:   p.Data.SmType + "_" + p.Data.ID.String(),
				Amount:      v,
				Description: "Refund",
			},
		},
		RefundsUUid: uuid.New(),
	}

	_, err = ctx.adyen.RefundTransaction(&req)
	if err != nil {

		e := tx.Rollback()
		if e != nil {
			sce.LogWithError(ctx, fmt.Errorf("couldnt rollback transaction, err was: %v", e), nil)
		}

		sce.LogWithError(ctx, err, nil)

		if onErrorReturn {
			return err
		}

		return ctx.MoveStateToError(p, sd)
	}

	err = tx.Commit()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	sce.Log(ctx, nil)

	return nil
}

func (ctx *Ctx) MoveStateToRefund(p *Pass, sd StateChangeEventSource) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(Refund),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	tx, err := ctx.dal.Db.Begin()
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}

	args := map[string]interface{}{
		ctx.Model.Columns.IsActive:    false,
		ctx.Model.Columns.IsConfirmed: false,
	}

	if ctx.BeforeRefund != nil {
		if c, err := ctx.BeforeRefund(p, args, tx, sd); err != nil {
			tx.Rollback()
			sce.LogWithError(ctx, err, nil)
			return err
		} else if !c {
			tx.Rollback()
			sce.Success = false
			sce.Error = "interrupted by callback"
			ctx.LogStateChangeEvent(&sce, nil)
			return nil
		}
	}

	err = ctx.UpdateState(
		p.Data.ID, p.Data.State,
		Refund, args, tx,
	)
	if err != nil {
		tx.Rollback()
		sce.LogWithError(ctx, err, nil)
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		sce.LogWithError(ctx, err, nil)
		return err
	}

	sce.Log(ctx, nil)

	if ctx.EmailUserAboutCancel != nil {
		ctx.EmailUserAboutCancel(p)
	}

	if ctx.EmailInstrAboutCancel != nil {
		ctx.EmailInstrAboutCancel(p)
	}

	return nil
}

func (ctx *Ctx) MoveStateToDispute(p *Pass, sd StateChangeEventSource) error {

	sce := StateChangeEvent{
		ID:            uuid.New(),
		ObjKey:        p.Data.ID,
		PreviousState: string(p.Data.State),
		NextState:     string(Dispute),
		Timestamp:     time.Now().In(time.UTC),
		Success:       true,
		Source:        sd,
	}

	err := ctx.UpdateState(
		p.Data.ID, p.Data.State,
		Dispute,
		map[string]interface{}{
			ctx.Model.Columns.IsConfirmed: false,
		}, nil,
	)
	if err != nil {
		sce.LogWithError(ctx, err, nil)
		return err
	}
	if ctx.EmailUserAboutDispute != nil {
		ctx.EmailUserAboutDispute(p)
	}
	if ctx.EmailInstrAboutDispute != nil {
		ctx.EmailInstrAboutDispute(p)
	}

	sce.Log(ctx, nil)
	return nil
}
