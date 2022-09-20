package masking

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func simpleMaskFuncBuilder() maskFuncBuilder {
	return createStringMaskFuncBuilder("simple", func(s *string, args ...string) error {
		if len(args) < 1 {
			return simpleMaskerWithRune('X')(s)
		}

		// take first argument as mask character
		runeArg := args[0]
		if len(runeArg) != 1 {
			return fmt.Errorf("first argument to simple mask must be a single character")
		}
		maskChar := []rune(runeArg)[0]

		return simpleMaskerWithRune(maskChar)(s, args[1:]...)
	})
}

func simpleMaskFuncBuilderWithRune(maskChar rune) maskFuncBuilder {
	return createStringMaskFuncBuilder(string(maskChar), simpleMaskerWithRune(maskChar))
}

func simpleMaskerWithRune(maskChar rune) func(s *string, args ...string) error {
	return func(s *string, args ...string) error {
		return maskSimpleWithArgs(s, maskChar, args...)
	}
}

func maskSimpleWithArgs(s *string, maskChar rune, args ...string) error {
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

		var err error
		switch argName {
		case "alphanumeric":
			alnumOnly = true
			if argVal != "" {
				return fmt.Errorf("alphanumeric specifier does not take an argument")
			}
		case "showfront":
			showFront, err = strconv.Atoi(argVal)
			if err != nil {
				return fmt.Errorf("unable to parse showfront value")
			}
		case "showback":
			showBack, err = strconv.Atoi(argVal)
			if err != nil {
				return fmt.Errorf("unable to parse showback value")
			}
		}
	}

	maskSimple(s, maskChar, showFront, showBack, alnumOnly)
	return nil
}

func maskSimple(s *string, maskChar rune, showFront, showBack int, alnumOnly bool) {
	var charMasker func(rune) rune
	if alnumOnly {
		charMasker = func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				return maskChar
			}
			return r
		}
	} else {
		charMasker = func(_ rune) rune {
			return maskChar
		}
	}

	oldS := *s
	lenS := len(oldS)

	prefix := oldS[0:showFront]
	suffix := oldS[lenS-showBack:]
	midMasked := strings.Map(charMasker, oldS[showFront:lenS-showBack])

	// construct new masked string
	newS := prefix + midMasked + suffix

	*s = newS
}
