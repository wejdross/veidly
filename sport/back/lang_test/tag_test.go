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

func TestExplainTag(t *testing.T) {
	method := "GET"
	path := "/api/tag/explain"

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
			RequestUrl:         path + "?t=pl&s=" + url.QueryEscape(`["orienteering"]`),
			ExpectedStatusCode: 200,
			ExpectedBody: ptr(helpers.JsonMustSerializeStr([]lang.TagWithCategory{
				{
					Tag: lang.Tag{
						Name: "orienteering",
						Translations: map[string]string{
							"pl": "biegi na orientację",
						},
					},
				},
			})),
		}, {
			RequestMethod:      method,
			RequestUrl:         path + "?t=pl&s=" + url.QueryEscape(`["curling", "basketball"]`),
			ExpectedStatusCode: 200,
			ExpectedBody: ptr(helpers.JsonMustSerializeStr([]lang.TagWithCategory{
				{
					Tag: lang.Tag{
						Name: "curling",
						Translations: map[string]string{
							"pl": "curling",
						},
					},
				}, {
					Tag: lang.Tag{
						Name: "basketball",
						Translations: map[string]string{
							"pl": "koszykówka",
						},
					},
				},
			})),
		},
	})
}

func TestSearchTag(t *testing.T) {
	method := "GET"
	path := "/api/tag"

	langCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         path,
			ExpectedStatusCode: 200,
		}, {
			RequestMethod:      method,
			RequestUrl:         path + "?q=" + url.QueryEscape("katare"),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []lang.ScoredTagWithCategory
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				if res[0].Tag.Name != "karate" {
					return fmt.Errorf("1: unexpected tag 0")
				}
				return nil
			},
		}, {
			RequestMethod:      method,
			RequestUrl:         path + "?q=" + url.QueryEscape("may thai") + "&l=pl",
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var res []lang.ScoredTagWithCategory
				if err := json.Unmarshal(b, &res); err != nil {
					return err
				}
				if res[0].Tag.Name != "muay thai" {
					return fmt.Errorf("2: unexpected tag 0")
				}
				if res[1].Tag.Name != "massage" {
					return fmt.Errorf("2: unexpected tag 1")
				}
				return nil
			},
		},
	})
}
