package adyen_sm

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type AgentErrorLocation int

const (
	LocDaemon     AgentErrorLocation = 1
	LocDaemonIter AgentErrorLocation = 2
)

func (ctx *Ctx) notifyAdminAboutAgentError(
	where AgentErrorLocation, err error, pass *Pass,
) {
	msg := ""
	switch where {
	case LocDaemonIter:
		msg = fmt.Sprintf(
			"For rsv: %v, DaemonDoRsv returned error: %v\n", pass.Data.ID, err)
	case LocDaemon:
		msg = fmt.Sprintf(
			"In rsv daemon process iter returned error: %v\n", err)
	default:
		msg = fmt.Sprintf(
			"For unkown location error occurred: %v\n", err)
	}
	go ctx.MailSupportAboutAgentErr(msg)
	fmt.Fprintln(os.Stderr, msg)
}

/*
	this routine is used to capture / cancel / confirm reservations
*/

const Day time.Duration = 24 * time.Hour

func (ctx *Ctx) AgentDoOne(p *Pass) error {
	now := time.Now().In(time.UTC)
	var sd = StateChangeEventSource{Agent, ""}
	switch p.Data.State {
	case Link:

		// if now.After(p.Data.SmTimeout) {
		// 	return ctx.MoveStateToLinkExpire(p, sd)
		// }

		break

	case LinkExpress:

		// if now.After(p.Data.SmTimeout) {
		// 	return ctx.MoveStateToLinkExpire(p, sd)
		// }

		break

	case Hold:

		// if r.InstructorDecision == Reject {
		// 	return MoveStateToPreCancelOrRefund(r)
		// }

		// Note that since Hold state is only achievable from Link state
		// we still should have AT LEAST 18 hours before training
		// so there is plenty of time to perform capture

		// now this may not be true in case of service going offline for long time

		// var err error

		// if now.After(p.Data.SmTimeout) {

		// 	if p.Data.ManualConfirm {
		// 		if p.Data.Decision == Approve {
		// 			err = ctx.MoveStateToWaitCapture(p, sd, p.Data.PspReference, false)
		// 		} else {
		// 			// instructor had time for confirm.
		// 			// but he didnt - move to cancel
		// 			err = ctx.MoveStateToWaitCancelOrRefund(p, sd, nil)
		// 		}
		// 	} else {
		// 		if p.Data.Decision == Reject {
		// 			err = ctx.MoveStateToWaitCancelOrRefund(p, sd, nil)
		// 		} else {
		// 			err = ctx.MoveStateToWaitCapture(p, sd, p.Data.PspReference, false)
		// 		}
		// 	}

		// }

		// return err

		break

	case WaitCapture:

		if now.After(p.Data.SmTimeout) {
			// timeout on capture so move to error

			ctx.LogStateChangeEvent(
				&StateChangeEvent{
					ID:            uuid.New(),
					ObjKey:        p.Data.ID,
					PreviousState: string(WaitCapture),
					NextState:     string(Capture),
					Timestamp:     time.Now().In(time.UTC),
					Success:       false,
					Source:        sd,
					Error:         "timeout while waiting for webhook notification",
				}, nil,
			)

			if p.Data.SmRetries >= ctx.Config.SmMaxRetries {
				return ctx.MoveStateToError(p, sd)
			}

			return ctx.UpdateRetries(p.Data.SmRetries+1, p.Data.ID, p.Data.State)
		}

		break

	case Capture:

		if now.After(p.Data.SmTimeout) {
			return ctx.MoveStateToWaitPayout(p, sd)
		}

		break

	case WaitPayout:

		if now.After(p.Data.SmTimeout) {

			ctx.LogStateChangeEvent(
				&StateChangeEvent{
					ID:            uuid.New(),
					ObjKey:        p.Data.ID,
					PreviousState: string(WaitPayout),
					NextState:     string(Payout),
					Timestamp:     time.Now().In(time.UTC),
					Success:       false,
					Source:        sd,
					Error:         "timeout while waiting for webhook notification",
				}, nil,
			)

			if p.Data.SmRetries >= ctx.Config.SmMaxRetries {
				return ctx.MoveStateToError(p, sd)
			}

			return ctx.UpdateRetries(p.Data.SmRetries+1, p.Data.ID, p.Data.State)
		}

		break

	case Payout:

		break

	case LinkExpire:

		break

	case Error:

		break
	case CancelOrRefund:

		break

	case WaitCancelOrRefund:

		if now.After(p.Data.SmTimeout) {
			// timeout on capture so move to error

			ctx.LogStateChangeEvent(
				&StateChangeEvent{
					ID:            uuid.New(),
					ObjKey:        p.Data.ID,
					PreviousState: string(WaitCancelOrRefund),
					NextState:     string(CancelOrRefund),
					Timestamp:     time.Now().In(time.UTC),
					Success:       false,
					Source:        sd,
					Error:         "timeout while waiting for webhook notification",
				}, nil,
			)

			if p.Data.SmRetries >= ctx.Config.SmMaxRetries {
				return ctx.MoveStateToError(p, sd)
			}

			return ctx.UpdateRetries(p.Data.SmRetries+1, p.Data.ID, p.Data.State)
		}

		break

	case RetryCancelOrRefund:

		if now.After(p.Data.SmTimeout) {
			return ctx.MoveStateToWaitCancelOrRefund(p, sd, nil)
		}

		break

	case WaitRefund:

		if now.After(p.Data.SmTimeout) {
			// timeout on capture so move to error

			ctx.LogStateChangeEvent(
				&StateChangeEvent{
					ID:            uuid.New(),
					ObjKey:        p.Data.ID,
					PreviousState: string(WaitRefund),
					NextState:     string(Refund),
					Timestamp:     time.Now().In(time.UTC),
					Success:       false,
					Source:        sd,
					Error:         "timeout while waiting for webhook notification",
				}, nil,
			)

			if p.Data.SmRetries >= ctx.Config.SmMaxRetries {
				return ctx.MoveStateToError(p, sd)
			}

			return ctx.UpdateRetries(p.Data.SmRetries+1, p.Data.ID, p.Data.State)
		}

		break

	case Refund:

		break

	case Dispute:

		break

	default:

		return fmt.Errorf("unsupported obj state: %s", p.Data.State)

	}

	return nil
}

func (ctx *Ctx) IterAgent() error {
	ps, err := ctx.GetPassInAgent()
	if err != nil {
		return err
	}

	for i := range ps {
		// its pretty good idea to do this in semaphore.
		// no race conditions should occur so it should be fine
		err := ctx.AgentDoOne(&ps[i])
		if err != nil {
			ctx.notifyAdminAboutAgentError(
				LocDaemonIter, err, &ps[i],
			)
		}
	}

	return nil
}

func (ctx *Ctx) RunAgent() {
	for {
		if err := ctx.IterAgent(); err != nil {
			ctx.notifyAdminAboutAgentError(
				LocDaemon, err, nil,
			)
		}
		time.Sleep(1 * time.Minute)
	}
}
