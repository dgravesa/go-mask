package masking_test

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/dgravesa/go-mask/masking"
)

func init() {
	masking.RegisterStringMasker("sponge", spongecase)
}

type Person struct {
	Name  string
	Quote string `mask:"sponge"`
}

func ExampleRegisterStringMasker() {
	person := Person{
		Name:  "Dan",
		Quote: "I have a really great idea.",
	}

	err := masking.Apply(&person)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf(`%s: "%s"`, person.Name, person.Quote)
	// Output: Dan: "i HaVe A rEaLlY gReAt IdEa."
}

func spongecase(s string) (string, error) {
	var sb strings.Builder
	uppercase := false
	rs := []rune(s)
	for _, r := range rs {
		if uppercase {
			sb.WriteRune(unicode.ToUpper(r))
		} else {
			sb.WriteRune(unicode.ToLower(r))
		}
		if unicode.IsLetter(r) {
			uppercase = !uppercase
		}
	}
	return sb.String(), nil
}
