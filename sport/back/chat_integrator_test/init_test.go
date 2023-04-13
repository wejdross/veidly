package chat_test

import (
	"sport/api"
	"sport/chat"
	"sport/chat_integrator"
	"sport/config"
	"sport/dal"
	"sport/static"
	"sport/user"
)

var chatIntegratorCtx *chat_integrator.Ctx
var userCtx *user.Ctx

func init() {
	db := "sportdb_test"

	c := config.NewLocalCtx()
	apiCtx := api.NewApi(c)

	dal.DeployDb(c, false, db, false)
	dalCtx := dal.NewDal(c, db)
	staticCtx := static.NewCtx(apiCtx)
	userCtx = user.NewCtx(apiCtx, dalCtx, staticCtx, nil)

	chatCtx := chat.NewCtx(apiCtx, userCtx.LangCtx, nil, db, false)
	chatIntegratorCtx = chat_integrator.NewCtx(apiCtx, chatCtx, userCtx)
}
