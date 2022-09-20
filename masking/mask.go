package masking

import (
	"fmt"
	"reflect"
)

// Mask applies masking to public fields of v based on struct tagging.
//
// Mask will apply masking to any non-pointer and non-slice fields of the struct and its nested
// structs. In cases where pointers or slices of the struct should also be masked, use DeepMask
// instead.
func Mask(v interface{}) error {
	return mask(reflect.ValueOf(v), false)
}

// DeepMask applies masking to all public fields of v, including pointers and slices, based on struct tagging.
func DeepMask(v interface{}) error {
	return mask(reflect.ValueOf(v), true)
}

func mask(ptr reflect.Value, maskPointedVals bool) error {
	ptrKind := ptr.Kind()
	if ptrKind != reflect.Pointer && ptrKind != reflect.Interface {
		return fmt.Errorf("mask: expected pointer or interface argument")
	}

	val := ptr.Elem()

	switch val.Kind() {
	case reflect.Struct:
		t := val.Type()
		for i := 0; i < t.NumField(); i++ {
			fieldPtr, isValPointer := getPointer(val.Field(i))
			if isValPointer && !maskPointedVals {
				continue
			}

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
				err := mask(fieldPtr, maskPointedVals)
				if err != nil {
					return err
				}
			}
		}

	case reflect.Slice:
		if maskPointedVals {
			for i := 0; i < val.Len(); i++ {
				itemPtr, isValPointer := getPointer(val.Index(i))
				if isValPointer && !maskPointedVals {
					continue
				}
				err := mask(itemPtr, maskPointedVals)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// getPointer returns val and true if val is a pointer, otherwise a pointer to val and false.
func getPointer(val reflect.Value) (reflect.Value, bool) {
	if val.Kind() == reflect.Pointer {
		return val, true
	}
	return val.Addr(), false
}
