package collections

import (
	"fmt"
	"reflect"
)

type (
	TypeNotCompatible struct {
		Estimate string
		Actually reflect.Type
	}

	MethodHasNoImplement struct {
		Method string
		Type   reflect.Type
	}
)

func (tnc *TypeNotCompatible) Error() string {
	return fmt.Sprintf(
		"type '%s' is not compatible with '%s'",
		tnc.Actually.String(), tnc.Estimate,
	)
}

func (mhni *MethodHasNoImplement) Error() string {
	return fmt.Sprintf(
		"type '%s' does have any '%s' method",
		mhni.Type.String(), mhni.Method,
	)
}

func throwTypeNotCompatiable(target string, actually reflect.Type) error {
	return &TypeNotCompatible{Estimate: target, Actually: actually}
}

func throwMethodHasNoImplement(method string, t reflect.Type) error {
	return &MethodHasNoImplement{Method: method, Type: t}
}
