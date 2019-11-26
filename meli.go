package meli

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type MeLi struct {
	http.Client

	Credentials creds
}

func (ml *MeLi) RouteTo(path string, id string, params url.Values) (string, error) {
	base := "https://api.mercadolibre.com"
	switch path {
	case "auth":
		base += "/oauth/token"
	case "product":
		base += "/items"
	case "category_predict":
		base += "/sites/MLA/category_predictor/predict"
	default:
		return "", errNonexistantPath
	}
	if id != "" {
		base += fmt.Sprintf("/%s", id)
	}
	if query := params.Encode(); query != "" {
		base = base + "?" + query
	}
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
