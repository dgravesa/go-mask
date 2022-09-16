package masking

import (
	"fmt"
	"reflect"
)

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
			fieldPtr := getPointer(val.Field(i))

			if fieldMaskTag := t.Field(i).Tag.Get("mask"); fieldMaskTag != "" {
				// apply masking if tag is specified
				maskFieldFunc, err := getMaskFunc(fieldMaskTag)
				if err != nil {
					return err
				}
				err = maskFieldFunc(fieldPtr)
				if err != nil {
					return err
				}
			} else {
				// perform masking recursively
				err := mask(fieldPtr)
				if err != nil {
					return err
				}
			}
		}

	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			itemPtr := getPointer(val.Index(i))
			err := mask(itemPtr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getPointer(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Pointer {
		return val
	}
	return val.Addr()
}
