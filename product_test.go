package meli

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/meli/melitest"
)

func TestMeLi_GetProduct(t *testing.T) {
	type args struct {
		id ProductId
	}
	tests := []struct {
		name     string
		args     args
		wantProd *Product
		wantErr  error
		stub     *melitest.Stub
	}{
		{
			name:    "REMOTE returns an ERR",
			wantErr: svErrFooBar,
			args:    args{id: "foo"},
			stub: &melitest.Stub{Status: 404,
				Body: svErrFooBar,
			},
		},
		{
			name:     "REMOTE returns the CORRECTly product",
			wantErr:  nil,
			wantProd: &Product{Id: "foo", Title: "bar"},
			args:     args{id: "foo"},
			stub: &melitest.Stub{Status: 200,
				Body: &Product{Id: "foo", Title: "bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{}
			svClose := tt.stub.Serve(t, ml)
			defer svClose()

			gotProd, err := ml.GetProduct(tt.args.id)

			if fmt.Sprintf("%v", tt.wantErr) != fmt.Sprintf("%v", err) {
				t.Errorf("MeLi.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.wantProd, gotProd); diff != "" {
				t.Errorf("MeLi.GetProduct() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestMeLi_SetProduct(t *testing.T) {
	tests := []struct {
		name     string
		creds    creds
		prod     *Product
		stub     *melitest.Stub
		wantProd *Product
		wantErr  error
	}{
		{
			name:    "but given NIL PRODUCT",
			wantErr: errNilProduct,
			creds:   creds{Access: "baz"},
		},
		{
			name:    "but given NIL CREDENTIALS",
			wantErr: errNilAccessToken,
			prod:    &Product{Id: "foo", Title: "bar"},
			creds:   creds{},
		},
		{
			name:    "but given NIL CREDENTIALS",
			wantErr: errNilAccessToken,
			prod:    &Product{Id: "foo", Title: "bar"},
			creds:   creds{},
		},
		{
			name:     "while EDITing product REMOTE returns CORRECTly",
			prod:     &Product{Id: "foo", Title: "bar"},
			wantProd: &Product{Id: "foo", Title: "bar", Price: 0},
			stub: &melitest.Stub{Status: 200,
				WantBodyReceive: JSONMarshal(t, &Product{Title: "bar"}), // The body sent lacks of id since its in the route
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: &Product{Id: "foo", Title: "bar", Price: 0},
			},
			creds: creds{Access: "baz"},
		},
		{
			name:     "while CREATing product, REMOTE returns CORRECTly",
			prod:     &Product{Title: "bar"},
			wantProd: &Product{Id: "quux", Title: "bar"},
			stub: &melitest.Stub{Status: 200,
				WantBodyReceive: JSONMarshal(t, &Product{Title: "bar"}),
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: &Product{Id: "quux", Title: "bar"},
			},
			creds: creds{Access: "baz"},
		},
		{
			name:    "while CREATing product, REMOTE returns CORRECTly",
			prod:    &Product{Title: "bar"},
			wantErr: svErrFooBar,
			stub: &melitest.Stub{Status: 400,
				WantBodyReceive: JSONMarshal(t, &Product{Title: "bar"}),
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: svErrFooBar,
			},
			creds: creds{Access: "baz"},
		},
		{
			name:    "while EDITing product, REMOTE returns an ERROR",
			prod:    &Product{Id: "foo", Title: "bar"},
			wantErr: svErrFooBar,
			stub: &melitest.Stub{Status: 400,
				WantBodyReceive: JSONMarshal(t, &Product{Title: "bar"}),
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: svErrFooBar,
			},
			creds: creds{Access: "baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{Credentials: tt.creds}
			svClose := tt.stub.Serve(t, ml)
			defer svClose()

			gotProd, err := ml.SetProduct(tt.prod)
			if fmt.Sprintf("%v", tt.wantErr) != fmt.Sprintf("%v", err) {
				t.Errorf("MeLi.SetProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.wantProd, gotProd); diff != "" {
				t.Errorf("MeLi.SetProduct() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestMeLi_DeleteProduct(t *testing.T) {
	tests := []struct {
		name    string
		creds   creds
		prod    *Product
		stub    *melitest.Stub
		wantErr error
	}{
		{
			name:    "but given NIL CREDENTIALS",
			wantErr: errNilAccessToken,
			prod:    &Product{Id: "foo", Title: "bar"},
			creds:   creds{},
		},
		{
			name: "while EDITing product REMOTE returns CORRECTly",
			prod: &Product{Id: "foo"},
			stub: &melitest.Stub{Status: 200,
				WantBodyReceive: JSONMarshal(t, &Product{Deleted: true}), // The body sent lacks of id since its in the route
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: &Product{Id: "foo", Title: "bar", Price: 0, Deleted: true},
			},
			creds: creds{Access: "baz"},
		},
		{
			name:    "while EDITing product, REMOTE returns an ERROR",
			prod:    &Product{Id: "foo"},
			wantErr: svErrFooBar,
			stub: &melitest.Stub{Status: 400,
				WantBodyReceive: JSONMarshal(t, &Product{Deleted: true}),
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: svErrFooBar,
			},
			creds: creds{Access: "baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{Credentials: tt.creds}
			svClose := tt.stub.Serve(t, ml)
			defer svClose()

			tt.prod.Delete()
			updProd, updErr := ml.updateProduct(tt.prod)

			delProd, delErr := ml.DeleteProduct(tt.prod.Id)

			if fmt.Sprintf("%v", updErr) != fmt.Sprintf("%v", delErr) {
				t.Errorf("MeLi.DeleteProduct() error = %v, wantErr %v", delProd, updProd)
			}

			if diff := cmp.Diff(delProd, updProd); diff != "" {
				t.Errorf("MeLi.DeleteProduct() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestProduct_RemoveCombination(t *testing.T) {
	type args struct {
		attName string
	}
	tests := []struct {
		name    string
		prod    *Product
		newProd *Product
		args    args
	}{
		{
			name: "combination is already set up in vars",
			prod: &Product{
				Variants: []*Variant{
					{AttributeCombinations: []*AttributeCombination{
						{Name: "foo", ValueId: "bar", ValueName: "baz"},
						{Name: "quux", ValueId: "quz", ValueName: "quuz"},
					}},
				},
			},
			newProd: &Product{
				Variants: []*Variant{
					{AttributeCombinations: []*AttributeCombination{
						{Name: "foo", ValueId: "", ValueName: ""},
						{Name: "quux", ValueId: "quz", ValueName: "quuz"},
					}},
				},
			},
			args: args{attName: "foo"},
		},
		{
			name: "combination is not set up in vars",
			prod: &Product{
				Variants: []*Variant{
					{AttributeCombinations: []*AttributeCombination{
						{Name: "quux", ValueId: "quz", ValueName: "quuz"},
					}},
				},
			},
			newProd: &Product{
				Variants: []*Variant{
					{AttributeCombinations: []*AttributeCombination{
						{Name: "quux", ValueId: "quz", ValueName: "quuz"},
					}},
				},
			},
			args: args{attName: "foo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prod.RemoveCombination(tt.args.attName)
			if diff := cmp.Diff(tt.prod, tt.newProd); diff != "" {
				t.Errorf("Product.RemoveCombination() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestProduct_ManageVarStocks(t *testing.T) {
	type args struct {
		stockById map[VariantId]int
	}
	tests := []struct {
		name    string
		prod    *Product
		newProd *Product
		args    args
	}{
		{
			name: "vars exists",
			prod: &Product{
				Variants: []*Variant{
					{Id: 1, AvailableQuantity: 2},
					{Id: 5, AvailableQuantity: 10},
				},
			},
			newProd: &Product{
				Variants: []*Variant{
					{Id: 1, AvailableQuantity: 1},
					{Id: 5, AvailableQuantity: 12},
				},
			},
			args: args{stockById: map[VariantId]int{1: -1, 5: 2}},
		},
		{
			name: "vars does notexists",
			prod: &Product{
				Variants: []*Variant{
					{Id: 5, AvailableQuantity: 10},
				},
			},
			newProd: &Product{
				Variants: []*Variant{
					{Id: 5, AvailableQuantity: 12},
				},
			},
			args: args{stockById: map[VariantId]int{1: -1, 5: 2}},
		},
		{
			name: "vars does not exists",
			prod: &Product{
				Variants: []*Variant{
					{Id: 5, AvailableQuantity: 10},
				},
			},
			newProd: &Product{
				Variants: []*Variant{
					{Id: 5, AvailableQuantity: 12},
				},
			},
			args: args{stockById: map[VariantId]int{1: -1, 5: 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prod.ManageVarStocks(tt.args.stockById)
			if diff := cmp.Diff(tt.prod, tt.newProd); diff != "" {
				t.Errorf("Product.ManageVarStocks() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestProduct_AddVariant(t *testing.T) {
	tests := []struct {
		name    string
		prod    *Product
		newProd *Product
		wantErr error
		v       *Variant
	}{
		{
			name:    "given var has NO COMBINATIONS",
			prod:    gProducts.Foo.None.copy(t),
			wantErr: errNilCombinations,
			v:       gVariants.Bar.AttributeCombinations.Zero.copy(t),
		},
		{
			name:    "given VAR IS ALREADY in prod",
			prod:    gProducts.Bar.None.copy(t),
			wantErr: errIncompatibleVar,
			v:       gVariants.Bar.None.copy(t),
		},
		{
			name:    "given var is INCOMPATIBLE due repeated attr combinations",
			prod:    gProducts.Foo.None.copy(t),
			wantErr: errIncompatibleVar,
			v:       gVariants.Bar.AttributeCombinations.Alt.copy(t),
		},
		{
			name:    "given var has NIL PRICE",
			prod:    gProducts.Foo.None.copy(t),
			wantErr: errNilVarPrice,
			v:       gVariants.Bar.Price.Zero.copy(t),
		},
		{
			name:    "given var has NIL STOCK",
			prod:    gProducts.Foo.None.copy(t),
			wantErr: errNilVarStock,
			v:       gVariants.Bar.AvailableQuantity.Zero.copy(t),
		},
		{
			name:    "given var has NIL PICS",
			prod:    gProducts.Foo.None.copy(t),
			wantErr: errNilVarPictures,
			v:       gVariants.Bar.PictureIds.Zero.copy(t),
		},
		{
			name:    "given var is SUCCESSFULLY ADDED",
			prod:    gProducts.Foo.None.copy(t),
			newProd: gProducts.Foo.None.copy(t).addVariantUnsafe(gVariants.Bar.None),
			v:       gVariants.Bar.None.copy(t),
		},
		{
			name:    "given var is SUCCESSFULLY ADDED on EMPTY PROD",
			prod:    gProducts.Foo.Variants.Zero.copy(t),
			newProd: gProducts.Foo.Variants.Zero.copy(t).addVariantUnsafe(gVariants.Bar.None),
			v:       gVariants.Bar.None.copy(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalProd := tt.prod.copy(t)
			gotErr := tt.prod.AddVariant(tt.v)
			if fmt.Sprintf("%v", gotErr) != fmt.Sprintf("%v", tt.wantErr) {
				t.Errorf("Product.AddVariant() errors mismatch; got: %v; want: %v", gotErr, tt.wantErr)
			}

			if tt.wantErr != nil || gotErr != nil {
				tt.newProd = originalProd
				return
			}

			if diff := cmp.Diff(tt.newProd, tt.prod); diff != "" {
				t.Errorf("Product.AddVariant() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func (p *Product) addVariantUnsafe(v *Variant) *Product {
	p.Variants = append(p.Variants, v)
	return p
}
