package adyen

import (
	"fmt"
	"sport/api"
	"sport/dal"
	"sport/notify"
)

const FreeInstructorShots = 2

// config
type Config struct {
	//AdyenWHPubUrl string `yaml:"adyen_wh_pub_url"`
	NotifyEnabled bool   `yaml:"notify_enabled"`
	NotifyEmail   string `yaml:"notify_email"`
	NotifyVer     string `yaml:"notify_ver"`
	Auth          struct {
		Username int
		Password string
		// api key to our system
		Apikey string `yaml:"apikey"`
	}
	// for example: "https://secure.przelewy24.pl"
	BaseUrl            string `yaml:"base_url"`
	Crc                string `yaml:"crc"`
	ApiPubUrl          string `yaml:"api_pub_url"`
	NotifyEndpointPath string `yaml:"notify_endpoint_path"`
	DonateReturnUrl    string `yaml:"donate_return_url"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Validate CtxRequest: " + m)
}

func (r *Config) Validate() error {
	if r.ApiPubUrl == "" {
		return r.ValidationErr("invalid api_pub_url")
	}
	if r.NotifyEndpointPath == "" {
		return r.ValidationErr("invalid notify_endpoint_path")
	}
	if r.DonateReturnUrl == "" {
		return r.ValidationErr("invalid donate_return_url")
	}
	if r.NotifyEnabled {
		if r.NotifyVer == "" {
			return r.ValidationErr("invalid notify_ver")
		}
		if r.NotifyEmail == "" {
			return r.ValidationErr("invalid notify_email")
		}
	}
	return nil
}

type HandlerFunc func(ni *NotificationRequestItem) error

type Ctx struct {
	Config     *Config
	Mockup     bool
	handlers   map[string]HandlerFunc
	NoReplyCtx notify.EmailSender
	Dal        *dal.Ctx
}

func (ctx *Ctx) whUrl() string {
	return fmt.Sprintf(
		"%s/api%s?apikey=%s",
		ctx.Config.ApiPubUrl,
		ctx.Config.NotifyEndpointPath,
		ctx.Config.Auth.Apikey,
	)
}

// isnt thread safe
func (ctx *Ctx) AddAdyenHandler(hdr string, f HandlerFunc) {
	if _, ok := ctx.handlers[hdr]; ok {
		panic(fmt.Sprintf("%s already implemented", hdr))
	}
	ctx.handlers[hdr] = f
}

func NewCtx(
	api *api.Ctx,
	noReplyCtx notify.EmailSender,
	dal *dal.Ctx,
) *Ctx {
	ctx := new(Ctx)
	ctx.Config = new(Config)
	api.Config.UnmarshalKeyPanic("p24", ctx.Config, ctx.Config.Validate)
	if ctx.Config.NotifyEnabled {
		ctx.NoReplyCtx = noReplyCtx
	}
	ctx.handlers = make(map[string]HandlerFunc)

	api.AnonGroup.POST(ctx.Config.NotifyEndpointPath, ctx.HandlerAdyenWH())
	api.AnonGroup.POST("/donate", ctx.DonateHandler())

	ctx.Dal = dal

	return ctx
}

func NewMockupCtx(api *api.Ctx) *Ctx {
	c := NewCtx(api, nil, nil)
	c.Mockup = true
	return c
}
