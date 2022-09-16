package masking

import (
	"fmt"
	"reflect"
	"strings"
)

var maskFuncBuilderRegistry = map[string]maskFuncBuilder{}

type maskFunc func(ptr reflect.Value) error

type maskFuncBuilder func(args ...string) (maskFunc, error)

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

	if len(funcName) == 1 {
		// special case, perform simple masking using the first character as the mask character
		simpleMasker, err := newSimpleMaskerFromStructTag(tag)
		if err != nil {
			return nil, err
		}
		return simpleMasker.mask, nil
	}

	builder, found := maskFuncBuilderRegistry[funcName]
	if !found {
		return nil, fmt.Errorf("unrecognized mask func: \"%s\"", funcName)
	}

	return builder(args[1:]...)
}
