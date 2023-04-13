package user_test

import (
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/static"
	"sport/user"
)

var userCtx *user.Ctx

func init() {
	testdb := "sportdb_test"
	apiCtx := api.NewApi(config.NewLocalCtx())
	dal.DeployDb(apiCtx.Config, false, testdb, true)
	dalCtx := dal.NewDal(apiCtx.Config, testdb)
	staticCtx := static.NewCtx(apiCtx)
	userCtx = user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
}
