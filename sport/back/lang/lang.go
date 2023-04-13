package lang

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"sport/helpers"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ctx *Ctx) ValidateUserLang(lang string) bool {
	if _, ok := ctx.userLangIndex[lang]; ok {
		return true
	} else {
		return false
	}
}

func (ctx *Ctx) ValidateApiLang(lang string) bool {
	if ok := ctx.Config.ApiLang[lang]; ok {
		return true
	} else {
		return false
	}
}

type Lang struct {
	Endonym string
	En      string
	//De      string
	//Fr           string
	ISO_639_1    string
	Translations map[string]string
}

type LangIndex map[string]Lang

func NewLangIndex(path string) (LangIndex, error) {
	fc, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var ls []Lang
	if err := json.Unmarshal(fc, &ls); err != nil {
		return nil, err
	}
	ret := make(LangIndex, len(ls))
	for i := range ls {
		comma := strings.Index(ls[i].Endonym, ",")
		if comma >= 0 {
			ls[i].Endonym = strings.TrimSpace(ls[i].Endonym[:comma])
		}
		comma = strings.Index(ls[i].En, ",")
		if comma >= 0 {
			ls[i].En = strings.TrimSpace(ls[i].En[:comma])
		}
		ret[ls[i].ISO_639_1] = ls[i]
	}
	return ret, nil
}

const MaxLangs = 59

func (ctx *Ctx) ExplainLangHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var langs []string

		l := g.Query("l")
		if l == "" || l == "[]" {
			g.AbortWithError(400, fmt.Errorf("empty langs"))
			return
		}
		if err := json.Unmarshal([]byte(l), &langs); err != nil {
			g.AbortWithError(400, err)
			return
		}

		if len(langs) > MaxLangs {
			g.AbortWithError(400, fmt.Errorf("too many langs"))
			return
		}

		t := g.Query("t")

		var ret = make([]Lang, len(langs))

		for i := range langs {
			iso := langs[i]
			if f, e := ctx.userLangIndex[iso]; !e {
				g.AbortWithError(404, fmt.Errorf("lang: %s not found", iso))
				return
			} else {
				ret[i] = Lang{
					Endonym:   f.Endonym,
					En:        f.En,
					ISO_639_1: f.ISO_639_1,
				}
				if x := f.Translations[t]; x != "" {
					ret[i].Translations = map[string]string{
						t: x,
					}
				}
			}
		}

		g.AbortWithStatusJSON(200, ret)
	}
}

type LangWithScore struct {
	Lang
	Score float32
}

const langMulThr = 0.8
const nopointThr = 0.3

/*
	find matches for requested language
*/
func (ctx *Ctx) GetLang() gin.HandlerFunc {
	return func(g *gin.Context) {
		q := strings.ToLower(g.Query("q"))
		if len(q) > 100 {
			g.AbortWithError(400, fmt.Errorf("too long query"))
			return
		}

		l := g.Query("l")

		var scores [4]float32
		var sum float32
		var res = make([]LangWithScore, len(ctx.userLangIndex))
		var i int
		var t string
		withTranslation := false
		for lang := range ctx.userLangIndex {
			if q == "" {
				continue
			}

			entry := ctx.userLangIndex[lang]
			t, withTranslation = entry.Translations[l]
			if withTranslation {
				scores[1] = helpers.NormLevenshtein(t, q)
			} else {
				scores[1] = 0
			}

			scores[0] = helpers.NormLevenshtein(entry.Endonym, q)
			scores[2] = helpers.NormLevenshtein(entry.En, q)
			scores[3] = helpers.NormLevenshtein(entry.ISO_639_1, q)

			sum = 0
			for i := range scores {
				if scores[i] < nopointThr {
					continue
				}
				if scores[i] >= langMulThr {
					sum += 5 * scores[i]
				} else {
					sum += scores[i]
				}
			}

			res[i] = LangWithScore{
				Lang: Lang{
					Endonym:   entry.Endonym,
					En:        entry.En,
					ISO_639_1: entry.ISO_639_1,
				},
				Score: sum,
			}
			if entry.Translations[l] != "" {
				res[i].Translations = map[string]string{
					l: entry.Translations[l],
				}
			}
			i++
		}
		sort.Slice(res, func(i, j int) bool {
			return res[i].Score > res[j].Score
		})
		g.AbortWithStatusJSON(200, res[:5])
	}
}

func (ctx *Ctx) ApiLangOrDefault(l string) string {
	if ctx.Config.ApiLang[l] {
		return l
	} else {
		return ctx.Config.DefaultLang
	}
}
