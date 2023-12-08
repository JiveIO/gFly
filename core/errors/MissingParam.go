package errors

import (
	"fmt"
	"strings"
)

// MissingParam is used when a required parameter is not present in the request
type MissingParam struct {
	Param []string `json:"param"`
}

// Error returns an error message regarding required parameters for a request
func (e MissingParam) Error() string {
	switch {
	case len(e.Param) > 1:
		return fmt.Sprintf("Parameters " + strings.Join(e.Param, ", ") + " are required for this request")
	case len(e.Param) == 1:
		return fmt.Sprintf("Parameter " + e.Param[0] + " is required for this request")
	default:
		return "This request is missing parameters"
	}
}
