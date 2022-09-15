package mask

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func alnumToMaskCharMapper(maskChar rune) func(rune) rune {
	return func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return maskChar
		}
		return r
	}
}

func allToMaskCharMapper(maskChar rune) func(rune) rune {
	return func(_ rune) rune {
		return maskChar
	}
}

func buildSimpleMaskFunc(maskChar rune, args ...string) (maskFunc, error) {
	var err error
	showFront := 0
	showBack := 0
	alnumOnly := false

	for _, arg := range args {
		var argName, argVal string
		argSplit := strings.SplitN(arg, "=", 2)
		if len(argSplit) == 1 {
			argName = arg
		} else {
			argName, argVal = argSplit[0], argSplit[1]
		}

		switch argName {
		case "alphanumeric":
			alnumOnly = true
			if argVal != "" {
				return nil, fmt.Errorf("alphanumeric specifier does not take an argument")
			}
		case "showfront":
			showFront, err = strconv.Atoi(argVal)
			if err != nil {
				return nil, fmt.Errorf("unable to parse showfront value")
			}
		case "showback":
			showBack, err = strconv.Atoi(argVal)
			if err != nil {
				return nil, fmt.Errorf("unable to parse showback value")
			}
		}
	}

	return func(ptr reflect.Value) error {
		val := ptr.Elem()
		if val.Kind() != reflect.String {
			return fmt.Errorf("mask only supports string types")
		}

		// construct masked string
		maskedStr := maskStringSimple(val.String(), maskChar, showFront, showBack, alnumOnly)
		newVal := reflect.ValueOf(maskedStr).Convert(val.Type())
		val.Set(newVal)

		return nil
	}, nil
}

func maskStringSimple(s string, maskChar rune, showFront, showBack int, alnumOnly bool) string {
	if showFront+showBack >= len(s) {
		return s
	}

	prefix := s[:showFront]
	suffix := s[len(s)-showBack:]
	midStr := s[showFront : len(s)-showBack]
	if alnumOnly {
		midStr = strings.Map(alnumToMaskCharMapper(maskChar), midStr)
	} else {
		midStr = strings.Map(allToMaskCharMapper(maskChar), midStr)
	}
	return prefix + midStr + suffix
}
