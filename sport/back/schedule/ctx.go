package schedule

import (
	"fmt"
	"sport/api"
	"sport/instr"
	"sport/rsv"
	"sport/sub"
	"sport/train"
)

// config
type Config struct {
	MaxSchedulePeriodInDays int `yaml:"max_schedule_period_in_days"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Validate Config: " + m)
}

func (r *Config) Validate() error {
	if r.MaxSchedulePeriodInDays <= 0 {
		return r.ValidationErr("invalid max_schedule_period_in_days")
	}
	return nil
}

type Ctx struct {
	Api    *api.Ctx
	Train  *train.Ctx
	Instr  *instr.Ctx
	Sub    *sub.Ctx
	Rsv    *rsv.Ctx
	Config *Config
}

func NewCtx(
	apiCtx *api.Ctx,
	instrCtx *instr.Ctx,
	trainCtx *train.Ctx,
	subCtx *sub.Ctx,
	rsvCtx *rsv.Ctx) *Ctx {

	_ = *apiCtx
	_ = *instrCtx
	_ = *trainCtx
	_ = *subCtx
	_ = *rsvCtx

	ctx := new(Ctx)

	ctx.Config = new(Config)
	apiCtx.Config.UnmarshalKeyPanic("schedule", ctx.Config, ctx.Config.Validate)

	ctx.Api = apiCtx
	ctx.Train = trainCtx
	ctx.Instr = instrCtx
	ctx.Sub = subCtx
	ctx.Rsv = rsvCtx

	ctx.Api.AnonGroup.GET("/schedule/rsv/t/:type", ctx.HandlerGetRsv())
	ctx.Api.AnonGroup.GET("/schedule", ctx.HandlerGetSched())

	return ctx
}
