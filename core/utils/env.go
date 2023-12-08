package utils

import (
	"os"
	"reflect"
	"strconv"
)

// Getenv Get environment parameter (From .env file).
func Getenv[V any](key string, init V) V {
	// Create out any type and assign initial value
	var out any = init

	// Try to get parameter
	if value, ok := os.LookupEnv(key); ok {
		switch reflect.TypeOf(init).Name() {
		case "string":
			out = value
		case "int":
			if num, err := strconv.Atoi(value); err == nil {
				out = num
			}
		case "float64":
			if num, err := strconv.ParseFloat(value, 64); err == nil {
				out = num
			}
		case "bool":
			if boolean, err := strconv.ParseBool(value); err == nil {
				out = boolean
			}
		case "float32":
			if num, err := strconv.ParseFloat(value, 32); err == nil {
				out = num
			}
		}
	}

	// Important! Type assertion
	return out.(V)
}
