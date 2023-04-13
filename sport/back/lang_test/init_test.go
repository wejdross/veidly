package lang_test

import (
	"sport/api"
	"sport/config"
	"sport/lang"
)

var langCtx *lang.Ctx

func init() {
	apiCtx := api.NewApi(config.NewLocalCtx())
	langCtx = lang.NewCtx(apiCtx)
}
