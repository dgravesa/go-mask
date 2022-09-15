package mask

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type simpleMasker struct {
	MaskChar         rune
	ShowFront        int
	ShowBack         int
	AlphanumericOnly bool
}

func newSimpleMaskerFromStructTag(tag string) (*simpleMasker, error) {
	var err error
	showFront := 0
	showBack := 0
	alphanumericOnly := false

	args := strings.Split(tag, ",")
	maskChar := []rune(args[0])[0]

	for _, arg := range args[1:] {
		var argName, argVal string
		argSplit := strings.SplitN(arg, "=", 2)
		if len(argSplit) == 1 {
			argName = arg
		} else {
			argName, argVal = argSplit[0], argSplit[1]
		}

		switch argName {
		case "alphanumeric":
			alphanumericOnly = true
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

	return &simpleMasker{
		MaskChar:         maskChar,
		ShowFront:        showFront,
		ShowBack:         showBack,
		AlphanumericOnly: alphanumericOnly,
	}, nil
}

func (m simpleMasker) mask(ptr reflect.Value) error {
	stringMaskFunc := makeStringMaskFunc(func(str string) (string, error) {
		return m.maskString(str), nil
	})
	return stringMaskFunc(ptr)
}

func (m simpleMasker) maskString(s string) string {
	if m.ShowFront+m.ShowBack >= len(s) {
		return s
	}

	prefix := s[:m.ShowFront]
	suffix := s[len(s)-m.ShowBack:]
	midStr := s[m.ShowFront : len(s)-m.ShowBack]

	var charMasker func(rune) rune
	if m.AlphanumericOnly {
		charMasker = func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				return m.MaskChar
			}
			return r
		}
	} else {
		charMasker = func(_ rune) rune {
			return m.MaskChar
		}
	}

	return prefix + strings.Map(charMasker, midStr) + suffix
}
