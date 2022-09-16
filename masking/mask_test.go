package masking_test

import (
	"testing"

	"github.com/dgravesa/go-mask/masking"
	"github.com/stretchr/testify/assert"
)

func Test_Mask_OnNonPointer_ReturnsError(t *testing.T) {
	// arrange
	type S struct {
		Secret string `mask:"X"`
	}
	type TestCase struct {
		V    interface{}
		Name string
	}
	testCases := []TestCase{
		{
			V:    "this is a string",
			Name: "string",
		},
		{
			V: S{
				Secret: "this is a secret",
			},
			Name: "struct",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// act
			err := masking.Apply(tc.V)

			// assert
			assert.Error(t, err)
		})
	}
}

func Test_Mask_UnrecognizedMaskFunc_ReturnsError(t *testing.T) {
	// arrange
	type S struct {
		Secret string `mask:"idk"`
	}
	s := S{
		Secret: "this is a secret",
	}

	// act
	err := masking.Apply(&s)

	// assert
	assert.Error(t, err)
}
