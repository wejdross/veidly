package adyen

import (
	"encoding/json"
	"fmt"
	"os"
	"sport/notify"
)

func (ctx *Ctx) SendMailToSupport(hdr, msg string) error {
	return notify.SendEmailToSupport(
		ctx.NoReplyCtx,
		ctx.Config.NotifyEmail,
		ctx.Config.NotifyVer,
		hdr, msg,
	)
}

func (ctx *Ctx) MailSupportAboutWHfail(hdr string, err error, req interface{}) {

	var j []byte

	if req != nil {
		var _e error
		j, _e = json.MarshalIndent(req, "", "  ")
		if _e != nil {
			j = []byte(_e.Error())
		}
	} else {
		j = []byte("<not provided>")
	}

	content := fmt.Sprintf("%s\n\n--- REQ ---\n%s", err.Error(), string(j))

	if err = ctx.SendMailToSupport("WH: "+hdr, content); err != nil {
		fmt.Fprintf(os.Stderr, "EmailSupportAboutWHfail SendEmailToSupport: %v", err)
		return
	}
}
