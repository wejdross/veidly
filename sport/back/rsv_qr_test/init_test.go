package rsv_qr_test

import (
	"sport/adyen"
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/dc"
	"sport/instr"
	"sport/rsv"
	"sport/rsv_qr"
	"sport/static"
	"sport/sub"
	"sport/train"
	"sport/user"
)

var qrCtx *rsv_qr.Ctx
var trainCtx *train.Ctx

func init() {
	db := "sportdb_test"
	apiCtx := api.NewApi(config.NewLocalCtx())
	dal.DeployDb(apiCtx.Config, false, db, false)
	dalCtx := dal.NewDal(apiCtx.Config, db)
	staticCtx := static.NewCtx(apiCtx)
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, nil)
	dcCtx := dc.NewCtx(apiCtx, dalCtx, instrCtx)
	trainCtx = train.NewCtx(apiCtx, dalCtx, userCtx, staticCtx, instrCtx)
	adyenCtx := adyen.NewMockupCtx(apiCtx)
	subCtx := sub.NewCtx(apiCtx, dalCtx, userCtx, instrCtx, adyenCtx, nil, nil)
	rsvCtx := rsv.NewCtx(
		apiCtx, dalCtx,
		userCtx, instrCtx,
		trainCtx, adyenCtx, nil,
		nil, dcCtx, subCtx, nil)
	qrCtx = rsv_qr.NewCtx(apiCtx, dalCtx, rsvCtx)
}
