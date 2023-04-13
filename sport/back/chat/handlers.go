package chat

import (
	"fmt"
	"os"
	"sport/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateChatRoomRequest struct {
	ChatRoomRequest
	UserData MemberData
}

func (c *CreateChatRoomRequest) Validate(ctx *Ctx) error {

	const _fmt = "Validate CreateChatRoomRequest: %s"

	if err := c.ChatRoomRequest.Validate(); err != nil {
		return fmt.Errorf(_fmt, err.Error())
	}

	if err := c.UserData.Validate(ctx); err != nil {
		return fmt.Errorf(_fmt, err.Error())
	}

	return nil
}

type CreateChatRoomResponse struct {
	ChatRoomID uuid.UUID
}

func (ctx *Ctx) ValidateTokenHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		_, err := ctx.AuthorizeRequest(g)
		if err != nil {
			g.AbortWithError(401, err)
			return
		}
		g.AbortWithStatus(204)
	}
}

func (ctx *Ctx) CreateChatRoomHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		userID, err := ctx.AuthorizeRequest(g)
		if err != nil {
			g.AbortWithError(401, err)
			return
		}

		var apiReq CreateChatRoomRequest

		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &apiReq, func() error {
				return apiReq.Validate(ctx)
			},
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		chatroomID := uuid.New()

		member := ChatRoomMember{
			ChatRoomID: chatroomID,
			ServerID:   ctx.Config.ServerID,
			UserID:     userID,
			Data:       apiReq.UserData,
		}

		if err := ctx.DalCreateChatRoomMember(&member, All); err != nil {
			_ = ctx.DalDeleteChatRoom(chatroomID)
			g.AbortWithError(500, err)
			return
		}

		cr := ChatRoom{
			ChatRoomID:      chatroomID,
			ChatRoomRequest: apiReq.ChatRoomRequest,
		}

		if err := ctx.CreateChatRoomWithNotify(&cr, &member); err != nil {
			g.AbortWithError(500, err)
			return
		}

		res := CreateChatRoomResponse{
			ChatRoomID: chatroomID,
		}

		g.AbortWithStatusJSON(200, res)
	}
}

type CreateAccessTokenResponse struct {
	*AccessToken
	JoinLink string
}

func (ctx *Ctx) CreateAccessTokenHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		userID, err := ctx.AuthorizeRequest(g)
		if err != nil {
			g.AbortWithError(401, err)
			return
		}

		var apiReq AccessTokenRequest

		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &apiReq, apiReq.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		// TODO validate that user is admin
		members, err := ctx.DalReadChatRoomMembers(
			apiReq.ChatRoomID, userID)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}
		if len(members) != 1 {
			g.AbortWithError(401, fmt.Errorf("member not found"))
			return
		}

		at := AccessToken{
			CreatorID:          userID,
			TokenValue:         uuid.New(),
			AccessTokenRequest: apiReq,
		}

		if err := ctx.DalCreateAccessToken(&at); err != nil {
			g.AbortWithError(500, err)
			return
		}

		res := CreateAccessTokenResponse{
			AccessToken: &at,
			JoinLink: fmt.Sprintf(ctx.Config.UiJoinChatUrlFmt,
				at.TokenValue, at.ChatRoomID),
		}

		g.AbortWithStatusJSON(200, &res)
	}
}

func (ctx *Ctx) GetAccessTokensHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		userID, err := ctx.AuthorizeRequest(g)
		if err != nil {
			g.AbortWithError(401, err)
			return
		}

		chatRoomIDstr := g.Query("chatRoomID")
		if chatRoomIDstr == "" {
			g.AbortWithError(400, fmt.Errorf("invalid chatRoomID"))
			return
		}

		chatRoomID, err := uuid.Parse(chatRoomIDstr)
		if err != nil {
			g.AbortWithError(400, fmt.Errorf("invalid chatRoomID"))
			return
		}

		its, err := ctx.DalReadAccessTokens(chatRoomID, uuid.Nil)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		var filteredTokens = make([]AccessToken, 0, 4)
		for i := range its {
			if its[i].CreatorID == userID {
				filteredTokens = append(filteredTokens, its[i])
			}
		}

		g.AbortWithStatusJSON(200, filteredTokens)
	}
}

type JoinChatRoomRequest struct {
	ChatRoomID uuid.UUID
	UserData   MemberData
	Token      uuid.UUID
}

func (c *JoinChatRoomRequest) Validate(ctx *Ctx, skipUserDataValidation bool) error {

	const _fmt = "Validate JoinChatRoomRequest: %s"

	if c.ChatRoomID == uuid.Nil {
		return fmt.Errorf(_fmt, "invalid ChatRoomID")
	}

	if c.Token == uuid.Nil {
		return fmt.Errorf(_fmt, "invalid InviteToken")
	}

	if !skipUserDataValidation {
		if err := c.UserData.Validate(ctx); err != nil {
			return fmt.Errorf(_fmt, err.Error())
		}
	}

	return nil
}

type TokenResponse struct {
	Token string
}

func (ctx *Ctx) JoinChatRoomHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		var userID uuid.UUID
		var authorized bool
		var err error

		if ctx.HasAuthorization(g) {
			userID, err = ctx.AuthorizeRequest(g)
			if err != nil {
				g.AbortWithError(401, err)
				return
			}
			authorized = true
		} else {
			userID = uuid.New()
			authorized = false
		}

		var apiReq JoinChatRoomRequest

		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &apiReq, func() error {
				return apiReq.Validate(ctx, true)
			},
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		ats, err := ctx.DalReadAccessTokens(
			apiReq.ChatRoomID, apiReq.Token)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		if len(ats) != 1 {
			g.AbortWithError(404, fmt.Errorf("token not found"))
			return
		}

		at := &ats[0]

		if at.ExpiresOn.Before(time.Now()) {
			g.AbortWithError(410, fmt.Errorf("token expired"))
			_ = ctx.DalDeleteAccessToken(at.ChatRoomID, at.TokenValue)
			return
		}

		if at.UserID != uuid.Nil {
			if !authorized {
				chatToken, err := ctx.NewToken(at.UserID)
				if err != nil {
					g.AbortWithError(500, err)
					return
				}
				g.AbortWithStatusJSON(200, TokenResponse{
					Token: chatToken,
				})
				return
			}
			if at.UserID != userID {
				g.AbortWithError(401, fmt.Errorf("token already used up"))
				return
			}
			/*
			 * if user is logged in, they already have access to the chat room,
			 * so i'm returning 'Not Modified' response
			 */
			g.AbortWithStatus(304)
			return
		}

		if err := apiReq.UserData.Validate(ctx); err != nil {
			g.AbortWithStatus(202)
			return
		}

		member := ChatRoomMember{
			ChatRoomID: apiReq.ChatRoomID,
			ServerID:   ctx.Config.ServerID,
			UserID:     userID,
			Data:       apiReq.UserData,
		}

		if err := ctx.DalUpdateAccessTokenUser(
			at.ChatRoomID, at.TokenValue, userID); err != nil {
			g.AbortWithError(500, err)
			return
		}

		if err := ctx.DalCreateChatRoomMember(&member, All&^(LastNotified|LastReadMsg)); err != nil {
			g.AbortWithError(500, err)
			return
		}

		if authorized {
			g.AbortWithStatus(204)
		} else {
			chatToken, err := ctx.NewToken(userID)
			if err != nil {
				g.AbortWithError(500, err)
				return
			}
			g.AbortWithStatusJSON(200, TokenResponse{
				Token: chatToken,
			})
		}

	}
}

type OpenChatRoomRequest struct {
	ChatRoomID         uuid.UUID
	FeedOpt            MsgFeedOptions
	ForceRedirectPeers bool
}

func (c *OpenChatRoomRequest) Validate() error {

	const _fmt = "Validate OpenChatRoomRequest: %s"

	if c.ChatRoomID == uuid.Nil {
		return fmt.Errorf(_fmt, "Invalid ChatRoomID")
	}

	return c.FeedOpt.Validate()
}

/* return false if redirect failed and should proceed on current server */
func (ctx *Ctx) RedirectToServer(
	g *gin.Context, serverID string) {

	_url := fmt.Sprintf(ctx.Config.OpenUrlFmt, serverID)

	g.Redirect(307, _url)
}

func (ctx *Ctx) ObtainWsTokenHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		userID, err := ctx.AuthorizeRequest(g)
		if err != nil {
			g.AbortWithError(401, err)
			return
		}

		var apiReq OpenChatRoomRequest

		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &apiReq, apiReq.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}

		members, err := ctx.DalReadChatRoomMembers(apiReq.ChatRoomID, uuid.Nil)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		member := FindChatMember(members, userID)

		if member == nil {
			g.AbortWithError(401, fmt.Errorf("member not found"))
			return
		}

		room, err := ctx.DalReadChatRoom(apiReq.ChatRoomID)
		if err != nil {
			g.AbortWithError(403, err)
			return
		}

		if room.Flags&ForceRedirectEnabled == 0 {
			apiReq.ForceRedirectPeers = false
		}
		dstServerID, err := ctx.NegotiateServer(members, userID, apiReq.ForceRedirectPeers)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		if dstServerID != ctx.Config.ServerID {
			ctx.RedirectToServer(g, dstServerID)
			return
		}

		wst := ctx.wsTokenCache.NewToken(apiReq.FeedOpt, members, member)

		g.AbortWithStatusJSON(200, TokenResponse{
			Token: wst.Token.String(),
		})
	}
}

func (ctx *Ctx) OpenChatRoomHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		token, err := uuid.Parse(g.Query("t"))
		if err != nil {
			g.AbortWithError(403, err)
			return
		}

		wst := ctx.wsTokenCache.LookupAndRm(token)
		if wst == nil {
			g.AbortWithError(403, err)
			return
		}

		conn, err := wssUpgrader.Upgrade(g.Writer, g.Request, g.Writer.Header())
		if err != nil {
			fmt.Fprintf(os.Stderr, "upgrade failed with err: %v\n", err)
			return
		}

		ctx.HandleNewConn(&wst.FeedOpt, conn, wst.Member, wst.Members)

		g.Abort()
	}
}

func (ctx *Ctx) GetChatRoomsHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID, err := ctx.AuthorizeRequest(g)
		if err != nil {
			g.AbortWithError(401, err)
			return
		}

		if err != nil {
			g.AbortWithError(401, err)
			return
		}

		members, err := ctx.DalReadChatRoomMembers(uuid.Nil, userID)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		g.AbortWithStatusJSON(200, members)
	}
}

func (ctx *Ctx) OpenNotificationChan() gin.HandlerFunc {
	return func(g *gin.Context) {

		userID, err := ctx.AuthorizeRequest(g)
		if err != nil {
			g.AbortWithError(401, err)
			return
		}

		conn, err := wssUpgrader.Upgrade(g.Writer, g.Request, g.Writer.Header())
		if err != nil {
			fmt.Fprintf(os.Stderr, "upgrade failed with err: %v\n", err)
			return
		}

		ctx.HandleNewNotifyConn(userID, conn)
		g.Abort()
	}
}
