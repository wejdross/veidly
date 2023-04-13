package notify

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"sport/config"
)

// SmtpConfig smtp protocol configuration
type SmtpConfig struct {
	Host string
	//Identity string
	Username string
	Password string
	Port     string
	From     string
}

func NewNoReplySmtpConfig(config *config.Ctx) *SmtpConfig {
	var ctx SmtpConfig
	config.UnmarshalKeyPanic("no_reply_smtp", &ctx, ctx.Validate)
	return &ctx
}

func (ctx *SmtpConfig) Validate() error {
	const hdr = "validate SMTPCtx: "
	if ctx.Host == "" {
		return fmt.Errorf("%shost was empty", hdr)
	}
	if ctx.Username == "" {
		return fmt.Errorf("%susername was empty", hdr)
	}
	if ctx.Password == "" {
		return fmt.Errorf("%spassword was empty", hdr)
	}
	if ctx.Port == "" {
		return fmt.Errorf("%sport was empty", hdr)
	}
	if ctx.From == "" {
		return fmt.Errorf("%sfrom was empty", hdr)
	}
	return nil
}

// will astablish TLS connection and then call AUTH and return client in this state
func (ctx *SmtpConfig) CreateSmtpClient() (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%s", ctx.Host, ctx.Port)
	auth := smtp.PlainAuth("", ctx.Username, ctx.Password, ctx.Host)

	tlsconfig := &tls.Config{
		ServerName: ctx.Host,
	}

	/*
		toCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		dialer := tls.Dialer{
			Config: tlsconfig,
		}
		conn, err := dialer.DialContext(toCtx, "tcp", addr)
		cancel()
		if err != nil {
			return nil, err
		}

		c, err := smtp.NewClient(conn, addr)
		if err != nil {
			return nil, err
		}

	*/

	c, err := smtp.Dial(addr)
	if err != nil {
		return nil, err
	}

	if err := c.StartTLS(tlsconfig); err != nil {
		return nil, err
	}

	if err = c.Auth(auth); err != nil {
		return nil, err
	}

	return c, nil
}

func (ctx *SmtpConfig) CreateHtmlData(dst mail.Address, subject string, html string) []byte {
	from := mail.Address{
		Name:    ctx.From,
		Address: ctx.From + "@veidly.com",
	}
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = dst.String()
	headers["Subject"] = subject
	headers["Mime-version"] = "1.0"
	headers["Content-Type"] = "text/html"
	headers["charset"] = "UTF-8"
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + html
	return []byte(message)
}

// SMTPSend send email to dst with given message
func (ctx *SmtpConfig) SendHtmlMail(dst mail.Address, subject string, html string) error {
	// var a smtp.Auth
	// const header = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// _subject := fmt.Sprintf("Subject: %s\n", subject)
	// msg := []byte(_subject + header + html)
	// a = smtp.PlainAuth( /*ctx.Identity*/ "", ctx.Username, ctx.Password, ctx.Host)
	// addr := fmt.Sprintf("%s:%s", ctx.Host, ctx.Port)
	// return smtp.SendMail(addr, a, ctx.From, dst, msg)

	c, err := ctx.CreateSmtpClient()
	if err != nil {
		return err
	}

	if err = c.Mail(ctx.Username); err != nil {
		return err
	}

	if err = c.Rcpt(dst.Address); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	d := ctx.CreateHtmlData(dst, subject, html)

	_, err = w.Write(d)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}

	return c.Close()
}

func (ctx *SmtpConfig) Test() error {
	return ctx.SendHtmlMail(
		mail.Address{Address: "sulucus@gmail.com"},
		"test",
		"<h1>testing</h1>")
}
