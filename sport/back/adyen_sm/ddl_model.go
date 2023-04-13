package adyen_sm

/*
	columns which must be present in the table for this module to work
*/
type RequiredColumnNames struct {
	InstructorDecision string
	InstructorID       string
	ID                 string
	State              string
	SmCache            string
	SmRetries          string
	SmTimeout          string
	IsConfirmed        string
	IsActive           string
	OrderID            string
	MethodID           string
}

var DefaultColumns = RequiredColumnNames{
	InstructorDecision: "instructor_decision",
	ID:                 "id",
	InstructorID:       "instructor_id",
	State:              "state",
	SmCache:            "sm_cache",
	SmRetries:          "sm_retries",
	SmTimeout:          "sm_timeout",
	IsConfirmed:        "is_confirmed",
	IsActive:           "is_active",
	OrderID:            "order_id",
	MethodID:           "method_id",
}

type DDLModel struct {
	TableName string
	Columns   RequiredColumnNames
}
