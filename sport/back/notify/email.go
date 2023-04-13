package notify

import (
	"fmt"
	"net/mail"
)

type EmailSender interface {
	SendHtmlMail(to mail.Address, subj string, html string) error
}

func SendEmailToSupport(
	e EmailSender,
	dstEmailAddr, ver string,
	hdr, msg string) error {
	if e == nil {
		return nil
	}
	return e.SendHtmlMail(
		mail.Address{
			Name:    "",
			Address: dstEmailAddr,
		},
		fmt.Sprintf("API:%s: %s", ver, hdr),
		fmt.Sprintf("<pre>%s</pre>", msg),
	)
}
