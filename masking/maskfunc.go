package masking

import (
	"fmt"
	"reflect"
	"strings"
)

type maskFunc func(ptr reflect.Value) error

type maskFuncBuilder func(args ...string) maskFunc

var maskFuncBuilderRegistry = map[string]maskFuncBuilder{
	"X":      simpleMaskFuncBuilderWithRune('X'),
	"x":      simpleMaskFuncBuilderWithRune('x'),
	"*":      simpleMaskFuncBuilderWithRune('*'),
	"-":      simpleMaskFuncBuilderWithRune('-'),
	"_":      simpleMaskFuncBuilderWithRune('_'),
	".":      simpleMaskFuncBuilderWithRune('.'),
	"simple": simpleMaskFuncBuilder(),
}

func getMaskFunc(tag string) (maskFunc, error) {
	args := strings.Split(tag, ",")
	funcName := args[0]

	builder, found := maskFuncBuilderRegistry[funcName]
	if !found {
		return nil, fmt.Errorf("unrecognized mask func: \"%s\"", funcName)
	}

	return builder(args[1:]...), nil
}

func registerMaskFuncBuilder(name string, builder maskFuncBuilder) error {
	if strings.Contains(name, ",") {
		return fmt.Errorf("commas not permitted in mask func names")
	}

	_, found := maskFuncBuilderRegistry[name]
	if found {
		return fmt.Errorf("mask func with name already exists: \"%s\"", name)
	}

	maskFuncBuilderRegistry[name] = builder
	return nil
}

func createStringMaskFuncBuilder(name string, masker func(*string, ...string) error) maskFuncBuilder {
	return func(args ...string) maskFunc {
		return func(ptr reflect.Value) error {
			// get the current value
			val := ptr.Elem()
			if val.Kind() != reflect.String {
				return fmt.Errorf("%s: mask func only supports string types", name)
			}

			// mask the string in place
			var sptr *string // convert to *string to enable type aliases
			sptr = ptr.Convert(reflect.TypeOf(sptr)).Interface().(*string)
			return masker(sptr, args...)
		}
	}
}

func createStructMaskFuncBuilder(name string, masker func(interface{}, ...string) error) maskFuncBuilder {
	return func(args ...string) maskFunc {
		return func(ptr reflect.Value) error {
			return masker(ptr.Interface(), args...)
		}
	}
}
