package rsv

import (
	"database/sql"
	"fmt"
	"os"
	"sport/adyen"
	"sport/adyen_sm"
	"sport/helpers"
	"sport/invoicing"
	"sport/review"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetRsvFromPass(d *adyen_sm.Pass) *DDLRsvWithInstr {
	return d.Args.(*DDLRsvWithInstr)
}

/*
	this function is idempotent - you may call it multiple times on single rsv without risk of multiple rollbacks
	however like many functions used in sm - this is not thread safe.
*/
func (ctx *Ctx) RollbackDiscountCodeUse(d *adyen_sm.Pass, args map[string]interface{}, tx *sql.Tx, sd adyen_sm.StateChangeEventSource) (bool, error) {
	// update db colun - indicate that discount code has been rollbacked
	args["dc_rollback"] = true
	r := GetRsvFromPass(d)
	if r.Dc != nil && !r.DcRollback {
		if err := ctx.Dc.DalUpdateDcUsesTx(r.Dc.ID, -1, tx); err != nil {
			if err != sql.ErrNoRows {
				return false, err
			}
		}
	}
	return true, nil
}

func (ctx *Ctx) VerifyIfRsvAvailable(d *adyen_sm.Pass, args map[string]interface{}, tx *sql.Tx, sd adyen_sm.StateChangeEventSource) (bool, error) {
	r := GetRsvFromPass(d)

	// this check is here to ensure that this is first capture, and not retry
	cd := d.Data.SmCache[adyen_sm.WaitPayout]
	if cd.Retries > 0 {
		return true, nil
	}

	c, err := ctx.ValidateCapacity(r.Training, r.DateStart, tx)
	if err != nil {
		return false, err
	}

	if !c {
		// note that i dont use Tx here since tx will be rollbacked once i return false
		err := ctx.AdyenSm.MoveStateToCancelOrRefund(d, adyen_sm.ManualSrc, nil)
		return false, err
	}

	args["sm_timeout"] = r.DateEnd.Add(time.Duration(ctx.Config.PayoutDelay))

	return true, err
}

func (ctx *Ctx) AdyenSmCreateReviewToken(d *adyen_sm.Pass) {
	r := GetRsvFromPass(d)

	// to make opinion review must be enabled
	if ctx.Review == nil {
		return
	}

	opinionReq := review.ReviewRequest{
		TrainingID: r.TrainingID,
		RsvID:      r.ID,
		// warning: this may be uuid.nil
		UserID:   r.UserID,
		Email:    r.UserContactData.Email,
		UserInfo: r.UserInfo,
	}
	o := opinionReq.NewOpinion()

	if err := ctx.Review.DalCreateReviewToken(o); err != nil {
		fmt.Fprintf(os.Stderr, "DalCreateReviewToken: %v\n", err)
		// TODO: notify admin about fuckup
	}
}

func (ctx *Ctx) ValidateIfUserCanRefund(d *adyen_sm.Pass, args map[string]interface{}, tx *sql.Tx, sd adyen_sm.StateChangeEventSource) (bool, error) {
	// check if refund is made from handler and is requested by user
	if sd.Type != adyen_sm.Handler || !strings.Contains(sd.Data, string(adyen_sm.User)) {
		return true, nil
	}

	rsv := GetRsvFromPass(d)

	// training happened already,
	// only viable option is dispute and manual resolution
	if rsv.DateStart.Before(time.Now()) {
		// err = ctx.MoveStateToDispute(
		// 	&rsv.Reservation, Handler, t+" refund")
		// if err != nil {
		// 	g.AbortWithError(500, err)
		// 	return
		// }
		// break
		return false, helpers.NewHttpError(
			412,
			"",
			fmt.Errorf("user cant request refund if training has started already, Use dispute instead"))
	}

	return true, nil
}

func (ctx *Ctx) SetTimeoutBeforeCap(d *adyen_sm.Pass, args map[string]interface{}, tx *sql.Tx, sd adyen_sm.StateChangeEventSource) (bool, error) {

	// timeout can be set to retry, in that case dont override it
	if _, ok := args["sm_timeout"]; ok {
		return true, nil
	}

	now := time.Now().In(time.UTC)
	rsv := GetRsvFromPass(d)

	var capTs time.Time
	df := rsv.DateStart.Sub(now)

	const day = time.Hour * 24

	// this scenario can only occurs if system enters fail state,
	// gets disconnected for long time, etc...
	// right now my solution is to follow the flow and capture payment asap
	if df <= 0 {
		capTs = now
	} else if df >= day {
		capTs = now.Add(day)
	} else { // df > 0 && df < Day
		capTs = now
	}

	args["sm_timeout"] = capTs

	return true, nil
}

func (ctx *Ctx) EmailUserAboutHoldAsync(d *adyen_sm.Pass) {
	rsv := GetRsvFromPass(d)
	if rsv.UserContactData.Email == "" || ctx.NoReplyCtx == nil {
		return
	}
	go func() {
		err := ctx.emailUserAboutHoldSync(&rsv.DDLRsv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "EmailUserAboutHold: %s", err)
		}
	}()
}

func (ctx *Ctx) EmailUserAboutCaptureAsync(d *adyen_sm.Pass) {
	rsv := GetRsvFromPass(d)
	if rsv.UserContactData.Email == "" || ctx.NoReplyCtx == nil {
		return
	}
	go func() {
		err := ctx.emailUserAboutCaptureSync(&rsv.DDLRsv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "EmailUserAboutCapture: %s", err)
		}
	}()
}

func (ctx *Ctx) EmailUserAboutDisputeAsync(d *adyen_sm.Pass) {
	rsv := GetRsvFromPass(d)
	if rsv.UserContactData.Email == "" || ctx.NoReplyCtx == nil {
		return
	}
	go func() {
		err := ctx.emailUserAboutDisputeSync(&rsv.DDLRsv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "EmailUserAboutDispute: %s", err)
		}
	}()
}

func (ctx *Ctx) EmailInstrAboutDisputeAsync(d *adyen_sm.Pass) {
	rsv := GetRsvFromPass(d)
	if ctx.NoReplyCtx == nil ||
		rsv.Instructor == nil ||
		rsv.Instructor.ContactData.Email == "" {
		return
	}
	go func() {
		err := ctx.emailInstrAboutDisputeSync(rsv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "EmailInstrAboutDispute: %s", err)
		}
	}()
}

func (ctx *Ctx) EmailUserAboutCancelAsync(d *adyen_sm.Pass) {
	rsv := GetRsvFromPass(d)
	if ctx.NoReplyCtx == nil ||
		rsv.UserContactData.Email == "" {
		return
	}
	go func() {
		err := ctx.emailUserAboutCancelSync(&rsv.DDLRsv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "EmailUserAboutCancel: %s", err)
		}
	}()
}

func (ctx *Ctx) EmailInstrAboutCancelAsync(d *adyen_sm.Pass) {
	rsv := GetRsvFromPass(d)
	if ctx.NoReplyCtx == nil ||
		rsv.Instructor == nil ||
		rsv.Instructor.ContactData.Email == "" {
		return
	}
	go func() {
		err := ctx.emailInstrAboutCancelSync(rsv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "emailInstrAboutCancel: %s", err)
		}
	}()
}

func (ctx *Ctx) EmailUserAboutFailedCaptureAsync(d *adyen_sm.Pass) {
	rsv := GetRsvFromPass(d)
	if ctx.NoReplyCtx == nil ||
		rsv.UserContactData.Email == "" {
		return
	}
	go func() {
		err := ctx.emailUserAboutFailedCaptureSync(&rsv.DDLRsv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "EmailUserAboutFailedCapture: %s", err)
		}
	}()
}

func (ctx *Ctx) EmailInstrAboutFailedPayoutAsync(d *adyen_sm.Pass) {
	rsv := GetRsvFromPass(d)
	if ctx.NoReplyCtx == nil ||
		rsv.Instructor == nil ||
		rsv.Instructor.ContactData.Email == "" {
		return
	}
	go func() {
		err := ctx.emailInstrAboutFailedPayoutSync(rsv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "emailUserAboutFailedPayout: %s", err)
		}
	}()
}

func (ctx *Ctx) RsvResponseToSmPass(r *DDLRsvWithInstr) adyen_sm.Pass {
	pi := ctx.GetRsvPricing(r.ToTrainingWithJoins())
	var iid uuid.UUID
	if r.InstructorID == nil || r.Instructor == nil {
		panic("cant have rsv in adyen flow without instructor")
	}
	iid = *r.InstructorID
	return adyen_sm.Pass{
		Args: r,
		Data: adyen_sm.SmData{
			ID:               r.ID,
			InstructorID:     iid,
			State:            r.State,
			LinkID:           r.LinkID,
			SmRetries:        r.SmRetries,
			OrderID:          r.OrderID,
			SmCache:          r.SmCache,
			IsConfirmed:      r.IsConfirmed,
			SmTimeout:        r.SmTimeout,
			ManualConfirm:    r.Training.ManualConfirm,
			Decision:         r.InstructorDecision,
			TotalPrice:       pi.TotalPrice,
			PayoutValue:      pi.SplitPayout,
			RefundValue:      pi.RefundAmount,
			InstrPrice:       pi.InstrPrice,
			Currency:         r.Training.Currency,
			ShopperReference: r.Instructor.CardInfo.CardRefID,
			SmType:           AdyenSmKey,
		},
	}
}

func (ctx *Ctx) RsvResponseToSmPassPtr(r *DDLRsvWithInstr) *adyen_sm.Pass {
	p := ctx.RsvResponseToSmPass(r)
	return &p
}

func (ctx *Ctx) GetRsvByNotificationRequestItem(ni *adyen.NotificationRequestItem) (*adyen_sm.Pass, error) {
	id, err := uuid.Parse(ni.SessionID)
	if err != nil {
		return nil, err
	}
	r, err := ctx.ReadSingleRsv(
		ReadRsvsArgs{
			// RefID:          &ni.SessionID,
			ID:             &id,
			WithInstructor: true,
		})
	if err != nil {
		return nil, err
	}
	if r.InstructorID == nil {
		return nil, fmt.Errorf("GetRsvByNotificationRequestItem: InstructorID is nil")
	}
	p := ctx.RsvResponseToSmPass(r)
	return &p, nil
}

func (ctx *Ctx) GetAgentRsvs() ([]adyen_sm.Pass, error) {
	rsvs, err := ctx.ReadDDLRsvs(ReadRsvsArgs{
		WithInstructor: true,
	})
	if err != nil {
		return nil, err
	}
	p := make([]adyen_sm.Pass, 0, len(rsvs.Rsv))
	for i := range rsvs.Rsv {
		r := &rsvs.Rsv[i]
		if r.InstructorID == nil {
			continue
		}
		x := ctx.RsvResponseToSmPass(r)
		if err != nil {
			return nil, err
		}
		p = append(p, x)
	}
	return p, nil
}

type SmApiRequest struct {
	ReservationID uuid.UUID
	AccessToken   uuid.UUID
}

func (req *SmApiRequest) Validate() error {
	if req.ReservationID == uuid.Nil {
		return fmt.Errorf("validate SmApiRequest: invalid rsv id")
	}
	return nil
}

func (ctx *Ctx) GetRsvFromSmApiReq(g *gin.Context, c adyen_sm.SmHandlerCaller, rawBody []byte) (*adyen_sm.Pass, error) {
	var rsv *DDLRsvWithInstr
	var err error
	var userID uuid.UUID
	var req SmApiRequest

	if err = helpers.ParseAndValidateJson(rawBody, &req, req.Validate); err != nil {
		g.AbortWithError(400, err)
		return nil, err
	}

	if g.GetHeader("Authorization") != "" {
		userID, err = ctx.Api.AuthorizeUserFromCtx(g)
		if err != nil {
			g.AbortWithError(401, err)
			return nil, err
		}
		switch c {
		case adyen_sm.User:
			rsv, err = ctx.ReadSingleRsv(
				ReadRsvsArgs{
					WithInstructor: true,
					UserID:         &userID,
					ID:             &req.ReservationID,
				})
		case adyen_sm.Instr:
			rsv, err = ctx.ReadSingleRsv(
				ReadRsvsArgs{
					WithInstructor: true,
					InstrUserID:    &userID,
					ID:             &req.ReservationID,
				})
		default:
			g.AbortWithError(400, fmt.Errorf("invalid type: %s", g.Param("type")))
			return nil, err
		}
	} else {
		if c == adyen_sm.Instr {
			g.AbortWithError(401, fmt.Errorf("instructor must authorize to use this endpoint"))
			return nil, err
		}
		if req.AccessToken == uuid.Nil {
			g.AbortWithError(401, fmt.Errorf("no access_token provided"))
			return nil, err
		}
		rsv, err = ctx.ReadSingleRsv(ReadRsvsArgs{AccessToken: &req.AccessToken, WithInstructor: true})
	}

	if err != nil {
		if err == sql.ErrNoRows {
			g.AbortWithError(404, err)
		} else {
			g.AbortWithError(500, err)
		}
		return nil, err
	}

	p := ctx.RsvResponseToSmPass(rsv)

	return &p, nil
}

func (ctx *Ctx) SetAdyenSm(invoicing *invoicing.Ctx) {

	cb := adyen_sm.SmCallbacks{
		BeforeHold:           ctx.SetTimeoutBeforeCap,
		BeforeLinkExpire:     ctx.RollbackDiscountCodeUse,
		BeforeCancelOrRefund: ctx.RollbackDiscountCodeUse,
		AfterPayout:          ctx.AdyenSmCreateReviewToken,
		BeforeRefund:         ctx.RollbackDiscountCodeUse,
		BeforeCapture:        ctx.VerifyIfRsvAvailable,
		BeforeWaitRefund:     ctx.ValidateIfUserCanRefund,
		BeforeWaitPayout:     nil,

		EmailUserAboutHold:          ctx.EmailUserAboutHoldAsync,
		EmailInstrAboutCancel:       ctx.EmailInstrAboutCancelAsync,
		EmailUserAboutCancel:        ctx.EmailUserAboutCancelAsync,
		EmailUserAboutCapture:       ctx.EmailUserAboutCaptureAsync,
		EmailInstrAboutDispute:      ctx.EmailInstrAboutDisputeAsync,
		EmailUserAboutDispute:       ctx.EmailUserAboutDisputeAsync,
		EmailUserAboutFailedCapture: ctx.EmailUserAboutFailedCaptureAsync,
		EmailInstrAboutFailedPayout: ctx.EmailInstrAboutFailedPayoutAsync,

		GetPassInWebhook: ctx.GetRsvByNotificationRequestItem,
		GetPassInAgent:   ctx.GetAgentRsvs,
		GetPassInHandler: ctx.GetRsvFromSmApiReq,
	}

	asm := adyen_sm.NewCtx(
		ctx.Api, ctx.Dal,
		ctx.Instr, ctx.Adyen,
		ctx.NoReplyCtx, cb,
		adyen_sm.DDLModel{
			TableName: "reservations",
			Columns:   adyen_sm.DefaultColumns,
		},
		&ctx.Config.AdyenSmConfig,
		invoicing)

	ctx.Api.JwtGroup.POST("/rsv/decision", asm.HandlerDecision())
	ctx.Api.AnonGroup.POST("/rsv/cancel", asm.HandlerCancel())
	// ctx.Api.AnonGroup.POST("/rsv/expire", asm.HandlerExpireLink())
	ctx.Api.AnonGroup.POST("/rsv/dispute/:type", asm.HandlerDispute())
	ctx.Api.JwtGroup.POST("/rsv/refund/:type", asm.HandlerRefund())

	ctx.Adyen.AddAdyenHandler(AdyenSmKey, asm.AdyenWh())

	ctx.AdyenSm = asm
}
