package adyen_test

import (
	"sport/adyen"
	"sport/api"
	"sport/config"
	"sport/dal"
)

var adyenCtx *adyen.Ctx
var apiCtx *api.Ctx

func init() {
	apiCtx = api.NewApi(config.NewLocalCtx())
	dal.DeployDb(apiCtx.Config, false, "sportdb_test", false)
	adyenCtx = adyen.NewMockupCtx(apiCtx)
}
