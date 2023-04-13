package instr

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// this is base request fields which instructor can modify whenever he wants
type VacationRequest struct {
	DateStart time.Time
	DateEnd   time.Time
}

func (vr *VacationRequest) Validate() error {
	if vr.DateStart.IsZero() {
		return fmt.Errorf("validate VacationRequest: invalid DateStart")
	}
	if vr.DateEnd.IsZero() {
		return fmt.Errorf("validate VacationRequest: invalid DateEnd")
	}
	if !vr.DateStart.Before(vr.DateEnd) {
		return fmt.Errorf("validate VacationRequest: DateStart must be before DateEnd")
	}
	return nil
}

// this is table representation
type Vacation struct {
	ID           uuid.UUID
	InstructorID uuid.UUID
	VacationRequest
}

// this gets returned to users
type VacationInfo struct {
	ID uuid.UUID
	VacationRequest
}

func (vr *VacationRequest) NewVacation(instrID uuid.UUID) *Vacation {
	return &Vacation{
		ID:              uuid.New(),
		InstructorID:    instrID,
		VacationRequest: *vr,
	}
}

type UpdateVacationRequest struct {
	VacationRequest
	ID uuid.UUID
}

func (uvr *UpdateVacationRequest) Validate() error {
	if uvr.ID == uuid.Nil {
		return fmt.Errorf("validate UpdateVacationRequest: invalid ID")
	}
	return uvr.VacationRequest.Validate()
}

type DeleteVacationRequest struct {
	ID uuid.UUID
}

func (dvr *DeleteVacationRequest) Validate() error {
	if dvr.ID == uuid.Nil {
		return fmt.Errorf("validate DeleteVacationRequest: invalid id")
	}
	return nil
}
