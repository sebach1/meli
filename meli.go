package meli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type MeLi struct {
	http.Client

	Credentials creds
}

func (ml *MeLi) SetClient(c http.Client) {
	ml.Client = c
}

func (ml *MeLi) SetCredentials(ctx context.Context, access, refresh, appId, secret string) error {
	var creds creds
	creds.Access = accessToken(access)
	creds.Refresh = refreshToken(refresh)
	creds.Secret = token(secret)
	creds.ApplicationId = applicationId(appId)
	return creds.validate()
}

func (ml *MeLi) SetAndValidateCredentials(ctx context.Context, access, refresh, appId, secret string) error {
	err := ml.SetCredentials(ctx, access, refresh, appId, secret)
	if err != nil {
		return err
	}
	return ml.RefreshToken()
}

// RouteTo retrieves a route given a path alias to the desired resource
// Notice: it returns a trailing slash on the return value in case it can contain childs.
// For example, in the case of /items, it'll return /items/ instead (alerting it is a sort of dir of sub-nodes)
// Path can be "auth", "product", "category_predict", "category", "category_attributes"
func (ml *MeLi) RouteTo(path string, params url.Values, ids ...interface{}) (string, error) {
	base := "https://api.mercadolibre.com"
	if ids != nil {
		base += fmt.Sprintf(path, ids...)
	} else {
		base += strings.ReplaceAll(path, "%v", "")
	}
	URL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	URL.RawQuery = params.Encode()
	base = URL.String()
	return base, nil
}

func (ml *MeLi) paramsWithToken() (url.Values, error) {
	if ml.Credentials.Access == "" {
		return nil, errNilAccessToken
	}
	params := url.Values{}
	params.Set("access_token", string(ml.Credentials.Access))
	return params, nil
}

func (ml *MeLi) Post(url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return ml.Do(req)
}

func (ml *MeLi) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return ml.Do(req)
}

func (ml *MeLi) Put(url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return ml.Do(req)
}
