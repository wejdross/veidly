package train

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"sport/dc"
	"sport/sub"

	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// used as base for update and create request
// dont provide seconds or lower parts of dates
type TrainingRequest struct {
	Title       string
	Description string
	DateStart   time.Time
	DateEnd     time.Time
	Capacity    int

	//
	RequiredGear    []string
	RecommendedGear []string
	InstructorGear  []string
	MinAge          int
	MaxAge          int

	LocationText    string
	LocationLat     float64
	LocationLng     float64
	LocationCountry string
	Price           int
	Currency        string // currency code
	Tags            []string
	Diff            []int32
	ManualConfirm   bool
	AllowExpress    bool
	Disabled        bool

	TrainingSupportsDisabled bool
	PlaceSupportsDisabled    bool
}

func (itr *TrainingRequest) ValidationErr(msg string) error {
	return fmt.Errorf("Validate CreateInstructorTrainingRequest: %s", msg)
}

func (itr *TrainingRequest) ValidateGearArray(ctx *Ctx, a []string) error {
	if len(a) > ctx.Config.MaxRequirementArrLen {
		return itr.ValidationErr("Invalid requirement array length")
	}
	for i := range a {
		if len(a[i]) > ctx.Config.MaxRequirementStrLen {
			return itr.ValidationErr("Invalid array item length")
		}
	}
	return nil
}

func (itr *TrainingRequest) ValidateDiffs() error {
	df := itr.Diff
	if len(df) == 0 {
		return nil
	}
	// come on its up to 3 difficulties O(n^2) is ok
	for i := range df {
		if df[i] < 1 || df[i] > 3 {
			return fmt.Errorf("invalid diff")
		}
		for j := range df {
			if i != j && df[i] == df[j] {
				return fmt.Errorf("duplicate diff")
			}
		}
	}
	return nil
}

func (itr *TrainingRequest) Validate(ctx *Ctx) error {

	if err := itr.ValidateDiffs(); err != nil {
		return err
	}

	if err := itr.ValidateGearArray(ctx, itr.RequiredGear); err != nil {
		return err
	}
	if err := itr.ValidateGearArray(ctx, itr.RecommendedGear); err != nil {
		return err
	}
	if err := itr.ValidateGearArray(ctx, itr.InstructorGear); err != nil {
		return err
	}

	if itr.MinAge < 0 || itr.MaxAge < 0 ||
		itr.MinAge > itr.MaxAge || itr.MaxAge > ctx.Config.MaxAge {
		return itr.ValidationErr("invalid age")
	}

	if itr.AllowExpress && itr.ManualConfirm {
		return itr.ValidationErr(
			"training cant have manual confirm with express payments enabled")
	}

	if itr.Title == "" || len(itr.Title) > 128 {
		return itr.ValidationErr("Invalid title")
	}
	if !itr.DateStart.IsZero() && !itr.DateEnd.IsZero() {
		if itr.DateEnd.Before(itr.DateStart) {
			return itr.ValidationErr("Invalid date range")
		}
	}
	if itr.DateStart.IsZero() != itr.DateEnd.IsZero() {
		return itr.ValidationErr("Invalid date range: either both start_date and end_date must be specified, or none")
	}
	if itr.Capacity <= 0 {
		return itr.ValidationErr("Invalid capacity")
	}
	/*
		right now i decided not to support free training sessions.
		we may introduce support for it in future,
		if we adapt reservation flow to omit payments
		'TODO'
	*/
	if itr.Price < 100 {
		return itr.ValidationErr("Invalid price")
	}
	if itr.Currency == "" {
		return itr.ValidationErr("Invalid currency code")
	}
	if ok := ctx.Config.AllowedCurrencies[itr.Currency]; !ok {
		return itr.ValidationErr("currency code not supported: " + itr.Currency)
	}

	if len(itr.Tags) > 5 {
		return itr.ValidationErr("invalid tags")
	}

	tx := make(map[string]struct{})
	for _, t := range itr.Tags {
		if len(t) == 0 || len(t) > 96 {
			return itr.ValidationErr("invalid tag")
		}
		if _, exists := tx[t]; exists {
			return itr.ValidationErr("duplicate tag in request")
		}
		tx[t] = struct{}{}
	}

	if itr.LocationCountry != "" {
		if len(itr.LocationCountry) != 2 {
			return itr.ValidationErr("invalid location")

		}
	}

	return nil
}

func (itr *TrainingRequest) NewTraining(
	ctx *Ctx,
	userID uuid.UUID,
) (*Training, error) {
	instructorID, err := ctx.Instr.DalReadInstructorID(userID)
	if err != nil {
		return nil, err
	}
	return &Training{
		ID:              uuid.New(),
		InstructorID:    instructorID,
		TrainingRequest: *itr,
		CreatedOn:       time.Now().In(time.UTC),
		MainImgID:       "",
		SecondaryImgIDs: nil,
	}, nil
}

// used as update request
type UpdateTrainingRequest struct {
	ID uuid.UUID
	TrainingRequest
}

func (itr *UpdateTrainingRequest) Validate(ctx *Ctx) error {
	return itr.TrainingRequest.Validate(ctx)
}

type Training struct {
	ID            uuid.UUID
	InstructorID  uuid.UUID
	NumberReviews int
	AvgMark       int
	CreatedOn     time.Time
	// those will be filled based on db relpaths
	// url = concat(basepath, relpath)
	MainImgUrl       string
	SecondaryImgUrls []string
	MainImgID        string
	SecondaryImgIDs  []string
	TrainingRequest
}

func TrainingSelectFields() string {
	return `
		t.id
		,t.instructor_id
		,t.title
		,t.description
		,t.date_start
		,t.date_end
		,t.capacity

		,t.required_gear
		,t.recommended_gear
		,t.instructor_gear
		,t.min_age
		,t.max_age

		,t.location_text
		,t.location_lat
		,t.location_lng
		,t.location_country
		,t.price
		,t.currency
		,t.tags
		,t.diff
		,t.number_reviews
		,t.avg_mark
		,t.manual_confirm
		,t.allow_express
		,t.created_on
		,t.main_img_relpath
		,t.secondary_img_relpaths
		,t.disabled
		,t.training_supports_disabled
		,t.place_supports_disabled
	`
}

// must correspond to return value from TrainingSelectFields() string
func (t *Training) ScanFields() []interface{} {
	return []interface{}{
		&t.ID,
		&t.InstructorID,
		&t.Title,
		&t.Description,
		&t.DateStart,
		&t.DateEnd,
		&t.Capacity,

		(*pq.StringArray)(&t.RequiredGear),
		(*pq.StringArray)(&t.RecommendedGear),
		(*pq.StringArray)(&t.InstructorGear),
		&t.MinAge,
		&t.MaxAge,

		&t.LocationText,
		&t.LocationLat,
		&t.LocationLng,
		&t.LocationCountry,
		&t.Price,
		&t.Currency,
		(*pq.StringArray)(&t.Tags),
		(*pq.Int32Array)(&t.Diff),
		&t.NumberReviews,
		&t.AvgMark,
		&t.ManualConfirm,
		&t.AllowExpress,
		&t.CreatedOn,
		&t.MainImgID,
		(*pq.StringArray)(&t.SecondaryImgIDs),
		&t.Disabled,
		&t.TrainingSupportsDisabled,
		&t.PlaceSupportsDisabled,
	}
}

func (t *Training) PostprocessAfterDbScan(ctx *Ctx) {
	if t.MainImgID != "" {
		t.MainImgUrl =
			ctx.Static.Config.AppendToBaseUrl(t.MainImgID)
	} else {
		t.MainImgUrl = ""
	}
	t.SecondaryImgUrls = make([]string, len(t.SecondaryImgIDs))
	for i := range t.SecondaryImgIDs {
		t.SecondaryImgUrls[i] =
			ctx.Static.Config.AppendToBaseUrl(t.SecondaryImgIDs[i])
	}
}

type TrainingWithJoins struct {
	Training    Training
	Occurrences []OccWithSecondary
	Groups      []Group
	Dcs         []dc.Dc
	Sms         []sub.SubModel
}

func (t *TrainingWithJoins) ValidateOccurrenceErr(msg string) error {
	return fmt.Errorf("Validate TrainingTemplateResponse Occurrence: %s", msg)
}

/*
	validates that specified time matches training template and returns end boundary on success
*/
func (t *TrainingWithJoins) ValidateOccurrence(occTime time.Time) (occ *OccWithSecondary, end time.Time, err error) {
	// if len(t.Occurrences) == 0 {
	// 	return time.Time{}, t.ValidateOccurrenceErr("cant validate empty occurrences")
	// }
	if !t.Training.DateStart.IsZero() && t.Training.DateStart.After(occTime) {
		return nil, time.Time{}, t.ValidateOccurrenceErr("target training didnt start yet")
	}
	if !t.Training.DateEnd.IsZero() && t.Training.DateEnd.Before(occTime) {
		return nil, time.Time{}, t.ValidateOccurrenceErr("target training ended already")
	}
	for i := 0; i < len(t.Occurrences); i++ {
		if t.Occurrences[i].RepeatDays <= 0 {
			if t.Occurrences[i].DateStart.Equal(occTime) {
				return &t.Occurrences[i], t.Occurrences[i].DateEnd, nil
			}
		} else {
			df := occTime.Sub(t.Occurrences[i].DateStart)
			if df < 0 {
				//return time.Time{}, t.ValidateOccurrenceErr("target training didnt start yet")
				continue
			}
			if df == 0 {
				return &t.Occurrences[i], t.Occurrences[i].DateEnd, nil
			}
			days := df.Hours() * 24
			if math.Trunc(days) != days {
				continue // specified occ record doesnt match
			}
			if (int(days) % t.Occurrences[i].RepeatDays) == 0 {
				return &t.Occurrences[i], occTime.Add(t.Occurrences[i].DateEnd.Sub(t.Occurrences[i].DateStart)), nil
			}
		}
	}
	//fmt.Println(occTime, helpers.JsonMustSerializeFormatStr(t.Occurrences))
	return nil, time.Time{}, fmt.Errorf("couldnt validate occurrence against template")
}

func (tt *Training) ScanErr(m string) error {
	return fmt.Errorf("Scan TrainingTemplate: " + m)
}

func (tt *Training) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), tt)
}

func (tt *Training) Value() (driver.Value, error) {
	return json.Marshal(tt)
}

type CreateTrainingRequest struct {
	Training    TrainingRequest
	Occurrences []CreateOccRequest
	ReturnID    bool
}

func (r *CreateTrainingRequest) ValidationErr(msg string) error {
	return fmt.Errorf("Validate PostInstructorTrainingRequest: %s", msg)
}

func (r *CreateTrainingRequest) Validate(ctx *Ctx) error {
	if err := r.Training.Validate(ctx); err != nil {
		return err
	}
	for i := 0; i < len(r.Occurrences); i++ {
		if err := r.Occurrences[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

type ObjectKey struct {
	ID uuid.UUID
}

func (o *ObjectKey) Validate() error {
	if o.ID == uuid.Nil {
		return fmt.Errorf("Validate ObjectKey: invalid id")
	}
	return nil
}

type DeleteTrainingOccRequest struct {
	ObjectKey
	TrainingID uuid.UUID
}

func (o *DeleteTrainingOccRequest) Validate() error {
	if o.TrainingID == uuid.Nil {
		return fmt.Errorf("Validate DeleteTrainingOccRequest: invalid TrainingID")
	}
	return o.ObjectKey.Validate()
}
