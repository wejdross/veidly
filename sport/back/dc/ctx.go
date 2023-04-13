package dc

import (
	"fmt"
	"sport/api"
	"sport/dal"
	"sport/instr"
)

type Config struct {
	MaxCodesPerInstr int `yaml:"max_codes_per_instr"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Validate Dc Config: " + m)
}

func (r *Config) Validate() error {

	if r.MaxCodesPerInstr <= 0 {
		return fmt.Errorf("Validate Config: invalid max_codes_per_instr")
	}

	return nil
}

type Ctx struct {
	Api    *api.Ctx
	Dal    *dal.Ctx
	Instr  *instr.Ctx
	Config *Config
}

func NewCtx(api *api.Ctx, dal *dal.Ctx, instrCtx *instr.Ctx) *Ctx {

	_ = *api
	_ = *dal
	_ = *instrCtx

	ctx := new(Ctx)
	ctx.Config = new(Config)
	api.Config.UnmarshalKeyPanic("dc", ctx.Config, ctx.Config.Validate)
	ctx.Api = api
	ctx.Dal = dal
	ctx.Instr = instrCtx

	ctx.Api.JwtGroup.POST("/dc", ctx.HandlerPostDc())
	ctx.Api.JwtGroup.GET("/dc", ctx.HandlerGetDc())
	ctx.Api.AnonGroup.GET("/dc/redeem", ctx.HandlerRedeemDc())
	ctx.Api.JwtGroup.PATCH("/dc", ctx.HandlerPatchDc())
	ctx.Api.JwtGroup.DELETE("/dc", ctx.HandlerDeleteDc())
	ctx.Api.JwtGroup.POST("/dc/binding", ctx.HandlerPostDcBinding())
	ctx.Api.JwtGroup.DELETE("/dc/binding", ctx.HandlerDeleteDcBinding())

	return ctx
}
