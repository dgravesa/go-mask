package mask

import (
	"fmt"
	"reflect"
	"strings"
)

func buildSimpleMaskFunc(args ...string) (Func, error) {
	var maskChar string
	switch len(args) {
	case 0:
		// if not specified, use "X"
		maskChar = "X"
	case 1:
		maskChar = args[0]
	default:
		return nil, fmt.Errorf("simple mask expects no more than 1 argument")
	}

	return func(ptr reflect.Value) error {
		val := ptr.Elem()
		if val.Kind() != reflect.String {
			return fmt.Errorf("simple mask only supports string type")
		}
		maskedStr := strings.Repeat(maskChar, val.Len())
		val.Set(reflect.ValueOf(maskedStr))
		return nil
	}, nil
}
