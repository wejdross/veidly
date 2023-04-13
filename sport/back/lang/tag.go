package lang

import (
	"encoding/json"
	"fmt"
	"sort"
	"sport/api"
	"sport/helpers"
	"strings"

	"github.com/gin-gonic/gin"
)

type TagWithCategory struct {
	Tag Tag
}

type Tag struct {
	Name         string
	Translations map[string]string
}

type TagCtx struct {
	_tagArr          []Tag
	TagMap           map[string]*Tag
	SupportedLangs   []string
	SupportedLangMap map[string]struct{}
}

func (ctx *TagCtx) ValidateTag(tag string) bool {
	_, e := ctx.TagMap[tag]
	return e
}

func (ctx *TagCtx) ExplainTags(tags []string, langs []string, noError bool) ([]TagWithCategory, error) {
	var ret = make([]TagWithCategory, len(tags))

	for i := range tags {
		_translations := map[string]string{}
		t := tags[i]
		if f, e := ctx.TagMap[t]; !e {
			if noError {
				for iter := range langs {
					if f != nil {
						if _, ok := f.Translations[langs[iter]]; ok {
							_translations[langs[iter]] = f.Translations[langs[iter]]
						}
					}
				}
				ret[i] = TagWithCategory{
					Tag: Tag{
						Name:         t,
						Translations: _translations,
					},
				}
			} else {
				return nil, fmt.Errorf("tag: %s not found", t)
			}
		} else {
			for iter := range langs {
				if f != nil {
					if _, ok := f.Translations[langs[iter]]; ok {
						_translations[langs[iter]] = f.Translations[langs[iter]]
					}
				}
			}
			t := Tag{
				Name:         t,
				Translations: _translations,
			}

			ret[i] = TagWithCategory{
				Tag: t,
			}
		}
	}

	return ret, nil
}

// GET {api_url}/api/tag?t=["tag", "other tag"]
func (ctx *TagCtx) ExplainTagHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		t := g.Query("t")
		s := g.Query("s")

		var tags []string
		if err := json.Unmarshal([]byte(s), &tags); err != nil {
			g.AbortWithError(400, err)
			return
		}

		if len(tags) == 0 {
			g.AbortWithError(400, fmt.Errorf("empty tags"))
			return
		}

		ret, _ := ctx.ExplainTags(tags, []string{t}, true)

		g.AbortWithStatusJSON(200, ret)
	}
}

func (ctx *TagCtx) GetTagLangHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.AbortWithStatusJSON(200, ctx.SupportedLangs)
	}
}

func NewTagCtx(a *api.Ctx, tagPath, categoryPath string) *TagCtx {

	ret := new(TagCtx)

	var tags []Tag

	if err := helpers.JsonDeserializeFile(tagPath, &tags); err != nil {
		panic(err)
	}

	ret.SupportedLangMap = make(map[string]struct{})
	ret.SupportedLangs = make([]string, len(tags[0].Translations))
	ret.TagMap = make(map[string]*Tag)
	ret._tagArr = tags

	for i := range tags {
		if i == 0 {
			for t := range tags[0].Translations {
				ret.SupportedLangs[i] = t
				i++
				ret.SupportedLangMap[t] = struct{}{}
			}
		}
		ret.TagMap[tags[i].Name] = &tags[i]
	}

	return ret
}

type ScoredTagWithCategory struct {
	TagWithCategory
	Score float32
}

func (ctx *TagCtx) GetTagHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		q := strings.ToLower(g.Query("q"))
		if len(q) > 100 {
			g.AbortWithError(400, fmt.Errorf("too long query"))
			return
		}

		l := g.Query("l")

		if q == "" {
			g.AbortWithStatusJSON(200, []interface{}{})
			return
		}

		var scores [2]float32
		var sum float32
		var res = make([]ScoredTagWithCategory, len(ctx._tagArr))
		var t string
		var i = 0
		withTranslation := false
		for key := range ctx.TagMap {
			entry := ctx.TagMap[key]
			t, withTranslation = entry.Translations[l]
			if withTranslation {
				scores[1] = helpers.NormLevenshtein(t, q)
				if scores[1] > 0.5 {
					scores[1] = 5 * scores[1]
				}
			} else {
				scores[1] = 0
			}

			const secMatchThr = 0.6

			scores[0] = helpers.NormLevenshtein(entry.Name, q)
			if scores[0] < secMatchThr {
				scores[0] = 0
			}

			sum = 0
			for z := range scores {
				sum += scores[z]
			}

			et, _ := ctx.ExplainTags([]string{entry.Name}, []string{l}, true)

			res[i] = ScoredTagWithCategory{
				TagWithCategory: et[0],
				Score:           sum,
			}
			i++
		}

		sort.Slice(res, func(i, j int) bool {
			return res[i].Score > res[j].Score
		})
		g.AbortWithStatusJSON(200, res[:10])
	}
}
