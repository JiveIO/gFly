package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// Sha256 hash from list of arguments
func Sha256(args ...any) string {
	var strSlice []string
	for i := 0; i < len(args); i++ {
		strSlice = append(strSlice, "%v")
	}
	// Create a new SHA256 hash.
	hash := sha256.New()
	code := fmt.Sprintf(strings.Join(strSlice, "-"), args)
	hash.Write([]byte(code))

	return hex.EncodeToString(hash.Sum(nil))
}
