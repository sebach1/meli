package meli

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMeLi_GetProduct(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantProd *Product
		wantErr  error
		stub     *stub
	}{
		{
			name:    "REMOTE returns an ERR",
			wantErr: meliErrFooBar,
			args:    args{id: "foo"},
			stub: &stub{status: 404,
				body: meliErrFooBar,
			},
		},
		{
			name:     "REMOTE returns the CORRECTly product",
			wantErr:  nil,
			wantProd: &Product{Id: "foo", Title: "bar"},
			args:     args{id: "foo"},
			stub: &stub{status: 200,
				body: &Product{Id: "foo", Title: "bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{}
			svClose := tt.stub.serve(t, ml)
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
		stub     *stub
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
			stub: &stub{status: 200,
				wantBodyReceive: JSONMarshal(t, &Product{Title: "bar"}), // The body sent lacks of id since its in the route
				wantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				body: &Product{Id: "foo", Title: "bar", Price: 0},
			},
			creds: creds{Access: "baz"},
		},
		{
			name:     "while CREATing product, REMOTE returns CORRECTly",
			prod:     &Product{Title: "bar"},
			wantProd: &Product{Id: "quux", Title: "bar"},
			stub: &stub{status: 200,
				wantBodyReceive: JSONMarshal(t, &Product{Title: "bar"}),
				wantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				body: &Product{Id: "quux", Title: "bar"},
			},
			creds: creds{Access: "baz"},
		},
		{
			name:    "while CREATing product, REMOTE returns CORRECTly",
			prod:    &Product{Title: "bar"},
			wantErr: meliErrFooBar,
			stub: &stub{status: 400,
				wantBodyReceive: JSONMarshal(t, &Product{Title: "bar"}),
				wantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				body: meliErrFooBar,
			},
			creds: creds{Access: "baz"},
		},
		{
			name:    "while EDITing product, REMOTE returns an ERROR",
			prod:    &Product{Id: "foo", Title: "bar"},
			wantErr: meliErrFooBar,
			stub: &stub{status: 400,
				wantBodyReceive: JSONMarshal(t, &Product{Title: "bar"}),
				wantParamsReceive: url.Values{
					"access_token": []string{"baz"},
				},
				body: meliErrFooBar,
			},
			creds: creds{Access: "baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{Credentials: tt.creds}
			svClose := tt.stub.serve(t, ml)
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
