package masking

import (
	"fmt"
	"reflect"
)

type Masker interface {
	func(input string) (output string) |
		func(input string) (output string, err error) |
		func(input string, args ...string) (output string, err error) |
		func(s *string) |
		func(s *string) (err error) |
		func(s *string, args ...string) (err error) |
		func(v interface{}) (err error) |
		func(v interface{}, args ...string) (err error)
}

// RegisterMasker registers a new masker function for use in struct tagging.
func RegisterMasker[M Masker](name string, masker M) error {
	// maskerType := reflect.TypeOf(masker)
	var mfb maskFuncBuilder

	switch m := reflect.ValueOf(masker).Interface().(type) {
	case func(string) string:
		mfb = createStringMaskFuncBuilder(name, func(s *string, _ ...string) error {
			*s = m(*s)
			return nil
		})
	case func(string) (string, error):
		mfb = createStringMaskFuncBuilder(name, func(s *string, _ ...string) error {
			var err error
			*s, err = m(*s)
			return err
		})
	case func(string, ...string) (string, error):
		mfb = createStringMaskFuncBuilder(name, func(s *string, args ...string) error {
			var err error
			*s, err = m(*s, args...)
			return err
		})
	case func(*string):
		mfb = createStringMaskFuncBuilder(name, func(s *string, _ ...string) error {
			m(s)
			return nil
		})
	case func(*string) error:
		mfb = createStringMaskFuncBuilder(name, func(s *string, _ ...string) error {
			return m(s)
		})
	case func(*string, ...string) error:
		mfb = createStringMaskFuncBuilder(name, m)
	default:
		return fmt.Errorf("unsupported masker signature")
	}

	return registerMaskFuncBuilder(name, mfb)
}

func createStringMaskFuncBuilder(name string, masker func(*string, ...string) error) maskFuncBuilder {
	return func(args ...string) (maskFunc, error) {
		return func(ptr reflect.Value) error {
			// get the current value
			val := ptr.Elem()
			if val.Kind() != reflect.String {
				return fmt.Errorf("%s: mask func only supports string types", name)
			}

			// mask the string in place
			return masker(val.Interface().(*string))
		}, nil
	}
}
