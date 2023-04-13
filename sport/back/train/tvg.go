package train

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var grpl = strings.NewReplacer(
	"user_id", "UserID",
	"max_people", "MaxPeople",
	"max_trainings", "MaxTrainings")

// var rgrpl = strings.NewReplacer(
// 	"UserID", "user_id",
// 	"MaxPeople", "max_people",
// 	"MaxTrainings", "max_trainings")

type GroupRequest struct {
	Name         string
	MaxPeople    int
	MaxTrainings int
}

func (t *GroupRequest) Validate() error {
	if t.Name == "" || len(t.Name) > 256 {
		return fmt.Errorf("validate GroupRequest: invalid group")
	}
	if t.MaxPeople <= 0 {
		return fmt.Errorf("validate GroupRequest: invalid MaxPeople")
	}
	if t.MaxTrainings <= 0 {
		return fmt.Errorf("validate GroupRequest: invalid MaxTrainings")
	}
	return nil
}

type Group struct {
	GroupRequest
	ID     uuid.UUID
	UserID uuid.UUID
}

type UpdateGroupRequest struct {
	ID uuid.UUID
	GroupRequest
}

func (u *UpdateGroupRequest) Validate() error {
	if u.ID == uuid.Nil {
		return fmt.Errorf("validate UpdateTrainingGroupRequest: invalid ID")
	}
	return u.GroupRequest.Validate()
}

type DeleteGroupRequest struct {
	ID uuid.UUID
}

func (d *DeleteGroupRequest) Validate() error {
	if d.ID == uuid.Nil {
		return fmt.Errorf("validate DeleteTrainingGroupRequest: invalid ID")
	}
	return nil
}

func (r *GroupRequest) NewGroup(userID uuid.UUID) Group {
	return Group{
		ID:           uuid.New(),
		UserID:       userID,
		GroupRequest: *r,
	}
}

type GroupArr []Group

// this is used for reading json groups from db
// since they may be direwct from db this replacer is nccessary
func (g *GroupArr) Scan(value interface{}) error {
	x := value.([]byte)
	xs := grpl.Replace(string(x))
	return json.Unmarshal([]byte(xs), g)
}

// no replacer needed for inserting - inserts are always in correct form
func (g *GroupArr) Value() (driver.Value, error) {
	return json.Marshal(g)
}

type GroupBindingRequest struct {
	TrainingID uuid.UUID
	GroupID    uuid.UUID
}

func (d *GroupBindingRequest) Validate() error {
	if d.TrainingID == uuid.Nil {
		return fmt.Errorf("validate GroupBindingRequest: invalid TrainingID")
	}
	if d.GroupID == uuid.Nil {
		return fmt.Errorf("validate GroupBindingRequest: invalid GroupID")
	}
	return nil
}
