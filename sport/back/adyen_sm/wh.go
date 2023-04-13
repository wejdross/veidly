package adyen_sm

import (
	"encoding/json"
	"fmt"
	"sport/adyen"
	"time"

	"github.com/google/uuid"
)

func (ctx *Ctx) AdyenWh() adyen.HandlerFunc {
	return func(ni *adyen.NotificationRequestItem) error {

		p, err := ctx.GetPassInWebhook(ni)

		lr := WhLog{
			ID:        uuid.New(),
			Timestamp: time.Now().In(time.UTC),
		}

		sd := StateChangeEventSource{
			Type: Webhook,
			Data: lr.ID.String(),
		}

		lmsg, e := json.Marshal(ni)
		if e != nil {
			return e
		}

		lr.Message = lmsg

		if err != nil {
			lr.Error = err.Error()
			ctx.LogWhNotification(&lr)
			return nil
		}

		switch p.Data.State {
		case Link:
			fallthrough
		case WaitCancelOrRefund:
			fallthrough
		case CancelOrRefund:
			err = fmt.Errorf("rsv state: %s, is illegal state for p24", p.Data.State)
		case LinkExpress:
			err = ctx.MoveStateToCapture(p, sd, ni.OrderID, nil)
		case WaitRefund:
			err = ctx.MoveStateToRefund(p, sd)
		case WaitPayout:
			// err = ctx.WhSuccessPayoutThirdParty(p, sd)
		default:
			err = fmt.Errorf("unexpected state: %s", p.Data.State)
		}

		if err == nil {
			lr.IsSuccess = true
		} else {
			lr.Error = err.Error()
			lr.IsSuccess = false
		}

		ctx.LogWhNotification(&lr)

		return nil
	}
}

// func (ctx *Ctx) WhSuccessAuth(
// 	ni *adyen.NotificationRequestItem, p *Pass, sd StateChangeEventSource,
// ) error {
// 	switch p.Data.State {
// 	case LinkExpress:

// 		// proceed to capture payment right now
// 		return ctx.MoveStateToWaitCapture(p, sd, ni.PspReference, true)
// 		// return ctx.MoveStateToWaitCapture(
// 		// 	rsv, Webhook, whLogID.String(),
// 		// 	ni.PspReference, true)

// 	case Link:

// 		return ctx.MoveStateToHold(
// 			p, sd,
// 			map[string]interface{}{
// 				ctx.Model.Columns.PspReference: ni.PspReference,
// 			})

// 	// auth came after link expired, which is not cool -
// 	// may indicate:
// 	//  - error in the system
// 	//  - race condition
// 	//  - some unusual delay
// 	case LinkExpire:

// 		return ctx.MoveStateToError(p, sd)

// 	default:

// 		// Abort with ok.
// 		// adyen may send us duplicate AUTH responses

// 		return nil
// 	}

// }

// func (ctx *Ctx) WhSuccessCancelOrRefund(
// 	p *Pass, sd StateChangeEventSource,
// ) error {
// 	switch p.Data.State {
// 	case WaitCancelOrRefund:

// 		return ctx.MoveStateToCancelOrRefund(p, sd)

// 	// those situations shouldnt happen,
// 	// but they might if admin manages state by hand
// 	// so just accept notification
// 	case Error:
// 		fallthrough
// 	case WaitCapture:
// 		fallthrough
// 	case Capture:
// 		fallthrough
// 	case Hold:
// 		fallthrough
// 	case CancelOrRefund:
// 		fallthrough
// 	case RetryCancelOrRefund:

// 		return ctx.MoveStateToCancelOrRefund(p, sd)

// 		// those are really bad situations, and even if admin was responsible something may break
// 	case WaitPayout:
// 		fallthrough
// 	case Payout:
// 		return ctx.MoveStateToError(p, sd)

// 	default:
// 		return fmt.Errorf(
// 			"Unexpected rsv state encountered for CANCEL_OR_REFUND notification: %s",
// 			p.Data.State)
// 	}
// }

// func (ctx *Ctx) WhFailCancelOrRefund(
// 	p *Pass, sd StateChangeEventSource,
// ) error {
// 	// on failure notification we will retry and after that fails move to
// 	// error

// 	switch p.Data.State {
// 	case WaitCancelOrRefund:

// 		// get current number of retries for this state
// 		cd := p.Data.SmCache[WaitCancelOrRefund]

// 		if cd.Retries >= ctx.Config.SmMaxRetries {
// 			return ctx.MoveStateToError(p, sd)
// 		} else {

// 			cd.Retries++
// 			if p.Data.SmCache == nil {
// 				p.Data.SmCache = make(map[State]SmCacheEntry)
// 			}
// 			p.Data.SmCache[WaitCancelOrRefund] = cd

// 			return ctx.MoveStateToRetryCancelOrRefund(
// 				p, sd, map[string]interface{}{
// 					ctx.Model.Columns.SmCache:   &p.Data.SmCache,
// 					ctx.Model.Columns.SmTimeout: time.Now().In(time.UTC).Add(time.Duration(ctx.Config.SmRetryTimeout)),
// 				})

// 		}

// 	// those situations really shouldnt happen
// 	case CancelOrRefund:
// 		fallthrough
// 	case RetryCancelOrRefund:

// 		return ctx.MoveStateToCancelOrRefund(p, sd)

// 	default:

// 		return fmt.Errorf(
// 			"Unexpected rsv state encountered for CANCEL_OR_REFUND notification: %s",
// 			p.Data.State)

// 	}
// }

// func (ctx *Ctx) WhFailCapture(
// 	p *Pass, sd StateChangeEventSource,
// ) error {

// 	switch p.Data.State {
// 	case WaitCapture:

// 		// if retries allow - move to hold and try to redo capture after some time
// 		// otherwise move to cancel transaction

// 		cd := p.Data.SmCache[WaitCapture]
// 		if cd.Retries >= ctx.Config.SmMaxRetries {

// 			if err := ctx.MoveStateToWaitCancelOrRefund(p, sd, nil); err != nil {
// 				return err
// 			}

// 		} else {

// 			to := time.Now().In(time.UTC).Add(time.Duration(ctx.Config.SmRetryTimeout))

// 			// if its first failed capture then we will notify user via email - but only then
// 			// we dont want to spam their inbox every time retry iteration happen
// 			if cd.Retries == 0 && ctx.EmailUserAboutFailedCapture != nil {
// 				ctx.EmailUserAboutFailedCapture(p)
// 			}

// 			cd.Retries++
// 			if p.Data.SmCache == nil {
// 				p.Data.SmCache = make(map[State]SmCacheEntry)
// 			}
// 			p.Data.SmCache[WaitCapture] = cd

// 			return ctx.MoveStateToHold(
// 				p, sd,
// 				map[string]interface{}{
// 					ctx.Model.Columns.SmCache:   &p.Data.SmCache,
// 					ctx.Model.Columns.SmTimeout: to,
// 				},
// 			)
// 		}

// 		break

// 	case Error:

// 		// just abort. this request got logged anyways

// 		break

// 	case Capture:

// 		// this may indicate that because of error we sent multiple capture requests
// 		// move state to error ?

// 		break

// 	default:
// 		return fmt.Errorf(
// 			"Unexpected rsv state encountered for CAPTURE notification: %s",
// 			p.Data.State)
// 	}

// 	return nil
// }

// // https://docs.adyen.com/online-payments/online-payouts/payout-notifications

// func (ctx *Ctx) WhSuccessPayoutThirdParty(
// 	p *Pass, sd StateChangeEventSource,
// ) error {
// 	switch p.Data.State {
// 	case WaitPayout:
// 		return ctx.MoveStateToPayout(p, sd)
// 	case Payout:
// 		return nil
// 	default:
// 		return fmt.Errorf(
// 			"Unexpected rsv state encountered for PAYOUT_THIRDPARTY notification: %s",
// 			p.Data.State)

// 	}
// }

// func (ctx *Ctx) WhFailPayoutThirdParty(
// 	p *Pass, sd StateChangeEventSource,
// ) error {
// 	switch p.Data.State {
// 	case WaitPayout:

// 		cd := p.Data.SmCache[WaitPayout]
// 		if cd.Retries >= ctx.Config.SmMaxRetries {

// 			return ctx.MoveStateToError(p, sd)

// 		} else {
// 			to := time.Now().In(time.UTC).Add(time.Duration(ctx.Config.SmRetryTimeout))

// 			// if its first failed capture then we will notify user via email - but only then
// 			// we dont want to spam their inbox every time retry iteration happen
// 			if cd.Retries == 0 && ctx.EmailInstrAboutFailedPayout != nil {
// 				ctx.EmailInstrAboutFailedPayout(p)
// 			}
// 			cd.Retries++
// 			if p.Data.SmCache == nil {
// 				p.Data.SmCache = make(map[State]SmCacheEntry)
// 			}
// 			p.Data.SmCache[WaitPayout] = cd

// 			return ctx.MoveStateToCapture(
// 				p, sd,
// 				map[string]interface{}{
// 					ctx.Model.Columns.SmCache:   &p.Data.SmCache,
// 					ctx.Model.Columns.SmTimeout: to,
// 				})

// 		}

// 	case Payout:

// 		return ctx.MoveStateToError(p, sd)

// 	default:
// 		return fmt.Errorf(
// 			"Unexpected rsv state encountered for PAYOUT_THIRDPARTY notification: %s",
// 			p.Data.State)
// 	}
// }

// func (ctx *Ctx) WhPaidoutReversed(p *Pass, sd StateChangeEventSource) error {

// 	if err := ctx.MoveStateToError(p, sd); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (ctx *Ctx) WhSuccessRefund(p *Pass, sd StateChangeEventSource) error {
// 	switch p.Data.State {
// 	case WaitRefund:
// 		return ctx.MoveStateToRefund(p, sd)
// 	default:
// 		return fmt.Errorf(
// 			"Unexpected rsv state encountered for PAYOUT_THIRDPARTY notification: %s",
// 			p.Data.State)

// 	}
// }

// func (ctx *Ctx) WhFailRefund(p *Pass, sd StateChangeEventSource) error {
// 	switch p.Data.State {
// 	case WaitPayout:
// 		return ctx.MoveStateToError(p, sd)
// 	default:
// 		return fmt.Errorf(
// 			"Unexpected rsv state encountered for PAYOUT_THIRDPARTY notification: %s",
// 			p.Data.State)

// 	}
// }
