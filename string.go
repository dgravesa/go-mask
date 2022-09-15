package mask

import (
	"fmt"
	"reflect"
)

type stringMasker func(string) (string, error)

func makeStringMaskFunc(masker stringMasker) maskFunc {
	return func(ptr reflect.Value) error {
		// get the current value
		val := ptr.Elem()
		if val.Kind() != reflect.String {
			return fmt.Errorf("mask only supports string types")
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
