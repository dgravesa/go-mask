package mask_test

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/dgravesa/go-mask"
)

func SpongeCase(s string) (string, error) {
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

func init() {
	mask.RegisterStringMasker("sponge", SpongeCase)
}

type Person struct {
	Name  string
	Quote string `mask:"sponge"`
}

func ExampleRegisterStringMaster() {
	person := Person{
		Name:  "Dan",
		Quote: "I have a really great idea.",
	}

	err := mask.Apply(&person)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf(`%s: "%s"`, person.Name, person.Quote)
	// Output: Dan: "i HaVe A rEaLlY gReAt IdEa."
}
