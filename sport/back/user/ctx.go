package user

import (
	"fmt"
	"sport/api"
	"sport/dal"
	"sport/lang"
	"sport/notify"
	"sport/static"
	"strings"
)

const UserNameMaxLength = 128
const UserPassMaxLength = 128

type OauthProviderRequest struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	Scopes       []string `yaml:"scopes"`
	RedirectURL  string   `yaml:"redirect_url"`
}

func (oauth *OauthProviderRequest) Validate() error {
	const hdr = "Validate OauthConfig: "
	if oauth.ClientID == "" {
		return fmt.Errorf("%sclient_id was empty", hdr)
	}
	if oauth.ClientSecret == "" {
		return fmt.Errorf("%sclient_secret was empty", hdr)
	}
	if len(oauth.Scopes) < 1 {
		return fmt.Errorf("%sscopes was empty", hdr)
	}
	return nil
}

type Config struct {
	EnableRetentionLoop           bool `yaml:"enable_retention_loop"`
	RetentionLoopRecordTTLseconds int  `yaml:"retention_loop_record_ttl_seconds"`
	RetentionLoopIntervalSeconds  int  `yaml:"retention_loop_interval_seconds"`

	// public service url needed to route return messages
	// example format: your.host.com:8712
	PublicURL string `yaml:"public_url"`

	// used to redirect after successful register confirmation
	RegisterRedirectURL string `yaml:"register_redirect_url"`

	// used to redirect to UI after password reset confirm
	UIPassResetUrl string `yaml:"ui_pass_reset_url"`

	Use2Fa bool `yaml:"use_2fa"`

	UseOauth bool `yaml:"use_oauth"`

	Google OauthProviderRequest
}

func (conf *Config) ValidateCountry(c string) error {
	if c == "" {
		return nil
	}
	if strings.ToUpper(c) != c {
		return fmt.Errorf("ValidateCountry: must be all upper case")
	}
	if len(c) != 2 {
		return fmt.Errorf("ValidateCountry: must be 2 letters")
	}
	return nil
}

func (conf *Config) Validate() error {

	const hdr = "Validate CtxRequest: "

	if conf.EnableRetentionLoop {
		if conf.RetentionLoopIntervalSeconds <= 0 {
			return fmt.Errorf("%sretention_loop_interval_seconds <= 0", hdr)
		}
		if conf.RetentionLoopRecordTTLseconds <= 0 {
			return fmt.Errorf("%sretention_loop_record_ttl_seconds <= 0", hdr)
		}
	}

	if conf.Use2Fa {
		if conf.PublicURL == "" {
			return fmt.Errorf("%spublic_url was empty", hdr)
		}
		if conf.RegisterRedirectURL == "" {
			return fmt.Errorf("%slogin_url was empty", hdr)
		}
		if conf.UIPassResetUrl == "" {
			return fmt.Errorf("%sui_pass_reset_url was empty", hdr)
		}
	}

	if conf.UseOauth {
		if err := conf.Google.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Ctx struct {
	Api     *api.Ctx
	Dal     *dal.Ctx
	Static  *static.Ctx
	Config  *Config
	Email   notify.EmailSender
	Oauth   *OauthCtx
	LangCtx *lang.Ctx
	//
	//configPath  string
	IsMockup2Fa bool
}

func (ctx *Ctx) EnableMockup2FA() error {
	if ctx.IsMockup2Fa {
		return fmt.Errorf("2fa is not disabled")
	}
	ctx.Config.Use2Fa = true
	ctx.IsMockup2Fa = true
	ctx.Email = &notify.MockupSmtpConfig{}
	return nil
}

func (ctx *Ctx) DisableMockup2FA(e notify.EmailSender) error {
	if !ctx.IsMockup2Fa {
		return fmt.Errorf("2fa is not enabled")
	}
	ctx.IsMockup2Fa = false
	if e != nil {
		// if err := ctx.createSmtpConfig(ctx.Config.configPath); err != nil {
		// 	return err
		// }
		ctx.Email = e
	} else {
		ctx.Config.Use2Fa = false
	}
	return nil
}

//var global *Ctx

func (ctx *Ctx) RegisterAllHandlers() {
	// oauth
	ctx.Api.AnonGroup.GET("/user/oauth/url/:provider", ctx.OauthUrlHandler())
	ctx.Api.AnonGroup.POST("/user/oauth/code/:provider", ctx.OauthCodeHandler())

	ctx.Api.AnonGroup.POST("/user/register", ctx.RegisterHandler())
	// 2FA
	ctx.Api.AnonGroup.GET(ConfirmRegisterPath, ctx.ConfirmRegisterHandler())
	ctx.Api.AnonGroup.POST("/user/register/resend", ctx.ResendRegisterEmailHandler())

	ctx.Api.AnonGroup.POST("/user/login", ctx.LoginHandler())

	ctx.Api.JwtGroup.DELETE("/user", ctx.DeleteHandler())
	ctx.Api.JwtGroup.GET("/user", ctx.GetHandler())
	ctx.Api.JwtGroup.PATCH("/user/data", ctx.PatchUserDataHandler())
	ctx.Api.JwtGroup.PATCH("/user/contact", ctx.PatchUserContactDataHandler())
	ctx.Api.JwtGroup.PATCH("/user/password", ctx.PatchPasswordHandler())

	ctx.Api.JwtGroup.PUT("/user/avatar", ctx.PutUserAvatar())
	ctx.Api.JwtGroup.DELETE("/user/avatar", ctx.DeleteUserAvatar())

	ctx.Api.JwtGroup.GET("/user/stat", ctx.Stat())

	// its hard to test 2FA or following handlers
	// maybe should introduce testing with imap protocol to ensure 2FA and password reset works
	ctx.Api.AnonGroup.GET("/user/password/forgot", ctx.ForgotPasswordHandler())
	ctx.Api.AnonGroup.POST("/user/password/reset", ctx.ResetPasswordHandler())
	ctx.Api.AnonGroup.GET("/user/password/validate", ctx.ValidatePassword())
}

func NewCtx(
	a *api.Ctx,
	d *dal.Ctx,
	staticCtx *static.Ctx,
	noReply notify.EmailSender) *Ctx {

	_ = *a
	_ = *d
	_ = *staticCtx

	c := new(Ctx)
	c.Config = new(Config)
	if err := a.Config.UnmarshalKey("user", c.Config, c.Config.Validate); err != nil {
		panic(err)
	}

	c.LangCtx = lang.NewCtx(a)

	c.Api = a
	c.Dal = d
	c.Static = staticCtx

	if c.Config.Use2Fa && noReply != nil {
		c.Email = noReply
		// if err := c.createSmtpConfig(conf.configPath); err != nil {
		// 	panic(err)
		// }
	} else {
		c.Config.Use2Fa = false
	}

	c.Static.RegisterDir("user")

	if c.Config.UseOauth {
		c.Oauth = c.Config.NewOauthCtx()
	}

	c.RegisterAllHandlers()

	return c
}
