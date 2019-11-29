package meli

import (
	"encoding/json"
)

type Category struct {
	Id                    CategoryId   `json:"id,omitempty"`
	Name                  string       `json:"name,omitempty"`
	PredictionProbability float64      `json:"prediction_probability,omitempty"`
	ShippingModes         []string     `json:"shipping_modes,omitempty"`
	PathFromRoot          []*Category  `json:"path_from_root,omitempty"`
	Variations            []*Attribute `json:"variations,omitempty"`
}

type CategoryId string

func (ml *MeLi) CategoryAttributes(catId CategoryId) ([]*Attribute, error) {
	if catId == "" {
		return nil, errNilCategoryId
	}
	URL, err := ml.RouteTo("/categories/%s/attributes", nil, catId)
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
	atts := []*Attribute{}
	err = json.NewDecoder(resp.Body).Decode(&atts)
	if err != nil {
		return nil, err
	}
	return atts, nil
}

func (ml *MeLi) CategoryVariableAttributes(catId CategoryId) ([]*Attribute, error) {
	attrs, err := ml.CategoryAttributes(catId)
	if err != nil {
		return nil, err
	}
	var varAttrs []*Attribute
	for _, att := range attrs {
		if att.tagValue("allow_variations") {
			varAttrs = append(varAttrs, att)
		}
	}
	return varAttrs, nil
}
