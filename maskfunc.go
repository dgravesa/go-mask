package mask

import (
	"fmt"
	"reflect"
	"strings"
)

type maskFunc func(ptr reflect.Value) error

func getMaskFunc(tag string) (maskFunc, error) {
	args := strings.Split(tag, ",")
	funcName := args[0]

	if len(funcName) == 1 {
		// special case, perform simple masking using the first character as the mask character
		return buildSimpleMaskFunc([]rune(funcName)[0], args[1:]...)
	}

	// TODO: implement other maskers, including ability to create custom maskers
	return nil, fmt.Errorf("unrecognized mask func: \"%s\"", funcName)
}
