package meli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type MeLi struct {
	http.Client

	creds *creds
}

func (ml *MeLi) SetClient(c http.Client) {
	ml.Client = c
}

func (ml *MeLi) SetServerCredentials(ctx context.Context, appId, secret string) error {
	ml.creds = &creds{}
	ml.creds.Secret = token(secret)
	ml.creds.ApplicationId = applicationId(appId)
	return ml.creds.validateServer()
}

func (ml *MeLi) SetCredentials(ctx context.Context, access, refresh, appId, secret string) error {
	ml.creds = &creds{}
	ml.creds.Access = accessToken(access)
	ml.creds.Refresh = refreshToken(refresh)
	ml.creds.Secret = token(secret)
	ml.creds.ApplicationId = applicationId(appId)
	return ml.creds.validateClient()
}

func (ml *MeLi) SetAndValidateCredentials(ctx context.Context, access, refresh, appId, secret string) error {
	err := ml.SetCredentials(ctx, access, refresh, appId, secret)
	if err != nil {
		return err
	}
	return ml.RefreshToken()
}

func (ml *MeLi) AuthRouteTo(path string, params url.Values, site SiteId) (string, error) {
	urls := map[SiteId]string{
		"MLA": "https://auth.mercadolibre.com.ar",
		"MLB": "https://auth.mercadolivre.com.br",
		"MCO": "https://auth.mercadolibre.com.co",
		"MCR": "https://auth.mercadolibre.com.cr",
		"MEC": "https://auth.mercadolibre.com.ec",
		"MLC": "https://auth.mercadolibre.cl",
		"MLM": "https://auth.mercadolibre.com.mx",
		"MLU": "https://auth.mercadolibre.com.uy",
		"MLV": "https://auth.mercadolibre.com.ve",
		"MPA": "https://auth.mercadolibre.com.pa",
		"MPE": "https://auth.mercadolibre.com.pe",
		"MPT": "https://auth.mercadolivre.pt",
		"MRD": "https://auth.mercadolibre.com.do",
	}
	base, ok := urls[site]
	if !ok {
		return "", errInvalidSiteId
	}
	URL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	URL.RawQuery = params.Encode()
	base = URL.String()

	return base, nil
}

var errInvalidSiteId = errors.New("the given siteId does not exists")

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
	if ml.creds.Access == "" {
		return nil, errNilAccessToken
	}
	params := url.Values{}
	params.Set("access_token", string(ml.creds.Access))
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
