package meli

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/copystructure"
	"github.com/sebach1/httpstub"
)

func TestMeLi_GetProduct(t *testing.T) {
	t.Parallel()
	type args struct {
		id ProductId
	}
	tests := []struct {
		name     string
		args     args
		wantProd *Product
		wantErr  error
		stub     *httpstub.Stub
	}{

		{
			name:    "REMOTE returns an ERR",
			wantErr: svErrFooBar,
			args:    args{id: "foo"},
			stub: &httpstub.Stub{Status: 404,
				URL:  "/items/foo",
				Body: svErrFooBar,
			},
		},
		{
			name:     "REMOTE returns the CORRECTly product",
			wantErr:  nil,
			wantProd: gProducts.Bar.None,
			args:     args{id: "foo"},
			stub: &httpstub.Stub{Status: 200,
				URL:  "/items/foo",
				Body: gProducts.Bar.None,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ml := &MeLi{}
			stubber := httpstub.Stubber{Stubs: []*httpstub.Stub{tt.stub}, Client: ml}
			cleanup := stubber.Serve(t)
			defer cleanup()

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
	t.Parallel()
	tests := []struct {
		name     string
		creds    creds
		prod     *Product
		stub     *httpstub.Stub
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
			prod:    gProducts.Foo.None,
			creds:   creds{},
		},
		{
			name:    "but given NIL CREDENTIALS",
			wantErr: errNilAccessToken,
			prod:    gProducts.Foo.None.copy(t),
			creds:   creds{},
		},
		{
			name:     "while EDITing product REMOTE returns CORRECTly",
			prod:     gProducts.Foo.None.copy(t),
			wantProd: gProducts.Foo.Title.Alt.copy(t),
			stub: &httpstub.Stub{Status: 200,
				URL:             "/items/" + string(gProducts.Foo.None.Id),
				WantBodyReceive: JSONMarshal(t, gProducts.Foo.Id.Zero), // The body sent lacks of id since its in the route
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: gProducts.Foo.Title.Alt,
			},
			creds: creds{Access: "baz"},
		},
		{
			name:     "while CREATing product, REMOTE returns CORRECTly",
			prod:     gProducts.Bar.Id.Zero.copy(t),
			wantProd: gProducts.Bar.None.copy(t),
			stub: &httpstub.Stub{Status: 200,
				URL:             "/items/",
				WantBodyReceive: JSONMarshal(t, gProducts.Bar.Id.Zero),
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: gProducts.Bar.None,
			},
			creds: creds{Access: "baz"},
		},
		{
			name:    "while CREATing product, REMOTE returns ERRored",
			prod:    gProducts.Bar.Id.Zero.copy(t),
			wantErr: svErrFooBar,
			stub: &httpstub.Stub{Status: 400,
				URL:             "/items/",
				WantBodyReceive: JSONMarshal(t, gProducts.Bar.Id.Zero),
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: svErrFooBar,
			},
			creds: creds{Access: "baz"},
		},
		{
			name:    "while EDITing product, REMOTE returns an ERROR",
			prod:    gProducts.Bar.None.copy(t),
			wantErr: svErrFooBar,
			stub: &httpstub.Stub{Status: 400,
				URL:             "/items/" + string(gProducts.Bar.None.Id),
				WantBodyReceive: JSONMarshal(t, gProducts.Bar.Id.Zero),
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: svErrFooBar,
			},
			creds: creds{Access: "baz"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ml := &MeLi{Credentials: tt.creds}
			stubber := httpstub.Stubber{Stubs: []*httpstub.Stub{tt.stub}, Client: ml}
			cleanup := stubber.Serve(t)
			defer cleanup()

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
	t.Parallel()
	tests := []struct {
		name    string
		creds   creds
		prod    *Product
		stub    *httpstub.Stub
		wantErr error
	}{
		{
			name:    "but given NIL CREDENTIALS",
			wantErr: errNilAccessToken,
			prod:    gProducts.Foo.None.copy(t),
			creds:   creds{},
		},
		{
			name: "while EDITing product REMOTE returns CORRECTly",
			prod: gProducts.Foo.None.copy(t),
			stub: &httpstub.Stub{Status: 200,
				WantBodyReceive: JSONMarshal(t, &Product{Deleted: true}), // The body sent lacks of id since its in the route
				WantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				Body: gProducts.Foo.None,
			},
			creds: creds{Access: "baz"},
		},
		{
			name:    "while EDITing product, REMOTE returns an ERROR",
			prod:    gProducts.Foo.None.copy(t),
			wantErr: svErrFooBar,
			stub: &httpstub.Stub{Status: 400,
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ml := &MeLi{Credentials: tt.creds}
			stubber := httpstub.Stubber{Stubs: []*httpstub.Stub{tt.stub}, Client: ml}
			cleanup := stubber.Serve(t)
			defer cleanup()

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

func TestProduct_ManageVarStocks(t *testing.T) {
	t.Parallel()
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
			name: "vars EXISTS",
			prod: &Product{
				Variants: []*Variant{
					{Id: 1, AvailableQuantity: pointerToInt(2)},
					{Id: 5, AvailableQuantity: pointerToInt(10)},
				},
			},
			newProd: &Product{
				Variants: []*Variant{
					{Id: 1, AvailableQuantity: pointerToInt(1)},
					{Id: 5, AvailableQuantity: pointerToInt(12)},
				},
			},
			args: args{stockById: map[VariantId]int{1: -1, 5: 2}},
		},
		{
			name: "vars does NOT EXISTS",
			prod: &Product{
				Variants: []*Variant{
					{Id: 5, AvailableQuantity: pointerToInt(10)},
				},
			},
			newProd: &Product{
				Variants: []*Variant{
					{Id: 5, AvailableQuantity: pointerToInt(12)},
				},
			},
			args: args{stockById: map[VariantId]int{1: -1, 5: 2}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.prod.ManageVarStocks(tt.args.stockById)
			if diff := cmp.Diff(tt.prod, tt.newProd); diff != "" {
				t.Errorf("Product.ManageVarStocks() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestProduct_RemoveCombination(t *testing.T) {
	t.Parallel()
	type args struct {
		attName string
	}
	tests := []struct {
		name     string
		prod     *Product
		newCombs []*AttributeCombination
		args     args
	}{
		{
			name:     "combination is REMOVED SUCCESSFULLY",
			prod:     gProducts.Foo.None.copy(t),
			newCombs: rmValueAndReturn(gVariants.Foo.None.copy(t).AttributeCombinations),
			args:     args{attName: gVariants.Foo.None.AttributeCombinations[0].Name},
		},
		{
			name:     "combination is NOT IN VARS",
			prod:     gProducts.Foo.None.copy(t),
			newCombs: gVariants.Foo.None.AttributeCombinations,
			args:     args{attName: gVariants.Bar.None.AttributeCombinations[0].Name},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.prod.RemoveCombination(tt.args.attName)
			for _, v := range tt.prod.Variants {
				if diff := cmp.Diff(tt.newCombs, v.AttributeCombinations); diff != "" {
					t.Errorf("Product.RemoveCombination() mismatch (-want +got): %s", diff)
				}
			}
		})
	}
}

func (v *Variant) withoutId() *Variant { v.Id = 0; return v }

func TestProduct_AddVariant(t *testing.T) {
	t.Parallel()
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
			newProd: gProducts.Foo.None.copy(t).appendVariantAndReturn(gVariants.Bar.None),
			v:       gVariants.Bar.None.copy(t),
		},
		{
			name:    "given var is SUCCESSFULLY ADDED on EMPTY PROD",
			prod:    gProducts.Foo.Variants.Zero.copy(t),
			newProd: gProducts.Foo.Variants.Zero.copy(t).appendVariantAndReturn(gVariants.Bar.None),
			v:       gVariants.Bar.None.copy(t),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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

func TestProduct_ManageStock(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		stockArg  int
		stockWant int
		stockHave int
	}{
		{
			name:      "rest stock",
			stockArg:  -1,
			stockWant: 3 - 1,
			stockHave: 3,
		},
		{
			name:      "add stock",
			stockArg:  3,
			stockWant: 3,
			stockHave: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stockHave, stockArg, stockWant := pointerToInt(tt.stockHave), pointerToInt(tt.stockArg), pointerToInt(tt.stockWant)
			p := &Product{AvailableQuantity: stockHave}
			p.ManageStock(*stockArg)
			if diff := cmp.Diff(stockWant, p.AvailableQuantity); diff != "" {
				t.Errorf("Product.ManageStock() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestProduct_removeVariant(t *testing.T) {
	t.Parallel()
	type args struct {
		varId VariantId
	}
	tests := []struct {
		name    string
		prod    *Product
		args    args
		newProd *Product
	}{
		{
			name:    "prod without vars",
			prod:    gProducts.Zero.copy(t),
			args:    args{varId: gVariants.Bar.None.copy(t).Id},
			newProd: gProducts.Zero,
		},
		{
			name:    "deletes successfully the var",
			prod:    gProducts.Foo.None.copy(t),
			args:    args{varId: gVariants.Foo.None.copy(t).Id},
			newProd: gProducts.Foo.Variants.Zero,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.prod.removeVariant(tt.args.varId)
			if diff := cmp.Diff(tt.prod, tt.newProd); diff != "" {
				t.Errorf("Product.removeVariant() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func (prod *Product) copy(t *testing.T) *Product {
	t.Helper()
	new, err := copystructure.Copy(prod)
	if err != nil {
		t.Fatalf("Couldnt be able to copy struct: %v", err)
	}
	newProd, ok := new.(*Product)
	if !ok {
		t.Fatalf("Couldnt be able to convert copied struct to native: %v", err)
	}
	return newProd
}

func (p *Product) appendVariantAndReturn(v *Variant) *Product {
	p.Variants = append(p.Variants, v)
	return p
}

func rmValueAndReturn(combs []*AttributeCombination) []*AttributeCombination {
	for _, attC := range combs {
		attC.ValueId = ""
		attC.ValueName = ""
	}
	return combs
}

// func (prod *Product) modAndReturn(mod func(*Product)) *Product {
// 	mod(prod)
// 	return prod
// }

func pointerToInt(integer int) *int {
	return &integer
}
