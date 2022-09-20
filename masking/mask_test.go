package masking_test

import (
	"testing"

	"github.com/dgravesa/go-mask/masking"
	"github.com/stretchr/testify/assert"
)

func Test_Mask_OnStructWithSlice_DoesNotMaskSliceItems(t *testing.T) {
	// arrange
	type InnerS struct {
		Secret string `mask:"X"`
	}
	type OuterS struct {
		Slice []InnerS
	}
	slice := []InnerS{
		{
			Secret: "Hello, World!",
		},
		{
			Secret: "These should not be masked",
		},
	}
	s := OuterS{
		Slice: slice,
	}

	// act
	err := masking.Mask(&s)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, []InnerS{
		{
			Secret: "Hello, World!",
		},
		{
			Secret: "These should not be masked",
		},
	}, slice)
}

func Test_Mask_OnStructWithNestedStruct_MasksNestedStruct(t *testing.T) {
	// arrange
	type InnerS struct {
		Secret string `mask:"X"`
	}
	type OuterS struct {
		Nested InnerS
	}
	s := OuterS{
		Nested: InnerS{
			Secret: "This is a secret",
		},
	}

	// act
	err := masking.Mask(&s)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, OuterS{
		Nested: InnerS{
			Secret: "XXXXXXXXXXXXXXXX",
		},
	}, s)
}

func Test_Mask_OnStructWithPointedField_DoesNotMaskPointedField(t *testing.T) {
	// arrange
	type InnerS struct {
		Secret string `mask:"X"`
	}
	type OuterS struct {
		Pointed *InnerS
	}
	pointed := InnerS{
		Secret: "This should not be changed",
	}
	s := OuterS{
		Pointed: &pointed,
	}

	// act
	err := masking.Mask(&s)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, InnerS{
		Secret: "This should not be changed",
	}, pointed)
}

func Test_DeepMask_OnStructWithSlice_MasksSliceItems(t *testing.T) {
	// arrange
	type InnerS struct {
		Secret string `mask:"X"`
	}
	type OuterS struct {
		Slice []InnerS
	}
	s := OuterS{
		Slice: []InnerS{
			{
				Secret: "Hello, World!",
			},
			{
				Secret: "These should not be masked",
			},
		},
	}

	// act
	err := masking.DeepMask(&s)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, OuterS{
		Slice: []InnerS{
			{
				Secret: "XXXXXXXXXXXXX",
			},
			{
				Secret: "XXXXXXXXXXXXXXXXXXXXXXXXXX",
			},
		},
	}, s)
}

func Test_DeepMask_OnStructWithNestedStruct_MasksNestedStruct(t *testing.T) {
	// arrange
	type InnerS struct {
		Secret string `mask:"X"`
	}
	type OuterS struct {
		Nested InnerS
	}
	s := OuterS{
		Nested: InnerS{
			Secret: "This is a secret",
		},
	}

	// act
	err := masking.DeepMask(&s)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, OuterS{
		Nested: InnerS{
			Secret: "XXXXXXXXXXXXXXXX",
		},
	}, s)
}

func Test_DeepMask_OnStructWithPointedField_MasksPointedField(t *testing.T) {
	// arrange
	type InnerS struct {
		Secret string `mask:"X"`
	}
	type OuterS struct {
		Pointed *InnerS
	}
	s := OuterS{
		Pointed: &InnerS{
			Secret: "This should be masked",
		},
	}

	// act
	err := masking.DeepMask(&s)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, &InnerS{
		Secret: "XXXXXXXXXXXXXXXXXXXXX",
	}, s.Pointed)
}

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
			err := masking.Mask(tc.V)

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
	err := masking.Mask(&s)

	// assert
	assert.Error(t, err)
}
