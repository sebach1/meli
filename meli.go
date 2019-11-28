package meli

import (
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

// RouteTo retrieves a route given a path alias to the desired resource
// Notice: it returns a trailing slash on the return value in case it can contain childs.
// For example, in the case of /items, it'll return /items/ instead (alerting it is a sort of dir of sub-nodes)
// Path can be "auth", "product", "category_predict", "category", "category_attributes"
func (ml *MeLi) RouteTo(path string, params url.Values, ids ...interface{}) (string, error) {
	base := "https://api.mercadolibre.com"
	switch path {
	case "auth":
		base += "/oauth/token"
	case "product":
		base += "/items/%s"
	case "category_predict":
		base += "/sites/%s/category_predictor/predict"
	case "category":
		base += "/categories/%s"
	case "category_attributes":
		base += "/categories/%s/attributes"
	default:
		return "", errNonexistantPath
	}
	if ids != nil {
		base = fmt.Sprintf(base, ids...)
	} else {
		base = strings.ReplaceAll(base, "%s", "")
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

func (ml *MeLi) Put(url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return ml.Do(req)
}
