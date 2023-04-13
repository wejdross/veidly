package sub

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type SubModelEditable struct {
	Name         string
	MaxEntrances int
	// days
	Duration          int
	Currency          string
	MaxActive         int
	IsFreeEntrance    bool
	AllTrainingsByDef bool
}

type SubModelRequest struct {
	SubModelEditable
	Price int
}

func SmSelectColumns() string {
	return `
		sm.id, 
		sm.instr_user_id, 
		sm.instructor_id, 
		sm.name, 
		sm.max_entrances, 
		sm.duration,
		sm.price,
		sm.processing_fee,
		sm.payout_value,
		sm.refund_value,
		sm.currency,
		sm.max_active,
		sm.is_free_entrance,
		sm.all_trainings_by_def`
}

var smr = strings.NewReplacer(
	"id", "ID",
	"instr_user_id", "InstrUserID",
	"instructor_id", "InstructorID",
	"name", "Name",
	"max_entrances", "MaxEntrances",
	"duration", "Duration",
	"price", "Price",
	"processing_fee", "ProcessingFee",
	"payout_value", "PayoutValue",
	"refund_value", "RefundValue",
	"currency", "Currency",
	"max_active", "MaxActive",
	"is_free_entrance", "IsFreeEntrance",
	"all_trainings_by_def", "AllTrainingsByDef")

func (sm *SubModel) ScanFields() []interface{} {
	return []interface{}{
		&sm.ID,
		&sm.InstrUserID,
		&sm.InstructorID,
		&sm.Name,
		&sm.MaxEntrances,
		&sm.Duration,
		&sm.PricingInfo.Price,
		&sm.PricingInfo.ProcessingFee,
		&sm.PricingInfo.PayoutValue,
		&sm.PricingInfo.RefundValue,
		&sm.Currency,
		&sm.MaxActive,
		&sm.IsFreeEntrance,
		&sm.AllTrainingsByDef,
	}
}

func (sm *SubModelRequest) Validate(ctx *Ctx) error {
	if sm.Name == "" || len(sm.Name) > 256 {
		return fmt.Errorf("Validate SubModelRequest: invalid Name")
	}
	if sm.MaxEntrances < -1 || sm.MaxEntrances == 0 {
		return fmt.Errorf("Validate SubModelRequest: invalid MaxEntrances")
	}

	if sm.Duration <= 0 {
		return fmt.Errorf("Validate SubModelRequest: Invalid Duration")
	}
	if sm.Price < 0 {
		return fmt.Errorf("Validate SubModelRequest: Invalid Price")
	}
	// if err := ctx.Train.ValidateCurrency(sm.Currency); err != nil {
	// 	return fmt.Errorf("Validate SubModelRequest: %v", err)
	// }
	if sm.MaxActive < -1 {
		return fmt.Errorf("")
	}
	return nil
}

type PricingInfo struct {
	Price         int
	ProcessingFee int
	PayoutValue   int
	RefundValue   int
}

func (c *PricingInfo) Total() int {
	return c.Price + c.ProcessingFee
}

func (ctx *Ctx) NewSubPricing(p int) PricingInfo {
	pv := ((100 - ctx.Config.ServiceFee) * p) / 100
	rv := (ctx.Config.RefundAmount * p) / 100
	return PricingInfo{
		Price:         p,
		ProcessingFee: (ctx.Config.ProcessingFee * p) / 100,
		PayoutValue:   pv,
		RefundValue:   rv,
	}
}

type SubModel struct {
	ID           uuid.UUID
	InstructorID uuid.UUID
	InstrUserID  uuid.UUID
	PricingInfo
	SubModelEditable
}

func (c *SubModel) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *SubModel) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), c)
}

func (sm *SubModelRequest) NewSubModel(ctx *Ctx, userID, instructorID uuid.UUID) SubModel {
	return SubModel{
		ID:               uuid.New(),
		InstructorID:     instructorID,
		InstrUserID:      userID,
		SubModelEditable: sm.SubModelEditable,
		PricingInfo:      ctx.NewSubPricing(sm.Price),
	}
}

type UpdateSubModelRequest struct {
	ID uuid.UUID
	SubModelRequest
}

func (sm *UpdateSubModelRequest) Validate(ctx *Ctx) error {
	if err := sm.SubModelRequest.Validate(ctx); err != nil {
		return err
	}
	if sm.ID == uuid.Nil {
		return fmt.Errorf("Validate UpdateSubModelRequest: invalid ID")
	}
	return nil
}

type DeleteSubModelRequest struct {
	ID uuid.UUID
}

func (sm *DeleteSubModelRequest) Validate() error {
	if sm.ID == uuid.Nil {
		return fmt.Errorf("Validate DeleteSubModelRequest: invalid ID")
	}
	return nil
}

type SmArr []SubModel

func (g *SmArr) Scan(value interface{}) error {
	x := value.([]byte)
	xs := smr.Replace(string(x))
	return json.Unmarshal([]byte(xs), g)
}
