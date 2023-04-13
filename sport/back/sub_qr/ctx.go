package sub_qr

import (
	"fmt"
	"sport/api"
	"sport/dal"
	"sport/sub"
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
	Sub    *sub.Ctx
	Config *Config
}

func NewCtx(
	api *api.Ctx,
	dal *dal.Ctx,
	sub *sub.Ctx,
) *Ctx {

	_ = *api
	_ = *dal
	_ = *sub

	ctx := new(Ctx)

	ctx.Config = new(Config)
	api.Config.UnmarshalKey("sub_qr", ctx.Config, ctx.Config.Validate)

	ctx.Api = api
	ctx.Dal = dal
	ctx.Sub = sub

	ctx.Api.JwtGroup.POST("/qr/sub", ctx.CreateQrHandler())
	ctx.Api.JwtGroup.GET("/qr/sub/eval", ctx.EvalQrHandler())
	ctx.Api.JwtGroup.GET("/qr/sub/confirm", ctx.ConfirmQrHandler())

	return ctx
}
