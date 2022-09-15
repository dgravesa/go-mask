package mask

import (
	"fmt"
	"strings"
)

// FuncBuilder is used to create a mask.Func given arguments specified in a struct tag.
type FuncBuilder func(args ...string) (Func, error)

// RegisterCustomMasker registers a new custom masking function which can be used via struct tags.
func RegisterCustomMasker(funcName string, builder FuncBuilder) error {
	_, found := funcRegistry[funcName]
	if found {
		return fmt.Errorf("mask func name already taken: \"%s\"", funcName)
	}
	funcRegistry[funcName] = builder
	return nil
}

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
