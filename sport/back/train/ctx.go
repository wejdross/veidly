package train

import (
	"fmt"
	"sport/api"
	"sport/dal"
	"sport/instr"
	"sport/static"
	"sport/user"
)

const TrainingMaxTitleLength = 128

// config
type Config struct {
	AllowedCurrencies    map[string]bool `yaml:"allowed_currencies"`
	MaxSecondaryTrImages int             `yaml:"max_tr_secondary_images"`
	MaxAge               int             `yaml:"max_age"`
	// Determines how long can requirement str be
	MaxRequirementStrLen int `yaml:"max_requirement_str_len"`
	// Determines how max amount of requirements
	MaxRequirementArrLen int `yaml:"max_requirement_arr_len"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Validate Config: " + m)
}

func (r *Config) Validate() error {
	if r.MaxSecondaryTrImages < 0 {
		return r.ValidationErr("max_tr_secondary_images")
	}
	if len(r.AllowedCurrencies) == 0 {
		return r.ValidationErr("allowed_currencies")
	}
	if r.MaxAge <= 0 {
		return r.ValidationErr("max_age")
	}
	if r.MaxRequirementStrLen <= 0 {
		return r.ValidationErr("max_requirement_str_len")
	}
	if r.MaxRequirementArrLen <= 0 {
		return r.ValidationErr("max_requirement_arr_len")
	}
	return nil
}

func (ctx *Ctx) ValidateCurrency(currency string) error {
	if currency == "" {
		return fmt.Errorf("ValidateCurrency: invalid currency")
	}
	if ok := ctx.Config.AllowedCurrencies[currency]; ok {
		return nil
	} else {
		return fmt.Errorf("ValidateCurrency: '%s' not suported", currency)
	}
}

type Ctx struct {
	Api    *api.Ctx
	Dal    *dal.Ctx
	User   *user.Ctx
	Instr  *instr.Ctx
	Static *static.Ctx
	Config *Config
}

func (ctx *Ctx) RegisterAllHandlers() {

	ctx.Api.JwtGroup.PATCH("/training", ctx.HandlerPatchTraining())
	ctx.Api.JwtGroup.POST("/training", ctx.HandlerPostTraining())
	ctx.Api.AnonGroup.GET("/training", ctx.HandlerGetTrainings())
	ctx.Api.JwtGroup.DELETE("/training", ctx.HandlerDeleteTraining())

	ctx.Api.JwtGroup.POST("/training/img", ctx.HandlerPostTrainingImage())
	ctx.Api.JwtGroup.DELETE("/training/img", ctx.HandlerDeleteTrainingImage())

	ctx.Api.JwtGroup.PUT("/training/occ", ctx.HandlerPutOccs())

	ctx.Api.JwtGroup.POST("/training/group", ctx.HandlerPostTrainingGroup())
	ctx.Api.JwtGroup.PATCH("/training/group", ctx.HandlerPatchTrainingGroup())
	ctx.Api.JwtGroup.DELETE("/training/group", ctx.HandlerDeleteTrainingGroup())
	ctx.Api.JwtGroup.GET("/training/group", ctx.HandlerGetTrainingGroups())

	ctx.Api.JwtGroup.PUT("/training/group/binding", ctx.HandlerPutGroupBinding())
	ctx.Api.JwtGroup.DELETE("/training/group/binding", ctx.HandlerDeleteGroupBinding())
}

const TrainingImgPath = "trainings"

func NewCtx(
	apiCtx *api.Ctx,
	dalCtx *dal.Ctx,
	userCtx *user.Ctx,
	staticCtx *static.Ctx,
	instrCtx *instr.Ctx,
) *Ctx {

	//
	_ = *apiCtx
	_ = *dalCtx
	_ = *userCtx
	_ = *staticCtx
	_ = *instrCtx

	ctx := new(Ctx)

	ctx.Config = new(Config)
	apiCtx.Config.UnmarshalKeyPanic("train", ctx.Config, ctx.Config.Validate)

	ctx.Api = apiCtx
	ctx.Dal = dalCtx
	ctx.User = userCtx
	ctx.Instr = instrCtx
	ctx.Static = staticCtx

	ctx.Static.RegisterDir(TrainingImgPath)
	ctx.RegisterAllHandlers()

	return ctx
}
