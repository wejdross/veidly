package dc

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sport/helpers"
	"strings"
	"time"

	"github.com/google/uuid"
)

/*
	golang cant parse dates if they dont contain [Z]one fragment.
	so i have to implement my own type.
	thanks golang
*/
type Time time.Time

func (j *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05", s)
		if err != nil {
			return err
		}
	}
	*j = Time(t)
	return nil
}

func (j Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

func (j Time) Format(s string) string {
	return j.Format(s)
}

type DcRequest struct {
	Quantity   int
	ValidStart Time
	ValidEnd   Time
	Name       string
	Discount   int
}

func (r *DcRequest) Validate() error {
	if r.Quantity < -1 {
		return fmt.Errorf("Validate DcRequest: invalid quantity")
	}
	if time.Time(r.ValidStart).IsZero() || time.Time(r.ValidStart).After(time.Time(r.ValidEnd)) {
		return fmt.Errorf("Validate DcRequest: invalid ValidStart")
	}
	if time.Time(r.ValidEnd).IsZero() {
		return fmt.Errorf("Validate DcRequest: invalid ValidEnd")
	}
	if r.Name == "" {
		return fmt.Errorf("Validate DcRequest: invalid Name")
	}
	if r.Discount <= 0 || r.Discount > 100 {
		return fmt.Errorf("Validate DcRequest: invalid Discount")
	}
	return nil
}

type Dc struct {
	ID, InstrID      uuid.UUID
	RedeemedQuantity int
	DcRequest
}

func (dc *Dc) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), dc)
}

func (dc *Dc) Value() (driver.Value, error) {
	return json.Marshal(dc)
}

func DcSelectColumns() string {
	return `dc.id, 
	dc.instr_id, 
	dc.redeemed_quantity, 
	dc.quantity, 
	dc.valid_start, 
	dc.valid_end,
	dc.name,
	dc.discount`
}

var dcl = strings.NewReplacer(
	"id", "ID",
	"instr_id", "InstrID",
	"redeemed_quantity", "RedeemedQuantity",
	"quantity", "Quantity",
	"valid_start", "ValidStart",
	"valid_end", "ValidEnd")

func DcNUllableSelectColumns() string {
	return fmt.Sprintf(`
		coalesce(dc.id, %s), 
		coalesce(dc.instr_id, %s), 
		coalesce(dc.redeemed_quantity, 0), 
		coalesce(dc.quantity, 0), 
		coalesce(dc.valid_start, %s), 
		coalesce(dc.valid_end, %s),
		coalesce(dc.name, ''),
		coalesce(dc.discount, '')`,
		helpers.ZeroUuidString,
		helpers.ZeroUuidString,
		helpers.ZeroTimeString,
		helpers.ZeroTimeString)
}

func (dc *Dc) ScanFields() []interface{} {
	return []interface{}{
		&dc.ID,
		&dc.InstrID,
		&dc.RedeemedQuantity,
		&dc.Quantity,
		(*time.Time)(&dc.ValidStart),
		(*time.Time)(&dc.ValidEnd),
		&dc.Name,
		&dc.Discount,
	}
}

func (r *DcRequest) NewDc(instrID uuid.UUID) *Dc {
	return &Dc{
		DcRequest:        *r,
		ID:               uuid.New(),
		InstrID:          instrID,
		RedeemedQuantity: 0,
	}
}

type UpdateDcRequest struct {
	DcRequest
	ID uuid.UUID
}

func (r *UpdateDcRequest) Validate() error {
	if r.ID == uuid.Nil {
		return fmt.Errorf("Validate UpdateDcRequest: invalid ID")
	}
	return r.DcRequest.Validate()
}

type DeleteDcRequest struct {
	ID uuid.UUID
}

func (r *DeleteDcRequest) Validate() error {
	if r.ID == uuid.Nil {
		return fmt.Errorf("Validate DeleteDcRequest: invalid ID")
	}
	return nil
}

type DcBindingRequest struct {
	TrainingID uuid.UUID
	DcID       uuid.UUID
}

func (r *DcBindingRequest) Validate() error {
	if r.TrainingID == uuid.Nil {
		return fmt.Errorf("Validate DcBinding: invalid TrainingID")
	}
	if r.DcID == uuid.Nil {
		return fmt.Errorf("Validate DcBinding: invalid DcID")
	}
	return nil
}

type DcArr []Dc

// this is used for reading json groups from db
// since they may be direwct from db this replacer is nccessary
func (g *DcArr) Scan(value interface{}) error {
	x := value.([]byte)
	xs := dcl.Replace(string(x))
	return json.Unmarshal([]byte(xs), g)
}
