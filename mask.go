package mask

import (
	"fmt"
	"reflect"
)

// Func is used for masking a value.
type Func func(ptr reflect.Value) error

// Apply applies masking based on struct tagging.
func Apply(v interface{}) error {
	ptr := reflect.ValueOf(v)
	if ptr.Kind() != reflect.Pointer {
		return fmt.Errorf("mask: expected pointer argument")
	}
	return mask(ptr)
}

func mask(ptr reflect.Value) error {
	val := ptr.Elem()

	switch val.Kind() {
	case reflect.Struct:
		t := val.Type()
		for i := 0; i < t.NumField(); i++ {
			maskFunc := mask // default to recursion

			// apply masking if tag is specified
			if maskTag := t.Field(i).Tag.Get("mask"); maskTag != "" {
				var err error
				maskFunc, err = getMaskFunc(maskTag)
				if err != nil {
					return err
				}
			}

			fieldPtr := val.Field(i)
			if fieldPtr.Kind() != reflect.Pointer {
				// if field is not a pointer, then use its address
				fieldPtr = fieldPtr.Addr()
			}

			err := maskFunc(fieldPtr)
			if err != nil {
				return err
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			itemPtr := val.Index(i)
			if itemPtr.Kind() != reflect.Pointer {
				// if item is not a pointer, then use its address
				itemPtr = itemPtr.Addr()
			}
			err := mask(itemPtr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
