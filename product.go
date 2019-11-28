package meli

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"
)

type Product struct {
	Id              ProductId `json:"id,omitempty"`
	SiteId          string    `json:"site_id,omitempty"`
	Title           string    `json:"title,omitempty"`
	Status          string    `json:"status,omitempty"`
	SellerId        int       `json:"seller_id,omitempty"`
	CategoryId      string    `json:"category_id,omitempty"`
	OfficialStoreId int       `json:"official_store_id,omitempty"`

	Price      int    `json:"price,omitempty"`
	BasePrice  int    `json:"base_price,omitempty"`
	CurrencyId string `json:"currency_id,omitempty"`

	AvailableQuantity int `json:"available_quantity,omitempty"`
	InitialQuantity   int `json:"initial_quantity,omitempty"`
	SoldQuantity      int `json:"sold_quantity,omitempty"`

	BuyingMode      string `json:"buying_mode,omitempty"`
	Condition       string `json:"condition,omitempty"`
	Permalink       string `json:"permalink,omitempty"`
	Thumbnail       string `json:"thumbnail,omitempty"`
	SecureThumbnail string `json:"secure_thumbnail,omitempty"`

	ListingTypeId  string `json:"listing_type_id,omitempty"`
	ListingSource  string `json:"listing_source,omitempty"`
	CatalogListing bool   `json:"catalog_listing,omitempty"`

	Shipping      *Shipping      `json:"shipping,omitempty"`
	SellerAddress *SellerAddress `json:"seller_address,omitempty"`
	SaleTerms     []*SaleTerm    `json:"sale_terms,omitempty"`
	Attributes    []*Attribute   `json:"attributes,omitempty"`
	Pictures      []*Picture     `json:"pictures,omitempty"`
	Variants      []*Variant     `json:"variations,omitempty"`

	Descriptions []struct {
		Id string `json:"id,omitempty"`
	} `json:"descriptions,omitempty"`

	StartTime   time.Time `json:"start_time,omitempty"`
	StopTime    time.Time `json:"stop_time,omitempty"`
	DateCreated time.Time `json:"date_created,omitempty"`
	LastUpdated time.Time `json:"last_updated,omitempty"`

	Warranty string `json:"warranty,omitempty"`

	AcceptsMercadopago        bool        `json:"accepts_mercadopago,omitempty"`
	InternationalDeliveryMode string      `json:"international_delivery_mode,omitempty"`
	SellerContact             interface{} `json:"seller_contact,omitempty"`
	Geolocation               struct {
		Latitude  float64 `json:"latitude,omitempty"`
		Longitude float64 `json:"longitude,omitempty"`
	} `json:"geolocation,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	DomainId        string   `json:"domain_id,omitempty"`
	AutomaticRelist bool     `json:"automatic_relist,omitempty"`

	Health                       float64       `json:"health,omitempty"`
	Location                     struct{}      `json:"location,omitempty"`
	VideoId                      string        `json:"video_id,omitempty"`
	SubStatus                    []interface{} `json:"sub_status,omitempty"`
	CatalogProductId             string        `json:"catalog_product_id,omitempty"`
	ParentItemId                 string        `json:"parent_item_id,omitempty"`
	DifferentialPricing          interface{}   `json:"differential_pricing,omitempty"`
	DealIds                      []interface{} `json:"deal_ids,omitempty"`
	CoverageAreas                []interface{} `json:"coverage_areas,omitempty"`
	Warnings                     []interface{} `json:"warnings,omitempty"`
	OriginalPrice                float64       `json:"original_price,omitempty"`
	NonMercadoPagoPaymentMethods []interface{} `json:"non_mercado_pago_payment_methods,omitempty"`
	Subtitle                     string        `json:"subtitle,omitempty"`

	Deleted bool `json:"deleted,omitempty"`
}

func (ml *MeLi) GetProduct(prodId ProductId) (*Product, error) {
	URL, err := ml.RouteTo("product", nil, prodId)
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
func (ml *MeLi) DeleteProduct(id ProductId) (*Product, error) {
	prod := &Product{Id: ProductId(id)}
	prod.Delete()
	return ml.updateProduct(prod)
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

// ManageStock adds to the product's stock the given stock.
// In case of giving a negative number, it rests the stock
func (p *Product) ManageStock(stock int) {
	p.AvailableQuantity += stock
}

func (p *Product) Delete() {
	p.Deleted = true
}

func (p *Product) AddVariant(v *Variant) error {
	err := v.validate()
	if err != nil {
		return err
	}
	if !p.varIsCompatible(v) {
		return errIncompatibleVar
	}
	p.Variants = append(p.Variants, v)
	return nil
}

func (prod *Product) RemoveVariant(varId VariantId) {
	for i, pV := range prod.Variants {
		if pV.Id == varId {
			prod.rmVariantByIdx(i)
		}
	}
}

func (prod *Product) ManageVarStocks(stockById map[VariantId]int) {
	for _, pV := range prod.Variants {
		if stock, ok := stockById[pV.Id]; ok {
			pV.ManageStock(stock)
		}
	}
}

func (prod *Product) RemoveCombination(attName string) {
	for _, v := range prod.Variants {
		for _, attC := range v.AttributeCombinations {
			if attC.Name == attName {
				attC.ValueId = ""
				attC.ValueName = ""
			}
		}
	}
}

func (p *Product) varIsCompatible(v *Variant) bool {
	for _, pV := range p.Variants {
		if !pV.isCompatible(v) {
			return false
		}
	}
	return true
}

func (prod *Product) rmVariantByIdx(i int) {
	var lock sync.Mutex // Avoid overlapping itself with a parallel call
	lock.Lock()
	lastIndex := len(prod.Variants) - 1
	prod.Variants[i] = prod.Variants[lastIndex]
	prod.Variants[lastIndex] = nil // Notices the GC to rm the last elem to avoid mem-leak
	prod.Variants = prod.Variants[:lastIndex]
	lock.Unlock()
}
