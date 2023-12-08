package utils

import (
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"unsafe"
)

// UnsafeBytes returns a byte pointer without allocation.
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString returns a string pointer without allocation.
func UnsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// CopyString copies a string to make it immutable
func CopyString(s string) string {
	return string(UnsafeBytes(s))
}

// CopyBytes copies a slice to make it immutable
func CopyBytes(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// QuoteString escape special characters in a given string
func QuoteString(raw string) string {
	bb := bytebufferpool.Get()
	quoted := UnsafeString(fasthttp.AppendQuotedArg(bb.B, UnsafeBytes(raw)))
	bytebufferpool.Put(bb)

	return quoted
}
