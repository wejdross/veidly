package chat_test

import (
	"sport/api"
	"sport/chat"
	"sport/config"
	"sport/lang"
)

var chatCtx *chat.Ctx
var apiCtx *api.Ctx

func init() {
	apiCtx = api.NewApi(config.NewLocalCtx())
	chatCtx = chat.NewCtx(apiCtx, lang.NewCtx(apiCtx), nil, "sportdb_test", false)
}
