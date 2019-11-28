package meli

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMeLi_Classify(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		args    args
		stub    *stub
		wantCat *Category
		wantErr error
	}{
		{
			name:    "correct behaviour",
			wantErr: nil,
			args:    args{title: "quux"},
			stub: &stub{status: 200,
				body: &Category{Id: "foo", PredictionProbability: 1, Name: "bar"},
				wantParamsReceive: url.Values{
					"title": []string{"quux"},
				},
			},
			wantCat: &Category{Id: "foo", PredictionProbability: 1, Name: "bar"},
		},
		{
			name:    "REMOTE returns an ERR",
			wantErr: svErrFooBar,
			args:    args{title: "quux"},
			stub: &stub{status: 400,
				body: svErrFooBar,
				wantParamsReceive: url.Values{
					"title": []string{"quux"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{}
			svClose := tt.stub.serve(t, ml)
			defer svClose()

			gotCat, err := ml.Classify(tt.args.title)

			if fmt.Sprintf("%v", tt.wantErr) != fmt.Sprintf("%v", err) {
				t.Errorf("MeLi.Classify() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.wantCat, gotCat); diff != "" {
				t.Errorf("MeLi.Classify() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestMeLi_ClassifyBatch(t *testing.T) {
	type args struct {
		titles []string
	}
	tests := []struct {
		name     string
		args     args
		stub     *stub
		wantCats []*Category
		wantErr  error
	}{
		{
			name:    "correct behaviour",
			wantErr: nil,
			args:    args{titles: []string{"a", "b"}},
			stub: &stub{status: 200,
				body: []*Category{
					{Id: "foo", PredictionProbability: 1, Name: "bar"},
					{Id: "baz", PredictionProbability: 1, Name: "quux"},
				},
				wantBodyReceive: JSONMarshal(t,
					[]map[string]string{{"title": "a"}, {"title": "b"}},
				),
			},
			wantCats: []*Category{
				{Id: "foo", PredictionProbability: 1, Name: "bar"},
				{Id: "baz", PredictionProbability: 1, Name: "quux"},
			},
		},
		{
			name:    "sends no body",
			wantErr: svErrFooBar,
			args:    args{titles: []string{}},
			stub: &stub{status: 400,
				body:            svErrFooBar,
				wantBodyReceive: JSONMarshal(t, []string{}),
			},
			wantCats: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MeLi{}
			svClose := tt.stub.serve(t, ml)
			defer svClose()

			gotCats, err := ml.ClassifyBatch(tt.args.titles)

			if fmt.Sprintf("%v", tt.wantErr) != fmt.Sprintf("%v", err) {
				t.Errorf("MeLi.ClassifyBatch() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.wantCats, gotCats); diff != "" {
				t.Errorf("MeLi.ClassifyBatch() mismatch (-want +got): %s", diff)
			}
		})
	}
}
