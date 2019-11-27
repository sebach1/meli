package meli

import (
	"bytes"
	"encoding/json"
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

func (ml *MeLi) GetProduct(id string) (*Product, error) {
	URL, err := ml.RouteTo("product", nil, id)
	if err != nil {
		return nil, err
	}
	resp, err := ml.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	prod := &Product{}
	err = json.NewDecoder(resp.Body).Decode(prod)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (ml *MeLi) SetProduct(prod *Product) (newProd *Product, err error) {
	if prod == nil {
		return nil, errNilProduct
	}
	if prod.Id == "" {
		newProd, err = ml.createProduct(prod)
	} else {
		newProd, err = ml.updateProduct(prod)
	}
	return
}

func (ml *MeLi) createProduct(prod *Product) (*Product, error) {
	params, err := ml.paramsWithToken()
	if err != nil {
		return nil, err
	}
	URL, err := ml.RouteTo("product", params)
	if err != nil {
		return nil, err
	}
	jsonProd, err := json.Marshal(prod)
	if err != nil {
		return nil, err
	}
	resp, err := ml.Post(URL, bytes.NewReader(jsonProd))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	newProd := &Product{}
	err = json.NewDecoder(resp.Body).Decode(newProd)
	if err != nil {
		return nil, err
	}
	return newProd, nil
}

func (ml *MeLi) updateProduct(prod *Product) (*Product, error) {
	params, err := ml.paramsWithToken()
	if err != nil {
		return nil, err
	}
	URL, err := ml.RouteTo("product", params, prod.Id)
	if err != nil {
		return nil, err
	}
	prod.Id = "" // Unset id since its in the route
	jsonProd, err := json.Marshal(prod)
	if err != nil {
		return nil, err
	}
	resp, err := ml.Put(URL, bytes.NewReader(jsonProd))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	newProd := &Product{}
	err = json.NewDecoder(resp.Body).Decode(newProd)
	if err != nil {
		return nil, err
	}
	return newProd, nil
}
