package meli

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/httpstub"
)

func TestMeLi_CategoryAttributes(t *testing.T) {
	t.Parallel()
	type args struct {
		catId CategoryId
	}
	tests := []struct {
		name    string
		args    args
		want    []*Attribute
		wantErr error
		stub    *httpstub.Stub
	}{
		{
			name:    "NIL cat ID",
			wantErr: ErrNilCategoryId,
			args:    args{catId: ""},
		},
		{
			name:    "REMOTE returns an ERR",
			wantErr: svErrFooBar,
			args:    args{catId: "foo"},
			stub: &httpstub.Stub{Status: 404,
				URL:  "/categories/foo/attributes",
				Body: svErrFooBar,
			},
		},
		{
			name: "REMOTE returns CORRECTly",
			args: args{catId: "foo"},
			stub: &httpstub.Stub{Status: 200,
				URL: "/categories/foo/attributes",
				Body: []*Attribute{
					{Id: "foo", Name: "bar"},
					{Id: "baz", Name: "quux"},
				},
			},
			want: []*Attribute{
				{Id: "foo", Name: "bar"},
				{Id: "baz", Name: "quux"},
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

			got, err := ml.CategoryAttributes(tt.args.catId)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.wantErr) {
				t.Errorf("MeLi.CategoryAttributes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("MeLi.CategoryAttributes() = (-want +got): %s", diff)
			}
		})
	}
}

func TestMeLi_CategoryVariableAttributes(t *testing.T) {
	type args struct {
		catId CategoryId
	}
	tests := []struct {
		name    string
		args    args
		want    []*Attribute
		wantErr error
		stub    *httpstub.Stub
	}{
		{
			name:    "NIL cat ID",
			wantErr: ErrNilCategoryId,
			args:    args{catId: ""},
		},
		{
			name:    "REMOTE returns an ERR",
			wantErr: svErrFooBar,
			args:    args{catId: "foo"},
			stub: &httpstub.Stub{Status: 404,
				URL:  "/categories/foo/attributes",
				Body: svErrFooBar,
			},
		},
		{
			name: "REMOTE returns CORRECTly",
			args: args{catId: "foo"},
			stub: &httpstub.Stub{Status: 200,
				URL: "/categories/foo/attributes",
				Body: []*Attribute{
					{Id: "foo", Name: "bar"},
					{Id: "baz", Name: "quux", Tags: []Tag{{"allow_variations": true}}},
				},
			},
			want: []*Attribute{
				{Id: "baz", Name: "quux", Tags: []Tag{{"allow_variations": true}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{}
			stubber := httpstub.Stubber{Stubs: []*httpstub.Stub{tt.stub}, Client: ml}
			cleanup := stubber.Serve(t)
			defer cleanup()

			got, err := ml.CategoryVariableAttributes(tt.args.catId)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.wantErr) {
				t.Errorf("MeLi.CategoryVariableAttributes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("MeLi.CategoryVariableAttributes() = (-want +got): %s", diff)
			}
		})
	}
}
