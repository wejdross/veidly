package api

import (
	"fmt"

	"sport/config"

	"github.com/gin-gonic/gin"
)

// ApiRequest is used to create main Api context.
// it mostly stores configuration data (user input)
// it may parsed from config
type ApiRequest struct {
	Addr    string
	Release bool
	/* on error delay response by value */
	Delay     uint32
	BodyLimit struct {
		B int64
		K int64
		M int64
		G int64
	} `yaml:"body_limit"`
	Jwt JwtRequest
	//
	WithTls  bool   `yaml:"with_tls"`
	CertPath string `yaml:"cert_path"`
	KeyPath  string `yaml:"key_path"`
}

func (r *ApiRequest) Validate() error {
	if len(r.Addr) == 0 {
		return fmt.Errorf("Invalid address length")
	}
	if r.BodyLimit.B == 0 && r.BodyLimit.K == 0 && r.BodyLimit.M == 0 && r.BodyLimit.G == 0 {
		return fmt.Errorf("Invalid body limit")
	}
	if r.WithTls {
		if r.CertPath == "" || r.KeyPath == "" {
			return fmt.Errorf("invalid cert_path or key_path")
		}
	}
	return r.Jwt.Validate()
}

// Get new context from configuration file
func NewApiRequestFromConfig(config *config.Ctx) (*ApiRequest, error) {
	var wrapper struct {
		Api ApiRequest `yaml:"api"`
	}
	if err := config.Unmarshal(&wrapper); err != nil {
		return nil, err
	}
	if err := wrapper.Api.Validate(); err != nil {
		return nil, err
	}
	return &wrapper.Api, nil
}

func (r *ApiRequest) NewApi(config *config.Ctx) (*Ctx, error) {

	_ = *config

	var err error

	api := new(Ctx)
	api.Request = r
	api.Config = config

	if api.Jwt, err = r.Jwt.NewJwt(); err != nil {
		return nil, err
	}

	if api.Request.Release {
		gin.SetMode(gin.ReleaseMode)
		api.engine = gin.New()
	} else {
		api.engine = gin.Default()
	}

	// if Delay == 0 then dont register delay middleware.
	// you could move this check to inside of the middleware
	// so any change in configuration will work.
	if api.Request.Delay > 0 {
		api.engine.Use(api.delayMiddleware)
	}

	// fuck cors
	// so much
	api.engine.Use(CorsMiddleware())

	api.engine.Use(LimitReaderMiddleware(api))

	api.AnonGroup = api.engine.Group(ApiBaseHref)
	api.JwtGroup = api.engine.Group(ApiBaseHref)
	api.EmptyGroup = api.engine.Group("")

	api.JwtGroup.Use(api.jwtMiddleware)

	//ApiAddCorsHandlers()

	return api, nil
}
