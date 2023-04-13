package lang

import (
	"fmt"
	"html/template"
	"sport/api"
)

type Config struct {
	// languages which are supported api langs, used for example in emails
	ApiLang         map[string]bool `yaml:"api_lang"`
	DefaultLang     string          `yaml:"default_lang"`
	TagPath         string          `yaml:"tag_path"`
	TagCategoryPath string          `yaml:"tag_category_path"`
	LangPath        string          `yaml:"lang_path"`
	PubUrl          string          `yaml:"pub_url"`
}

type Ctx struct {
	Config        *Config
	userLangIndex LangIndex
	Tag           *TagCtx
	Api           *api.Ctx
	TemplateCache map[string]*template.Template
}

func (conf *Config) Validate() error {
	const hdr = "Validate Config: "

	if conf.DefaultLang == "" {
		return fmt.Errorf("%sdefault_lang is empty", hdr)
	}

	if !conf.ApiLang[conf.DefaultLang] {
		return fmt.Errorf("%sdefault_lang not present in supported_lang", hdr)
	}

	if conf.TagPath == "" {
		return fmt.Errorf("%stag_path is empty", hdr)
	}
	if conf.TagCategoryPath == "" {
		return fmt.Errorf("%stag_category_path is empty", hdr)
	}

	if conf.LangPath == "" {
		return fmt.Errorf("%slang_path is empty", hdr)
	}

	if conf.PubUrl == "" {
		return fmt.Errorf("%spub_url is empty", hdr)
	}

	return nil
}

func NewCtx(a *api.Ctx) *Ctx {
	ctx := new(Ctx)
	ctx.Api = a
	ctx.Config = new(Config)
	a.Config.UnmarshalKeyPanic("lang", ctx.Config, ctx.Config.Validate)

	ctx.Tag = NewTagCtx(a, ctx.Config.TagPath, ctx.Config.TagCategoryPath)
	ctx.TemplateCache = make(map[string]*template.Template)

	var err error
	ctx.userLangIndex, err = NewLangIndex(ctx.Config.LangPath)
	if err != nil {
		panic(err)
	}

	a.AnonGroup.GET("/lang", ctx.GetLang())
	a.AnonGroup.GET("/lang/explain", ctx.ExplainLangHandler())
	a.AnonGroup.GET("/tag/explain", ctx.Tag.ExplainTagHandler())
	a.AnonGroup.GET("/tag/lang", ctx.Tag.GetTagLangHandler())
	a.AnonGroup.GET("/tag", ctx.Tag.GetTagHandler())

	return ctx
}
