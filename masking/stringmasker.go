package masking

import (
	"fmt"
	"reflect"
)

// StringMasker is a function that masks a string input.
type StringMasker func(input string) (output string, err error)

// RegisterStringMasker registers a new custom string masker that may be used in struct tags.
func RegisterStringMasker(name string, masker StringMasker) error {
	stringMaskFunc := makeStringMaskFunc(masker)
	builder := func(args ...string) (maskFunc, error) {
		if len(args) > 0 {
			return nil, fmt.Errorf("mask func \"%s\" takes no additional arguments", name)
		}
		return stringMaskFunc, nil
	}
	return registerMaskFuncBuilder(name, builder)
}

func makeStringMaskFunc(masker StringMasker) maskFunc {
	return func(ptr reflect.Value) error {
		// get the current value
		val := ptr.Elem()
		if val.Kind() != reflect.String {
			return fmt.Errorf("mask func only supports string types")
		}

		// create the masked string
		maskedStr, err := masker(val.String())
		if err != nil {
			return err
		}

		// set new value
		newVal := reflect.ValueOf(maskedStr).Convert(val.Type())
		val.Set(newVal)

		return nil
	}
}
