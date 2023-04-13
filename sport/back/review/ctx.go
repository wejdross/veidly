package review

import (
	"fmt"
	"sport/api"
	"sport/dal"
	"sport/helpers"
	"sport/train"
	"sport/user"
)

const MaxMark = 6
const MaxReviewLen = 256

// config
type Config struct {
	ReviewExp helpers.Duration `yaml:"review_exp"`
}

func (r *Config) ValidationErr(m string) error {

	return fmt.Errorf("Validate Config: " + m)
}

func (r *Config) Validate() error {
	if r.ReviewExp <= 0 {
		return r.ValidationErr("invalid review_exp")
	}
	return nil
}

type Ctx struct {
	Api    *api.Ctx
	Dal    *dal.Ctx
	Train  *train.Ctx
	User   *user.Ctx
	Config *Config
}

func (ctx *Ctx) RegisterAllHandlers() {
	ctx.Api.AnonGroup.POST("/review", ctx.HandlerPostReview())
	ctx.Api.JwtGroup.DELETE("/review", ctx.HandlerDeleteReview())
	ctx.Api.JwtGroup.GET("/review/user", ctx.HandlerGetUserReview())
	ctx.Api.AnonGroup.GET("/review/pub", ctx.HandlerGetPubReviews())
}

func NewCtx(
	apiCtx *api.Ctx,
	dalCtx *dal.Ctx,
	trainCtx *train.Ctx,
	userCtx *user.Ctx) *Ctx {

	_ = *apiCtx
	_ = *dalCtx
	_ = *trainCtx
	_ = *userCtx

	ctx := new(Ctx)
	ctx.Config = new(Config)

	apiCtx.Config.UnmarshalKeyPanic("review", ctx.Config, ctx.Config.Validate)

	ctx.Api = apiCtx
	ctx.Dal = dalCtx
	ctx.Train = trainCtx
	ctx.User = userCtx

	ctx.RegisterAllHandlers()

	return ctx
}
