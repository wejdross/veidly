package chat

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type NotificationItem struct {
	Author       string
	MsgSummary   string
	Timestamp    time.Time
	TimestampStr string
	AuthorID     uuid.UUID
}

type NotificationRing struct {
	Slice        []*NotificationItem
	index        int
	LastNotified map[uuid.UUID]time.Time
	Mu           sync.RWMutex
}

func NewNotificationRing(capacity int) *NotificationRing {
	if capacity <= 0 {
		panic("creating NewNotificationRing with capacity <= 0")
	}
	return &NotificationRing{
		Slice:        make([]*NotificationItem, capacity),
		LastNotified: make(map[uuid.UUID]time.Time),
	}
}

func (n *NotificationRing) Enqueue(m *NotificationItem) {
	n.Slice[n.index] = m
	n.index = (n.index + 1) % len(n.Slice)
}

// validityPerdiod is used to discard all messages which are too old to be relevant
func (n *NotificationRing) DequeueAll(lastNotified time.Time, nbf time.Duration, userID uuid.UUID) []*NotificationItem {
	ret := make([]*NotificationItem, 0, len(n.Slice))
	now := time.Now()
	for i := range n.Slice {
		el := n.Slice[i]
		if el != nil &&
			el.Timestamp.After(lastNotified) &&
			el.Timestamp.Add(nbf).Before(now) &&
			el.AuthorID != userID {
			ret = append(ret, el)
		}
	}
	return ret
}
