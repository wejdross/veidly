package chat

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ChatroomNotification struct {
	Count         int
	LastTimestamp int64
}

type ChatroomNotifications map[uuid.UUID]ChatroomNotification
type NbfMap map[uuid.UUID]int64

type NotificationCtx struct {
	nbf       NbfMap
	nbfLock   sync.Mutex
	chatCtx   *Ctx
	conn      *websocket.Conn
	userID    uuid.UUID
	rwNotChan chan struct{}
}

func (ctx *NotificationCtx) getAndWriteNotifications(force bool) error {

	ms, err := ctx.chatCtx.DalReadChatRoomMembers(uuid.Nil, ctx.userID)
	if err != nil {
		return err
	}

	if len(ms) == 0 {
		return nil
	}

	crids := make([]uuid.UUID, len(ms))
	cmap := make(map[uuid.UUID]*ChatRoomMember)
	for i := range ms {
		crids[i] = ms[i].ChatRoomID
		cmap[ms[i].ChatRoomID] = &ms[i]
	}

	rooms, err := ctx.chatCtx.DalReadChatRooms(crids)
	if err != nil {
		return err
	}

	var nots = make(ChatroomNotifications)

	for i := range rooms {
		crid := rooms[i].ChatRoomID
		lts := rooms[i].LastMsgTimestmap
		nbf, e := ctx.nbf[crid]
		cmem := cmap[crid]

		cn := ChatroomNotification{}

		// if new chatroom appears - force notification
		// i'm leaving it up to UI to determine if user should be notified about
		// 		new empty chatrooms or not
		if !e {
			force = true
		}

		if lts > nbf && (cmem == nil || lts > cmem.LastReadMsg) {
			cn.Count = 1
			cn.LastTimestamp = lts
		}

		nots[crid] = cn
		ctx.nbf[crid] = lts
	}

	if force {
		return ctx.conn.WriteJSON(nots)
	}

	for k := range nots {
		if nots[k].Count != 0 {
			return ctx.conn.WriteJSON(nots)
		}
	}

	return nil
}

func (ctx *NotificationCtx) mergeFeedOpts(uopts NbfMap) error {

	ctx.nbfLock.Lock()
	defer ctx.nbfLock.Unlock()

	for crid := range uopts {
		if uopts[crid] != 0 {
			ctx.nbf[crid] = uopts[crid]
		}
	}

	return nil
}

func (ctx *NotificationCtx) notifyRead() error {
	_, msg, err := ctx.conn.ReadMessage()
	if err != nil {
		return err
	}
	if len(msg) == 0 {
		return nil
	}
	var uopts NbfMap
	if err := json.Unmarshal(msg, &uopts); err != nil {
		return err
	}
	if err := ctx.mergeFeedOpts(uopts); err != nil {
		return err
	}
	ctx.rwNotChan <- struct{}{}
	return nil
}

func (ctx *NotificationCtx) notifyWrite() error {
	forceNotification := false
	select {
	case <-ctx.rwNotChan:
		forceNotification = true
		break
	case <-time.After(time.Second):
		if len(ctx.nbf) == 0 {
			return nil
		}
		break
	}
	if err := ctx.getAndWriteNotifications(forceNotification); err != nil {
		return err
	}
	return nil
}

func (ctx *Ctx) HandleNewNotifyConn(
	userID uuid.UUID,
	conn *websocket.Conn,
) {

	notCtx := &NotificationCtx{
		nbf:       make(NbfMap),
		nbfLock:   sync.Mutex{},
		chatCtx:   ctx,
		conn:      conn,
		userID:    userID,
		rwNotChan: make(chan struct{}),
	}

	go func() {
		for {
			if err := notCtx.notifyRead(); err != nil {
				fmt.Fprintf(os.Stderr, "notifyRead: %v\n", err)
				conn.Close()
				return
			}
		}
	}()

	go func() {

		for {
			if err := notCtx.notifyWrite(); err != nil {
				fmt.Fprintf(os.Stderr, "notifyWrite: %v\n", err)
				conn.Close()
				return
			}
		}

	}()
}
