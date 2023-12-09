package utils

import (
	"crypto/rand"
	"math/big"
)

// RandInt64 Generate a random number Int64 in rang [0, max)
//
// NOTE:
//
//	Get error `G404 (CWE-338): Use of weak random number generator (math/rand instead of crypto/rand)
//	(Confidence: MEDIUM, Severity: HIGH)` when use `rand.Intn(max)` from "math/rand".
//	Fixed Refer https://github.com/securego/gosec/issues/294#issuecomment-487452731
func RandInt64(max int64) int64 {
	// Get random Int64
	n, _ := rand.Int(rand.Reader, big.NewInt(max))

	return n.Int64()
}
