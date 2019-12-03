package meli

import (
	"bytes"
	"encoding/json"
	"net/url"
)

type Variation struct {
	AttributeGroupID   string   `json:"attribute_group_id,omitempty"`
	AttributeGroupName string   `json:"attribute_group_name,omitempty"`
	Hierarchy          string   `json:"hierarchy,omitempty"`
	ID                 string   `json:"id,omitempty"`
	Name               string   `json:"name,omitempty"`
	Relevance          int      `json:"relevance,omitempty"`
	Tag                []Tag    `json:"tag,omitempty"`
	ValueMaxLength     int      `json:"value_max_length,omitempty"`
	ValueType          string   `json:"value_type,omitempty"`
	Values             []*Value `json:"values,omitempty"`
}

type SiteId string

func (ml *MeLi) Classify(title string, siteId SiteId) (*Category, error) {
	params := url.Values{}
	params.Set("title", title)
	URL, err := ml.RouteTo("/sites/%v/categories/category_predictor/predict", params, siteId)
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
	cat := &Category{}
	err = json.NewDecoder(resp.Body).Decode(cat)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (ml *MeLi) ClassifyBatch(titles []string, siteId SiteId) ([]*Category, error) {
	URL, err := ml.RouteTo("/sites/%v/categories/category_predictor/predict", nil, siteId)
	if err != nil {
		return nil, err
	}
	jsonTitles := []map[string]string{}
	for _, title := range titles {
		jsonTitles = append(jsonTitles, map[string]string{"title": title})
	}
	reqBody, err := json.Marshal(jsonTitles)
	if err != nil {
		return nil, err
	}
	resp, err := ml.Post(URL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	cats := []*Category{}
	err = json.NewDecoder(resp.Body).Decode(&cats)
	if err != nil {
		return nil, err
	}
	return cats, nil
}
