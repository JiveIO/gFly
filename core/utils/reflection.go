package utils

import "reflect"

// ReflectType Get the name of a struct instance.
func ReflectType(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
