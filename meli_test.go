package meli

import (
	"fmt"
	"net/url"
	"testing"
)

func TestMeLi_RouteTo(t *testing.T) {
	type args struct {
		path   string
		id     fmt.Stringer
		params url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "all products",
			args: args{path: "product"},
			want: "https://api.mercadolibre.com/items",
		},
		{
			name: "single product",
			args: args{path: "product", id: ProductId("foo")},
			want: "https://api.mercadolibre.com/items/foo",
		},
		{
			name: "auth",
			args: args{path: "auth"},
			want: "https://api.mercadolibre.com/oauth/token",
		},
		{
			name: "category predict",
			args: args{path: "category_predict", id: SiteId("MLA")},
			want: "https://api.mercadolibre.com/sites/MLA/category_predictor/predict",
		},
		{
			name: "category predict with params",
			args: args{path: "category_predict", id: SiteId("MLA"), params: url.Values{"foo": []string{"bar"}}},
			want: "https://api.mercadolibre.com/sites/MLA/category_predictor/predict?foo=bar",
		},
		{
			name: "single product with params",
			args: args{path: "product", id: ProductId("foo"), params: url.Values{"bar": []string{"baz"}}},
			want: "https://api.mercadolibre.com/items/foo?bar=baz",
		},
		{
			name: "all products with params",
			args: args{path: "product", params: url.Values{"bar": []string{"baz"}}},
			want: "https://api.mercadolibre.com/items?bar=baz",
		},
		{
			name:    "nonexistant path",
			args:    args{path: "invalid", params: url.Values{"bar": []string{"baz"}}},
			wantErr: errNonexistantPath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{}
			got, err := ml.RouteTo(tt.args.path, tt.args.params, tt.args.id)
			if err != tt.wantErr {
				t.Errorf("MeLi.RouteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MeLi.RouteTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
