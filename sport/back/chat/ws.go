package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// CTS stands for client to server
type WsCTSMsgType int

const (
	FeedReq WsCTSMsgType = 1
	Msg     WsCTSMsgType = 2
	ReadMsg WsCTSMsgType = 3
)

type WsCTSPayload struct {
	Type             WsCTSMsgType
	FeedOpts         *MsgFeedOptions
	Msg              string
	ReadMsgTimestamp int64
}

type WsSTCMsgType int

const (
	Feed    WsSTCMsgType = 1
	Members WsSTCMsgType = 2
)

type WsMember struct {
	MemberPubData
	UserID      uuid.UUID
	IsConnected bool
	You         bool
}

// STC stands for server to client
type WsSTCPayload struct {
	Type    WsSTCMsgType
	Msgs    []ChatMsgData `json:",omitempty"`
	Members []WsMember    `json:",omitempty"`
}

type WsConn struct {
	conn          *websocket.Conn
	isRunning     bool
	PendingMsgs   chan ChatMsgData
	LastTimestamp int64

	member *ChatRoomMember
}

func (ctx *Ctx) __sendMemberReport(
	chatRoomID uuid.UUID, allMembers []ChatRoomMember, lockTaken bool) {

	if len(allMembers) == 0 {
		return
	}

	members := make([]WsMember, len(allMembers))

	if !lockTaken {
		ctx.wsConnMapLock.RLock()
		defer ctx.wsConnMapLock.RUnlock()
	}

	wsConns := ctx.wsConnMap[chatRoomID]

	h := make(map[uuid.UUID]*WsConn)
	for i := range wsConns {
		h[wsConns[i].member.UserID] = wsConns[i]
	}

	for i := range allMembers {
		members[i] = WsMember{
			MemberPubData: allMembers[i].Data.MemberPubData,
			UserID:        allMembers[i].UserID,
			IsConnected:   false,
		}
		if c := h[allMembers[i].UserID]; c != nil {
			members[i].IsConnected = true
		}
	}

	for i := range members {
		if !members[i].IsConnected {
			continue
		}
		c := h[members[i].UserID]
		members[i].You = true
		c.conn.WriteJSON(WsSTCPayload{
			Type:    Members,
			Members: members,
		})
		members[i].You = false
	}
}

func (ctx *Ctx) SendMemberReport(chatRoomID uuid.UUID, lockTaken bool) {
	members, err := ctx.DalReadChatRoomMembers(chatRoomID, uuid.Nil)
	if err != nil {
		// TODO: proper error handling
		fmt.Println(err)
	}
	ctx.__sendMemberReport(chatRoomID, members, lockTaken)
}

func (ctx *Ctx) __closeWssConn(chatRoomID uuid.UUID, userID uuid.UUID, lockTaken bool) {
	if !lockTaken {
		ctx.wsConnMapLock.Lock()
		defer ctx.wsConnMapLock.Unlock()
	}
	ca := ctx.wsConnMap[chatRoomID]
	var c *WsConn
	for i := range ca {
		if ca[i].member.UserID != userID {
			continue
		}
		li := len(ca) - 1
		c = ca[i]
		ca[i] = ca[li]
		ctx.wsConnMap[chatRoomID] = ca[:li]
		break
	}

	// moved to runner...
	// if len(ctx.wsConnMap[chatRoomID]) == 0 {
	// 	delete(ctx.wsConnMap, chatRoomID)
	// }

	if c == nil {
		return
	}
	c.isRunning = false
	if err := c.conn.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "CloseWssConn: (warn) err during close: %v", err)
	}
}

func (ctx *Ctx) CloseWssConn(c *WsConn, lockTaken bool) {
	ctx.__closeWssConn(c.member.ChatRoomID, c.member.UserID, lockTaken)
	ctx.SendMemberReport(c.member.ChatRoomID, lockTaken)
}

func (ctx *Ctx) sendMsgToConn(msg ChatMsgData, c *WsConn) {

	select {
	case m := <-c.PendingMsgs:
		// todo: write pending msgs in bulk
		if err := c.conn.WriteJSON(WsSTCPayload{
			Type: Feed,
			Msgs: []ChatMsgData{m},
		}); err != nil {
			fmt.Fprintf(os.Stderr, "WriteJSON: failed to write msg: %v", err)
			c.PendingMsgs <- m
			c.PendingMsgs <- msg
			return
		}
	default:
		break
	}

	if err := c.conn.WriteJSON(WsSTCPayload{
		Type: Feed,
		Msgs: []ChatMsgData{msg},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "WriteJSON: failed to write msg: %v", err)
		c.PendingMsgs <- msg
	}
}

func TimeToMsgTimestmap(d time.Time) int64 {
	return d.UnixMicro()
}

func MsgTimestampToTime(t int64) time.Time {
	return time.UnixMicro(t)
}

func (ctx *Ctx) wssRead(c *WsConn) error {
	_, p, err := c.conn.ReadMessage()
	if err != nil {
		return err
	}

	var cts WsCTSPayload
	if err = json.Unmarshal(p, &cts); err != nil {
		return err
	}

	switch cts.Type {
	case ReadMsg:
		if cts.ReadMsgTimestamp > c.member.LastReadMsg {
			c.member.LastReadMsg = cts.ReadMsgTimestamp
			return ctx.DalCreateChatRoomMember(c.member, LastReadMsg)
		}
	case FeedReq:
		if cts.FeedOpts == nil {
			return fmt.Errorf("cannot request msg feed without providing options")
		}
		if err := cts.FeedOpts.Validate(); err != nil {
			return err
		}
		return ctx.SendMsgFeed(cts.FeedOpts, c.member.ChatRoomID, c.conn)
	case Msg:
		if cts.Msg == "" || len(cts.Msg) > 128 {
			return nil
		}
		ts := TimeToMsgTimestmap(time.Now())
		if ts == c.LastTimestamp {
			// throttling
			fmt.Println("throttling...")
			time.Sleep(time.Microsecond * 10)
			ts = TimeToMsgTimestmap(time.Now())
		}
		c.LastTimestamp = ts
		m := ChatMsg{
			ChatRoomID: c.member.ChatRoomID,
			ChatMsgData: ChatMsgData{
				Timestamp: ts,
				Content:   string(cts.Msg),
				UserID:    c.member.UserID,
			},
		}

		err := ctx.CreateChatMsgNotify(c.member, &m)

		if err != nil {
			return err
		}

		ctx.wsConnMapLock.RLock()
		defer ctx.wsConnMapLock.RUnlock()

		chatroomConnections := ctx.wsConnMap[c.member.ChatRoomID]
		for i := range chatroomConnections {
			ctx.sendMsgToConn(m.ChatMsgData, chatroomConnections[i])
		}
	}

	return nil
}

func (ctx *Ctx) RunConnectionLoop(c *WsConn) {

	for {
		if !c.isRunning {
			return
		}

		if err := ctx.wssRead(c); err != nil {
			fmt.Fprintf(os.Stderr, "wssRead: %v\n", err)
			ctx.CloseWssConn(c, false)
		}
	}
}

type WssConnMap map[uuid.UUID][]*WsConn

var wssUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type MsgFeedOptions struct {
	Limit int
	Start int64
	End   int64
}

const MaxMsgFeed = 100

func (fo *MsgFeedOptions) Validate() error {
	if fo.Limit < 0 {
		return fmt.Errorf("invalid Limit")
	}
	if fo.Limit > MaxMsgFeed {
		return fmt.Errorf("too many msgs requested in Limit")
	}
	if fo.Limit == 0 {
		fo.Limit = MaxMsgFeed
	}
	return nil
}

func (ctx *Ctx) SendMsgFeed(
	opt *MsgFeedOptions,
	chatroomID uuid.UUID,
	conn *websocket.Conn) error {

	msgs, err := ctx.DalReadChatMsgs(chatroomID, opt, false)
	if err != nil {
		return err
	}

	if len(msgs) > 0 {
		if err := conn.WriteJSON(WsSTCPayload{
			Type: Feed,
			Msgs: msgs,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *Ctx) HandleNewConnErr(
	feedOpt *MsgFeedOptions,
	conn *websocket.Conn,
	member *ChatRoomMember,
	allMembers []ChatRoomMember) error {

	wsConn := &WsConn{
		conn:        conn,
		isRunning:   true,
		PendingMsgs: make(chan ChatMsgData, 100),
		member:      member,
	}

	ctx.wsConnMapLock.Lock()

	if ctx.wsConnMap[member.ChatRoomID] == nil {
		ctx.wsConnMap[member.ChatRoomID] = make([]*WsConn, 0, 4)
	}
	ctx.wsConnMap[member.ChatRoomID] = append(ctx.wsConnMap[member.ChatRoomID], wsConn)
	ctx.__sendMemberReport(member.ChatRoomID, allMembers, true)

	ctx.wsConnMapLock.Unlock()

	if err := ctx.SendMsgFeed(feedOpt, member.ChatRoomID, conn); err != nil {
		return err
	}

	go ctx.RunConnectionLoop(wsConn)
	// go ctx.RunRetentionLoop(wsConn)

	return nil
}

func (ctx *Ctx) HandleNewConn(
	feedOpt *MsgFeedOptions,
	conn *websocket.Conn,
	member *ChatRoomMember,
	allMembers []ChatRoomMember) {

	if err := ctx.HandleNewConnErr(feedOpt, conn, member, allMembers); err != nil {
		fmt.Fprintf(os.Stderr, "HandleNewConnErr: %v\n", err)
		return
	}
}
