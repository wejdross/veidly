package chat

import (
	"fmt"
	"sport/helpers"
	"time"

	"github.com/google/uuid"
)

type AccessTokenRequest struct {
	ChatRoomID uuid.UUID
	ExpiresOn  time.Time
}

func (itr *AccessTokenRequest) Validate() error {
	const _fmt = "Validate InviteTokenRequest: %s"
	if itr.ChatRoomID == uuid.Nil {
		return fmt.Errorf(_fmt, "invalid ChatRoomID")
	}
	if itr.ExpiresOn.Before(time.Now()) {
		return fmt.Errorf(_fmt, "invalid ExpiresOn")
	}
	return nil
}

type AccessToken struct {
	CreatorID  uuid.UUID
	UserID     uuid.UUID
	TokenValue uuid.UUID
	AccessTokenRequest
}

type MemberPubData struct {
	DisplayName string
	IconRelpath string
}

type MemberData struct {
	Email        string
	ChatRoomName string
	Language     string
	MemberPubData
}

func (r *MemberData) Validate(ctx *Ctx) error {
	if len(r.Email) > 0 {
		if err := helpers.ValidateEmail(r.Email); err != nil {
			return err
		}
	}

	if len(r.ChatRoomName) > 48 || r.ChatRoomName == "" {
		return fmt.Errorf("invalid ChatRoomName")
	}

	if len(r.DisplayName) > 48 || r.DisplayName == "" {
		return fmt.Errorf("invalid DisplayName")
	}

	if !ctx.langCtx.ValidateApiLang(r.Language) {
		r.Language = ctx.langCtx.Config.DefaultLang
	}

	return nil
}

type ChatRoomMember struct {
	ChatRoomID   uuid.UUID
	ServerID     string
	UserID       uuid.UUID
	LastNotified time.Time
	LastReadMsg  int64
	Data         MemberData
}

type ChatMsgData struct {
	Timestamp int64
	UserID    uuid.UUID
	Content   string
}

type ChatMsg struct {
	ChatRoomID uuid.UUID `json:"-"`
	ChatMsgData
}

type ChatRoomFlags int

const (
	FreeJoin             ChatRoomFlags = 1
	ForceRedirectEnabled ChatRoomFlags = 2
)

type ChatRoomRequest struct {
	Flags ChatRoomFlags
}

func (crr *ChatRoomRequest) Validate() error {
	if crr.Flags < 0 {
		return fmt.Errorf("invalid flags")
	}
	return nil
}

type ChatRoom struct {
	ChatRoomID       uuid.UUID
	LastMsgTimestmap int64
	ChatRoomRequest
}
