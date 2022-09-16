package masking

import (
	"fmt"
	"reflect"
)

// StructMasker is a function that masks a struct.
// The type v will always be a pointer.
type StructMasker func(v interface{}) (err error)

// RegisterStructMasker registers a new custom struct masker that may be used in struct tags.
func RegisterStructMasker(name string, masker StructMasker) error {
	structMaskFunc := makeStructMaskFunc(masker)
	builder := func(args ...string) (maskFunc, error) {
		if len(args) > 0 {
			return nil, fmt.Errorf("mask func \"%s\" takes no additional arguments", name)
		}
		return structMaskFunc, nil
	}
	return registerMaskFuncBuilder(name, builder)
}

func makeStructMaskFunc(masker StructMasker) maskFunc {
	return func(ptr reflect.Value) error {
		// mask the struct in place
		return masker(ptr.Interface())
	}
}
