//revive:disable
package common

import "fmt"

var (
	ErrQueryRequest          = fmt.Errorf("error query")
	ErrServiceUnavailable    = fmt.Errorf("service unavailable")
	ErrEventNotFound         = fmt.Errorf("event not found")
	ErrIPSubnetAlreadyExists = fmt.Errorf("subnet exist")
	ErrIPSubnetNotFound      = fmt.Errorf("subnet not found")
)
