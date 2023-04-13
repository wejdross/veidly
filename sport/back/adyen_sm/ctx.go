package adyen_sm

import (
	"fmt"
	"sport/adyen"
	"sport/api"
	"sport/dal"
	"sport/helpers"
	"sport/instr"
	"sport/invoicing"
	"sport/notify"
)

type Config struct {
	NotifyEmail   string          `yaml:"notify_email"`
	NotifyOn      map[string]bool `yaml:"notify_on"`
	NotifyEnabled bool            `yaml:"notify_enabled"`
	NotifyVer     string          `yaml:"notify_ver"`
	SmLogDir      string          `yaml:"sm_log_dir"`

	//

	SmMaxRetries   int              `yaml:"sm_max_retries"`
	SmRetryTimeout helpers.Duration `yaml:"sm_retry_timeout"`
	// max time allowed to process a payout
	InstantPayoutTimeout helpers.Duration `yaml:"instant_payout_timeout"`
	// max time allowed to process a capture
	CaptureTimeout helpers.Duration `yaml:"capture_timeout"`
	// max time allowed to process a refund
	RefundTimeout helpers.Duration `yaml:"refund_timeout"`

	//
	FreeInstrShots   int  `yaml:"free_instr_shots"`
	NoFreeInstrShots bool `yaml:"no_instr_shots"`
	//
	InstrShotPenaltyPercent    int  `yaml:"instr_shot_penalty_percent"`
	NoInstrShotPenalty         bool `yaml:"no_instr_shot_penalty"`
	InstrShotNoMoreThanPercent int  `yaml:"instr_shot_no_more_than_percent"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Sm: Validate Config: " + m)
}

func (r *Config) Validate() error {

	if r.FreeInstrShots < 0 {
		return r.ValidationErr("invalid free_instr_shots")
	}
	if !r.NoFreeInstrShots && r.FreeInstrShots == 0 {
		return r.ValidationErr("invalid free_instr_shots")
	}

	if r.InstrShotPenaltyPercent < 0 {
		return r.ValidationErr("invalid instr_shot_penalty_percent")
	}
	if !r.NoInstrShotPenalty && r.InstrShotPenaltyPercent == 0 {
		return r.ValidationErr("invalid instr_shot_penalty_percent")
	}

	if r.InstrShotNoMoreThanPercent <= 0 {
		return r.ValidationErr("invalid instr_shot_no_more_than_percent")
	}

	if r.NotifyEnabled {
		if r.NotifyVer == "" {
			return r.ValidationErr("invalid notify_ver")
		}
		if r.NotifyEmail == "" {
			return r.ValidationErr("invalid notify_email")
		}
	}
	if r.SmLogDir == "" {
		return r.ValidationErr("invalid sm_log_dir")
	}
	if r.InstantPayoutTimeout <= 0 {
		return r.ValidationErr("invalid instant_payout_timeout")
	}
	if r.CaptureTimeout <= 0 {
		return r.ValidationErr("invalid capture_timeout")
	}
	if r.RefundTimeout <= 0 {
		return r.ValidationErr("invalid refund_timeout")
	}
	return nil
}

type Ctx struct {
	Config *Config

	dal        *dal.Ctx
	noReplyCtx notify.EmailSender
	adyen      *adyen.Ctx
	instr      *instr.Ctx
	invoicing  *invoicing.Ctx

	SmCallbacks
	Model DDLModel
}

func NewCtx(
	api *api.Ctx, dalCtx *dal.Ctx,
	instrCtx *instr.Ctx,
	adyenCtx *adyen.Ctx,
	noReplyCtx notify.EmailSender,
	cb SmCallbacks,
	model DDLModel,
	conf *Config,
	invoicing *invoicing.Ctx,
) *Ctx {

	_ = *api
	_ = *dalCtx
	_ = *instrCtx
	_ = *adyenCtx
	_ = *conf

	if err := conf.Validate(); err != nil {
		panic(err)
	}

	ctx := new(Ctx)
	ctx.Config = conf
	ctx.dal = dalCtx
	ctx.instr = instrCtx
	ctx.SmCallbacks = cb
	ctx.adyen = adyenCtx
	ctx.Model = model
	ctx.invoicing = invoicing

	if ctx.Config.NotifyEnabled {
		ctx.noReplyCtx = noReplyCtx
	}

	return ctx
}
