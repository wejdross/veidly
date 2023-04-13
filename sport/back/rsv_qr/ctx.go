package rsv_qr

import (
	"fmt"
	"sport/api"
	"sport/dal"
	"sport/rsv"
)

type Config struct {
	QrEvalUrlFmt string `yaml:"qr_eval_url_fmt"`
	MaxQrCodes   int    `yaml:"max_qr_codes"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Validate CtxRequest: " + m)
}

func (r *Config) Validate() error {

	if r.QrEvalUrlFmt == "" {
		return fmt.Errorf("Validate Config: invalid qr_eval_url")
	}

	if r.MaxQrCodes <= 0 {
		return fmt.Errorf("Validate Config: invalid max_qr_codes")
	}

	return nil
}

type Ctx struct {
	Api    *api.Ctx
	Dal    *dal.Ctx
	Rsv    *rsv.Ctx
	Config *Config
}

func NewCtx(
	api *api.Ctx,
	dal *dal.Ctx,
	rsv *rsv.Ctx,
) *Ctx {

	_ = *api
	_ = *dal
	_ = *rsv

	ctx := new(Ctx)
	ctx.Config = new(Config)
	api.Config.UnmarshalKeyPanic("rsv_qr", ctx.Config, ctx.Config.Validate)
	ctx.Api = api
	ctx.Dal = dal
	ctx.Rsv = rsv

	ctx.Api.AnonGroup.POST("/qr/rsv", ctx.HandlersPostQr())
	ctx.Api.JwtGroup.GET("/qr/rsv/eval", ctx.HandlersEvalQr())

	return ctx
}
