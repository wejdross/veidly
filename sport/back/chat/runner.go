package chat

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"
)

func Latest(ts ...time.Time) time.Time {
	latest := time.Time{}
	for i := range ts {
		if ts[i].After(latest) {
			latest = ts[i]
		}
	}
	return latest
}

func (ctx *Ctx) NotifyChanMemberAboutUnreadMsgs(cm *ChatRoomMember) {

	cring := ctx.chanNotificationMap[cm.ChatRoomID]
	if cring == nil {
		return
	}

	uid := cm.UserID

	cring.Mu.RLock()

	canNotifytAfter := cring.LastNotified[uid].Add(time.Duration(ctx.Config.EmailNotAfter))

	lastNotified := Latest(
		cring.LastNotified[uid],
		// only messages sent within last 10 minutes will trigger notification
		time.Now().Add(-10*time.Minute),
	)

	// message must be present in db for at least <nbf>
	// for it to be sent
	nbf := time.Minute * 5

	nots := cring.DequeueAll(lastNotified, nbf, cm.UserID)

	cring.Mu.RUnlock()

	if len(nots) == 0 {
		return
	}

	if !canNotifytAfter.Before(time.Now()) {
		return
	}

	dbm, err := ctx.DalReadChatRoomMembers(cm.ChatRoomID, cm.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(dbm) != 1 {
		return
	}
	crm := dbm[0]

	// because user could be notified from other server
	canNotifytAfter = crm.LastNotified.Add(time.Duration(ctx.Config.EmailNotAfter))
	if !canNotifytAfter.Before(time.Now()) {
		return
	}

	lastNotified = Latest(
		MsgTimestampToTime(crm.LastReadMsg),
		cring.LastNotified[uid],
		crm.LastNotified,
		// only messages sent within last 15 minutes will trigger notification
		time.Now().Add(-15*time.Minute),
	)

	cring.Mu.RLock()

	nots = cring.DequeueAll(lastNotified, nbf, cm.UserID)

	cring.Mu.RUnlock()

	if len(nots) == 0 {
		return
	}

	now := time.Now()

	cring.Mu.Lock()
	cring.LastNotified[uid] = now
	cring.Mu.Unlock()

	ctx.SendEmailAboutUnreadMsgsAsync(&dbm[0], nots)
	cm.LastNotified = now
	if err := ctx.DalCreateChatRoomMember(cm, LastNotified); err != nil {
		fmt.Println(err)
	}
}

func (ctx *Ctx) RunnerIter() error {

	ctx.wsConnMapLock.Lock()

	var chatroomMap = make(map[uuid.UUID][]ChatRoomMember, len(ctx.wsConnMap))
	for chatRoomID := range ctx.wsConnMap {
		chatroomMap[chatRoomID] = nil
	}
	var chatroomMapLock sync.Mutex

	ctx.wsConnMapLock.Unlock()

	const maxThreads = 16

	context := context.Background()

	var sem = semaphore.NewWeighted(maxThreads)

	for crid := range chatroomMap {
		if err := sem.Acquire(context, 1); err != nil {
			return err
		}
		go func(crid uuid.UUID) {
			members, err := ctx.DalReadChatRoomMembers(crid, uuid.Nil)
			if err == nil {
				chatroomMapLock.Lock()
				chatroomMap[crid] = members
				chatroomMapLock.Unlock()
			} else {
				// log error but dont interrupt the flow
				fmt.Fprintln(os.Stderr, "in RunnerIter, DalReadChatRoomMembers failed with err: "+err.Error())
			}
			sem.Release(1)
		}(crid)
	}

	if err := sem.Acquire(context, maxThreads); err != nil {
		return err
	}

	ctx.wsConnMapLock.Lock()
	defer ctx.wsConnMapLock.Unlock()
	//

	for chatRoomID := range ctx.wsConnMap {

		allMembers := chatroomMap[chatRoomID]
		if allMembers == nil {
			// either failed to download members for some reason
			// or this is new connection - in that case member will be evaluated next time
			continue
		}

		userIndex := make(map[uuid.UUID]*ChatRoomMember)
		for i := range ctx.wsConnMap[chatRoomID] {
			c := ctx.wsConnMap[chatRoomID][i]
			userIndex[c.member.UserID] = c.member
		}

		for i := range allMembers {

			userID := allMembers[i].UserID

			if _, e := userIndex[userID]; e {

				// member is connected - check if he needs to be redirected
				serverID, err := ctx.NegotiateServer(allMembers, userID, false)
				if err != nil {
					// this error should never happen
					panic(err)
				}
				if serverID != ctx.Config.ServerID {
					for k := range ctx.wsConnMap[chatRoomID] {
						if ctx.wsConnMap[chatRoomID][k].member.UserID == userID {
							ctx.CloseWssConn(ctx.wsConnMap[chatRoomID][k], true)
						}
					}
				}
				continue
			}

			// member is not connected

			ctx.NotifyChanMemberAboutUnreadMsgs(&allMembers[i])
		}

		// for ow chatrooms are being persisted on servers
		// this could cause problems for huge amount of chatrooms
		// TODO: introduce some retention
		// if len(ctx.wsConnMap[chatRoomID]) == 0 {
		// 	delete(ctx.wsConnMap, chatRoomID)
		// }
	}

	return nil
}

func (ctx *Ctx) StartChatroomRunner() {
	for {
		if err := ctx.RunnerIter(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 4)
	}
}
