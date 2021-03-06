package meli

import (
	"bytes"
	"encoding/json"
)

type Variant struct {
	Id    VariantId `json:"id,omitempty"`
	Price float64   `json:"price,omitempty"`

	AvailableQuantity *int `json:"available_quantity,omitempty"`
	SoldQuantity      int  `json:"sold_quantity,omitempty"`

	AttributeCombinations []*Attribute `json:"attribute_combinations,omitempty"`
	SaleTerms             []*SaleTerm  `json:"sale_terms,omitempty"`
	PictureIds            []string     `json:"picture_ids,omitempty"`
	CatalogProductId      interface{}  `json:"catalog_product_id,omitempty"`
}

type VariantId int

func (ml *MeLi) GetVariant(varId VariantId, prodId ProductId) (*Variant, error) {
	if prodId == "" {
		return nil, ErrNilProductId
	}
	URL, err := ml.RouteTo("/items/%v/variations/%v", nil, prodId, varId)
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
		return nil, ErrNilProductId
	}
	if v == nil {
		return nil, ErrNilVariant
	}
	if err != nil {
		return nil, err
	}
	exists := v.Id == 0
	if exists {
		newVar, err = ml.createVariant(v, prodId)
	} else {
		err = v.validate()
		if err != nil {
			return nil, err
		}
		newVar, err = ml.updateVariant(v, prodId)
	}
	return
}

func (ml *MeLi) DeleteVariant(varId VariantId, prodId ProductId) (*Variant, error) {
	prod, err := ml.GetProduct(prodId)
	if err != nil {
		return nil, err
	}
	v := prod.removeVariant(varId)
	if v == nil {
		return nil, ErrVariantNotFound
	}
	_, err = ml.SetProduct(prod)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// ManageStock adds to the variant's stock the given stock.
// In case of giving a negative number, it rests the stock
func (v *Variant) ManageStock(stock int) {
	if v.AvailableQuantity == nil {
		v.AvailableQuantity = &stock
		return
	}
	*v.AvailableQuantity += stock
}

func NewVariant(
	attrs []*Attribute,
	price float64,
	stock *int,
	picIds []string,
) (*Variant, error) {
	v := &Variant{AttributeCombinations: attrs, Price: price, AvailableQuantity: stock, PictureIds: picIds}
	err := v.validate()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func NewExistantVariant(
	attrs []*Attribute,
	price float64,
	stock *int,
	picIds []string,
) (*Variant, error) {
	v := &Variant{AttributeCombinations: attrs, Price: price, AvailableQuantity: stock, PictureIds: picIds}
	return v, nil
}

func (v *Variant) validate() error {
	if v.Price == 0 {
		return ErrNilVarPrice
	}
	if v.AvailableQuantity == nil {
		return ErrNilVarStock
	}
	if len(v.AttributeCombinations) == 0 {
		return ErrNilCombinations
	}
	if len(v.PictureIds) == 0 {
		return ErrNilVarPictures
	}
	return nil
}

func (ml *MeLi) updateVariant(v *Variant, prodId ProductId) (*Variant, error) {
	params, err := ml.paramsWithToken()
	if err != nil {
		return nil, err
	}

	URL, err := ml.RouteTo("/items/%v/variations/%v", params, prodId, v.Id)
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
	prod, err := ml.GetProduct(prodId)
	if err != nil {
		return nil, err
	}
	err = prod.AddVariant(v)
	if err != nil {
		return nil, err
	}
	prod, err = ml.SetProduct(prod)
	if err != nil {
		return nil, err
	}
	for _, pV := range prod.Variants {
		if !pV.isCompatible(v) { // Note: it checks for being the same var as `v` due
			// in .AddVariant it already errs in case of not being compat with another one
			v = pV
			break
		}
	}
	return v, nil
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

func (attC *Attribute) equals(otherC *Attribute) bool {
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
