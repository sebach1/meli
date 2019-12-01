package meli

import (
	"testing"

	"github.com/mitchellh/copystructure"
)

func (v *Variant) copy(t *testing.T) *Variant {
	t.Helper()
	new, err := copystructure.Copy(v)
	if err != nil {
		t.Fatalf("Couldnt be able to copy struct: %v", err)
	}
	newVar, ok := new.(*Variant)
	if !ok {
		t.Fatalf("Couldnt be able to convert copied struct to native: %v", err)
	}
	return newVar
}
