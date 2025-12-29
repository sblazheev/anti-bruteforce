//revive:disable
package common

import "fmt"

var (
	ErrQueryRequest          = fmt.Errorf("error query")
	ErrServiceUnavailable    = fmt.Errorf("service unavailable")
	ErrIPSubnetAlreadyExists = fmt.Errorf("subnet exist")
	ErrIPSubnetNotFound      = fmt.Errorf("subnet not found")
	ErrIPSubnetOverlapped    = fmt.Errorf("subnet overlapped")
	ErrStorageUnknownType    = fmt.Errorf("storage unknown type")
	ErrFormatIp              = fmt.Errorf("error format ip4")
	ErrDuplicateValue        = fmt.Errorf("duplicate value")
)
