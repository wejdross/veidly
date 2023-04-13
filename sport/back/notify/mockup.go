package notify

import "net/mail"

type MockupSmtpConfig struct {
}

func (ctx *MockupSmtpConfig) SendHtmlMail(dst mail.Address, subject string, html string) error {
	return nil
}
