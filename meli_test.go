package meli

import (
	"net/url"
	"testing"
)

func TestMeLi_RouteTo(t *testing.T) {
	type args struct {
		path   string
		ids    []interface{}
		params url.Values
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantsErr bool
	}{
		{
			name:     "gives more ids than formattable",
			args:     args{path: "%s/%s", ids: []interface{}{"foo", "bar", "baz"}},
			wantsErr: true,
		},
		{
			name: "all resource",
			args: args{path: "/items/%s"},
			want: "https://api.mercadolibre.com/items/",
		},
		{
			name: "multiple embedding",
			args: args{path: "/items/%s/variations/%s/sth/%s", ids: []interface{}{"foo", "bar", "baz"}},
			want: "https://api.mercadolibre.com/items/foo/variations/bar/sth/baz",
		},
		{
			name: "single embedding",
			args: args{path: "/items/%s", ids: []interface{}{"foo"}},
			want: "https://api.mercadolibre.com/items/foo",
		},
		{
			name: "wout formatting",
			args: args{path: "/oauth/token"},
			want: "https://api.mercadolibre.com/oauth/token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{}
			got, err := ml.RouteTo(tt.args.path, tt.args.params, tt.args.ids...)
			if (err != nil) != tt.wantsErr {
				t.Errorf("MeLi.RouteTo() error = %v, wantErr %v", err, tt.wantsErr)
				return
			}
			if got != tt.want {
				t.Errorf("MeLi.RouteTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
