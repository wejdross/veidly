package invoicing_test

import (
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/instr"
	"sport/invoicing"
	"sport/static"
	"sport/user"
)

var invoicingCtx *invoicing.Ctx

func init() {
	db := "sportdb_test"
	apiCtx := api.NewApi(config.NewLocalCtx())
	dal.DeployDb(apiCtx.Config, false, db, false)
	dalCtx := dal.NewDal(apiCtx.Config, db)
	staticCtx := static.NewCtx(apiCtx)
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, nil)
	invoicingCtx = invoicing.NewCtx(apiCtx, dalCtx, instrCtx)
}
