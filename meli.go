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

type ProductId string

func (pId ProductId) String() string { return string(pId) }

func (ml *MeLi) RouteTo(path string, params url.Values, ids ...fmt.Stringer) (string, error) {
	base := "https://api.mercadolibre.com"
	ids = rmEmptyStringers(ids)
	if len(ids) == 0 {
		ids = []fmt.Stringer{ProductId("")}
	}
	switch path {
	case "auth":
		base += "/oauth/token"
	case "product":
		base += fmt.Sprintf("/items/%s", ids[0])
	case "variant":
		base += fmt.Sprintf("/items/%s/variations/%s", ids[0], ids[1])
	case "category_predict":
		base += fmt.Sprintf("/sites/%s/category_predictor/predict", ids[0])
	case "category":
		base += fmt.Sprintf("/categories/%s", ids[0])
	case "category_attributes":
		base += fmt.Sprintf("/categories/%s/attributes", ids[0])
	default:
		return "", errNonexistantPath
	}
	base = strings.TrimRight(base, "/")
	if query := params.Encode(); query != "" {
		base = base + "?" + query
	}
	return base, nil
}

func rmEmptyStringers(stringers []fmt.Stringer) []fmt.Stringer {
	var fullStringers []fmt.Stringer
	for _, str := range stringers {
		if str != nil {
			fullStringers = append(fullStringers, str)
		}
	}
	return fullStringers
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
