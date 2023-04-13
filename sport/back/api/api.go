package api

import (
	"path"
	"sport/config"
	"sport/helpers"

	"github.com/gin-gonic/gin"
)

/*
	this will be added to all urls
*/
const ApiBaseHref = "/api"

type Ctx struct {
	Config *config.Ctx

	/*
		primary api context.
	*/
	engine *gin.Engine

	Request *ApiRequest

	EmptyGroup *gin.RouterGroup

	/*
		anon context. used to register public endpoints
	*/
	AnonGroup *gin.RouterGroup

	/*
		this is Jwt token authorization group
		meant to be used to obtain master token from user credentials via portal
		which then may be used to either access resources, or generate apikey to delegate permissions.
	*/
	Jwt      *Jwt
	JwtGroup *gin.RouterGroup
}

//var Global *Api

func NewApi(config *config.Ctx) *Ctx {
	ctx, err := NewApiRequestFromConfig(config)
	if err != nil {
		panic(err)
	}
	api, err := ctx.NewApi(config)
	if err != nil {
		panic(err)
	}
	return api
}

// func NewApiNoLog(path string) *Ctx {
// 	ctx, err := NewApiRequestFromConfig(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	ctx.Release = true
// 	api, err := ctx.NewApi()
// 	api.Config = config
// 	if err != nil {
// 		panic(err)
// 	}
// 	return api
// }

/*
	initialize global api context or panic
*/
// func NewGlobalFromConfig(path string, auth *auth.Auth) {
// 	api, err := NewFromConfig(path, auth)
// 	if err != nil {
// 		panic(err)
// 	}
// 	Global = api
// }

func (api *Ctx) delayMiddleware(g *gin.Context) {
	g.Next()
	if len(g.Errors) > 0 {
		helpers.CRNG_Delay(api.Request.Delay)
	}
}

/*
	start listening using global Api context
*/
func (api *Ctx) Run() error {
	if api.Request.WithTls {
		return api.RunTls()
	}
	return api.engine.Run(api.Request.Addr)
}

/*
	this will return path =
		{ApiBaseHref}/{rel}
*/
func ApiGetFullHref(rel string) string {
	return path.Join(ApiBaseHref, rel)
}
