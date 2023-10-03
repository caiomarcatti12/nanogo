package util

import "reflect"

func ZeroReflectGeneric[T any](obj T) T {
	return reflect.Zero(reflect.TypeOf(obj)).Interface().(T)
}
