package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (ctx *Ctx) DetectLanguage(query string) (string, error) {
	var payload struct {
		Q string `json:"q"`
	}
	payload.Q = query
	var j []byte

	j, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// httpReq, err := http.NewRequest(
	// 	"GET",
	// 	fmt.Sprintf("%s?q=%s", ctx.Conf.LinguaApiUrl, url.QueryEscape(query)), nil)
	httpReq, err := http.NewRequest("POST", ctx.Conf.LinguaApiUrl, bytes.NewBuffer(j))

	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}

	defer httpRes.Body.Close()

	jres, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return "", err
	}

	if httpRes.StatusCode != 200 {
		return "", fmt.Errorf("%s\n%s", httpRes.Status, string(jres))
	}

	var res []struct {
		Lang  string  `json:"lang"`
		Score float32 `json:"score"`
	}

	if err := json.Unmarshal(jres, &res); err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", fmt.Errorf("no language detections found")
	}

	// if relative point distance is smaller than threshold - its unreliable and ignore
	if len(res) > 1 {
		// results are sorted by score desc
		if res[0].Score-res[1].Score < 0.2 {
			return "", fmt.Errorf("unreliable language detection")
		}
	}

	return res[0].Lang, nil
}

func (ctx *Ctx) GenericGoogleTranslateRequest(url string, req, res interface{}) error {

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(
		"POST", url,
		bytes.NewReader(jsonReq))
	if err != nil {
		return err
	}

	if !ctx.Conf.GoogleTranslate.Enabled || ctx.Conf.GoogleTranslate.Token == "" {
		return fmt.Errorf("cannot perform translate without translate service being enabled")
	}

	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("Authorization", "Bearer "+ctx.Conf.GoogleTranslate.Token)

	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}

	defer httpRes.Body.Close()

	jres, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	if httpRes.StatusCode != 200 {
		return fmt.Errorf("%s\n%s", httpRes.Status, string(jres))
	}

	if err := json.Unmarshal(jres, res); err != nil {
		return err
	}

	return nil
}

type TranslateResponse struct {
	Data struct {
		Transalations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}

type TranslateFormat string

const (
	Text TranslateFormat = "text"
)

type TranslateRequest struct {
	Q      string          `json:"q"`
	Source string          `json:"source"`
	Target string          `json:"target"`
	Format TranslateFormat `json:"format"`
}

/*
	res must not be nil. (function will write response into res)
*/
func (ctx *Ctx) GoogleTranslate(req *TranslateRequest, res *TranslateResponse) error {
	return ctx.GenericGoogleTranslateRequest(
		"https://translation.googleapis.com/language/translate/v2",
		req, req,
	)
}

type DetectRequest struct {
	Q string `json:"q"`
}

type DetectResponse struct {
	Data struct {
		Detections []struct {
			Language string `json:"language"`
		} `json:"detections"`
	} `json:"data"`
}

func (ctx *Ctx) GoogleDetectLang(req *DetectRequest, res *DetectResponse) error {
	return ctx.GenericGoogleTranslateRequest(
		"https://translation.googleapis.com/language/translate/v2/detect",
		req, req,
	)
}
