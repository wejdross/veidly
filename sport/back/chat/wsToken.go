package chat

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type WsToken struct {
	FeedOpt MsgFeedOptions

	Members []ChatRoomMember
	Member  *ChatRoomMember

	Token     uuid.UUID
	CreatedOn time.Time
}

type WsTokenCache struct {
	tm map[uuid.UUID]*WsToken
	mu sync.Mutex
}

func (b *WsTokenCache) LookupAndRm(token uuid.UUID) *WsToken {
	b.mu.Lock()
	defer b.mu.Unlock()
	ret := b.tm[token]
	delete(b.tm, token)
	return ret
}

func (b *WsTokenCache) NewToken(
	feedOpt MsgFeedOptions,
	members []ChatRoomMember,
	member *ChatRoomMember,
) *WsToken {
	b.mu.Lock()
	defer b.mu.Unlock()
	t := &WsToken{
		FeedOpt:   feedOpt,
		Token:     uuid.New(),
		CreatedOn: time.Now(),
		Members:   members,
		Member:    member,
	}
	b.tm[t.Token] = t
	return t
}

func NewWsTokenCache() *WsTokenCache {
	ret := new(WsTokenCache)
	ret.tm = make(map[uuid.UUID]*WsToken)
	go func() {
		for {
			ret.mu.Lock()
			now := time.Now()
			for k := range ret.tm {
				v := ret.tm[k]
				if v.CreatedOn.Add(time.Minute).Before(now) {
					delete(ret.tm, k)
				}
			}
			ret.mu.Unlock()
			time.Sleep(time.Second * 15)
		}
	}()
	return ret
}
