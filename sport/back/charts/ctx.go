package charts

import (
	"sport/api"
	"sport/dal"
	"sport/instr"
)

type Ctx struct {
	Api   *api.Ctx
	Dal   *dal.Ctx
	Instr *instr.Ctx
}

func NewCtx(api *api.Ctx, dal *dal.Ctx, inst *instr.Ctx) *Ctx {

	_ = *api
	_ = *dal
	_ = *inst

	ctx := new(Ctx)
	ctx.Api = api
	ctx.Dal = dal
	ctx.Instr = inst

	ctx.Api.JwtGroup.POST("/charts", ctx.PostCharts())
	ctx.Api.JwtGroup.GET("/charts/user", ctx.GetChartsUserPerspective())

	return ctx
}
