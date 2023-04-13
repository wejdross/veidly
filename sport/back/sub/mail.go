package sub

import (
	"fmt"
	"net/mail"
	"os"
	"reflect"
	"runtime"
	"sport/adyen_sm"
	"sport/lang"
	"sport/user"

	"github.com/google/uuid"
)

const MailDateFmt = "02/01/2006 15:04"
const EmailBasePath = "../lang/email_templates/sub/"

func (ctx *Ctx) SendEmailToUser(
	ui *user.User,
	c string,
	subjHdr LocIx) error {

	if ui.ContactData.Email == "" {
		return nil
	}

	return ctx.noReplyCtx.SendHtmlMail(
		mail.Address{
			Name:    ui.Name,
			Address: ui.ContactData.Email,
		},
		Locale[ui.Language][subjHdr], c)
}

func (ctx *Ctx) GenerateAndSendEmailToUserID(
	userID uuid.UUID,
	getEmailData func(u *user.User) interface{},
	templateName string,
	emailTitleLocaleIx LocIx,
) error {
	u, err := ctx.User.DalReadUser(userID, user.KeyTypeID, true)
	if err != nil {
		return err
	}
	if u.ContactData.Email == "" {
		return nil
	}

	req := getEmailData(u)

	path := lang.CombineEmailPath(EmailBasePath, u.Language, templateName)

	html, err := ctx.User.LangCtx.ExecuteTemplate(path, &req)
	if err != nil {
		return err
	}

	return ctx.SendEmailToUser(u, html, emailTitleLocaleIx)
}

type CommonEmailData struct {
	Name    string
	SubName string
	SubUrl  string
}

func (ctx *Ctx) NewCommonEmailData(s *Sub, u *user.User) CommonEmailData {
	return CommonEmailData{
		Name:    u.Name,
		SubName: s.SubModel.Name,
		SubUrl:  ctx.SubUrl(s),
	}
}

func (ctx *Ctx) emailUserAboutLinkSync(s *Sub) error {
	return ctx.GenerateAndSendEmailToUserID(
		s.UserID,
		func(u *user.User) interface{} {
			return struct {
				CommonEmailData
				PaymentUrl string
			}{
				CommonEmailData: ctx.NewCommonEmailData(s, u),
				PaymentUrl:      s.LinkUrl,
			}
		},
		"sub_link",
		NewSub,
	)
}

type EmailHandler func(s *Sub) error

func (ctx *Ctx) SendEmailAsync(s *Sub, h EmailHandler) {
	if ctx.noReplyCtx == nil {
		return
	}
	go func() {
		err := h(s)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s: SendHtmlEmail: %s",
				runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(),
				err)
		}
	}()
}

func (ctx *Ctx) EmailUserAboutLinkAsync(p *adyen_sm.Pass) {
	s := SmPassToSub(p)
	ctx.SendEmailAsync(s, ctx.emailUserAboutLinkSync)
}

func (ctx *Ctx) emailUserAboutCaptureSync(s *Sub) error {
	return ctx.GenerateAndSendEmailToUserID(
		s.UserID,
		func(u *user.User) interface{} {
			return ctx.NewCommonEmailData(s, u)
		},
		"sub_capture",
		SubConfirm,
	)
}

func (ctx *Ctx) EmailUserAboutCaptureAsync(p *adyen_sm.Pass) {
	s := SmPassToSub(p)
	ctx.SendEmailAsync(s, ctx.emailUserAboutCaptureSync)
}

func (ctx *Ctx) emailUserAboutDisputeSync(s *Sub) error {
	return ctx.GenerateAndSendEmailToUserID(
		s.UserID,
		func(u *user.User) interface{} {
			return ctx.NewCommonEmailData(s, u)
		},
		"sub_dispute",
		SubDispute,
	)
}

func (ctx *Ctx) EmailUserAboutDisputeAsync(p *adyen_sm.Pass) {
	s := SmPassToSub(p)
	ctx.SendEmailAsync(s, ctx.emailUserAboutDisputeSync)
}

func (ctx *Ctx) emailInstrAboutDisputeSync(s *Sub) error {
	return ctx.GenerateAndSendEmailToUserID(
		s.InstrUserID,
		func(u *user.User) interface{} {
			return ctx.NewCommonEmailData(s, u)
		},
		"sub_dispute",
		SubDispute,
	)
}

func (ctx *Ctx) EmailInstrAboutDisputeAsync(p *adyen_sm.Pass) {
	s := SmPassToSub(p)
	ctx.SendEmailAsync(s, ctx.emailInstrAboutDisputeSync)
}

func (ctx *Ctx) emailUserAboutCancelSync(s *Sub) error {
	return ctx.GenerateAndSendEmailToUserID(
		s.UserID,
		func(u *user.User) interface{} {
			return ctx.NewCommonEmailData(s, u)
		},
		"sub_cancel",
		SubCancelled,
	)
}

func (ctx *Ctx) EmailUserAboutCancelAsync(p *adyen_sm.Pass) {
	s := SmPassToSub(p)
	ctx.SendEmailAsync(s, ctx.emailUserAboutCancelSync)
}

func (ctx *Ctx) emailInstrAboutCancelSync(s *Sub) error {
	return ctx.GenerateAndSendEmailToUserID(
		s.InstrUserID,
		func(u *user.User) interface{} {
			return ctx.NewCommonEmailData(s, u)
		},
		"sub_cancel",
		SubCancelled,
	)
}

func (ctx *Ctx) EmailInstrAboutCancelAsync(p *adyen_sm.Pass) {
	s := SmPassToSub(p)
	ctx.SendEmailAsync(s, ctx.emailInstrAboutCancelSync)
}

func (ctx *Ctx) emailUserAboutFailedCaptureSync(s *Sub) error {
	return ctx.GenerateAndSendEmailToUserID(
		s.UserID,
		func(u *user.User) interface{} {
			return struct {
				CommonEmailData
				SupportEmail string
			}{
				CommonEmailData: ctx.NewCommonEmailData(s, u),
				SupportEmail:    ctx.AdyenSm.Config.NotifyEmail,
			}
		},
		"sub_fail_capture",
		SubFailCapture,
	)
}

func (ctx *Ctx) EmailUserAboutFailedCaptureAsync(p *adyen_sm.Pass) {
	s := SmPassToSub(p)
	ctx.SendEmailAsync(s, ctx.emailUserAboutFailedCaptureSync)
}

//

func (ctx *Ctx) emailInstrAboutFailedPayoutSync(s *Sub) error {
	return ctx.GenerateAndSendEmailToUserID(
		s.InstrUserID,
		func(u *user.User) interface{} {
			return struct {
				CommonEmailData
				SupportEmail  string
				PaymentConfig string
			}{
				CommonEmailData: ctx.NewCommonEmailData(s, u),
				SupportEmail:    ctx.AdyenSm.Config.NotifyEmail,
				PaymentConfig:   ctx.Config.PayoutConfigUrl,
			}
		},
		"rsv_fail_payout",
		SubFailPayout,
	)
}

func (ctx *Ctx) EmailInstrAboutFailedPayoutAsync(p *adyen_sm.Pass) {
	s := SmPassToSub(p)
	ctx.SendEmailAsync(s, ctx.emailInstrAboutFailedPayoutSync)
}
