package mask

import (
	"fmt"
	"strings"
)

// FuncBuilder is used to create a mask.Func given arguments specified in a struct tag.
type FuncBuilder func(args ...string) (Func, error)

var funcRegistry = map[string]FuncBuilder{
	"simple": buildSimpleMaskFunc,
}

func getMaskFunc(tag string) (Func, error) {
	args := strings.Split(tag, ",")
	funcName := args[0]

	builder, ok := funcRegistry[funcName]
	if !ok {
		return nil, fmt.Errorf("unrecognized mask func: \"%s\"", funcName)
	}

	return builder(args[1:]...)
}
