package adyen_sm

import (
	"encoding/json"
	"fmt"
	"os"
	"sport/notify"
)

func (ctx *Ctx) SendMailToSupport(hdr, msg string) error {
	return notify.SendEmailToSupport(
		ctx.noReplyCtx,
		ctx.Config.NotifyEmail,
		ctx.Config.NotifyVer,
		hdr, msg,
	)
}

func (ctx *Ctx) MailSupportAboutWHfail(hdr string, err error) {

	if err = ctx.SendMailToSupport("WH: "+hdr, err.Error()); err != nil {
		fmt.Fprintf(os.Stderr, "EmailSupportAboutWHfail SendEmailToSupport: %v", err)
		return
	}
}

// this should be executed async cus emails take fuck load of time to send
func (ctx *Ctx) MailSupportAboutRsvStateChanged(o *StateChangeEvent) {
	j, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "EmailSupportAboutRsvStateChanged MarshalIndent: %v", err)
	}
	if err = ctx.SendMailToSupport("rsv_state_change", string(j)); err != nil {
		fmt.Fprintf(os.Stderr, "EmailSupportAboutRsvStateChanged SendHtmlEmail: %v", err)
		return
	}
}

func (ctx *Ctx) MailSupportAboutAgentErr(
	msg string,
) {
	if err := ctx.SendMailToSupport("agent_err", msg); err != nil {
		fmt.Fprintf(os.Stderr, "EmailSupportAboutAgentErr SendHtmlEmail: %v", err)
	}
}

func (ctx *Ctx) MailSupportAboutWHNotification(o *WhLog) {
	j, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "EmailSupportAboutWHNotification MarshalIndent: %v", err)
	}
	if err = ctx.SendMailToSupport("wh_notify", string(j)); err != nil {
		fmt.Fprintf(os.Stderr, "EmailSupportAboutWHNotification SendHtmlEmail: %v", err)
		return
	}
}
