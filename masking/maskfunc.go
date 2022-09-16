package masking

import (
	"fmt"
	"reflect"
	"strings"
)

type maskFunc func(ptr reflect.Value) error

type maskFuncBuilder func(args ...string) (maskFunc, error)

var maskFuncBuilderRegistry = map[string]maskFuncBuilder{
	"X":    simpleMaskFuncBuilderWithChar('X'),
	"x":    simpleMaskFuncBuilderWithChar('x'),
	"*":    simpleMaskFuncBuilderWithChar('*'),
	"-":    simpleMaskFuncBuilderWithChar('-'),
	"_":    simpleMaskFuncBuilderWithChar('_'),
	".":    simpleMaskFuncBuilderWithChar('.'),
	"char": simpleMaskFuncBuilder(),
}

func registerMaskFuncBuilder(funcName string, builder maskFuncBuilder) error {
	if strings.Contains(funcName, ",") {
		return fmt.Errorf("commas not permitted in mask func names")
	}

	_, found := maskFuncBuilderRegistry[funcName]
	if found {
		return fmt.Errorf("mask func with name already exists: \"%s\"", funcName)
	}

	maskFuncBuilderRegistry[funcName] = builder
	return nil
}

func getMaskFunc(tag string) (maskFunc, error) {
	args := strings.Split(tag, ",")
	funcName := args[0]

	builder, found := maskFuncBuilderRegistry[funcName]
	if !found {
		return nil, fmt.Errorf("unrecognized mask func: \"%s\"", funcName)
	}

	return builder(args[1:]...)
}
