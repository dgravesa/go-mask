package mask_test

import (
	"testing"

	"github.com/dgravesa/go-mask"
	"github.com/stretchr/testify/assert"
)

func Test_MaskSimple_WithNoChar_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type AccountInfo struct {
		Name          string
		AccountNumber string `mask:"simple"`
	}
	accountInfo := AccountInfo{
		Name:          "Test Account",
		AccountNumber: "12345678",
	}
	expectedMask := AccountInfo{
		Name:          "Test Account",
		AccountNumber: "XXXXXXXX",
	}

	// act
	maskedAccountInfo := accountInfo
	mask.Apply(&maskedAccountInfo)

	// assert
	assert.Equal(t, expectedMask, maskedAccountInfo)
}

func Test_MaskSimple_WithCharArgument_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type UserPass struct {
		Username string
		Password string `mask:"simple,*"`
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
	mask.Apply(&up)

	// assert
	assert.Equal(t, expectedMask, up)
}

func Test_MaskSimple_WithArray_ReturnsCorrectResult(t *testing.T) {
	// arrange
	type UserPass struct {
		Username string
		Password string `mask:"simple,*"`
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
	mask.Apply(&ups)

	// assert
	assert.Equal(t, expectedMask, ups)
}
