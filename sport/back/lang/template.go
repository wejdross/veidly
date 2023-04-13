package lang

import (
	"html/template"
	"path"
	"strings"
)

const baseTemplatePath = "../lang/email_templates/base.html"

func (ctx *Ctx) ExecuteTemplate(path string, data interface{}) (string, error) {

	var err error

	t, exists := ctx.TemplateCache[path]
	if !exists {
		t, err = template.ParseFiles(baseTemplatePath, path)
		if err != nil {
			return "", err
		}
		ctx.TemplateCache[path] = t
	}

	var out strings.Builder

	err = t.ExecuteTemplate(&out, "base", data)

	res := strings.Replace(out.String(), "${LogoUrl}", ctx.Config.PubUrl, -1)

	return res, err
}

func CombineEmailPath(basePath, lang, filename string) string {
	return path.Join(basePath, lang+filename)
}
