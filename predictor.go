package meli

import (
	"bytes"
	"encoding/json"
	"net/url"
)

type CategoryPrediction struct {
	Id                    string
	Name                  string
	PredictionProbability float64
	ShippingModes         []string `json:"shipping_modes"`
	PathFromRoot          []*CategoryPrediction
	Variations            []*Variation
}

type Variation struct {
	AttributeGroupID   string `json:"attribute_group_id"`
	AttributeGroupName string `json:"attribute_group_name"`
	Hierarchy          string `json:"hierarchy"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Relevance          int    `json:"relevance"`
	Tags               struct {
		AllowVariations bool `json:"allow_variations"`
		DefinesPicture  bool `json:"defines_picture"`
	} `json:"tags"`
	ValueMaxLength int    `json:"value_max_length"`
	ValueType      string `json:"value_type"`
	Values         []*Value
}
type Category struct {
	Id   string
	Name string
}

func (ml *MeLi) Classify(title string) (*CategoryPrediction, error) {
	params := url.Values{}
	params.Set("title", title)
	URL, err := ml.RouteTo("category_predict", "", params)
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
	cat := &CategoryPrediction{}
	err = json.NewDecoder(resp.Body).Decode(cat)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (ml *MeLi) ClassifyBatch(titles []string) ([]*CategoryPrediction, error) {
	URL, err := ml.RouteTo("category_predict", "", nil)
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
	cats := []*CategoryPrediction{}
	err = json.NewDecoder(resp.Body).Decode(&cats)
	if err != nil {
		return nil, err
	}
	return cats, nil
}
