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
	case func(interface{}) error:
		mfb = createStructMaskFuncBuilder(name, func(v interface{}, _ ...string) error {
			return m(v)
		})
	case func(interface{}, ...string) error:
		mfb = createStructMaskFuncBuilder(name, m)
	default:
		return fmt.Errorf("unsupported masker signature")
	}

	return registerMaskFuncBuilder(name, mfb)
}
