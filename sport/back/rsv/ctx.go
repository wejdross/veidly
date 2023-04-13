package rsv

import (
	"fmt"
	"sport/adyen"
	"sport/adyen_sm"
	"sport/api"
	"sport/dal"
	"sport/dc"
	"sport/helpers"
	"sport/instr"
	"sport/invoicing"
	"sport/notify"
	"sport/review"
	"sport/sub"
	"sport/train"
	"sport/user"
)

const MaxMark = 6

const FreeInstructorShots = 2

// config
type Config struct {
	MaxPageSize      int    `yaml:"max_page_size"`
	RsvDetailsUrlFmt string `yaml:"rsv_details_url_fmt"`
	PayoutConfigUrl  string `yaml:"payout_config_url"`

	ServiceFee    int `yaml:"service_fee"`
	ProcessingFee int `yaml:"processing_fee"`
	RefundAmount  int `yaml:"refund_amount"`

	//
	LinkExpire        helpers.Duration `yaml:"link_expire"`
	LinkExpressExpire helpers.Duration `yaml:"link_express_expire"`
	LinkManualExpire  helpers.Duration `yaml:"link_manual_expire"`

	//
	LinkAtLeastBefore        helpers.Duration `yaml:"link_at_least_before"`
	LinkManualAtLeastBefore  helpers.Duration `yaml:"link_manual_at_least_before"`
	LinkExpressAtLeastBefore helpers.Duration `yaml:"link_express_at_least_before"`

	//
	AdyenSmConfig adyen_sm.Config `yaml:"adyen_sm"`

	// after end of rsv, payout will be delayed for this duration
	// this ensures that user has enough time to notify us about problems
	PayoutDelay helpers.Duration `yaml:"payout_delay_after_rsv_end"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Validate Config: " + m)
}

func (r *Config) Validate() error {
	if r.LinkExpire <= 0 {
		return r.ValidationErr("invalid link_expire")
	}
	if r.LinkExpressExpire <= 0 {
		return r.ValidationErr("invalid link_express_expire")
	}
	if r.LinkManualExpire <= 0 {
		return r.ValidationErr("invalid link_manual_expire")
	}
	if r.LinkAtLeastBefore <= 0 {
		return r.ValidationErr("invalid link_at_least_before")
	}
	if r.LinkManualAtLeastBefore <= 0 {
		return r.ValidationErr("invalid link_manual_at_least_before")
	}
	if r.ServiceFee > 100 || r.ServiceFee < 0 {
		return r.ValidationErr("invalid service_fee")
	}
	if r.ProcessingFee > 100 || r.ProcessingFee < 0 {
		return r.ValidationErr("invalid service_fee")
	}
	if r.RefundAmount > 100 || r.RefundAmount < 0 {
		return r.ValidationErr("invalid refund_amount")
	}
	if r.RsvDetailsUrlFmt == "" {
		return r.ValidationErr("invalid rsv_details_url_fmt")
	}
	if r.PayoutConfigUrl == "" {
		return r.ValidationErr("invalid payout_config_url")
	}
	return nil
}

type Ctx struct {
	Api           *api.Ctx
	Dal           *dal.Ctx
	User          *user.Ctx
	Adyen         *adyen.Ctx
	Train         *train.Ctx
	Instr         *instr.Ctx
	Review        *review.Ctx
	Dc            *dc.Ctx
	Config        *Config
	NoPaymentFlow bool
	Sub           *sub.Ctx

	NoReplyCtx notify.EmailSender

	AdyenSm *adyen_sm.Ctx
}

func (ctx *Ctx) RegisterAllHandlers() {
	ctx.Api.AnonGroup.POST("/rsv", ctx.HandlerPostRsv())
	ctx.Api.AnonGroup.POST("/rsv/pricing", ctx.HandlerPostRsvPricing())
	ctx.Api.AnonGroup.GET("/rsv/t/:type", ctx.HandlerGetRsv())

	ctx.Api.AnonGroup.GET("/rsv/instr/contact", ctx.HandlerGetInstrContact())
}

// noReplyCtx may be nil
func (ctx *Ctx) SetNoReplyCtx(noReplyCtx notify.EmailSender) {
	ctx.NoReplyCtx = noReplyCtx
}

const AdyenSmKey = "RSV"

var adyenWhPrefix = fmt.Sprintf("%s_", AdyenSmKey)

func (ctx *Ctx) RunAgent() {
	ctx.AdyenSm.RunAgent()
}

// reviewCtx can be nil
func NewCtx(
	apiCtx *api.Ctx,
	dalCtx *dal.Ctx,
	userCtx *user.Ctx,
	instrCtx *instr.Ctx,
	trainCtx *train.Ctx,
	adyenCtx *adyen.Ctx,
	reviewCtx *review.Ctx,
	noReplyCtx notify.EmailSender,
	dcCtx *dc.Ctx,
	subCtx *sub.Ctx,
	invoicing *invoicing.Ctx) *Ctx {

	_ = *apiCtx
	_ = *dalCtx
	_ = *userCtx
	_ = *instrCtx
	_ = *trainCtx
	_ = *adyenCtx
	_ = *dcCtx
	_ = *subCtx

	ctx := new(Ctx)

	ctx.Config = new(Config)
	apiCtx.Config.UnmarshalKeyPanic("rsv", ctx.Config, ctx.Config.Validate)

	ctx.Api = apiCtx
	ctx.Dal = dalCtx
	ctx.User = userCtx
	ctx.Adyen = adyenCtx
	if adyenCtx.Mockup {
		ctx.NoPaymentFlow = true
	}
	ctx.Train = trainCtx
	ctx.Instr = instrCtx
	ctx.Review = reviewCtx
	ctx.Dc = dcCtx
	ctx.Sub = subCtx

	if adyenCtx.Config.NotifyEnabled {
		ctx.NoReplyCtx = noReplyCtx
	}

	ctx.RegisterAllHandlers()

	ctx.SetAdyenSm(invoicing)

	return ctx
}
