package fluentsql

import (
	"fmt"
	"strings"
)

// joinSlice Join a slice with separator
func joinSlice[T any](values []T, separator string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(values)), fmt.Sprintf("%s ", separator)), "[]")
}
