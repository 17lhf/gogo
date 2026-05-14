package model

// LogStatus represents the status of an operation log entry.
type LogStatus int16

const (
	LogStatusSuccess LogStatus = 1
	LogStatusFailure LogStatus = 2
)

func (s LogStatus) String() string {
	switch s {
	case LogStatusSuccess:
		return "success"
	case LogStatusFailure:
		return "failure"
	default:
		return "unknown"
	}
}
