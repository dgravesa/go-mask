package masking_test

import (
	"fmt"

	"github.com/dgravesa/go-mask/masking"
)

type MyInnerType struct {
	ShouldStrBeMasked bool
	Str               string
}

type MyOuterType struct {
	Name  string
	Inner MyInnerType `mask:"xxx"`
}

func maskInnerType(v interface{}) error {
	innerType, ok := v.(*MyInnerType)
	if !ok {
		return fmt.Errorf("unable to convert to *MyInnerType")
	}
	if innerType.ShouldStrBeMasked {
		innerType.Str = "XXXXXXXX"
	}
	return nil
}

func init() {
	masking.RegisterMasker("xxx", maskInnerType)
}

func ExampleRegisterMasker_struct() {
	t1 := MyOuterType{
		Name: "unmasked",
		Inner: MyInnerType{
			ShouldStrBeMasked: false,
			Str:               "this should not be masked",
		},
	}

	err := masking.Mask(&t1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s: %s\n", t1.Name, t1.Inner.Str)

	t2 := MyOuterType{
		Name: "masked",
		Inner: MyInnerType{
			ShouldStrBeMasked: true,
			Str:               "this should be masked",
		},
	}

	err = masking.Mask(&t2)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s: %s\n", t2.Name, t2.Inner.Str)

	// Output:
	// unmasked: this should not be masked
	// masked: XXXXXXXX
}
