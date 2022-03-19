package reflectutil

import (
	"reflect"
)

func IndirectValueOf(o interface{}) *reflect.Value {
	v := reflect.Indirect(reflect.ValueOf(o))
	return &v
}
func IsSlice(o interface{}) bool {
	theValue := reflect.Indirect(reflect.ValueOf(o))
	if theValue.Kind() != reflect.Slice {
		return false
	}
	return true

}

func IsArrayOrSlice(o interface{}) bool {
	theValue := reflect.Indirect(reflect.ValueOf(o))
	if theValue.Kind() != reflect.Slice && theValue.Kind() != reflect.Array {
		return false
	}
	return true
}
