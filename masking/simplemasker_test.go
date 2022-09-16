package masking_test

import (
	"testing"

	"github.com/dgravesa/go-mask/masking"
	"github.com/stretchr/testify/assert"
)

func Test_MaskSimple_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type UserPass struct {
		Username string
		Password string `mask:"*"`
	}
	up := UserPass{
		Username: "John Smith",
		Password: "abcd 1234",
	}
	expectedMask := UserPass{
		Username: "John Smith",
		Password: "*********",
	}

	// act
	masking.Apply(&up)

	// assert
	assert.Equal(t, expectedMask, up)
}

func Test_MaskSimple_WithShowFront_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type AccountInfo struct {
		Name          string
		AccountNumber string `mask:"X,showfront=4"`
	}
	accountInfo := AccountInfo{
		Name:          "Test Account",
		AccountNumber: "123456789",
	}
	expectedMask := AccountInfo{
		Name:          "Test Account",
		AccountNumber: "1234XXXXX",
	}

	// act
	maskedAccountInfo := accountInfo
	masking.Apply(&maskedAccountInfo)

	// assert
	assert.Equal(t, expectedMask, maskedAccountInfo)
}

func Test_MaskSimple_WithArray_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type UserPass struct {
		Username string
		Password string `mask:"*"`
	}
	ups := []UserPass{
		{
			Username: "John Smith",
			Password: "abcd 1234",
		},
		{
			Username: "Jim Brown",
			Password: "verylongpassword123",
		},
	}
	expectedMask := []UserPass{
		{
			Username: "John Smith",
			Password: "*********",
		},
		{
			Username: "Jim Brown",
			Password: "*******************",
		},
	}

	// act
	masking.Apply(&ups)

	// assert
	assert.Equal(t, expectedMask, ups)
}

func Test_MaskSimple_WithShowBackAndAlphaNumeric_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type CreditCard struct {
		Number string `mask:"x,showback=4,alphanumeric"`
	}
	card := CreditCard{
		Number: "1234-5678-9012-3456",
	}
	expectedNumber := "xxxx-xxxx-xxxx-3456"

	// act
	masking.Apply(&card)

	// assert
	assert.Equal(t, expectedNumber, card.Number)
}

func Test_MaskSimple_WithNestedStruct_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type InnerInfo struct {
		SecretAnswer string `mask:"X,alphanumeric"`
	}
	type User struct {
		AccountNumber string `mask:"X"`
		PublicInfo    string
		Info          InnerInfo
	}
	user := User{
		PublicInfo:    "user is cool",
		AccountNumber: "12345",
		Info: InnerInfo{
			SecretAnswer: "the water is wet.",
		},
	}
	expectedMask := User{
		PublicInfo:    "user is cool",
		AccountNumber: "XXXXX",
		Info: InnerInfo{
			SecretAnswer: "XXX XXXXX XX XXX.",
		},
	}

	// act
	masking.Apply(&user)

	// assert
	assert.Equal(t, expectedMask, user)
}

func Test_MaskSimple_WithNestedStructPointer_ReturnsCorrectResult(t *testing.T) {
	// NOTE:
	// arrange
	type InnerInfo struct {
		SecretAnswer string `mask:"X,alphanumeric"`
	}
	type User struct {
		AccountNumber string `mask:"X"`
		PublicInfo    string
		Info          *InnerInfo
	}
	user := User{
		PublicInfo:    "user is cool",
		AccountNumber: "12345",
		Info: &InnerInfo{
			SecretAnswer: "the water is wet.",
		},
	}
	expectedMask := User{
		PublicInfo:    "user is cool",
		AccountNumber: "XXXXX",
		Info: &InnerInfo{
			SecretAnswer: "XXX XXXXX XX XXX.",
		},
	}

	// act
	userMasked := user
	masking.Apply(&userMasked)

	// assert
	assert.Equal(t, expectedMask, userMasked)

	// NOTE: pointers do not get deep-copied in this case. Because the original user struct
	// uses a pointer, it's the struct pointed to that gets modified. This may cause some
	// unintended behavior if not careful.
	assert.Equal(t, user.Info.SecretAnswer, "XXX XXXXX XX XXX.")
}

func Test_MaskSimple_WithStringAliasType_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type Password string
	type UserPass struct {
		Username string
		Password Password `mask:"*"`
	}
	up := UserPass{
		Username: "John Smith",
		Password: "abcd 1234",
	}
	expectedMask := UserPass{
		Username: "John Smith",
		Password: "*********",
	}

	// act
	masking.Apply(&up)

	// assert
	assert.Equal(t, expectedMask, up)
}

func Test_MaskSimple_WithNonStandardChar_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type UserInfo struct {
		PhoneNumber string `mask:"simple,#,alphanumeric,showback=4"`
	}
	ui := UserInfo{
		PhoneNumber: "(123)-456-7890",
	}
	expectedMask := "(###)-###-7890"

	// act
	masking.Apply(&ui)

	// assert
	assert.Equal(t, expectedMask, ui.PhoneNumber)
}
