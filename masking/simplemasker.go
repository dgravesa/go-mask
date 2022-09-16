package masking

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

func simpleMaskFuncBuilderWithChar(maskChar rune) maskFuncBuilder {
	return func(args ...string) (maskFunc, error) {
		masker := &simpleMasker{
			MaskChar: maskChar,
		}
		if err := masker.updateFromArgs(args...); err != nil {
			return nil, err
		}
		return masker.mask, nil
	}
}

func simpleMaskFuncBuilder() maskFuncBuilder {
	return func(args ...string) (maskFunc, error) {
		if len(args) == 0 {
			return simpleMaskFuncBuilderWithChar('*')()
		}

		// set mask character from first argument
		if len(args[0]) != 1 {
			return nil, fmt.Errorf("first argument to simple mask but be a single character")
		}
		maskChar := []rune(args[0])[0]
		// build mask func with remaining arguments
		return simpleMaskFuncBuilderWithChar(maskChar)(args[1:]...)
	}
}

func (m *simpleMasker) updateFromArgs(args ...string) error {
	for _, arg := range args {
		var argName, argVal string
		argSplit := strings.SplitN(arg, "=", 2)
		if len(argSplit) == 1 {
			argName = arg
		} else {
			argName, argVal = argSplit[0], argSplit[1]
		}

		var err error
		switch argName {
		case "alphanumeric":
			m.AlphanumericOnly = true
			if argVal != "" {
				return fmt.Errorf("alphanumeric specifier does not take an argument")
			}
		case "showfront":
			m.ShowFront, err = strconv.Atoi(argVal)
			if err != nil {
				return fmt.Errorf("unable to parse showfront value")
			}
		case "showback":
			m.ShowBack, err = strconv.Atoi(argVal)
			if err != nil {
				return fmt.Errorf("unable to parse showback value")
			}
		}
	}

	return nil
}

func (m simpleMasker) mask(ptr reflect.Value) error {
	return makeStringMaskFunc(m.maskString)(ptr)
}

func (m simpleMasker) maskString(s string) (string, error) {
	if m.ShowFront+m.ShowBack >= len(s) {
		return s, nil
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

	return prefix + strings.Map(charMasker, midStr) + suffix, nil
}
