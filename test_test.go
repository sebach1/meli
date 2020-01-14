package meli

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/sebach1/meli/internal/test/assist"
)

var (
	errFoo = errors.New("foo")
	errBar = errors.New("bar")

	svErrFooBar = &Error{ResponseErr: errFoo.Error(), Message: errBar.Error()}
)

var (
	gVariants goldenVariants
	gProducts goldenProducts
)

func init() {
	assist.DecodeJsonnet("variants", &gVariants)
	assist.DecodeJsonnet("products", &gProducts)
}

type goldenVariants struct {
	Foo  *variadicVariants `json:"foo,omitempty"`
	Bar  *variadicVariants `json:"bar,omitempty"`
	Zero *Variant          `json:"zero,omitempty"`
}

type variadicVariants struct {
	None                  *Variant     `json:"none,omitempty"`
	Id                    *optVariants `json:"id,omitempty"`
	Price                 *optVariants `json:"price,omitempty"`
	AvailableQuantity     *optVariants `json:"available_quantity,omitempty"`
	Attributes *optVariants `json:"attribute_combinations,omitempty"`
	PictureIds            *optVariants `json:"picture_ids,omitempty"`
}

type optVariants struct {
	Alt  *Variant `json:"alt,omitempty"`
	Zero *Variant `json:"zero,omitempty"`
}

type goldenProducts struct {
	Foo  *variadicProducts `json:"foo,omitempty"`
	Bar  *variadicProducts `json:"bar,omitempty"`
	Zero *Product          `json:"zero,omitempty"`
}

type variadicProducts struct {
	None              *Product     `json:"none,omitempty"`
	Id                *optProducts `json:"id,omitempty"`
	Variants          *optProducts `json:"variations,omitempty"`
	CategoryId        *optProducts `json:"category_id,omitempty"`
	AvailableQuantity *optProducts `json:"available_quantity,omitempty"`
	Title             *optProducts `json:"title,omitempty"`
	Price             *optProducts `json:"price,omitempty"`
}

type optProducts struct {
	Alt  *Product `json:"alt,omitempty"`
	Zero *Product `json:"zero,omitempty"`
}

func JSONMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()
	bytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("couldn't marshal given value: %v", err)
	}
	return bytes
}
