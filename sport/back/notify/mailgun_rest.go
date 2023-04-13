package notify

// import (
// 	"encoding/base64"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"net/url"
// 	"sport/helpers"
// 	"strings"
// )

// // MailgunRestCtx mailgun rest api configuration
// type MailgunRestCtx struct {
// 	APIKey string `yaml:"api_key"`
// 	From   string
// 	Domain string
// 	Host   string
// }

// func (ctx *MailgunRestCtx) Validate() error {
// 	const hdr = "validate MailgunRestCtx: "
// 	if ctx.APIKey == "" {
// 		return fmt.Errorf("%sapi_key was empty", hdr)
// 	}
// 	if ctx.From == "" {
// 		return fmt.Errorf("%sfrom was empty", hdr)
// 	}
// 	if ctx.Domain == "" {
// 		return fmt.Errorf("%sdomain was empty", hdr)
// 	}
// 	if ctx.Host == "" {
// 		return fmt.Errorf("%shost was empty", hdr)
// 	}
// 	return nil
// }

// // create and validate ctx from config
// func NewMailgunRestCtx(configPath string) (*MailgunRestCtx, error) {
// 	var wrapper struct {
// 		Ctx MailgunRestCtx `yaml:"mailgun_rest"`
// 	}
// 	if err := helpers.YmlParseFile(configPath, &wrapper); err != nil {
// 		return nil, err
// 	}
// 	return &wrapper.Ctx, wrapper.Ctx.Validate()
// }

// // MailgunSend send email via mailgun rest api
// func (ctx *MailgunRestCtx) MailgunSend(to []string, subject, html string) error {

// 	var data = make(url.Values)
// 	if len(to) == 0 {
// 		return fmt.Errorf("no dest specified - nothing to do")
// 	}
// 	for _, dest := range to {
// 		data.Set("to", dest)
// 	}
// 	data.Set("from", ctx.From+" <mailgun@"+ctx.Domain+">")
// 	data.Set("subject", subject)
// 	data.Set("html", html)

// 	url := fmt.Sprintf(
// 		"https://%s/v3/%s/messages",
// 		ctx.Host,
// 		ctx.Domain)

// 	req, err := http.NewRequest(
// 		"POST",
// 		url,
// 		strings.NewReader(data.Encode()))
// 	if err != nil {
// 		return err
// 	}

// 	bCred := "api:" + ctx.APIKey
// 	bCred64 := base64.StdEncoding.EncodeToString([]byte(bCred))

// 	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
// 	req.Header.Add("Authorization", "Basic "+bCred64)

// 	var res *http.Response
// 	if res, err = http.DefaultClient.Do(req); err != nil {
// 		return err
// 	}

// 	if res.StatusCode != 200 {
// 		var resbody []byte
// 		var rs string
// 		if resbody, err = ioutil.ReadAll(res.Body); err == nil {
// 			rs = string(resbody)
// 		}
// 		return fmt.Errorf("Unexepected status code: %d\n%s", res.StatusCode, rs)
// 	}
// 	return nil
// }
