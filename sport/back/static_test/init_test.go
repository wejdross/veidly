package static

import (
	"sport/api"
	"sport/config"
	"sport/static"
)

var staticCtx *static.Ctx

func init() {
	apiCtx := api.NewApi(config.NewLocalCtx())
	staticCtx = static.NewCtx(apiCtx)
}
