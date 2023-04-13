package chat_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sport/api"
	"sport/chat"
	"sport/helpers"
	"sport/user"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func createTestMsg(t *testing.T, userID, chatroomID uuid.UUID) {
	req := chat.ChatMsg{
		ChatRoomID: chatroomID,
		ChatMsgData: chat.ChatMsgData{
			Timestamp: chat.TimeToMsgTimestmap(time.Now()),
			UserID:    userID,
			Content:   helpers.CRNG_stringPanic(20),
		},
	}
	if err := chatCtx.CreateChatMsgNotify(&chat.ChatRoomMember{
		ChatRoomID: chatroomID,
		UserID:     userID,
	}, &req); err != nil {
		t.Fatal(err)
	}
}

func AssertCreateChatRoomResponse(
	t *testing.T,
	chatUserID uuid.UUID,
	req *chat.CreateChatRoomRequest,
	res *chat.CreateChatRoomResponse,
) {

	if res.ChatRoomID == uuid.Nil {
		t.Fatal("invalid res.ChatRoomID")
	}

	room, err := chatCtx.DalReadChatRoom(res.ChatRoomID)
	if err != nil {
		t.Fatalf("DalReadChatRoom: %v", err)
	}

	if room.Flags != req.Flags || room.ChatRoomID != res.ChatRoomID {
		t.Fatal("DalReadChatRoom: invalid response")
	}

	members, err := chatCtx.DalReadChatRoomMembers(res.ChatRoomID, chatUserID)
	if err != nil {
		t.Fatalf("DalReadChatRoomMembers: %v", err)
	}

	if len(members) != 1 {
		t.Fatal("DalReadChatRoomMembers: invalid response")
	}

	member := members[0]

	helpers.AssertJson(t, "DalReadChatRoomMembers", member.Data, req.UserData)

	if member.UserID != chatUserID ||
		member.ServerID != chatCtx.Config.ServerID ||
		!member.LastNotified.IsZero() ||
		member.ChatRoomID != res.ChatRoomID {

		t.Fatal("DalReadChatRoomMembers: invalid response")
	}
}

func newTestMemberData() chat.MemberData {
	return chat.MemberData{
		MemberPubData: chat.MemberPubData{
			DisplayName: helpers.CRNG_stringPanic(32),
		},
		Email:        helpers.CRNG_stringPanic(32) + "@test.test",
		ChatRoomName: helpers.CRNG_stringPanic(32),
		Language:     "en",
	}
}

func newTestCreateChatroomReqest() chat.CreateChatRoomRequest {
	return chat.CreateChatRoomRequest{
		ChatRoomRequest: chat.ChatRoomRequest{
			Flags: chat.ForceRedirectEnabled | chat.FreeJoin,
		},
		UserData: newTestMemberData(),
	}
}

func TestTokenCreate(t *testing.T) {

	var userID = uuid.New()
	token, err := chatCtx.NewToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	req := newTestCreateChatroomReqest()

	var res chat.CreateChatRoomResponse

	m := "POST"
	u := "/api/chat/room"

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
			RequestHeaders:     user.GetAuthHeader(helpers.CRNG_stringPanic(30)),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(req),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &res)
			},
		},
	})

	AssertCreateChatRoomResponse(t, userID, &req, &res)
}

func GinUnwrap(gh gin.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		g, _ := gin.CreateTestContext(w)
		g.Request = r
		gh(g)
		for i := range g.Errors {
			fmt.Println("GinUnwrap: err: " + g.Errors[i].Error())
		}
	}
}

func TestWs(t *testing.T) {
	var userID = uuid.New()
	token, err := chatCtx.NewToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	req := newTestCreateChatroomReqest()

	var res chat.CreateChatRoomResponse

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(req),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &res)
			},
		},
	})

	var wst chat.TokenResponse

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/token/ws",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(chat.OpenChatRoomRequest{
				ChatRoomID:         res.ChatRoomID,
				ForceRedirectPeers: false,
			}),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &wst)
			},
		},
	})

	path := chatCtx.GetOpenPath() + "?t=" + wst.Token

	serv := httptest.NewServer(GinUnwrap(chatCtx.OpenChatRoomHandler()))
	u := "ws" + strings.TrimPrefix(serv.URL, "http") + path
	defer serv.Close()
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	var payload chat.WsSTCPayload

	ws.SetReadDeadline(time.Now().Add(time.Second))
	_, wsBytes, err := ws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(wsBytes, &payload); err != nil {
		t.Fatal(err)
	}

	if payload.Type != chat.Members || len(payload.Members) != 1 ||
		!payload.Members[0].IsConnected {
		t.Fatal("invalid member report")
	}

	msgContent := "hello world"
	p := chat.WsCTSPayload{
		Type: chat.Msg,
		Msg:  msgContent,
	}

	if err := ws.WriteMessage(websocket.TextMessage,
		[]byte(helpers.JsonMustSerializeStr(p))); err != nil {
		t.Fatal(err)
	}

	ws.SetReadDeadline(time.Now().Add(time.Second))
	_, wsBytes, err = ws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(wsBytes, &payload); err != nil {
		t.Fatal(err)
	}

	if payload.Type != chat.Feed {
		t.Fatal("invalid payload type")
	}

	if len(payload.Msgs) != 1 {
		t.Fatal("invalid msgs")
	}

	msg := payload.Msgs[0]
	if msg.UserID != userID {
		t.Fatal("invalid msg author")
	}
	if msg.Content != msgContent {
		t.Fatal("invalid msg content")
	}

	members, err := chatCtx.DalReadChatRoomMembers(res.ChatRoomID, userID)
	if err != nil {
		t.Fatal(err)
	}

	if len(members) != 1 {
		t.Fatal("invalid amount of members in chatroom")
	}

	member := &members[0]

	if member.ChatRoomID != res.ChatRoomID ||
		!member.LastNotified.IsZero() ||
		member.UserID != userID ||
		member.ServerID != chatCtx.Config.ServerID {

		t.Fatal("invalid chat room member")
	}

	if err := helpers.AssertErr("assert member", member.Data, req.UserData); err != nil {
		t.Fatal(err)
	}

	ddlMsgs, err := chatCtx.DalReadChatMsgs(res.ChatRoomID, nil, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(ddlMsgs) != 1 {
		t.Fatal("invalid ddlMsgs")
	}

	ddlMsg := &ddlMsgs[0]

	if ddlMsg.UserID != userID ||
		ddlMsg.Content != msgContent ||
		ddlMsg.Timestamp != msg.Timestamp {

		t.Fatal("invalid ddlMsg")
	}

	rooms, err := chatCtx.DalReadChatRoomMembers(uuid.Nil, userID)
	if err != nil {
		t.Fatal(err)
	}

	if len(rooms) != 1 {
		t.Fatal("invalid rooms len")
	}

	room := rooms[0]
	helpers.AssertJson(t, "invalid room", &room, member)

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/chat/room?t=" + token,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(req),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []chat.ChatRoomMember
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}

				helpers.AssertJson(t, "[api] invalid room", res, rooms)
				return nil
			},
		},
	})
}

func TestJoinRoomWithExistingAccount(t *testing.T) {
	var userID = uuid.New()
	token, err := chatCtx.NewToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	req := newTestCreateChatroomReqest()

	var res chat.CreateChatRoomResponse

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(req),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &res)
			},
		},
	})

	var userID2 = uuid.New()
	token2, err := chatCtx.NewToken(userID2)
	if err != nil {
		t.Fatal(err)
	}

	createTokenReq := chat.AccessTokenRequest{
		ChatRoomID: res.ChatRoomID,
		ExpiresOn:  time.Now().Add(time.Hour * 3),
	}

	var inviteRes []chat.AccessToken
	var chMems []chat.ChatRoomMember

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room/access_token",
			ExpectedStatusCode: 401,
			RequestHeaders:     user.GetAuthHeader(token2),
			RequestReader:      helpers.JsonMustSerializeReader(createTokenReq),
		},
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room/access_token",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(createTokenReq),
		},
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/chat/room/access_token?chatRoomID=" + res.ChatRoomID.String(),
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token2),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				helpers.JsonMustDeserialize(b, &inviteRes)
				if len(inviteRes) != 0 {
					return fmt.Errorf("0: invalid number of invites")
				}
				return nil
			},
		},
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/chat/room/access_token?chatRoomID=" + res.ChatRoomID.String(),
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				helpers.JsonMustDeserialize(b, &inviteRes)
				if len(inviteRes) != 1 {
					return fmt.Errorf("1: invalid number of invites")
				}
				return nil
			},
		},
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/chat/room",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				helpers.JsonMustDeserialize(b, &chMems)
				if len(chMems) != 1 || chMems[0].UserID != userID {
					return fmt.Errorf("0: invalid chMems")
				}
				return nil
			},
		},
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/chat/room",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token2),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				helpers.JsonMustDeserialize(b, &chMems)
				if len(chMems) != 0 {
					return fmt.Errorf("1: invalid chMems")
				}
				return nil
			},
		},
	})

	jre := &chat.JoinChatRoomRequest{
		UserData:   newTestMemberData(),
		ChatRoomID: res.ChatRoomID,
		Token:      inviteRes[0].TokenValue,
	}

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room/join",
			ExpectedStatusCode: 204,
			RequestHeaders:     user.GetAuthHeader(token2),
			RequestReader:      helpers.JsonMustSerializeReader(jre),
		},
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room/join",
			ExpectedStatusCode: 304,
			RequestHeaders:     user.GetAuthHeader(token2),
			RequestReader:      helpers.JsonMustSerializeReader(jre),
		},
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/chat/room",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token2),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				helpers.JsonMustDeserialize(b, &chMems)
				if len(chMems) != 1 {
					return fmt.Errorf("3: invalid chMems")
				}
				return nil
			},
		},
	})

	mems, err := chatCtx.DalReadChatRoomMembers(res.ChatRoomID, uuid.Nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(mems) != 2 {
		t.Fatal("invalid amount of chat members")
	}
}

func TestJoinRoomNewUser(t *testing.T) {
	var userID = uuid.New()
	token, err := chatCtx.NewToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	req := newTestCreateChatroomReqest()

	var res chat.CreateChatRoomResponse

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(req),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &res)
			},
		},
	})

	createTokenReq := chat.AccessTokenRequest{
		ChatRoomID: res.ChatRoomID,
		ExpiresOn:  time.Now().Add(time.Hour * 3),
	}

	var inviteRes chat.CreateAccessTokenResponse

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room/access_token",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &inviteRes)
			},
			RequestReader: helpers.JsonMustSerializeReader(createTokenReq),
		},
	})

	jre := &chat.JoinChatRoomRequest{
		UserData:   newTestMemberData(),
		ChatRoomID: res.ChatRoomID,
		Token:      inviteRes.TokenValue,
	}

	mems, err := chatCtx.DalReadChatRoomMembers(res.ChatRoomID, uuid.Nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(mems) != 1 {
		t.Fatal("0: invalid amount of chat members")
	}

	var tokenRes chat.TokenResponse

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room/join",
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				helpers.JsonMustDeserialize(b, &tokenRes)
				return nil
			},
			RequestReader: helpers.JsonMustSerializeReader(jre),
		},
	})

	mems, err = chatCtx.DalReadChatRoomMembers(res.ChatRoomID, uuid.Nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(mems) != 2 {
		t.Fatal("1: invalid amount of chat members")
	}

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room/join",
			ExpectedStatusCode: 304,
			RequestHeaders:     user.GetAuthHeader(tokenRes.Token),
			RequestReader:      helpers.JsonMustSerializeReader(jre),
		},
	})
}

func TestNotifyChan(t *testing.T) {
	var userID = uuid.New()
	token, err := chatCtx.NewToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	var ccrr chat.CreateChatRoomResponse

	apiCtx.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/chat/room",
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(newTestCreateChatroomReqest()),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &ccrr)
			},
		},
	})

	createTestMsg(t, userID, ccrr.ChatRoomID)

	path := "/api/chat/notify/open?t=" + token

	serv := httptest.NewServer(GinUnwrap(chatCtx.OpenNotificationChan()))
	u := "ws" + strings.TrimPrefix(serv.URL, "http") + path
	defer serv.Close()
	nws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer nws.Close()

	var req = make(chat.NbfMap)
	req[ccrr.ChatRoomID] = 0 // chat.TimeToMsgTimestmap(time.Now().Add(-time.Second))

	if err := nws.WriteJSON(req); err != nil {
		t.Fatal(err)
	}

	var res chat.ChatroomNotifications

	_, msg, err := nws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	helpers.JsonMustDeserialize(msg, &res)

	if len(res) != 1 || res[ccrr.ChatRoomID].Count != 1 {
		t.Fatal("0: invalid notification")
	}

	// trigger notification
	if err := nws.WriteJSON(req); err != nil {
		t.Fatal(err)
	}
	_, msg, err = nws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	helpers.JsonMustDeserialize(msg, &res)
	if len(res) != 1 || res[ccrr.ChatRoomID].Count != 0 {
		t.Fatal("1: invalid notification")
	}

	createTestMsg(t, userID, ccrr.ChatRoomID)
	createTestMsg(t, userID, ccrr.ChatRoomID)

	if err := nws.WriteJSON(req); err != nil {
		t.Fatal(err)
	}
	_, msg, err = nws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	helpers.JsonMustDeserialize(msg, &res)
	if len(res) != 1 || res[ccrr.ChatRoomID].Count != 1 {
		t.Fatal("2: invalid notification")
	}
}
