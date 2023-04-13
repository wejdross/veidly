package chat_integrator

import (
	"context"
	"fmt"
	"sport/chat"
	"sport/helpers"
	"sport/user"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"
)

// return first non empty string
func coalesce(strs ...string) string {
	for i := range strs {
		if strs[i] != "" {
			return strs[i]
		}
	}
	return ""
}

func newChatroomMember(
	chatRoomID uuid.UUID, _user, peerUser *user.User,
) chat.ChatRoomMember {
	return chat.ChatRoomMember{
		ChatRoomID:   chatRoomID,
		UserID:       _user.ID,
		LastNotified: time.Time{},
		Data: chat.MemberData{
			Email:        _user.Email,
			ChatRoomName: coalesce(peerUser.Name, peerUser.Email),
			MemberPubData: chat.MemberPubData{
				DisplayName: coalesce(_user.Name, _user.Email),
				IconRelpath: _user.AvatarRelpath,
			},
		},
	}
}

type CreateUserChatroomRes struct {
	ChatRoomID  uuid.UUID
	AccessToken string
}

type AnonData struct {
	DisplayName string
	Email       string
}

func (ad *AnonData) Validate() error {
	const _fmt = "validate AnonData: %s"
	if ad.DisplayName == "" || len(ad.DisplayName) > 40 {
		return fmt.Errorf(_fmt, "invalid AnonData")
	}
	if err := helpers.ValidateEmail(ad.Email); err != nil {
		return fmt.Errorf(_fmt, "invalid email: "+err.Error())
	}
	return nil
}

type CreateUserChatroomReq struct {
	PeerUserID  uuid.UUID
	AnonData    *AnonData
	InitContent string
}

func (r *CreateUserChatroomReq) Validate() error {
	const _fmt = "validate CreateUserChatroomRequest: %s"
	if r.PeerUserID == uuid.Nil {
		return fmt.Errorf(_fmt, "invalid PeerUserID")
	}
	if len(r.InitContent) > 512 {
		return fmt.Errorf(_fmt, "invalid InitContent")
	}
	return nil
}

/*
 * find common chatroom between users and abort with its ID if found
 * first user ID must be caller
 */
func (ctx *Ctx) abortIfChatExists(g *gin.Context, isLoggedIn bool, userIDs ...uuid.UUID) bool {

	if len(userIDs) < 1 {
		g.AbortWithError(500, fmt.Errorf("abortIfChatExists: invalid userIDs"))
		return true
	}

	var userChatrooms = make(map[uuid.UUID]int, len(userIDs))
	var lock sync.Mutex
	const semCount = 8
	var sem = *semaphore.NewWeighted(semCount)
	var _ctx = context.Background()
	var lerr error

	for i := range userIDs {
		sem.Acquire(_ctx, 1)
		go func(i int) {
			defer sem.Release(1)
			members, err := ctx.Chat.DalReadChatRoomMembers(uuid.Nil, userIDs[i])
			if err != nil {
				lerr = err
				return
			}
			lock.Lock()
			defer lock.Unlock()
			for j := range members {
				userChatrooms[members[j].ChatRoomID]++
			}
		}(i)
	}

	sem.Acquire(_ctx, semCount)

	if lerr != nil {
		g.AbortWithError(500, lerr)
		return true
	}

	for k := range userChatrooms {
		if userChatrooms[k] == len(userIDs) {

			at := ""

			if !isLoggedIn {
				at, lerr = ctx.Chat.NewToken(userIDs[0])
				if lerr != nil {
					g.AbortWithError(500, lerr)
					return true
				}
			}

			g.AbortWithStatusJSON(200, CreateUserChatroomRes{
				ChatRoomID:  k,
				AccessToken: at,
			})
			return true
		}
	}

	return false
}

/*
 create [if not exists] chatroom between 2 veidly users,
 return chatroom_id
*/
func (ctx *Ctx) CreateUserChatroom() gin.HandlerFunc {
	return func(g *gin.Context) {

		var req CreateUserChatroomReq

		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}

		var reqUser *user.User
		var isLoggedIn bool

		if g.GetHeader("Authorization") == "" {
			if req.AnonData == nil {
				g.AbortWithError(400,
					fmt.Errorf("user must provide AnonData if no authorization is set"))
				return
			}
			if err := req.AnonData.Validate(); err != nil {
				g.AbortWithError(400, err)
				return
			}

			/*
			 * create dummy user which can be then used in newChatroomMember
			 *			 instead of real user
			 */
			reqUser = &user.User{
				ID: uuid.New(),
				PrivUserInfo: user.PrivUserInfo{
					Email: req.AnonData.Email,
					PubUserInfo: user.PubUserInfo{
						UserData: user.UserData{
							Name: req.AnonData.DisplayName,
						},
					},
				},
			}

		} else {

			reqUserID, err := ctx.Api.AuthorizeUserFromCtx(g)
			if err != nil {
				g.AbortWithError(401, err)
				return
			}

			reqUser, err = ctx.User.DalReadUser(reqUserID, user.KeyTypeID, true)
			if err != nil {
				g.AbortWithError(401, err)
				return
			}
			isLoggedIn = true
		}

		peerUser, err := ctx.User.DalReadUser(req.PeerUserID, user.KeyTypeID, true)
		if err != nil {
			g.AbortWithError(404, err)
			return
		}

		if ctx.abortIfChatExists(g, isLoggedIn, reqUser.ID, peerUser.ID) {
			return
		}

		cr := chat.ChatRoom{
			ChatRoomID: uuid.New(),
			ChatRoomRequest: chat.ChatRoomRequest{
				Flags: chat.ForceRedirectEnabled,
			},
		}
		reqMember := newChatroomMember(cr.ChatRoomID, reqUser, peerUser)
		peerMember := newChatroomMember(cr.ChatRoomID, peerUser, reqUser)

		if err := ctx.Chat.DalCreateChatRoomMember(&reqMember, chat.All); err != nil {
			g.AbortWithError(500, err)
			return
		}
		if err := ctx.Chat.DalCreateChatRoomMember(&peerMember, chat.All); err != nil {
			g.AbortWithError(500, err)
			return
		}

		if err := ctx.Chat.CreateChatRoomWithNotify(&cr, &reqMember); err != nil {
			g.AbortWithError(500, err)
			return
		}

		if len(req.InitContent) > 0 {
			err = ctx.Chat.CreateChatMsgNotify(&reqMember, &chat.ChatMsg{
				ChatRoomID: cr.ChatRoomID,
				ChatMsgData: chat.ChatMsgData{
					Timestamp: chat.TimeToMsgTimestmap(time.Now()),
					UserID:    reqMember.UserID,
					Content:   req.InitContent,
				},
			})
			if err != nil {
				fmt.Println("couldnt send init message, err was: ", err)
			}
		}

		res := CreateUserChatroomRes{
			ChatRoomID: cr.ChatRoomID,
		}

		if !isLoggedIn {
			res.AccessToken, err = ctx.Chat.NewToken(reqUser.ID)
			if err != nil {
				g.AbortWithError(500, err)
				return
			}
		}

		g.AbortWithStatusJSON(200, res)

	}
}

/*
 * generate chat access token from veidly token
 * if no veidly token provided, generate new anonymous access token
 */
func (ctx *Ctx) GetChatTokenHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		var userID uuid.UUID
		var err error

		if g.GetHeader("Authorization") != "" {
			userID, err = ctx.Api.AuthorizeUserFromCtx(g)
			if err != nil {
				g.AbortWithError(401, err)
				return
			}
		} else {
			userID = uuid.New()
		}
		chatToken, err := ctx.Chat.NewToken(userID)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}
		g.AbortWithStatusJSON(200, chat.TokenResponse{
			Token: chatToken,
		})
	}
}

// type AddSupportToChatroomReq struct {
// 	ChatRoomID uuid.UUID
// }

// func (a *AddSupportToChatroomReq) Validate() error {
// 	if a.ChatRoomID != uuid.Nil {
// 		return fmt.Errorf("Validate AddSupportToChatroomReq: invalid ChatRoomID")
// 	}
// 	return nil
// }

// func (ctx *Ctx) AddSupportToChatroomHandler() gin.HandlerFunc {
// 	return func(g *gin.Context) {
// 		userID, err := ctx.Chat.AuthorizeRequest(g)
// 		if err != nil {
// 			g.AbortWithError(401, err)
// 			return
// 		}

// 		var apiReq AddSupportToChatroomReq

// 		if err := helpers.ReadJsonBodyFromReader(
// 			g.Request.Body, &apiReq, apiReq.Validate,
// 		); err != nil {
// 			g.AbortWithError(400, err)
// 			return
// 		}

// 		crs, err := ctx.Chat.DalReadChatRoomMembers(apiReq.ChatRoomID, uuid.Nil)
// 		if err != nil {
// 			g.AbortWithError(500, err)
// 			return
// 		}

// 		var member *chat.ChatRoomMember

// 		for i := range crs {
// 			cr := &crs[i]
// 			/* support is already present in the chatroom */
// 			if cr.Data.Email == ctx.Config.SupportEmail {
// 				g.AbortWithStatus(304)
// 				return
// 			}
// 			if cr.UserID == userID {
// 				member = cr
// 			}
// 		}

// 		if member == nil {
// 			g.AbortWithError(404, fmt.Errorf("user not found"))
// 			return
// 		}

// 		supportMember := chat.ChatRoomMember{
// 			ChatRoomID:   apiReq.ChatRoomID,
// 			ServerID:     member.ServerID,
// 			UserID:       uuid.Nil,
// 			LastNotified: time.Time{},
// 			Data: chat.MemberData{
// 				Email:        ctx.Config.SupportEmail,
// 				ChatRoomName: member.Data.Email,
// 				MemberPubData: chat.MemberPubData{
// 					DisplayName: "veidly",
// 					IconRelpath: "",
// 				},
// 			},
// 		}

// 		if err := ctx.Chat.DalCreateChatRoomMember(&supportMember, chat.All); err != nil {
// 			g.AbortWithError(500, err)
// 			return
// 		}

// 		g.AbortWithStatus(204)
// 	}
// }
