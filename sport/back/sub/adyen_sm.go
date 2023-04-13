package sub

import (
	"database/sql"
	"fmt"
	"sport/adyen"
	"sport/adyen_sm"
	"sport/helpers"
	"sport/invoicing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const AdyenSmKey = "SUB"
const AdyenSmPrefix = "SUB_"

func SmPassToSub(d *adyen_sm.Pass) *Sub {
	return d.Args.(*Sub)
}

func (ctx *Ctx) SubToSmPass(r *Sub) adyen_sm.Pass {
	return adyen_sm.Pass{
		Args: r,
		Data: adyen_sm.SmData{
			ID:            r.ID,
			InstructorID:  r.InstructorID,
			State:         r.State,
			LinkID:        r.LinkID,
			SmRetries:     r.SmRetries,
			OrderID:       r.OrderID,
			RefID:         r.RefID.String(),
			SmCache:       r.SmCache,
			IsConfirmed:   r.IsConfirmed,
			SmTimeout:     r.SmTimeout,
			ManualConfirm: false,
			Decision:      adyen_sm.Unset,
			TotalPrice:    r.SubModel.Total(),
			PayoutValue:   r.SubModel.PayoutValue,
			RefundValue:   r.SubModel.RefundValue,
			InstrPrice:    r.SubModel.Price,
			Currency:      r.SubModel.Currency,
			// this is needed only when making payout (MoveStateToWaitPayout) request
			// so i will set it in beforePayout callback to avoid making needless joins,
			// or querying extra data during scans
			ShopperReference: "",
			SmType:           AdyenSmKey,
		},
	}
}

func (ctx *Ctx) SubToSmPassPtr(r *Sub) *adyen_sm.Pass {
	p := ctx.SubToSmPass(r)
	return &p
}

func (ctx *Ctx) SetTimeoutForPayout(d *adyen_sm.Pass, args map[string]interface{}, tx *sql.Tx, sd adyen_sm.StateChangeEventSource) (bool, error) {

	// timeout can be set to retry, in that case dont override it
	if _, ok := args["sm_timeout"]; ok {
		return true, nil
	}

	// delay before payout is made
	args["sm_timeout"] = time.Now().In(time.UTC).Add(time.Duration(ctx.Config.PayoutDelay))

	return true, nil
}

func (ctx *Ctx) SetShopperReference(
	d *adyen_sm.Pass,
	args map[string]interface{},
	tx *sql.Tx,
	sd adyen_sm.StateChangeEventSource) (bool, error) {
	/*
		TODO:
		optionally, on error or when pi.CardRefID == ""
		send mail to instructor with notificaiton that he didnt set up payment data
	*/
	s := SmPassToSub(d)
	//
	pi, err := ctx.Instr.DalReadCardInfo(s.InstructorID)
	if err != nil {
		return false, err
	}
	d.Data.ShopperReference = pi.CardRefID

	return true, nil
}

func (ctx *Ctx) GetPassByNotificationRequestItem(ni *adyen.NotificationRequestItem) (*adyen_sm.Pass, error) {
	var refid uuid.UUID
	var err error
	if refid, err = uuid.Parse(ni.SessionID); err != nil {
		return nil, err
	}
	s, err := ctx.DalReadSingleSubWithJoins(
		ReadSubRequest{
			RefID: &refid,
		})
	if err != nil {
		return nil, err
	}
	p := ctx.SubToSmPass(&s.Sub)
	return &p, nil
}

func (ctx *Ctx) GetAgentSubs() ([]adyen_sm.Pass, error) {
	sbs, err := ctx.DalReadSubWithJoins(EmptyReadSubRequest)
	if err != nil {
		return nil, err
	}
	p := make([]adyen_sm.Pass, 0, len(sbs))
	for i := range sbs {
		s := &sbs[i]
		x := ctx.SubToSmPass(&s.Sub)
		if err != nil {
			return nil, err
		}
		p = append(p, x)
	}
	return p, nil
}

type SmApiRequest struct {
	SubID uuid.UUID
}

func (req *SmApiRequest) Validate() error {
	if req.SubID == uuid.Nil {
		return fmt.Errorf("validate SmApiRequest: invalid sub id")
	}
	return nil
}

type DisputeSubRequest struct {
	SmApiRequest
	adyen_sm.DisputeRequest
}

func (ctx *Ctx) GetPassFromSmApiReq(g *gin.Context, c adyen_sm.SmHandlerCaller, rawBody []byte) (*adyen_sm.Pass, error) {
	var s *SubWithJoins
	var err error
	var userID uuid.UUID
	var req SmApiRequest

	if err = helpers.ParseAndValidateJson(rawBody, &req, req.Validate); err != nil {
		g.AbortWithError(400, err)
		return nil, err
	}

	userID, err = ctx.Api.AuthorizeUserFromCtx(g)
	if err != nil {
		g.AbortWithError(401, err)
		return nil, err
	}

	switch c {
	case adyen_sm.User:
		s, err = ctx.DalReadSingleSubWithJoins(
			ReadSubRequest{
				UserID: &userID,
				ID:     &req.SubID,
			})
		break
	case adyen_sm.Instr:
		s, err = ctx.DalReadSingleSubWithJoins(
			ReadSubRequest{
				InstrUserID: &userID,
				ID:          &req.SubID,
			})
		break
	default:
		g.AbortWithError(400, fmt.Errorf("Invalid type: %s", g.Param("type")))
		return nil, err
	}

	if err != nil {
		if err == sql.ErrNoRows {
			g.AbortWithError(404, err)
		} else {
			g.AbortWithError(500, err)
		}
		return nil, err
	}

	p := ctx.SubToSmPass(&s.Sub)

	return &p, nil
}

func (ctx *Ctx) SetAdyenSm(invoicing *invoicing.Ctx) {

	cb := adyen_sm.SmCallbacks{
		BeforeHold:           nil,
		BeforeLinkExpire:     nil,
		BeforeCancelOrRefund: nil,
		AfterPayout:          nil,
		BeforeRefund:         nil,
		BeforeCapture:        ctx.SetTimeoutForPayout,
		BeforeWaitRefund:     nil,
		BeforeWaitPayout:     ctx.SetShopperReference,

		// hold is not used in subscriptions
		EmailUserAboutHold:          nil,
		EmailInstrAboutCancel:       ctx.EmailInstrAboutCancelAsync,
		EmailUserAboutCancel:        ctx.EmailUserAboutCancelAsync,
		EmailUserAboutCapture:       ctx.EmailUserAboutCaptureAsync,
		EmailInstrAboutDispute:      ctx.EmailInstrAboutDisputeAsync,
		EmailUserAboutDispute:       ctx.EmailUserAboutDisputeAsync,
		EmailUserAboutFailedCapture: ctx.EmailUserAboutFailedCaptureAsync,
		EmailInstrAboutFailedPayout: ctx.EmailInstrAboutFailedPayoutAsync,

		GetPassInWebhook: ctx.GetPassByNotificationRequestItem,
		GetPassInAgent:   ctx.GetAgentSubs,
		GetPassInHandler: ctx.GetPassFromSmApiReq,
	}

	asm := adyen_sm.NewCtx(
		ctx.Api, ctx.Dal,
		ctx.Instr, ctx.Adyen,
		ctx.noReplyCtx, cb,
		adyen_sm.DDLModel{
			TableName: "subs",
			Columns:   adyen_sm.DefaultColumns,
		},
		&ctx.Config.AdyenSmConfig,
		invoicing)

	// ctx.Api.AnonGroup.POST("/sub/expire", asm.HandlerExpireLink())
	ctx.Api.AnonGroup.POST("/sub/dispute/:type", asm.HandlerDispute())
	ctx.Api.AnonGroup.POST("/sub/dispute", asm.HandlerDispute())

	ctx.Adyen.AddAdyenHandler(AdyenSmKey, asm.AdyenWh())

	ctx.AdyenSm = asm
}
