package search

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sport/helpers"
	"sport/instr"
	"sport/lang"
	"sport/rsv"
	"sport/schedule"
	"sport/train"
	"sport/user"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/btree"
	"github.com/google/uuid"
)

type Weekday int

const (
	Monday Weekday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// func (w Weekday) ToTimeWeekday() time.Weekday {
// }

func TimeWeekdayToWeekday(w time.Weekday) Weekday {
	return Weekday((w-1)%7) + 1
}

type SearchRequestOpts struct {
	MinTagPoints       float32
	OmitEmptySchedules bool

	/*  */
	SugEnabled bool
	SugDistKm  int

	/*
		if schedule records is not available then omit it in response,
		note that if this flag is set along with OmitEmptySchedules then trainings without available schedules will be omitted completely
	*/
	OnlyAvailable bool

	// either lat lng distKm must be provided
	// or country
	Lat, Lng float64
	DistKm   int

	Country string // ISO 3166-1 alpha 2 country code

	Pagination helpers.PaginationRequest
	Langs      []string // ISO 639-1 lang codes

	PriceMin    int
	PriceMax    int
	CapacityMin int
	CapacityMax int

	Age int

	// difficulties
	Diffs []int

	// every training with duration between those 2
	// MINUTES
	DurationMin int
	DurationMax int

	// those 2 conditions if specified will affect all the following ones
	// Will match every training session which overlaps with specified dates
	//
	// for example:
	//
	//       | - SESSION 1 - |           | - SESSION 2 - |
	// --------------------------------------------------> TIME
	//                 /-------------/
	//           DateStart        DateEnd
	//
	// SESSION 1 WILL MATCH
	// SESSION 2 WILL NOT MATCH
	DateStart time.Time
	DateEnd   time.Time

	// 1 = monday
	Days []Weekday

	// every training starting between those 2 hrs
	// format: 08:00
	HrStart string
	HrEnd   string

	Sort []SearchSortRecord

	TrainingSupportsDisabled bool
	PlaceSupportsDisabled    bool
}

type SearchRequest struct {
	Query []string
	SearchRequestOpts
}

func (ctx *Ctx) NewSearchRequest(ar *ApiSearchRequest, g *gin.Context) *SearchRequest {
	parts := strings.Split(ar.Query, " ")

	fparts := make([]string, 0, len(parts))
	for i := range parts {
		p := strings.ToLower(strings.TrimSpace(parts[i]))
		if p != "" {
			fparts = append(fparts, p)
		}
	}
	parts = fparts
	if len(parts) > 5 {
		parts = parts[:5]
	}

	ar.Langs = ctx.GetSearchLangs(ar, g)

	return &SearchRequest{
		Query:             parts,
		SearchRequestOpts: ar.SearchRequestOpts,
	}
}

// every time training passes a filter fitlerScore will be implemented by this value
// note that this is only used when request flag 'UseLooseFilter' is set to true
// othwerise every training will get the same filter score
//const FilterScoreIncrement = 1

// if we already saw same instructor in results then decrement score
const NegativeInstrScoreDecrementBy float64 = 0.5

type recordMeta struct {
	isSuggestion bool
}

type Score struct {
	// for dist elimination
	DistScore float64
	// for filtering
	//FilterScore float64
	// for fuzzy string sort and elimination
	InstrScore         float64
	TagScore           float64
	TitleScore         float64
	NegativeInstrScore float64
}

func (ts Score) TotalScore() float64 {
	return ts.DistScore + ts.InstrScore + ts.TagScore + ts.TitleScore - NegativeInstrScoreDecrementBy
}

type TagWithScore struct {
	Tag           lang.Tag
	Score         float32
	NumberMatches int
}

type TagWithScoreMap map[string]*TagWithScore

type ScheduleWithScore struct {
	schedule.TrainingSchedule
	Score      Score
	UserInfo   *user.PubUserInfo
	Instructor *instr.PubInstructorInfo
	meta       recordMeta
}

func deleteFromScheduleWScore(s *[]ScheduleWithScore, ix int) {
	if len(*s) == 0 {
		return
	}

	// delete from trainings
	lastElemIx := len(*s) - 1
	(*s)[ix] = (*s)[lastElemIx]
	(*s) = (*s)[:lastElemIx]
}

func ApplyPagination(
	p *helpers.PaginationRequest, ss *[]ScheduleWithScore) helpers.PaginationResponse {

	size := p.Size
	offset := p.Page * size
	l := len(*ss)
	sliceSize := l - offset
	if sliceSize < 0 {
		// nothing remaining
		(*ss) = []ScheduleWithScore{}
		return helpers.PaginationResponse{
			NPages:            0,
			PaginationRequest: *p,
		}
	}
	last := offset + size
	if last > l {
		last = l
	}

	(*ss) = (*ss)[offset:last]

	return helpers.PaginationResponse{
		NPages:            int(math.Ceil(float64(l) / float64(size))),
		PaginationRequest: *p,
	}
}

func ApplySort(req *SearchRequest, ss []ScheduleWithScore) {

	if len(req.Sort) == 0 {
		req.Sort = []SearchSortRecord{
			{
				Column: Accuracy,
				IsDesc: true,
			},
		}
	}

	// for now we only perform single sort
	s := req.Sort[0]
	isDesc := s.IsDesc
	var sortFunc func(i, j int) bool
	switch s.Column {
	case Accuracy:
		sortFunc = func(i, j int) bool {
			if ss[i].Score.TotalScore() > ss[j].Score.TotalScore() {
				return isDesc
			}
			return !isDesc
		}
	case Price:
		sortFunc = func(i, j int) bool {
			if ss[i].Training.Price > ss[j].Training.Price {
				return isDesc
			}
			if ss[i].Training.Price < ss[j].Training.Price {
				return !isDesc
			}
			// secondary sort by total score desc
			return ss[i].Score.TotalScore() > ss[j].Score.TotalScore()
		}
	case NumberReviews:
		sortFunc = func(i, j int) bool {
			if ss[i].Training.NumberReviews > ss[j].Training.NumberReviews {
				return isDesc
			}
			if ss[i].Training.NumberReviews < ss[j].Training.NumberReviews {
				return !isDesc
			}
			return ss[i].Score.TotalScore() > ss[j].Score.TotalScore()
		}
	case AvgMark:
		sortFunc = func(i, j int) bool {
			if ss[i].Training.AvgMark > ss[j].Training.AvgMark {
				return isDesc
			}
			if ss[i].Training.AvgMark < ss[j].Training.AvgMark {
				return !isDesc
			}
			return ss[i].Score.TotalScore() > ss[j].Score.TotalScore()
		}
	case Capacity:
		sortFunc = func(i, j int) bool {
			if ss[i].Training.Capacity > ss[j].Training.Capacity {
				return isDesc
			}
			if ss[i].Training.Capacity < ss[j].Training.Capacity {
				return !isDesc
			}
			return ss[i].Score.TotalScore() > ss[j].Score.TotalScore()
		}
	}

	// first sort
	sort.Slice(ss, sortFunc)

	// adjust instructor scores

	// ascending accuracy sort is not adjusted - there is no point
	if s.Column != Accuracy || s.IsDesc {
		tc := make(map[uuid.UUID]float64)

		for i := 0; i < len(ss); i++ {
			iid := ss[i].Instructor.ID
			c := tc[iid]
			ss[i].Score.NegativeInstrScore += c
			tc[iid] += NegativeInstrScoreDecrementBy
		}

		if len(tc) > 0 {
			// must sort again since we modified score
			sort.Slice(ss, sortFunc)
		}
	}

	// sort schedules
	for i := range ss {
		sort.Slice(ss[i].Schedule, func(j, z int) bool {
			return ss[i].Schedule[j].Start.Before(ss[i].Schedule[z].Start)
		})
	}
}

/*
	sets score along with instructor info
*/
func postprocessSched(
	sched *ScheduleWithScore,
	tagScore, instrScore, titleScore float32,
	instr *instr.InstructorWithUser) {

	score := &sched.Score
	score.TagScore = float64(tagScore)
	score.TitleScore = float64(titleScore)
	score.InstrScore = float64(instrScore)
	sched.Instructor = &instr.PubInstructorInfo
	sched.UserInfo = &instr.UserInfo
}

func ApplyPostFilterWSchedule(
	req *SearchRequest, cache *SearchCache,
	sched *[]ScheduleWithScore,
	tagScores TagWithScoreMap) {

	// if similarity is lower than this value
	// score will be completely ignored
	const tagSimilarityThr = 0.5

	pointCache := make(map[string]float32)

	for i := 0; i < len(*sched); i++ {
		t := &(*sched)[i].Training
		instr := cache.InstrIndex.Cluster[t.InstructorID]
		// tags score
		var tagScore float32
		// training title score
		var titleScore float32
		// instr name score
		var instrScore float32

		if len(req.Days) > 0 || req.OnlyAvailable {
			dm := make(map[Weekday]struct{})
			for j := range req.Days {
				dm[req.Days[j]] = struct{}{}
			}
			schedsToBeRemoved := make([]int, 0, 2)
			for j := range (*sched)[i].Schedule {
				sched := &(*sched)[i].Schedule[j]

				if len(req.Days) > 0 {
					daySpan := int(math.Ceil(sched.End.Sub(sched.Start).Hours() / 24))
					if daySpan > 7 {
						daySpan = 7
					}
					startDay := TimeWeekdayToWeekday(sched.Start.Weekday())
					found := true
					var _d Weekday
					for d := startDay; d < startDay+Weekday(daySpan); d++ {
						_d = d
						if _d > 7 {
							_d = (_d % 8) + 1
						}
						if _, exists := dm[_d]; !exists {
							found = false
							break
						}
					}
					if !found {
						schedsToBeRemoved = append(schedsToBeRemoved, j)
						continue
					}
				}

				if req.OnlyAvailable && !sched.IsAvailable {
					schedsToBeRemoved = append(schedsToBeRemoved, j)
				}
			}

			if len(schedsToBeRemoved) > 0 {
				sort.Ints(schedsToBeRemoved)
				for z := len(schedsToBeRemoved) - 1; z >= 0; z-- {
					lastIndex := len((*sched)[i].Schedule) - 1
					(*sched)[i].Schedule[schedsToBeRemoved[z]] = (*sched)[i].Schedule[lastIndex]
					(*sched)[i].Schedule = (*sched)[i].Schedule[:lastIndex]
				}
			}
		}

		if (req.OmitEmptySchedules && len((*sched)[i].Schedule) == 0) || instr == nil {
			deleteFromScheduleWScore(sched, i)
			i--
			continue
		}

		if len(req.Query) != 0 {

			tags, _ := cache.Tag.ExplainTags(t.Tags, req.Langs, true)
			// for every translation try to find best one
			var maxTranslationScore float32
			var translationScore float32 = 0

			// for every tag
			for z := range tags {

				if _p, ok := pointCache[tags[z].Tag.Name]; ok {
					tagScore += _p
					continue
				}

				maxTranslationScore = 0

				for lang := range tags[z].Tag.Translations {
					translation := tags[z].Tag.Translations[lang]
					translationScore = 0
					for j := range req.Query {
						lev := helpers.NormLevenshtein(req.Query[j], translation)
						if lev < tagSimilarityThr {
							lev = 0
						}
						translationScore += lev
					}
					if translationScore > maxTranslationScore {
						maxTranslationScore = translationScore
					}
				}

				// extra 'en' translation
				translationScore = 0
				for j := range req.Query {
					lev := helpers.NormLevenshtein(req.Query[j], tags[z].Tag.Name)
					if lev < tagSimilarityThr {
						lev = 0
					}
					translationScore += lev
				}
				if translationScore > maxTranslationScore {
					maxTranslationScore = translationScore
				}

				pointCache[tags[z].Tag.Name] = maxTranslationScore
				tagScores[tags[z].Tag.Name] = &TagWithScore{
					Tag:           tags[z].Tag,
					Score:         maxTranslationScore,
					NumberMatches: 0,
				}
				tagScore += maxTranslationScore
			}

			// normalizing tag score
			if tagScore > 0 {
				tagScore = tagScore / float32(len(tags)*len(req.Query))
			}

			if tagScore < req.MinTagPoints {
				// didnt match tags
				goto DEL
			}

			if _p, ok := pointCache[t.Title]; ok {
				titleScore = _p
			} else {
				for j := range req.Query {
					_p += helpers.NormLevenshtein(req.Query[j], t.Title)
				}
				pointCache[t.Title] = _p
				titleScore = _p
			}

			if len(instr.UserInfo.Name) > 0 {
				if _p, ok := pointCache[instr.UserInfo.Name]; ok {
					instrScore = _p
				} else {
					_p = helpers.NormLevenshtein(strings.Join(req.Query, " "), instr.UserInfo.Name)
					pointCache[t.Title] = _p
					instrScore = _p
				}
			}
		}

		postprocessSched(
			&(*sched)[i], tagScore, instrScore, titleScore, instr)

		if !(*sched)[i].meta.isSuggestion {
			/* set tag scores */
			for j := 0; j < len(t.Tags); j++ {
				n := t.Tags[j]
				c := tagScores[n]
				if c == nil {
					continue
				}
				c.NumberMatches++
			}
		}

		continue

	DEL:
		if req.SugEnabled {
			postprocessSched(
				&(*sched)[i], tagScore, instrScore, titleScore, instr)
			if !(*sched)[i].meta.isSuggestion {
				(*sched)[i].meta.isSuggestion = true
			}
			continue
		}

		deleteFromScheduleWScore(sched, i)
		i--
	}
}

func isValidInstructor(c *SearchCache, instrId uuid.UUID) bool {
	if c.ignoreInvalidInstructors {
		return true
	}
	_instr := c.InstrIndex.Cluster[instrId]
	if _instr == nil {
		return false
	}
	// ignore instructors which are not configured correctly
	if instr.GetInstrConfig(_instr) != 0 {
		return false
	}
	return true
}

func EliminateCountry(req *SearchRequest, c *SearchCache) []ScheduleWithScore {
	trainings := c.TrainingIndex.CountryHash[req.Country].Slice
	cpy := make([]ScheduleWithScore, 0, len(trainings))
	for i := 0; i < len(trainings); i++ {
		if !isValidInstructor(c, trainings[i].Training.InstructorID) {
			continue
		}
		cpy = append(cpy, ScheduleWithScore{
			TrainingSchedule: schedule.TrainingSchedule{
				Schedule:          nil,
				TrainingWithJoins: trainings[i],
			},
			Score: Score{
				DistScore: 1,
			},
		})
	}
	return cpy
}

const kmDistMulSq = 110.25 * 110.25

func LatLngDistanceSquared(lat1, lng1, lat2, lng2 float64) float64 {
	latDistSquared := math.Pow(lat1-lat2, 2)
	lngDistSquared := math.Pow(lng1-lng2, 2)
	latCorrectionOnLng := math.Cos(helpers.DgToRad(lat1))
	actDistKm := kmDistMulSq * (latDistSquared + (lngDistSquared * latCorrectionOnLng))
	return actDistKm
}

func KmToLng(km int, lat float64) float64 {
	return float64(km) / (111.32 * math.Cos(helpers.DgToRad(lat)))
}

func EliminateLatLng(req *SearchRequest, c *SearchCache) []ScheduleWithScore {

	ss := make([]ScheduleWithScore, 0, 100)

	//reqDistKmSq := float64(req.DistKm * req.DistKm)
	reqDistKm := float64(req.DistKm)
	if req.SugEnabled {
		reqDistKm = math.Max(float64(req.DistKm), float64(req.SugDistKm))
	}
	lngDist := KmToLng(int(reqDistKm), req.Lat)
	minLng := req.Lng - lngDist
	maxLng := req.Lng + lngDist

	c.TrainingIndex.Lng.AscendRange(&BTreeItem{Dg: minLng}, &BTreeItem{Dg: maxLng},
		func(i btree.Item) bool {
			bt := i.(*BTreeItem)
			for i := range bt.Val {
				t := &bt.Val[i].Training
				distSq := LatLngDistanceSquared(req.Lat, req.Lng, t.LocationLat, t.LocationLng)
				dist := math.Sqrt(distSq)
				if dist > reqDistKm {
					continue
				}
				var isSuggestion bool
				if req.SugEnabled && dist > float64(req.DistKm) {
					isSuggestion = true
				}
				if !isValidInstructor(c, t.InstructorID) {
					continue
				}
				ss = append(ss, ScheduleWithScore{
					TrainingSchedule: schedule.TrainingSchedule{
						Schedule:          nil,
						TrainingWithJoins: bt.Val[i],
					},
					meta: recordMeta{
						isSuggestion: isSuggestion,
					},
					Score: Score{
						DistScore: math.Abs((dist - float64(req.DistKm)) / float64(req.DistKm)),
					},
				})
			}
			return true
		})

	return ss
}

func HrMinToDec(hr, min int) float32 {
	return float32(hr) + (float32(min) / 60)
}

func ParseHrMin(tm string) (int, int, error) {
	parts := strings.Split(tm, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid fmt")
	}
	p1, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	p2, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return p1, p2, nil
}

func ParseHrMinDec(tmp string) (float32, error) {
	h, m, err := ParseHrMin(tmp)
	if err != nil {
		return 0, err
	}
	return HrMinToDec(h, m), nil
}

func ApplyFilter(
	req *SearchRequest,
	ss *[]ScheduleWithScore,
	c *SearchCache,
) error {

	//ft := make([]*train.TrainingWithJoins, 0, min(len(filteredTrainings)/2, 2))
	// filters - note that most of the trainings have been eliminated at this point
	for i := 0; i < len(*ss); i++ {

		t := &(*ss)[i].Training
		removeOccs := make(map[int]struct{}, 4)
		occs := (*ss)[i].Occurrences

		if req.PriceMin != 0 || req.PriceMax != 0 {
			if t.Price < req.PriceMin || t.Price > req.PriceMax {
				goto DEL
			}
			//score.FilterScore += FilterScoreIncrement
		}

		if req.Age != 0 {
			if t.MinAge != 0 || t.MaxAge != 0 {
				if t.MinAge > req.Age || t.MaxAge < req.Age {
					goto DEL
				}
			}
			//score.FilterScore += FilterScoreIncrement
		}

		if req.TrainingSupportsDisabled && !t.TrainingSupportsDisabled {
			goto DEL
		}

		if req.PlaceSupportsDisabled && !t.PlaceSupportsDisabled {
			goto DEL
		}

		if req.CapacityMin != 0 || req.CapacityMax != 0 {
			if t.Capacity < req.CapacityMin || t.Capacity > req.CapacityMax {
				goto DEL
			}
			//score.FilterScore += FilterScoreIncrement
		}

		if len(req.Diffs) > 0 {
			var found bool
			for j := 0; j < len(req.Diffs); j++ {
				found = false
				for k := 0; k < len(t.Diff); k++ {
					if t.Diff[k] == int32(req.Diffs[j]) {
						found = true
						break
					}
				}
				if !found {
					break
				}
			}
			if !found {
				goto DEL
			}
			//score.FilterScore += FilterScoreIncrement
		}

		if req.HrStart != "" && req.HrEnd != "" {
			s, err := ParseHrMinDec(req.HrStart)
			if err != nil {
				return helpers.NewHttpError(400, "", err)
			}
			e, err := ParseHrMinDec(req.HrEnd)
			if err != nil {
				return helpers.NewHttpError(400, "", err)
			}
			for z := range occs {
				h, m, _ := occs[z].DateStart.Clock()
				x := HrMinToDec(h, m)
				if s > x || e < x {
					// this occurrence doesnt match
					removeOccs[z] = struct{}{}
				}
			}
		}

		if req.DurationMin != 0 || req.DurationMax != 0 {
			dmin := time.Duration(req.DurationMin) * time.Minute
			dmax := time.Duration(req.DurationMax) * time.Minute
			for j := 0; j < len(occs); j++ {
				o := occs[j]
				dur := o.DateEnd.Sub(o.DateStart)
				if dur > dmax || dur < dmin {
					removeOccs[j] = struct{}{}
				}
			}
		}

		// ? date start / end ? (its done in schedule but maybe we can quickly reduce amount of rows in here)

		if len(removeOccs) > 0 {

			// shallow copy of TrainingWithJoins so we will be able to replace occ slice
			// without interfering into cache
			tcpy := new(train.TrainingWithJoins)
			(*tcpy) = *(*ss)[i].TrainingWithJoins

			_occs := make([]train.OccWithSecondary, len(occs)-len(removeOccs))
			l := 0
			for z := range occs {
				_, f := removeOccs[z]
				if f {
					continue
				}
				_occs[l] = occs[z]
				l++
			}
			if len(_occs) == 0 {
				// no occurrences left - this training has been eliminated
				goto DEL
			}
			tcpy.Occurrences = _occs
			(*ss)[i].TrainingWithJoins = tcpy
			//score.FilterScore += FilterScoreIncrement
		}
		// else if filteredByOcc {
		// 	// filtered by occ and every occ passed, increment the score
		// 	score.FilterScore += FilterScoreIncrement
		// }

		continue

	DEL:

		if req.SugEnabled {
			if !(*ss)[i].meta.isSuggestion {
				(*ss)[i].meta.isSuggestion = true
			}
			continue
		}

		deleteFromScheduleWScore(ss, i)
		i--
	}

	return nil
}

type PrettyJsonDuration time.Duration

func (p PrettyJsonDuration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Duration(p).String() + "\""), nil
}

func (p *PrettyJsonDuration) UnmarshalJSON(data []byte) error {
	var x string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	d, err := time.ParseDuration(x)
	*p = PrettyJsonDuration(d)
	return err
}

type SearchPerformanceMetadata struct {
	CountryEliminationTime PrettyJsonDuration
	LatLngEliminationTime  PrettyJsonDuration
	TrainingFilterTime     PrettyJsonDuration
	ScheduleGenerationTime PrettyJsonDuration
	FuzzySearchTime        PrettyJsonDuration
	SortTime               PrettyJsonDuration
}

type SearchMetadata struct {
	helpers.PaginationResponse
	TotalCount int
	Perf       SearchPerformanceMetadata
	Langs      []string
	Tags       []*TagWithScore
}

type SearchResultsWithMetadata struct {
	Data    []ScheduleWithScore
	SugData []ScheduleWithScore
	Meta    SearchMetadata
}

func ProcessTagScores(tagScores TagWithScoreMap) []*TagWithScore {
	ret := make([]*TagWithScore, 0, len(tagScores))
	for t := range tagScores {
		if tagScores[t].Score > 0 {
			ret = append(ret, tagScores[t])
		}
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[j].Score > ret[j].Score
	})
	const maxLen = 5
	if len(ret) > maxLen {
		return ret[:maxLen]
	}
	return ret
}

/*
	search flow:

	1. eliminate lat / lng

	2. apply general filters, such as:
		- age, capacity, price
		O(n') <- n' is number of rows reduced in previous step

	3. apply occurrence filter
		O(n'*M) where M is average number of occurrence records

	4. levenstein distance order by and [maybe] elimination for
		- tags, instructor name, training name


	5. create schedule for those trainings.
		now we need:
		- instr vacations [DB request]
		- reservations    [DB request]
		may be done in parallel

	done
*/
func (ctx *Ctx) SearchWithCache(req *SearchRequest, c *SearchCache) (*SearchResultsWithMetadata, error) {

	var perfMeta = SearchPerformanceMetadata{}

	var ss []ScheduleWithScore
	if req.Country != "" {
		s := time.Now()
		ss = EliminateCountry(req, c)
		perfMeta.CountryEliminationTime = PrettyJsonDuration(time.Since(s))
	} else {
		s := time.Now()
		ss = EliminateLatLng(req, c)
		perfMeta.LatLngEliminationTime = PrettyJsonDuration(time.Since(s))
	}
	if len(ss) == 0 {
		return &SearchResultsWithMetadata{
			Meta: SearchMetadata{
				Perf: perfMeta,
			},
		}, nil
	}

	var err error
	s := time.Now()
	err = ApplyFilter(req, &ss, c)
	perfMeta.TrainingFilterTime = PrettyJsonDuration(time.Since(s))
	if err != nil {
		return nil, err
	}

	// prepare schedule
	s = time.Now()
	for i := 0; i < len(ss); i++ {
		ss[i].Schedule, err = ctx.Sched.NewTrainingScheduleRecords(
			ss[i].TrainingWithJoins,
			helpers.DateRange{
				Start: req.DateStart,
				End:   req.DateEnd,
			},
			schedule.ScheduleBuffers{
				CanRegisterBuffers: rsv.CanRegisterBuffers{
					Vacations:   c.InstrIndex.VacationHash,
					RsvCount:    c.RsvIndex.RsvCountHash,
					GroupedRsvs: c.RsvIndex.GroupHash,
				},
			},
		)
		if err != nil {
			return nil, err
		}
	}

	perfMeta.ScheduleGenerationTime = PrettyJsonDuration(time.Since(s))
	if err != nil {
		return nil, err
	}

	s = time.Now()

	tagScores := make(TagWithScoreMap)
	ApplyPostFilterWSchedule(req, c, &ss, tagScores)

	perfMeta.FuzzySearchTime = PrettyJsonDuration(time.Since(s))

	if err != nil {
		return nil, err
	}

	s = time.Now()

	var sug = make([]ScheduleWithScore, 0, 4)

	if req.SugEnabled {
		for i := 0; i < len(ss); i++ {
			if !ss[i].meta.isSuggestion {
				continue
			}
			sug = append(sug, ss[i])
			deleteFromScheduleWScore(&ss, i)
			i--
		}
	}

	ApplySort(req, ss)

	if len(sug) > 0 {
		ApplySort(req, sug)
	}

	perfMeta.SortTime = PrettyJsonDuration(time.Since(s))

	// pagination
	paginationRes := ApplyPagination(&req.Pagination, &ss)
	if len(sug) > 0 {
		_ = ApplyPagination(&helpers.PaginationRequest{
			Page: 0,
			Size: 5,
		}, &sug)
	}

	ts := ProcessTagScores(tagScores)

	return &SearchResultsWithMetadata{
		Data:    ss,
		SugData: sug,
		Meta: SearchMetadata{
			Perf:               perfMeta,
			PaginationResponse: paginationRes,
			Tags:               ts,
		},
	}, nil
}
