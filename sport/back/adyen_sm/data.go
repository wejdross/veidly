package adyen_sm

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type State string

const (
	// payment intialized and obj created
	Link State = "link"
	// payment initialized, obj created and capture request sent
	LinkExpress State = "link_express"
	// fund put on hold on clients account
	Hold State = "hold"
	// processing order cancellation or refund
	WaitCancelOrRefund State = "wait_cancel_or_refund"
	// processed cancel or refund
	// THIS IS END STATE
	CancelOrRefund State = "cancel_or_refund"
	// currently retrying cancel or refund requests due to errors
	RetryCancelOrRefund State = "retry_cancel_or_refund"
	// made capture request
	WaitCapture State = "wait_capture"
	// capture has been made and confirmed
	Capture State = "capture"
	// move to pay vendor for his service
	WaitPayout State = "wait_payout"
	// payout has been made and confirmed
	Payout State = "payout"
	// something unexpected happened during processing.
	Error State = "error"
	// link expired
	LinkExpire State = "link_expire"
	// refund requested
	WaitRefund State = "wait_refund"
	// refund confirmation received
	Refund State = "refund"
	// dispute between 2 sides
	Dispute State = "dispute"
)

type SmCacheEntry struct {
	Retries int
}

type SmCache map[State]SmCacheEntry

// used to convert golang struct into db composite type representation
func (c *SmCache) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// used to convert DB composite type representation into golang struct
func (c *SmCache) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), c)
}

type InstructorDecision string

const (
	Unset   InstructorDecision = "unset"
	Approve InstructorDecision = "approve"
	Reject  InstructorDecision = "reject"
)

type SmData struct {
	ID           uuid.UUID
	InstructorID uuid.UUID
	//
	State       State
	LinkID      string
	SmRetries   int
	OrderID     string
	RefID       string
	SmCache     SmCache
	IsConfirmed bool
	SmTimeout   time.Time
	//
	ManualConfirm bool
	Decision      InstructorDecision
	//
	TotalPrice  int
	InstrPrice  int
	PayoutValue int
	RefundValue int
	Currency    string
	//
	ShopperReference string
	//
	SmType string
}

//
type Pass struct {

	// custom argument
	Args interface{}

	// database data
	Data SmData
}
