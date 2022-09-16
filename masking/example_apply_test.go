package masking_test

import (
	"fmt"

	"github.com/dgravesa/go-mask/masking"
)

type UserAccount struct {
	Username       string
	Password       string `mask:"*"`
	AccountNumber  string `mask:"X,showback=4"`
	ActivationCode string `mask:"X,showfront=6,alphanumeric"`
}

func ExampleApply() {
	account := UserAccount{
		Username:       "John Smith",
		Password:       "thisisthepassword",
		AccountNumber:  "1234567890",
		ActivationCode: "ab13ea-12cb55fab125-3f3b97",
	}

	masking.Apply(&account)

	fmt.Printf("%s, %s, %s, %s\n",
		account.Username,
		account.Password,
		account.AccountNumber,
		account.ActivationCode)
	// Output: John Smith, *****************, XXXXXX7890, ab13ea-XXXXXXXXXXXX-XXXXXX
}
