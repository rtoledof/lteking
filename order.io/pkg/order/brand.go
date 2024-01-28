package order

import (
	"fmt"
	"io"
	"strconv"
)

type Brand string

const (
	BrandBmw Brand = "BMW"
)

var AllBrand = []Brand{
	BrandBmw,
}

func (e Brand) IsValid() bool {
	switch e {
	case BrandBmw:
		return true
	}
	return false
}

func (e Brand) String() string {
	return string(e)
}

func (e *Brand) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Brand(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Brand", str)
	}
	return nil
}

func (e Brand) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
