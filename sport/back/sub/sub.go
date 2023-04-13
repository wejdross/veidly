package sub

import (
	"fmt"
	"sport/adyen"
	"sport/adyen_sm"
	"time"

	"github.com/google/uuid"
)

type Sub struct {
	ID uuid.UUID

	SubModel     SubModel
	SubModelID   *uuid.UUID
	InstructorID uuid.UUID
	InstrUserID  uuid.UUID
	UserID       uuid.UUID
	RefID        uuid.UUID

	InstructorDecision adyen_sm.InstructorDecision
	State              adyen_sm.State

	SmCache   adyen_sm.SmCache
	SmRetries int
	SmTimeout time.Time

	IsConfirmed bool
	IsActive    bool
	OrderID     string

	CreatedOn time.Time

	//
	LinkID  string
	LinkUrl string

	DateStart time.Time
	DateEnd   time.Time

	RemainingEntries int
}

func SubColumns() string {
	return `s.id,

	s.sub_model,
	s.sub_model_id,
	s.instructor_id,
	s.instr_user_id,
	s.user_id,
	s.ref_id,

	s.instructor_decision,
	s.state,

	s.sm_cache,
	s.sm_retries,
	s.sm_timeout,

	s.is_confirmed,
	s.is_active,
	s.order_id,

	s.created_on,

	s.link_id,
	s.link_url,

	s.date_start,
	s.date_end,

	s.remaining_entries`
}

func (s *Sub) ScanFields() []interface{} {
	return []interface{}{
		&s.ID,

		&s.SubModel,
		&s.SubModelID,
		&s.InstructorID,
		&s.InstrUserID,
		&s.UserID,
		&s.RefID,

		&s.InstructorDecision,
		&s.State,

		&s.SmCache,
		&s.SmRetries,
		&s.SmTimeout,

		&s.IsConfirmed,
		&s.IsActive,
		&s.OrderID,

		&s.CreatedOn,

		&s.LinkID,
		&s.LinkUrl,

		&s.DateStart,
		&s.DateEnd,

		&s.RemainingEntries,
	}
}

// key is sub model id, val is assigned subs
type GroupedSm map[uuid.UUID][]SubWithJoins

func NewGroupedSm(s []SubWithJoins) GroupedSm {
	ret := make(GroupedSm)
	for i := 0; i < len(s); i++ {
		smID := s[i].SubModel.ID
		x, exists := ret[smID]
		if exists {
			ret[smID] = append(x, s[i])
		} else {
			ret[smID] = make([]SubWithJoins, 0, 2)
			ret[smID] = append(ret[smID], s[i])
		}
	}
	return ret
}

type SubRequest struct {
	ID                 uuid.UUID
	SubModel           *SubModel
	RefID              uuid.UUID
	AdyenRes           *adyen.CreatePaymentLinkResponse
	UserID             uuid.UUID
	DateStart, DateEnd time.Time
}

func (r *SubRequest) NewSub(ctx *Ctx) (*Sub, error) {
	if r.SubModel == nil {
		return nil, fmt.Errorf("NewSub: SubModel is nil")
	}
	if r.RefID == uuid.Nil {
		return nil, fmt.Errorf("NewSub: RefID is nil")
	}
	if r.ID == uuid.Nil {
		return nil, fmt.Errorf("NewSub: ID is nil")
	}
	if r.AdyenRes == nil {
		return nil, fmt.Errorf("NewSub: AdyenRes is nil")
	}
	if r.UserID == uuid.Nil {
		return nil, fmt.Errorf("NewSub: UserID is nil")
	}
	if r.DateStart.IsZero() || r.DateEnd.IsZero() || !r.DateEnd.After(r.DateStart) {
		return nil, fmt.Errorf("NewSub: invalid DateStart, DateEnd")
	}
	now := time.Now().In(time.UTC)

	var state adyen_sm.State
	var smto time.Time
	if ctx.Adyen.Mockup {
		state = adyen_sm.Hold
		smto = time.Date(5000, 12, 31, 24, 59, 59, 0, time.UTC)
	} else {
		state = adyen_sm.LinkExpress
		smto = now.Add(time.Duration(ctx.Config.LinkExpire))
	}

	s := &Sub{
		CreatedOn:          now,
		ID:                 r.ID,
		SubModel:           *r.SubModel,
		SubModelID:         &r.SubModel.ID,
		InstructorID:       r.SubModel.InstructorID,
		InstrUserID:        r.SubModel.InstrUserID,
		UserID:             r.UserID,
		RefID:              r.RefID,
		OrderID:            "",
		State:              state,
		SmTimeout:          smto,
		SmCache:            make(adyen_sm.SmCache),
		SmRetries:          0,
		InstructorDecision: adyen_sm.Unset,
		IsConfirmed:        false,
		IsActive:           true,
		LinkID:             r.AdyenRes.ID,
		LinkUrl:            r.AdyenRes.Url,
		DateStart:          r.DateStart,
		DateEnd:            r.DateEnd,
		RemainingEntries:   r.SubModel.MaxEntrances,
	}
	return s, nil
}

type SubValidStatus int

const (
	SubValid      SubValidStatus = 0
	Expired       SubValidStatus = 1
	NotYet        SubValidStatus = 2
	NotConfirmed  SubValidStatus = 4
	NoMoreEntries SubValidStatus = 8
)

func (sub *Sub) ValidStatus(t time.Time) SubValidStatus {

	var s SubValidStatus
	if !sub.IsConfirmed {
		s |= NotConfirmed
	}

	if t.Before(sub.DateStart) {
		s |= NotYet
	}

	if t.After(sub.DateEnd) {
		s |= Expired
	}

	if sub.RemainingEntries == 0 {
		s |= NoMoreEntries
	}

	return s
}
