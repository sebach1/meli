package meli

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/sebach1/httpstub"
)

func TestMeLi_RefreshToken(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		creds            *creds
		wantErr          error
		stub             *httpstub.Stub
		wantAccessToken  accessToken
		wantRefreshToken accessToken
	}{
		{
			name:    "but NIL CREDS",
			wantErr: errNilCredentials,
		},
		{
			name:    "but NO APP Id",
			creds:   &creds{Access: "foo", Refresh: "bar", Secret: "baz"},
			wantErr: errNilApplicationId,
		},
		{
			name:    "but NO SECRET",
			creds:   &creds{Access: "foo", Refresh: "bar", ApplicationId: "baz"},
			wantErr: errNilSecret,
		},
		{
			name:    "but NO ACCESS",
			creds:   &creds{Refresh: "bar", ApplicationId: "baz", Secret: "foo"},
			wantErr: errNilAccessToken,
		},
		{
			name:    "but NO REFRESH",
			creds:   &creds{Access: "bar", ApplicationId: "baz", Secret: "foo"},
			wantErr: errNilRefreshToken,
		},
		{
			name:    "REMOTE returns an ERR",
			creds:   &creds{Access: "bar", ApplicationId: "baz", Secret: "foo", Refresh: "asd"},
			wantErr: svErrFooBar,
			stub: &httpstub.Stub{Status: 404,
				URL:  "/oauth/token",
				Body: svErrFooBar,
				Receive: httpstub.Receive{
					Params: url.Values{
						"grant_type":    []string{"refresh_token"},
						"refresh_token": []string{"asd"},
						"client_secret": []string{"foo"},
						"client_id":     []string{"baz"},
					},
				},
			},
		},
		{
			name:    "REMOTE returns an INCONSISTENCY",
			creds:   &creds{Access: "bar", ApplicationId: "baz", Secret: "foo", Refresh: "asd"},
			wantErr: errRemoteInconsistency,
			stub: &httpstub.Stub{Status: 200,
				URL:  "/oauth/token",
				Body: svErrFooBar,
				Receive: httpstub.Receive{
					Params: url.Values{
						"grant_type":    []string{"refresh_token"},
						"refresh_token": []string{"asd"},
						"client_secret": []string{"foo"},
						"client_id":     []string{"baz"}},
				},
			},
		},
		{
			name:             "REMOTE returns CORRECTLY",
			creds:            &creds{Access: "bar", ApplicationId: "baz", Secret: "foo", Refresh: "asd"},
			wantAccessToken:  "qux",
			wantRefreshToken: "quux",
			stub: &httpstub.Stub{Status: 200,
				URL:  "/oauth/token",
				Body: &authBody{AccessToken: "qux", RefreshToken: "quux"},
				Receive: httpstub.Receive{
					Params: url.Values{
						"grant_type":    []string{"refresh_token"},
						"refresh_token": []string{"asd"},
						"client_secret": []string{"foo"},
						"client_id":     []string{"baz"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ml := &MeLi{creds: tt.creds}
			stubber := httpstub.Stubber{Stubs: []*httpstub.Stub{tt.stub}, Client: ml}
			cleanup := stubber.Serve(t)
			defer cleanup()

			var oldAccessToken, newAccessToken accessToken
			var oldRefreshToken, newRefreshToken refreshToken
			if ml.creds != nil {
				oldAccessToken, oldRefreshToken = ml.creds.Access, ml.creds.Refresh
			}
			err := ml.RefreshToken()
			if fmt.Sprintf("%v", tt.wantErr) != fmt.Sprintf("%v", err) {
				t.Errorf("MeLi.RefreshToken() error got = %v, want: %v", err, tt.wantErr)
			}

			if ml.creds != nil {
				newAccessToken, newRefreshToken = ml.creds.Access, ml.creds.Refresh
			}
			if err != nil {
				if oldAccessToken != newAccessToken || oldRefreshToken != newRefreshToken {
					t.Errorf("MeLi.RefreshToken() ASSIGNED a new TOKEN on err")
				}
				return
			}

			if tt.wantAccessToken != ml.creds.Access {
				t.Errorf("MeLi.RefreshToken() assign mismatch, got: %v, want: %v", ml.creds.Access, tt.wantAccessToken)
			}
		})
	}
}
