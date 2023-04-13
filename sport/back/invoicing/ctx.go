package invoicing

import (
	"fmt"
	"sport/api"
	"sport/dal"
	"sport/instr"
)

type Config struct {
	CompanyLines []string `yaml:"company_lines"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Invoicing: Validate Config: " + m)
}

func (r *Config) Validate() error {

	if len(r.CompanyLines) == 0 {
		return r.ValidationErr("invalid company_lines")
	}

	return nil
}

type Ctx struct {
	Config *Config
	Dal    *dal.Ctx
	Instr  *instr.Ctx
}

func NewCtx(
	api *api.Ctx,
	dalCtx *dal.Ctx,
	instrCtx *instr.Ctx,
) *Ctx {

	_ = *api
	_ = *dalCtx
	_ = *instrCtx

	ctx := new(Ctx)
	ctx.Config = new(Config)

	api.Config.UnmarshalKeyPanic("invoicing", ctx.Config, ctx.Config.Validate)

	ctx.Dal = dalCtx
	ctx.Instr = instrCtx

	api.JwtGroup.GET("/invoice", ctx.GetInvoiceHandler())
	api.AnonGroup.GET("/invoice/print", ctx.PrintInvoiceHandler())

	return ctx
}
