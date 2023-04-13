package user

import (
	"path"
)

const basePath = "../lang/email_templates/user"

type RegisterData struct {
	Name string
	Url  string
}

func (ctx *Ctx) GetRegisterHtml(lang string, d *RegisterData) (string, error) {
	path := path.Join(basePath, lang+".register.html")
	return ctx.LangCtx.ExecuteTemplate(path, d)
}

type ForgotPassData struct {
	Name string
	Url  string
}

func (ctx *Ctx) GetForgotPassHtml(lang string, d *ForgotPassData) (string, error) {
	path := path.Join(basePath, lang+".forgot_pass.html")
	return ctx.LangCtx.ExecuteTemplate(path, d)
}
