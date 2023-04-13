package chat

import (
	"fmt"
	"net/mail"
	"path"
	"sport/lang"
	"time"

	"github.com/google/uuid"
)

const EmailBasePath = "../lang/email_templates/chat/"

type NewChanTemplate struct {
	Target   string
	ChatUrl  string
	ChatName string
	ValidTo  string
}

func (ctx *Ctx) SendEmailAboutNewChan(crm *ChatRoomMember) error {

	if ctx.langCtx == nil || ctx.noReply == nil {
		return nil
	}

	_lang := ctx.langCtx.ApiLangOrDefault(crm.Data.Language)
	at := &AccessToken{
		CreatorID:  crm.UserID,
		UserID:     crm.UserID,
		TokenValue: uuid.New(),
		AccessTokenRequest: AccessTokenRequest{
			ChatRoomID: crm.ChatRoomID,
			ExpiresOn:  time.Now().Add(time.Hour * 24 * 30),
		},
	}
	if err := ctx.DalCreateAccessToken(at); err != nil {
		return err
	}
	path := path.Join(EmailBasePath, _lang+".new_chatroom.html")
	html, err := ctx.langCtx.ExecuteTemplate(path, UnreadMsgsTemplate{
		Target:   crm.Data.DisplayName,
		ChatUrl:  fmt.Sprintf(ctx.Config.UiJoinChatUrlFmt, at.TokenValue, crm.ChatRoomID),
		ChatName: crm.Data.ChatRoomName,
		ValidTo:  at.AccessTokenRequest.ExpiresOn.Format(lang.MailDateFmt),
	})
	if err != nil {
		return err
	}

	return ctx.noReply.SendHtmlMail(
		mail.Address{
			Name:    crm.Data.DisplayName,
			Address: crm.Data.Email,
		},
		fmt.Sprintf(lang.Locale.NewChatroomFmt[_lang], crm.Data.ChatRoomName),
		html)
}

func (ctx *Ctx) SendEmailAboutNewChanAsync(crm *ChatRoomMember) {
	go func() {
		if err := ctx.SendEmailAboutNewChan(crm); err != nil {
			fmt.Println(err)
		}
	}()
}

type UnreadMsgsTemplate struct {
	Target        string
	ChatUrl       string
	ChatName      string
	ValidTo       string
	Notifications []*NotificationItem
}

func (ctx *Ctx) SendEmailAboutUnreadMsgs(crm *ChatRoomMember, nots []*NotificationItem) error {

	if ctx.langCtx == nil || ctx.noReply == nil {
		return nil
	}

	_lang := ctx.langCtx.ApiLangOrDefault(crm.Data.Language)
	path := path.Join(EmailBasePath, _lang+".unread_messages.html")

	at := &AccessToken{
		CreatorID:  crm.UserID,
		UserID:     crm.UserID,
		TokenValue: uuid.New(),
		AccessTokenRequest: AccessTokenRequest{
			ChatRoomID: crm.ChatRoomID,
			ExpiresOn:  time.Now().Add(time.Hour * 24 * 30),
		},
	}

	if err := ctx.DalCreateAccessToken(at); err != nil {
		return err
	}

	html, err := ctx.langCtx.ExecuteTemplate(path, UnreadMsgsTemplate{
		Notifications: nots,
		Target:        crm.Data.DisplayName,
		ChatUrl:       fmt.Sprintf(ctx.Config.UiJoinChatUrlFmt, at.TokenValue, crm.ChatRoomID),
		ChatName:      crm.Data.ChatRoomName,
		ValidTo:       at.AccessTokenRequest.ExpiresOn.Format(lang.MailDateFmt),
	})
	if err != nil {
		return err
	}

	return ctx.noReply.SendHtmlMail(
		mail.Address{
			Name:    crm.Data.DisplayName,
			Address: crm.Data.Email,
		},
		fmt.Sprintf(lang.Locale.UnreadMsgsFmt[_lang], crm.Data.ChatRoomName),
		html)
}

func (ctx *Ctx) SendEmailAboutUnreadMsgsAsync(crm *ChatRoomMember, nots []*NotificationItem) {
	go func() {
		if err := ctx.SendEmailAboutUnreadMsgs(crm, nots); err != nil {
			fmt.Println(err)
		}
	}()
}
