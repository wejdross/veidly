package train

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sport/helpers"
	"time"

	"github.com/google/uuid"
)

// single occurrence record for training
type Occ struct {
	ID         uuid.UUID
	TrainingID uuid.UUID
	OccRequest
}

type OccWithSecondary struct {
	Occ
	SecondaryOccs []SecondaryOcc
}

func (tt *OccWithSecondary) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), tt)
}

func (tt *OccWithSecondary) Value() (driver.Value, error) {
	return json.Marshal(tt)
}

func OccurrenceSelectFields() string {
	return `
		o.id as occurrence_id
		,o.training_id
		,o.date_start
		,o.date_end
		,o.repeat_days
		,o.color
		,o.remarks
	`
}

// this function assumes that occ fields may be nulls in response
// and will coalesce them so scan wont fail
func NullableOccurrenceSelectFields() string {
	return fmt.Sprintf(`
		coalesce(o.id, %s) as occurrence_id
		,coalesce(o.training_id, %s) as training_id
		,coalesce(o.date_start, %s) as date_start
		,coalesce(o.date_end, %s) as date_end
		,coalesce(o.repeat_days, 0)
		,coalesce(o.color, '')
		,coalesce(o.remarks, '')
	`,
		helpers.ZeroUuidString,
		helpers.ZeroUuidString,
		helpers.ZeroTimeString,
		helpers.ZeroTimeString)
}

func (occ *Occ) ScanFields() []interface{} {
	return []interface{}{
		&occ.ID,
		&occ.TrainingID,
		&occ.DateStart,
		&occ.DateEnd,
		&occ.RepeatDays,
		&occ.Color,
		&occ.Remarks,
	}
}

// func (t *Occ) Validate() error {
// 	if t.ID == uuid.Nil {
// 		return fmt.Errorf("Validate TrainingOcc: invalid id")
// 	}
// 	if t.TrainingID == uuid.Nil {
// 		return fmt.Errorf("Validate TrainingOcc: invalid trainingID")
// 	}
// 	return nil
// }

func (tt *Occ) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), tt)
}

func (tt *Occ) Value() (driver.Value, error) {
	return json.Marshal(tt)
}

/*
   date_start and date_end represent when specific session happens
   for example:
   date_start = '2020-01-01 08:00'
   date_end = '2020-01-01 10:00'
   repeat_days = 1

   means that training will at morning (8 am GMT) since 1st january 2020
   and will repeat each 1 day 02 jan, 03 jan, ...
      (will repeat forever unless specified otherwise in instructor_trainings table, column "date_end")

   another example:
   date_start = '2020-09-22 22:00'
   date_end = '2020-09-23 2:00'
   repeat_days = 7

   means that training will start at 10 PM since 2020-09-22
       will last 4 hours until next day's 2 AM
   and will occur at same time each week (2020-09-29, ...)
      (will repeat forever unless specified otherwise in instructor_trainings table, column "date_end")
*/
type OccRequest struct {
	DateStart  time.Time
	DateEnd    time.Time
	RepeatDays int
	Color      string
	Remarks    string
}

func (i *OccRequest) ValidationErr(msg string) error {
	return fmt.Errorf("Validate OccRequest: %s", msg)
}

func (i *OccRequest) Validate() error {
	if i.DateStart.IsZero() {
		return i.ValidationErr("invalid DateStart, cannot be empty")
	}
	if i.DateEnd.IsZero() {
		return i.ValidationErr("invalid DateEnd, cannot be empty")
	}
	if i.DateEnd.IsZero() != i.DateStart.IsZero() {
		return i.ValidationErr("DateStart xor DateEnd failed")
	}
	// if both are zero this yields true anyways
	if i.DateEnd.Before(i.DateStart) {
		return i.ValidationErr("invalid DateEnd, must be bigger or equal to DateStart")
	}
	if i.RepeatDays < 0 {
		return i.ValidationErr("invalid RepeatDays, cannot be < 0")
	}
	if len(i.Remarks) > 512 {
		return i.ValidationErr("invalid remarks")
	}
	if len(i.Color) > 10 {
		return i.ValidationErr("invalid color")
	}
	return nil
}

type CreateOccRequest struct {
	OccRequest
	SecondaryOccs []SecondaryOccRequest
}

func (c *CreateOccRequest) ValidationErr(msg string) error {
	return fmt.Errorf("Validate CreateOccRequest: %s", msg)
}

func (c *CreateOccRequest) Validate() error {

	min := 1<<32 - 1
	max := -1 << 31

	for i := range c.SecondaryOccs {
		if err := c.SecondaryOccs[i].Validate(); err != nil {
			return err
		}
		// offsets overlap <=> dates overlap
		for j := range c.SecondaryOccs {
			if i != j && helpers.OverlapsInt(
				c.SecondaryOccs[i].OffsetStart, c.SecondaryOccs[i].OffsetEnd,
				c.SecondaryOccs[j].OffsetStart, c.SecondaryOccs[j].OffsetEnd) {
				return fmt.Errorf("invalid overlap")
			}
		}
		if c.SecondaryOccs[i].OffsetStart < min {
			min = c.SecondaryOccs[i].OffsetStart
		}
		if c.SecondaryOccs[i].OffsetEnd > max {
			max = c.SecondaryOccs[i].OffsetEnd
		}
	}

	if len(c.SecondaryOccs) != 0 {
		if min != 0 {
			return fmt.Errorf("secondary occs dont start with dateStart")
		}
		end := c.DateStart.Add(time.Duration(max) * time.Minute)
		if !end.Equal(c.DateEnd) {
			return fmt.Errorf("secondary occs dont end with dateEnd")
		}
	}

	return c.OccRequest.Validate()
}

func (r *OccRequest) NewOcc(
	trainingID uuid.UUID,
) Occ {
	or := *r
	or.DateStart = or.DateStart.Round(time.Minute)
	or.DateEnd = or.DateEnd.Round(time.Minute)
	return Occ{
		OccRequest: or,
		TrainingID: trainingID,
		ID:         uuid.New(),
	}
}

type PutTrainingOccsRequest struct {
	Occurrences []CreateOccRequest
	TrainingID  uuid.UUID
}

func (r *PutTrainingOccsRequest) Validate() error {
	if r.TrainingID == uuid.Nil {
		return fmt.Errorf("Validate PostTrainingOccRequest: invalid trainingID")
	}
	for i := range r.Occurrences {
		if err := r.Occurrences[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}
