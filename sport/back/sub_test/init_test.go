package sub_test

import (
	"sport/adyen"
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/instr"
	"sport/static"
	"sport/sub"
	"sport/train"
	"sport/user"
)

var trainCtx *train.Ctx
var subCtx *sub.Ctx

func init() {
	apiCtx := api.NewApi(config.NewLocalCtx())
	dal.DeployDb(apiCtx.Config, false, "sportdb_test", false)
	dalCtx := dal.NewDal(apiCtx.Config, "sportdb_test")
	staticCtx := static.NewCtx(apiCtx)
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, nil)
	trainCtx = train.NewCtx(apiCtx, dalCtx, userCtx, staticCtx, instrCtx)
	adyenCtx := adyen.NewMockupCtx(apiCtx)
	subCtx = sub.NewCtx(apiCtx, dalCtx, userCtx, instrCtx, adyenCtx, nil, nil)
}
