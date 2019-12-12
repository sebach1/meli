package meli

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/httpstub"
)

func TestMeLi_Classify(t *testing.T) {
	t.Parallel()
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		args    args
		stub    *httpstub.Stub
		wantCat *Category
		wantErr error
	}{
		{
			name:    "correct behaviour",
			wantErr: nil,
			args:    args{title: "quux"},
			stub: &httpstub.Stub{Status: 200,
				Body: &Category{Id: "foo", PredictionProbability: 1, Name: "bar"},
				WantParamsReceive: url.Values{
					"title": []string{"quux"},
				},
			},
			wantCat: &Category{Id: "foo", PredictionProbability: 1, Name: "bar"},
		},
		{
			name:    "REMOTE returns an ERR",
			wantErr: svErrFooBar,
			args:    args{title: "quux"},
			stub: &httpstub.Stub{Status: 400,
				Body: svErrFooBar,
				WantParamsReceive: url.Values{
					"title": []string{"quux"},
				},
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

			gotCat, err := ml.Classify(tt.args.title, "MLA")

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
	t.Parallel()
	type args struct {
		titles []string
	}
	tests := []struct {
		name     string
		args     args
		stub     *httpstub.Stub
		wantCats []*Category
		wantErr  error
	}{
		{
			name:    "correct behaviour",
			wantErr: nil,
			args:    args{titles: []string{"a", "b"}},
			stub: &httpstub.Stub{Status: 200,
				Body: []*Category{
					{Id: "foo", PredictionProbability: 1, Name: "bar"},
					{Id: "baz", PredictionProbability: 1, Name: "quux"},
				},
				WantBodyReceive: JSONMarshal(t,
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
			stub: &httpstub.Stub{Status: 400,
				Body:            svErrFooBar,
				WantBodyReceive: JSONMarshal(t, []string{}),
			},
			wantCats: nil,
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

			gotCats, err := ml.ClassifyBatch(tt.args.titles, "MLA")

			if fmt.Sprintf("%v", tt.wantErr) != fmt.Sprintf("%v", err) {
				t.Errorf("MeLi.ClassifyBatch() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.wantCats, gotCats); diff != "" {
				t.Errorf("MeLi.ClassifyBatch() mismatch (-want +got): %s", diff)
			}
		})
	}
}
