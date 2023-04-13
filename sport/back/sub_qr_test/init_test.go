package sub_qr_test

import (
	"sport/adyen"
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/instr"
	"sport/static"
	"sport/sub"
	"sport/sub_qr"
	"sport/user"
)

var qrCtx *sub_qr.Ctx
var subCtx *sub.Ctx
var uctx *user.Ctx

func init() {
	db := "sportdb_test"
	apiCtx := api.NewApi(config.NewLocalCtx())
	dal.DeployDb(apiCtx.Config, false, db, false)
	dalCtx := dal.NewDal(apiCtx.Config, db)
	staticCtx := static.NewCtx(apiCtx)
	uctx = user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, uctx, nil)
	adyenCtx := adyen.NewMockupCtx(apiCtx)
	subCtx = sub.NewCtx(apiCtx, dalCtx, uctx, instrCtx, adyenCtx, nil, nil)
	qrCtx = sub_qr.NewCtx(apiCtx, dalCtx, subCtx)
}
