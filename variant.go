package meli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sebach1/meli/copy"
)

type VariantId int

type Variant struct {
	Id    VariantId `json:"id,omitempty"`
	Price int       `json:"price,omitempty"`

	AvailableQuantity int `json:"available_quantity,omitempty"`
	SoldQuantity      int `json:"sold_quantity,omitempty"`

	Attributes            []*Attribute            `json:"attributes,omitempty"`
	AttributeCombinations []*AttributeCombination `json:"attribute_combinations,omitempty"`
	SaleTerms             []*SaleTerm             `json:"sale_terms,omitempty"`
	PictureIds            []PictureId             `json:"picture_ids,omitempty"`
	CatalogProductId      interface{}             `json:"catalog_product_id,omitempty"`
}

type AttributeCombination struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	ValueName string `json:"value_name,omitempty"`
	ValueId   string `json:"value_id,omitempty"`
}

func (ml *MeLi) GetVariant(varId VariantId, prodId ProductId) (*Variant, error) {
	if prodId == "" {
		return nil, errNilProductId
	}
	URL, err := ml.RouteTo("/items/%s/variations/%s", nil, prodId, varId)
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
	v := &Variant{}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (ml *MeLi) SetVariant(v *Variant, prodId ProductId) (newVar *Variant, err error) {
	if prodId == "" {
		return nil, errNilProductId
	}
	if v == nil {
		return nil, errNilVariant
	}
	err = v.validate()
	if err != nil {
		return nil, err
	}
	if v.Id == 0 {
		newVar, err = ml.createVariant(v, prodId)
	} else {
		newVar, err = ml.updateVariant(v, prodId)
	}
	return
}

// ManageStock adds to the variant's stock the given stock.
// In case of giving a negative number, it rests the stock
func (v *Variant) ManageStock(stock int) {
	v.AvailableQuantity += stock
}

func (v *Variant) validate() error {
	if len(v.AttributeCombinations) == 0 {
		return errNilCombinations
	}
	if v.Price == 0 {
		return errNilVarPrice
	}
	if v.AvailableQuantity == 0 {
		return errNilVarStock
	}
	if len(v.PictureIds) == 0 {
		return errNilVarPictures
	}

	return nil
}

func (ml *MeLi) updateVariant(v *Variant, prodId ProductId) (*Variant, error) {
	params, err := ml.paramsWithToken()
	if err != nil {
		return nil, err
	}

	URL, err := ml.RouteTo("/items/%s/variations/%s", params, prodId, v.Id)
	if err != nil {
		return nil, err
	}
	v.Id = 0 // Unset id since its in the route
	jsonVar, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	resp, err := ml.Put(URL, bytes.NewReader(jsonVar))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	newVar := &Variant{}
	err = json.NewDecoder(resp.Body).Decode(newVar)
	if err != nil {
		return nil, err
	}
	return newVar, nil
}

func (ml *MeLi) createVariant(v *Variant, prodId ProductId) (*Variant, error) {
	params, err := ml.paramsWithToken()
	if err != nil {
		return nil, err
	}
	URL, err := ml.RouteTo("/items/%s/variations/%s", params, prodId)
	if err != nil {
		return nil, err
	}
	jsonVar, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	resp, err := ml.Post(URL, bytes.NewReader(jsonVar))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	newVar := &Variant{}
	err = json.NewDecoder(resp.Body).Decode(newVar)
	if err != nil {
		return nil, err
	}
	return newVar, nil

}

func (v *Variant) isCompatible(otherV *Variant) bool {
	if v.Id == otherV.Id {
		return false
	}
	var equalsQt int
	for _, vAtt := range v.AttributeCombinations {
		for _, oAtt := range otherV.AttributeCombinations {
			if vAtt.equals(oAtt) {
				equalsQt += 1
			}
		}
	}
	if equalsQt == len(v.AttributeCombinations) || equalsQt == len(otherV.AttributeCombinations) {
		return false
	}
	return true
}

func (attC *AttributeCombination) equals(otherC *AttributeCombination) bool {
	if attC.Id != otherC.Id {
		return false
	}
	if attC.ValueId != otherC.ValueId {
		return false
	}
	if attC.ValueName != otherC.ValueName {
		return false
	}
	return true
}

func (v *Variant) copy(t *testing.T) *Variant {
	t.Helper()
	newVar := Variant{}
	err := copy.Copy(&newVar, v)
	if err != nil {
		t.Fatalf("Couldnt be able to copy struct: %v", err)
	}
	newVar.AttributeCombinations = nil
	for _, attC := range v.AttributeCombinations {
		newAttC := &AttributeCombination{}
		err = copy.Copy(newAttC, attC)
		if err != nil {
			t.Fatalf("Couldnt be able to copy struct: %v", err)
		}
		newVar.AttributeCombinations = append(newVar.AttributeCombinations, newAttC)
	}
	return &newVar
}
