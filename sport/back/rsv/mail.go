package rsv

import (
	"fmt"
	"net/mail"
	"os"
	"sport/lang"
)

const MailBasePath = "../lang/email_templates/rsv/"

func CombineEMailPath(language, name string) string {
	return lang.CombineEmailPath(MailBasePath, language, name)
}

func (ctx *Ctx) SendEmailToRsvUser(
	rsv *DDLRsv,
	c string,
	subjHdr LocIx) error {
	return ctx.NoReplyCtx.SendHtmlMail(
		mail.Address{
			Name:    rsv.UserInfo.Name,
			Address: rsv.UserContactData.Email,
		},
		Locale[rsv.UserInfo.Language][subjHdr], c)
}

func (ctx *Ctx) SendEmailToRsvInstr(
	rsv *DDLRsvWithInstr,
	c string,
	subjHdr LocIx,
) error {
	return ctx.NoReplyCtx.SendHtmlMail(
		mail.Address{
			Name:    rsv.Instructor.UserInfo.Name,
			Address: rsv.Instructor.ContactData.Email,
		},
		Locale[rsv.Instructor.UserInfo.Language][subjHdr],
		c)
}

func (ctx *Ctx) emailUserAboutLinkSync(
	rsv *DDLRsv,
) error {
	req := struct {
		Name       string
		PaymentUrl string
		RsvUrl     string
		Training   string
		RsvDate    string
	}{
		Name:       rsv.UserInfo.Name,
		PaymentUrl: rsv.LinkUrl,
		RsvUrl:     ctx.RsvTokenUrl(rsv),
		Training:   rsv.Training.Title,
		RsvDate:    rsv.DateStart.Format(lang.MailDateFmt),
	}

	path := lang.CombineEmailPath(MailBasePath, rsv.UserInfo.Language, ".rsv_link.html")
	html, err := ctx.User.LangCtx.ExecuteTemplate(path, req)
	if err != nil {
		return err
	}

	return ctx.SendEmailToRsvUser(rsv, html, NewRsv)
}

func (ctx *Ctx) EmailUserAboutLinkAsync(
	rsv *DDLRsv,
) {
	if rsv.UserContactData.Email == "" || ctx.NoReplyCtx == nil {
		return
	}
	go func() {
		err := ctx.emailUserAboutLinkSync(rsv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "EmailUserAboutCreatedRsv: SendHtmlEmail: %s", err)
		}
	}()
}

func (ctx *Ctx) emailUserAboutHoldSync(rsv *DDLRsv) error {
	req := struct {
		Name        string
		CaptureDate string
		RsvUrl      string
		Training    string
		RsvDate     string
	}{
		Name:        rsv.UserInfo.Name,
		CaptureDate: rsv.SmTimeout.Format(lang.MailDateFmt),
		RsvUrl:      ctx.RsvTokenUrl(rsv),
		Training:    rsv.Training.Title,
		RsvDate:     rsv.DateStart.Format(lang.MailDateFmt),
	}

	path := lang.CombineEmailPath(MailBasePath, rsv.UserInfo.Language, ".rsv_hold.html")
	html, err := ctx.User.LangCtx.ExecuteTemplate(path, req)
	if err != nil {
		return err
	}

	return ctx.SendEmailToRsvUser(rsv, html, RsvHold)
}

func (ctx *Ctx) emailUserAboutCaptureSync(rsv *DDLRsv) error {
	req := struct {
		Name         string
		Training     string
		RsvUrl       string
		TrainingDate string
	}{
		Name:         rsv.UserInfo.Name,
		Training:     rsv.Training.Title,
		RsvUrl:       ctx.RsvTokenUrl(rsv),
		TrainingDate: rsv.DateStart.Format(lang.MailDateFmt),
	}

	path := lang.CombineEmailPath(MailBasePath, rsv.UserInfo.Language, ".rsv_capture.html")
	html, err := ctx.User.LangCtx.ExecuteTemplate(path, req)
	if err != nil {
		return err
	}

	return ctx.SendEmailToRsvUser(rsv, html, RsvConfirm)
}

func (ctx *Ctx) emailUserAboutDisputeSync(rsv *DDLRsv) error {
	req := struct {
		Name     string
		Training string
		RsvUrl   string
	}{
		Name:     rsv.UserInfo.Name,
		Training: rsv.Training.Title,
		RsvUrl:   ctx.RsvTokenUrl(rsv),
	}

	path := lang.CombineEmailPath(MailBasePath, rsv.UserInfo.Language, ".rsv_dispute.html")
	html, err := ctx.User.LangCtx.ExecuteTemplate(path, req)
	if err != nil {
		return err
	}
	return ctx.SendEmailToRsvUser(rsv, html, RsvDispute)
}

func (ctx *Ctx) emailInstrAboutDisputeSync(rsv *DDLRsvWithInstr) error {
	req := struct {
		Name     string
		Training string
		RsvUrl   string
	}{
		Name:     rsv.Instructor.UserInfo.Name,
		Training: rsv.Training.Title,
		RsvUrl:   ctx.RsvTokenUrl(&rsv.DDLRsv),
	}

	path := lang.CombineEmailPath(
		MailBasePath, rsv.Instructor.UserInfo.Language, ".rsv_dispute.html")
	html, err := ctx.User.LangCtx.ExecuteTemplate(path, req)
	if err != nil {
		return err
	}
	return ctx.SendEmailToRsvInstr(rsv, html, RsvDispute)
}

func (ctx *Ctx) emailUserAboutCancelSync(rsv *DDLRsv) error {
	req := struct {
		Name      string
		DateStart string
		Training  string
		RsvUrl    string
	}{
		Name:      rsv.UserInfo.Name,
		Training:  rsv.Training.Title,
		RsvUrl:    ctx.RsvTokenUrl(rsv),
		DateStart: rsv.DateStart.Format(lang.MailDateFmt),
	}

	path := lang.CombineEmailPath(
		MailBasePath, rsv.UserInfo.Language, ".rsv_cancel.html")
	html, err := ctx.User.LangCtx.ExecuteTemplate(path, req)
	if err != nil {
		return err
	}

	return ctx.SendEmailToRsvUser(rsv, html, RsvCancelled)
}

func (ctx *Ctx) emailInstrAboutCancelSync(rsv *DDLRsvWithInstr) error {
	req := struct {
		Name      string
		DateStart string
		Training  string
		RsvUrl    string
	}{
		Name:      rsv.Instructor.UserInfo.Name,
		Training:  rsv.Training.Title,
		RsvUrl:    ctx.RsvTokenUrl(&rsv.DDLRsv),
		DateStart: rsv.DateStart.Format(lang.MailDateFmt),
	}

	path := lang.CombineEmailPath(
		MailBasePath, rsv.Instructor.UserInfo.Language, ".rsv_cancel.html")
	html, err := ctx.User.LangCtx.ExecuteTemplate(path, req)
	if err != nil {
		return err
	}

	return ctx.SendEmailToRsvInstr(rsv, html, RsvCancelled)
}

func (ctx *Ctx) emailUserAboutFailedCaptureSync(rsv *DDLRsv) error {
	req := struct {
		Name         string
		Training     string
		RsvUrl       string
		SupportEmail string
	}{
		Name:         rsv.UserInfo.Name,
		Training:     rsv.Training.Title,
		RsvUrl:       ctx.RsvTokenUrl(rsv),
		SupportEmail: ctx.AdyenSm.Config.NotifyEmail,
	}

	path := lang.CombineEmailPath(
		MailBasePath, rsv.UserInfo.Language, ".rsv_fail_capture.html")
	html, err := ctx.User.LangCtx.ExecuteTemplate(path, req)
	if err != nil {
		return err
	}

	return ctx.SendEmailToRsvUser(rsv, html, RsvFailCapture)
}

//

func (ctx *Ctx) emailInstrAboutFailedPayoutSync(rsv *DDLRsvWithInstr) error {
	req := struct {
		Name          string
		Training      string
		RsvUrl        string
		SupportEmail  string
		PaymentConfig string
	}{
		Name:          rsv.Instructor.UserInfo.Name,
		Training:      rsv.Training.Title,
		RsvUrl:        ctx.RsvTokenUrl(&rsv.DDLRsv),
		SupportEmail:  ctx.AdyenSm.Config.NotifyEmail,
		PaymentConfig: ctx.Config.PayoutConfigUrl,
	}

	path := lang.CombineEmailPath(
		MailBasePath, rsv.Instructor.UserInfo.Language, ".rsv_fail_payout.html")
	html, err := ctx.User.LangCtx.ExecuteTemplate(path, req)
	if err != nil {
		return err
	}

	return ctx.SendEmailToRsvInstr(rsv, html, RsvFailPayout)
}
