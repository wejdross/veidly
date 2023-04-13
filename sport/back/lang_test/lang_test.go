package lang_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sport/api"
	"sport/helpers"
	"sport/lang"
	"testing"
)

func TestExplainLang(t *testing.T) {
	method := "GET"
	path := "/api/lang/explain"

	ptr := func(s string) *string {
		x := new(string)
		*x = s
		return x
	}

	langCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         path,
			ExpectedStatusCode: 400,
		}, {
			RequestMethod:      method,
			RequestUrl:         path + "?t=pl",
			ExpectedStatusCode: 400,
		}, {
			RequestMethod:      method,
			RequestUrl:         path + "?t=pl&l=[\"zh\"]",
			ExpectedStatusCode: 200,
			ExpectedBody: ptr(helpers.JsonMustSerializeStr([]lang.Lang{
				{
					Endonym:   "中文",
					ISO_639_1: "zh",
					En:        "chinese",
					Translations: map[string]string{
						"pl": "chiński",
					},
				},
			})),
		}, {
			RequestMethod:      method,
			RequestUrl:         path + "?t=en&l=[\"zh\", \"pl\"]",
			ExpectedStatusCode: 200,
			ExpectedBody: ptr(helpers.JsonMustSerializeStr([]lang.Lang{
				{
					Endonym:   "中文",
					ISO_639_1: "zh",
					En:        "chinese",
					Translations: map[string]string{
						"en": "chinese",
					},
				}, {
					Endonym:   "język polski",
					ISO_639_1: "pl",
					En:        "polish",
					Translations: map[string]string{
						"en": "polish",
					},
				},
			})),
		},
	})
}

func TestSearchLang(t *testing.T) {
	method := "GET"
	path := "/api/lang"

	langCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         path,
			ExpectedStatusCode: 200,
		}, {
			RequestMethod:      method,
			RequestUrl:         path + "?q=" + url.QueryEscape("po"),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []lang.LangWithScore
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				if res[0].En != "lao" && res[0].Score != 1.5 {
					return fmt.Errorf("1: unexpected lang")
				}
				return nil
			},
		}, {
			RequestMethod:      method,
			RequestUrl:         path + "?q=" + url.QueryEscape("po") + "&l=pl",
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []lang.LangWithScore
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				if res[0].En != "lao" && res[0].Score != 1.5 && res[0].Translations["pl"] != "laotański" {
					return fmt.Errorf("2: unexpected lang")
				}
				return nil
			},
		},
	})
}
