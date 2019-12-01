package meli

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/meli/melitest"
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
		stub    *melitest.Stub
	}{
		{
			name:    "NIL cat ID",
			wantErr: errNilCategoryId,
			args:    args{catId: ""},
		},
		{
			name:    "REMOTE returns an ERR",
			wantErr: svErrFooBar,
			args:    args{catId: "foo"},
			stub: &melitest.Stub{Status: 404,
				Body: svErrFooBar,
			},
		},
		{
			name: "REMOTE returns CORRECTly",
			args: args{catId: "foo"},
			stub: &melitest.Stub{Status: 200,
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
			svClose := tt.stub.Serve(t, ml)
			defer svClose()

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
		stub    *melitest.Stub
	}{
		{
			name:    "NIL cat ID",
			wantErr: errNilCategoryId,
			args:    args{catId: ""},
		},
		{
			name:    "REMOTE returns an ERR",
			wantErr: svErrFooBar,
			args:    args{catId: "foo"},
			stub: &melitest.Stub{Status: 404,
				Body: svErrFooBar,
			},
		},
		{
			name: "REMOTE returns CORRECTly",
			args: args{catId: "foo"},
			stub: &melitest.Stub{Status: 200,
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
			svClose := tt.stub.Serve(t, ml)
			defer svClose()

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
