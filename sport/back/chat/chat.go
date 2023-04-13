package chat

import (
	"fmt"
	"sport/api"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctx *Ctx) HasAuthorization(g *gin.Context) bool {
	_, err := api.GetTokenFromHdr(g)
	return err == nil || g.Query("t") != ""
}

func (ctx *Ctx) AuthorizeToken(token string) (uuid.UUID, error) {
	return ctx.jwt.AuthorizeUserFromToken(token)
}

func (ctx *Ctx) AuthorizeRequest(g *gin.Context) (uuid.UUID, error) {

	/* try to get token from header, if that fails move to query string */

	token, err := api.GetTokenFromHdr(g)
	if err != nil {
		token = g.Query("t")
		if token == "" {
			return uuid.Nil, fmt.Errorf("no auth specified")
		}
	}

	return ctx.AuthorizeToken(token)
}

// IMPORTANT
// this function is only allowed to return error if forceRedirect is set to true
func (ctx *Ctx) NegotiateServer(
	roomMembers []ChatRoomMember, userID uuid.UUID, forceRedirect bool) (string, error) {

	for i := range roomMembers {
		m := &roomMembers[i]

		if m.UserID == userID {
			continue
		}

		if m.ServerID != "" {
			if forceRedirect {
				m.ServerID = ctx.Config.ServerID
				if err := ctx.DalCreateChatRoomMember(m, ServerID); err != nil {
					return "", err
				}
			}
			return m.ServerID, nil
		}
	}

	return ctx.Config.ServerID, nil
}

func FindChatMember(members []ChatRoomMember, userID uuid.UUID) *ChatRoomMember {
	for i := range members {
		if members[i].UserID == userID {
			return &members[i]
		}
	}
	return nil
}

func (ctx *Ctx) NewToken(userID uuid.UUID) (string, error) {
	return ctx.jwt.GenToken(&api.JwtPayload{
		UserID: userID,
	})
}
