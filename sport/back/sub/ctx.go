package sub

import (
	"fmt"
	"sport/adyen"
	"sport/adyen_sm"
	"sport/api"
	"sport/dal"
	"sport/helpers"
	"sport/instr"
	"sport/invoicing"
	"sport/notify"
	"sport/user"
	//"sport/train"
)

type Config struct {
	LinkExpire           helpers.Duration `yaml:"link_expire_after"`
	WaitCaptureTimeout   helpers.Duration `yaml:"wait_capture_timeout"`
	InstantPayoutTimeout helpers.Duration `yaml:"instant_payout_timeout"`

	ServiceFee    int `yaml:"service_fee"`
	ProcessingFee int `yaml:"processing_fee"`
	RefundAmount  int `yaml:"refund_amount"`

	SubDetailsUrlFmt string           `yaml:"sub_details_url_fmt"`
	PayoutDelay      helpers.Duration `yaml:"payout_after"`
	PayoutConfigUrl  string           `yaml:"payout_config_url"`

	AdyenSmConfig adyen_sm.Config `yaml:"adyen_sm"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Subs: Validate Config: " + m)
}

func (ctx *Ctx) SubUrl(s *Sub) string {
	return fmt.Sprintf(ctx.Config.SubDetailsUrlFmt, s.ID)
}

func (r *Config) Validate() error {

	if r.LinkExpire <= 0 {
		return r.ValidationErr("Invalid link_expire_after")
	}
	if r.WaitCaptureTimeout <= 0 {
		return r.ValidationErr("Invalid wait_capture_timeout")
	}
	if r.InstantPayoutTimeout <= 0 {
		return r.ValidationErr("Invalid instant_payout_timeout")
	}

	if r.ServiceFee <= 0 {
		return r.ValidationErr("Invalid service_fee")
	}
	if r.ProcessingFee <= 0 {
		return r.ValidationErr("Invalid processing_fee")
	}

	if r.SubDetailsUrlFmt == "" {
		return r.ValidationErr("Invalid sub_details_url_fmt")
	}
	if r.PayoutDelay <= 0 {
		return r.ValidationErr("Invalid payout_config_url")
	}
	if r.PayoutConfigUrl == "" {
		return r.ValidationErr("Invalid payout_after")
	}
	if r.RefundAmount > 100 || r.RefundAmount < 0 {
		return r.ValidationErr("invalid refund_amount")
	}

	return r.AdyenSmConfig.Validate()
}

type Ctx struct {
	Api        *api.Ctx
	Dal        *dal.Ctx
	User       *user.Ctx
	Instr      *instr.Ctx
	Adyen      *adyen.Ctx
	Config     *Config
	noReplyCtx notify.EmailSender
	AdyenSm    *adyen_sm.Ctx
}

func (ctx *Ctx) RunAgent() {
	ctx.AdyenSm.RunAgent()
}

func NewCtx(
	api *api.Ctx,
	dal *dal.Ctx,
	userCtx *user.Ctx,
	instrCtx *instr.Ctx,
	adyenCtx *adyen.Ctx,
	noReplyCtx notify.EmailSender,
	invoicing *invoicing.Ctx,
) *Ctx {

	_ = *api
	_ = *dal
	_ = *instrCtx
	_ = *adyenCtx
	_ = *userCtx

	ctx := new(Ctx)

	ctx.Config = new(Config)
	api.Config.UnmarshalKey("sub", ctx.Config, ctx.Config.Validate)

	ctx.Api = api
	ctx.Dal = dal
	ctx.Instr = instrCtx
	ctx.Adyen = adyenCtx
	ctx.User = userCtx

	if adyenCtx.Config.NotifyEnabled {
		ctx.noReplyCtx = noReplyCtx
	}

	ctx.Api.JwtGroup.POST("/sub/model", ctx.HandlerPostSubModel())
	ctx.Api.JwtGroup.PATCH("/sub/model", ctx.HandlerPatchSubModel())
	ctx.Api.JwtGroup.DELETE("/sub/model", ctx.HandlerDeleteSubModel())
	ctx.Api.AnonGroup.GET("/sub/model", ctx.HandlerGetSubModel())

	ctx.Api.JwtGroup.POST("/sub/model/binding", ctx.HandlerPostSubModelBinding())
	ctx.Api.JwtGroup.DELETE("/sub/model/binding", ctx.HandlerDeleteSubModelBinding())

	ctx.Api.JwtGroup.POST("/sub", ctx.HandlerPostSub())
	ctx.Api.JwtGroup.GET("/sub", ctx.HandlerGetSubs())
	ctx.Api.JwtGroup.GET("/sub/:type", ctx.HandlerGetSubs())

	ctx.SetAdyenSm(invoicing)

	return ctx
}
