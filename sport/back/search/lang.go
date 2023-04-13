package search

import (
	"sport/api"
	"sport/helpers"
	"sport/user"

	"github.com/gin-gonic/gin"
)

// return err on failure (must be handled)
// return "", nil when no language was found but no error occurred either
func (ctx *Ctx) GetSearchLangFromUserAccount(g *gin.Context) (string, error) {
	if g.GetHeader("Authorization") == "" {
		return "", nil
	}
	userID, err := ctx.Api.AuthorizeUserFromCtx(g)
	if err != nil {
		return "", err
	}
	u, err := ctx.User.DalReadUser(userID, user.KeyTypeID, true)
	if err != nil {
		return "", err
	}
	return u.Language, nil
}

func (ctx *Ctx) GetSearchLangFromGeo(g *gin.Context) (string, error) {
	ip, err := api.GetIP(g, true)
	if err != nil {
		return "", err
	}

	matches, err := GetPossibleIP4Localizations(ctx.Dal, ip)
	if err != nil {
		// could be ip6 - just ignore it for now
		return "", nil
	}

	loc := GetBestIP4Localization(matches)
	if loc.IsFilled && loc.Language != "" {
		return loc.Language, nil
	}
	return "", nil
}

func (ctx *Ctx) GetSearchLangs(apiReq *ApiSearchRequest, g *gin.Context) []string {

	var detects [3]string

	helpers.RunFunctionsInParallel([]func() error{
		func() error {
			lang, err := ctx.GetSearchLangFromUserAccount(g)
			if err == nil {
				detects[0] = lang
			}
			return nil
		},
		func() error {
			lang, err := ctx.GetSearchLangFromGeo(g)
			if err == nil {
				detects[1] = lang
			}
			return nil
		},
		func() error {
			lang, err := ctx.DetectLanguage(apiReq.Query)
			if err == nil {
				detects[2] = lang
			}
			return nil
		},
	})

	possibleLangs := append(detects[:], apiReq.Langs...)
	filteredLangs := make([]string, 0, len(possibleLangs))
	duplicates := make(map[string]struct{}, len(possibleLangs))
	for i := range possibleLangs {

		l := possibleLangs[i]

		// is lang supported?
		if _, e := ctx.User.LangCtx.Tag.SupportedLangMap[l]; !e {
			continue
		}

		// is lang duplicate
		if _, e := duplicates[l]; e {
			continue
		}

		filteredLangs = append(filteredLangs, l)
		duplicates[l] = struct{}{}
	}

	return filteredLangs
}
