package notify

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"time"
)

type SmtpClientWrapper struct {
	Client *smtp.Client
	Expire time.Time
}

type ReusableSmtpCtx struct {
	Config SmtpConfig
	Pool   chan *SmtpClientWrapper
}

func smtpExpire() time.Time {
	return time.Now().Add(time.Minute)
}

func (c *SmtpClientWrapper) Close(force bool) error {
	if c.Client != nil {
		if err := c.Client.Quit(); err != nil && !force {
			return err
		}
		c.Client = nil
	}
	c.Expire = time.Time{}
	return nil
}

// this agent is responsible for removing connections
// if they are not used for specified amount of time
func (ctx *ReusableSmtpCtx) RunAgent() {
	for {
		// each few seconds pop element from pool and test it.
		el := <-ctx.Pool
		if !el.Expire.IsZero() && el.Expire.Before(time.Now()) {
			el.Close(true)
		}
		ctx.Pool <- el
		time.Sleep(time.Second * 10)
	}
}

// config cant be nil, maxConnections cant be <= 0
func NewReusableSmtpCtx(config *SmtpConfig, maxConnections int) *ReusableSmtpCtx {
	res := new(ReusableSmtpCtx)
	res.Config = *config
	// create buffer
	res.Pool = make(chan *SmtpClientWrapper, maxConnections)
	// fill up buffer with wrappers
	for i := 0; i < maxConnections; i++ {
		res.Pool <- &SmtpClientWrapper{
			Client: nil, // we dont create connections right away
		}
	}
	go res.RunAgent()
	return res
}

func (ctx *ReusableSmtpCtx) SendHtmlMail(dst mail.Address, subject string, html string) error {

	var cw *SmtpClientWrapper

	select {
	case cw = <-ctx.Pool:
		break
	case <-time.After(time.Second * 10):
		return fmt.Errorf("cant Send email: too many active connections at the moment")
	}

	defer func() {
		// before returning cw into the pool update its expire timer
		cw.Expire = smtpExpire()
		ctx.Pool <- cw
	}()

	var err error

	repeat := 0

	// this loop is used to detect broken connections
	// if smtp connection taken from the pool is in broken state
	// then MAIL will fail and we will have to recreate said connection
	// if retry fails too, then we will return error
	for {

		if cw.Client == nil {
			cw.Client, err = ctx.Config.CreateSmtpClient()
			if err != nil {
				return err
			}
		}

		if err = cw.Client.Mail(ctx.Config.Username); err != nil {
			cw.Close(true)
			if repeat >= 1 {
				return err
			} else {
				repeat++
				continue
			}
		}

		break
	}

	if err = cw.Client.Rcpt(dst.Address); err != nil {
		return err
	}

	w, err := cw.Client.Data()
	if err != nil {
		return err
	}

	d := ctx.Config.CreateHtmlData(dst, subject, html)

	_, err = w.Write(d)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
