package search

import (
	"fmt"
	"sport/helpers"
	"time"

	"github.com/gin-gonic/gin"
)

type SortColumn string

const (
	Accuracy      SortColumn = "accuracy"
	Price         SortColumn = "price"
	NumberReviews SortColumn = "number_reviews"
	AvgMark       SortColumn = "avg_mark"
	Capacity      SortColumn = "capacity"
)

type SearchSortRecord struct {
	Column SortColumn
	IsDesc bool
}

// const defaultSort = "(tag_points + distance_points + title_points + ntc)"

// var SortExp = map[string]string{
// 	//
// 	"accuracy":       defaultSort,
// 	"price":          "price",
// 	"number_reviews": "number_reviews",
// 	"avg_mark":       "avg_mark",
// 	"capacity":       "capacity",
// }

type ApiSearchRequest struct {
	Query string
	SearchRequestOpts
}

func (s *ApiSearchRequest) ValidationErr(msg string) error {
	return fmt.Errorf("Validate SearchInstructorsRequest: %s", msg)
}

const maxDistKm = 200

const MaxPageSize = 100

func (s *ApiSearchRequest) Validate(conf *Config) error {

	if s.Age < 0 {
		return s.ValidationErr("Invalid age")
	}

	if len(s.Query) > conf.MaxQueryLength {
		return s.ValidationErr("Invalid query")
	}

	if s.Pagination.Page < 0 {
		s.Pagination.Page = 0
	}

	if s.Pagination.Size <= 0 || s.Pagination.Size > MaxPageSize {
		s.Pagination.Size = 10
	}

	// max 2 user provided langs: app lang and browser lang
	if len(s.Langs) > 2 {
		return s.ValidationErr("invalid langs")
	}

	if s.DistKm < 0 || s.DistKm > maxDistKm {
		s.DistKm = maxDistKm
	}

	if s.SugDistKm < 0 || s.SugDistKm > maxDistKm {
		s.SugDistKm = maxDistKm
	}

	if len(s.Country) > 0 && len(s.Country) != 2 {
		return s.ValidationErr("invalid country")
	}

	byDist := s.DistKm != 0
	byCountry := len(s.Country) != 0
	if !byDist && !byCountry {
		return s.ValidationErr("Must specify search either by distance or by country")
	}

	if s.PriceMin < 0 {
		return s.ValidationErr("invalid price min")
	}

	if s.PriceMax < 0 {
		return s.ValidationErr("invalid price max")
	}

	if s.PriceMax < s.PriceMin {
		return s.ValidationErr("invalid price range")
	}

	if s.CapacityMin < 0 {
		return s.ValidationErr("invalid min capacity")
	}

	if s.CapacityMax < 0 {
		return s.ValidationErr("invalid max capacity")
	}

	if s.CapacityMax < s.CapacityMin {
		return s.ValidationErr("invalid capacity range")
	}

	if len(s.Diffs) > 10 {
		return s.ValidationErr("invalid diffs")
	}

	var dx = map[int]struct{}{}

	for _, x := range s.Diffs {
		if _, ok := dx[x]; ok {
			return s.ValidationErr("duplicate diff")
		} else {
			dx[x] = struct{}{}
		}
	}

	if s.DurationMin < 0 {
		return s.ValidationErr("invalid duration min")
	}

	if s.DurationMax < 0 {
		return s.ValidationErr("invalid duration max")
	}

	if s.DurationMax < s.DurationMin {
		return s.ValidationErr("invalid duration range")
	}

	if (s.HrStart == "") != (s.HrEnd == "") {
		return s.ValidationErr("invalid hr")
	}

	if s.DateStart.IsZero() != s.DateEnd.IsZero() {
		return s.ValidationErr("DateStart xor DateEnd failed")
	}

	if !s.DateStart.IsZero() {
		if s.DateEnd.Before(s.DateStart) {
			return s.ValidationErr("invalid dates: DateStart after DateEnd")
		}
		iv := s.DateEnd.Sub(s.DateStart)
		const aiv = time.Hour * 24 * 365
		if iv > aiv {
			return s.ValidationErr(fmt.Sprintf(
				"invalid dates: Interval between DateStart and DateEnd (%v) exceeeded allowed interval (%v)", iv, aiv))
		}
	}

	var allowedDays = map[Weekday]bool{
		Monday:    true,
		Thursday:  true,
		Tuesday:   true,
		Wednesday: true,
		Friday:    true,
		Saturday:  true,
		Sunday:    true,
	}

	for _, d := range s.Days {
		if enabled, ok := allowedDays[d]; ok {
			if enabled {
				allowedDays[d] = false
			} else {
				return s.ValidationErr("duplicate day")
			}
		} else {
			return s.ValidationErr("invalid day")
		}
	}

	if len(s.Sort) > 5 {
		return s.ValidationErr("invalid sort")
	}

	return nil
}

/*
GET /api/instructor/search

DESC
	search for instructors

REQUEST
	- Authorization: Bearer <token> in header
	- query string:
		- query
			required, search query. too long search query will result in 400
		- lat, lng
			optional, either none or both must be provided
		- filter params:
			see SearchFilterRequest to get possible values
		- sort params:
			see SearchSortRequest to get possible values
		- page, size
			optional, page is 0-based
		- lang:
			optional - language of search.
				if not provided then user language will be used.
					if no user language exists then language will be obtained from request IP
					if that fails endpoint will terminate with 403 status code

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 400 when validation of query string fails
			- 500 on unexpected error
			- 401 if invalid authorization was provided
			- 403 if language / ip validation failed
	- on success application/json body of type []SearchResponse
*/
func (ctx *Ctx) SearchHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		var apiReq ApiSearchRequest
		var err error

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &apiReq, func() error {
				return apiReq.Validate(ctx.Conf)
			}); err != nil {
			g.AbortWithError(400, err)
			return
		}

		req := ctx.NewSearchRequest(&apiReq, g)

		ctx.RWLock.RLock()

		res, err := ctx.SearchWithCache(req, ctx.Cache)

		ctx.RWLock.RUnlock()

		if err != nil {
			if e := helpers.HttpErr(err); e != nil {
				e.WriteAndAbort(g)
				return
			}
			if helpers.PgIsInvalidDatetimeFormat(err) {
				g.AbortWithError(400, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		res.Meta.Langs = req.Langs

		g.AbortWithStatusJSON(200, res)
	}
}
