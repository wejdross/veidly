package chat_integrator

import (
	"sport/api"
	"sport/chat"
	"sport/user"
)

// type Config struct {
// 	SupportEmail string
// }

type Ctx struct {
	Api  *api.Ctx
	Chat *chat.Ctx
	User *user.Ctx
	//Config *Config
}

func NewCtx(
	apiCtx *api.Ctx,
	chatCtx *chat.Ctx,
	userCtx *user.Ctx) *Ctx {

	res := new(Ctx)
	res.Api = apiCtx
	res.Chat = chatCtx
	res.User = userCtx

	apiCtx.AnonGroup.GET("/chat_integrator/token", res.GetChatTokenHandler())
	apiCtx.AnonGroup.POST("/chat_integrator/room", res.CreateUserChatroom())

	return res
}
