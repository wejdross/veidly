package sub

import (
	"fmt"

	"github.com/google/uuid"
)

type SubModelBinding struct {
	SubModelID uuid.UUID
	TrainingID uuid.UUID
}

func (smb *SubModelBinding) Validate() error {
	if smb.SubModelID == uuid.Nil {
		return fmt.Errorf("Validate SubModelBinding: invalid SubModelID")
	}
	if smb.TrainingID == uuid.Nil {
		return fmt.Errorf("Validate SubModelBinding: invalid TrainingID")
	}
	return nil
}
