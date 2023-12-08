package errors

import (
	"fmt"
	"strings"
)

// InvalidParam is used when an invalid parameter is passed in the request
type InvalidParam struct {
	Param []string
}

// Error returns an error message regarding incorrect parameter values for a request
func (e InvalidParam) Error() string {
	switch {
	case len(e.Param) > 1:
		return fmt.Sprintf("Incorrect value for parameters: " + strings.Join(e.Param, ", "))
	case len(e.Param) == 1:
		return fmt.Sprintf("Incorrect value for parameter: " + e.Param[0])
	default:
		return "This request has invalid parameters"
	}
}
