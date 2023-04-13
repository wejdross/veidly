package train

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sport/helpers"

	"github.com/google/uuid"
)

/*
	secondary occurrences are not assigned directly
	to training but to other [primary] Occurrences
	you cant register for secondary occ, instead they are used to represent trainings
	consisting of several days and multiple sessions
*/
type SecondaryOcc struct {
	ID         uuid.UUID
	TrainingID uuid.UUID
	OccID      uuid.UUID
	SecondaryOccRequest
}

// this function assumes that occ fields may be nulls in response
// and will coalesce them so scan wont fail
func NullableSecondaryOccSelectFields() string {
	return fmt.Sprintf(`
		coalesce(so.id, %s)
		,coalesce(so.occ_id, %s)
		,coalesce(so.training_id, %s)
		,coalesce(so.offset_start, 0)
		,coalesce(so.offset_end, 0)
		,coalesce(so.color, '')
		,coalesce(so.remarks, '')
	`,
		helpers.ZeroUuidString,
		helpers.ZeroUuidString,
		helpers.ZeroUuidString)
}

func (occ *SecondaryOcc) ScanFields() []interface{} {
	return []interface{}{
		&occ.ID,
		&occ.OccID,
		&occ.TrainingID,
		&occ.OffsetStart,
		&occ.OffsetEnd,
		&occ.Color,
		&occ.Remarks,
	}
}

func (tt *SecondaryOcc) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), tt)
}

func (tt *SecondaryOcc) Value() (driver.Value, error) {
	return json.Marshal(tt)
}

type SecondaryOccRequest struct {
	// in minutes
	OffsetStart int
	OffsetEnd   int

	Color   string
	Remarks string
}

func (i *SecondaryOccRequest) ValidationErr(msg string) error {
	return fmt.Errorf("Validate SecondaryOccRequest: %s", msg)
}

func (i *SecondaryOccRequest) Validate() error {
	if i.OffsetStart < 0 {
		return i.ValidationErr("invalid OffsetStart, cannot be negative")
	}
	if i.OffsetEnd < 0 {
		return i.ValidationErr("invalid OffsetEnd, cannot be negative")
	}
	if i.OffsetEnd <= i.OffsetStart {
		return i.ValidationErr("invalid Offset, End must be bigger than start")
	}
	if len(i.Remarks) > 512 {
		return i.ValidationErr("invalid remarks")
	}
	if len(i.Color) > 10 {
		return i.ValidationErr("invalid color")
	}
	return nil
}

func (r *SecondaryOccRequest) NewSecondaryOcc(occID, trainingID uuid.UUID) SecondaryOcc {
	return SecondaryOcc{
		ID:                  uuid.New(),
		OccID:               occID,
		SecondaryOccRequest: *r,
		TrainingID:          trainingID,
	}
}
