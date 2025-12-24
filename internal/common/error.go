//revive:disable
package common

import "fmt"

var (
	ErrQueryRequest         = fmt.Errorf("error query")
	ErrServiceUnavailable   = fmt.Errorf("service unavailable")
	ErrEventNotFound        = fmt.Errorf("event not found")
	ErrEventDateBusy        = fmt.Errorf("the selected time is already busy")
	ErrEventInvalidEvent    = fmt.Errorf("invalid event data")
	ErrEventAlreadyExists   = fmt.Errorf("event already exists")
	ErrEventConflictOverlap = fmt.Errorf("event overlaps with another event")
	ErrStorageUnknownType   = fmt.Errorf("storage unknown type")
)
