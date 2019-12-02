package meli

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/copystructure"
	"github.com/sebach1/meli/melitest"
)

func TestMeLi_GetVariant(t *testing.T) {
	type args struct {
		varId  VariantId
		prodId ProductId
	}
	tests := []struct {
		name    string
		args    args
		wantVar *Variant
		wantErr error
		stub    *melitest.Stub
	}{
		{
			name:    "VAR is not in PROD",
			wantErr: svErrFooBar,
			args:    args{prodId: gProducts.Foo.None.Id, varId: gVariants.Bar.None.Id},
			stub: &melitest.Stub{Status: 404,
				Body: svErrFooBar,
			},
		},
		{
			name:    "REMOTE returns the CORRECTly variant",
			wantErr: nil,
			wantVar: gVariants.Foo.None,
			args:    args{prodId: gProducts.Foo.None.Id, varId: gVariants.Foo.None.Id},
			stub: &melitest.Stub{Status: 200,
				Body: gVariants.Foo.None,
			},
		},
		{
			name:    "NIL given PROD id",
			wantErr: errNilProductId,
			wantVar: gVariants.Foo.None,
			args:    args{prodId: "", varId: gVariants.Foo.None.Id},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{}
			svClose := tt.stub.Serve(t, ml)
			defer svClose()

			gotVar, err := ml.GetVariant(tt.args.varId, tt.args.prodId)

			if fmt.Sprintf("%v", tt.wantErr) != fmt.Sprintf("%v", err) {
				t.Errorf("MeLi.GetVariant() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.wantVar, gotVar); diff != "" {
				t.Errorf("MeLi.GetVariant() mismatch (-want +got): %s", diff)
			}
		})
	}
}

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
