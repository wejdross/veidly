package instr

import (
	"sport/adyen"
	"sport/api"
	"sport/dal"
	"sport/lang"
	"sport/static"
	"sport/user"
)

const FreeInstructorShots = 2

type Ctx struct {
	Api    *api.Ctx
	Dal    *dal.Ctx
	User   *user.Ctx
	Adyen  *adyen.Ctx
	Lang   *lang.Ctx
	Static *static.Ctx
	//
	//Config *Config

	//MockAdyen bool
}

func (ctx *Ctx) RegisterAllHandlers() {

	ctx.Api.JwtGroup.POST("/instructor", ctx.HandlerPostInstructor())
	ctx.Api.AnonGroup.GET("/instructor", ctx.HandlerGetInstructor())
	ctx.Api.JwtGroup.PATCH("/instructor", ctx.HandlerPatchInstructor())
	ctx.Api.JwtGroup.DELETE("/instructor", ctx.HandlerDeleteInstructor())

	ctx.Api.JwtGroup.GET("/instructor/info", ctx.HandlerInstructorInfo())

	ctx.Api.JwtGroup.GET("/instructor/can_delete", ctx.HandlerCanDeleteInstructor())

	ctx.Api.JwtGroup.PATCH("/instructor/payout", ctx.HandlerPatchInstructorPayments())
	ctx.Api.JwtGroup.GET("/instructor/payout", ctx.HandlerGetInstructorPayments())
	ctx.Api.JwtGroup.DELETE("/instructor/payout", ctx.HandlerDeletePayoutInfo())

	ctx.Api.JwtGroup.POST("/instructor/vacation", ctx.HandlerPostVacation())
	ctx.Api.AnonGroup.GET("/instructor/vacation", ctx.HandlerGetVacations())
	ctx.Api.JwtGroup.PATCH("/instructor/vacation", ctx.HandlerPatchVacation())
	ctx.Api.JwtGroup.DELETE("/instructor/vacation", ctx.HandlerDeleteVacation())

	ctx.Api.JwtGroup.DELETE("/instructor/profile/img", ctx.HandlerDeleteProfileImg())
	ctx.Api.JwtGroup.POST("/instructor/profile/img", ctx.HandlerPostProfileImg())
}

const InstrProfileImgDir = "instr_profile"

func NewCtx(
	apiCtx *api.Ctx,
	dalCtx *dal.Ctx,
	userCtx *user.Ctx,
	adyenCtx *adyen.Ctx,
) *Ctx {

	_ = *apiCtx
	_ = *dalCtx
	_ = *userCtx
	// conf, err := NewConfig(configPath)
	// if err != nil {
	// 	panic(err)
	// }

	ctx := new(Ctx)
	//ctx.Config = conf

	ctx.Api = apiCtx
	ctx.Dal = dalCtx
	ctx.User = userCtx
	ctx.Adyen = adyenCtx
	ctx.Static = userCtx.Static
	ctx.Lang = userCtx.LangCtx

	ctx.RegisterAllHandlers()

	ctx.Static.RegisterDir(InstrProfileImgDir)

	return ctx
}
