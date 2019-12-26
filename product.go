package meli

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

type ProductId string

type Product struct {
	Id              ProductId  `json:"id,omitempty"`
	SiteId          string     `json:"site_id,omitempty"`
	Title           string     `json:"title,omitempty"`
	Status          string     `json:"status,omitempty"`
	SellerId        int        `json:"seller_id,omitempty"`
	CategoryId      CategoryId `json:"category_id,omitempty"`
	OfficialStoreId int        `json:"official_store_id,omitempty"`

	Price      float64 `json:"price,omitempty"`
	BasePrice  float64 `json:"base_price,omitempty"`
	CurrencyId string  `json:"currency_id,omitempty"`

	AvailableQuantity *int `json:"available_quantity,omitempty"`
	InitialQuantity   int  `json:"initial_quantity,omitempty"`
	SoldQuantity      int  `json:"sold_quantity,omitempty"`
	//
	BuyingMode      BuyingMode `json:"buying_mode,omitempty"`
	Condition       Condition  `json:"condition,omitempty"`
	Permalink       string     `json:"permalink,omitempty"`
	Thumbnail       string     `json:"thumbnail,omitempty"`
	SecureThumbnail string     `json:"secure_thumbnail,omitempty"`

	ListingTypeId  ListingTypeId `json:"listing_type_id,omitempty"`
	ListingSource  string        `json:"listing_source,omitempty"`
	CatalogListing bool          `json:"catalog_listing,omitempty"`

	Shipping      *Shipping      `json:"shipping,omitempty"`
	SellerAddress *SellerAddress `json:"seller_address,omitempty"`
	SaleTerms     []*SaleTerm    `json:"sale_terms,omitempty"`
	Attributes    []*Attribute   `json:"attributes,omitempty"`
	Pictures      []*Picture     `json:"pictures,omitempty"`
	Variants      []*Variant     `json:"variations,omitempty"`

	Description struct {
		PlainText string `json:"plain_text,omitempty"`
	} `json:"description,omitempty"`

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

	EndTime           time.Time `json:"end_time"`
	ExpirationTime    time.Time `json:"expiration_time"`
	SellerCustomField string    `json:"seller_custom_field"`

	lock sync.Mutex
}

type Condition string

func (c Condition) validate() error {
	for _, validC := range []Condition{"new", "used"} {
		if c == validC {
			return nil
		}
	}
	return ErrInvalidCondition
}

type BuyingMode string

type ListingTypeId string

// TODO: test with the api responses of https://api.mercadolibre.com/sites/{Site_id}/listing_types
func (ltId ListingTypeId) validate(siteId SiteId) error {
	validListingTypeIds := map[SiteId][]ListingTypeId{
		"MLA": {"gold_pro", "gold_premium", "gold_special", "gold", "silver", "bronze", "free"},
	}
	for _, validLtId := range validListingTypeIds[siteId] {
		if ltId == validLtId {
			return nil
		}
	}
	return ErrInvalidListingTypeId
}

func (bM BuyingMode) validate() error {
	for _, validBM := range []BuyingMode{"buy_it_now", "auction", "classified"} {
		if bM == validBM {
			return nil
		}
	}
	return ErrInvalidBuyingMode
}

func NewProduct(
	title, condition, buyingMode, listingTypeId string,
	categoryId CategoryId,
	price float64,
	stock *int,
	picsSrcs []string,
) (*Product, error) {
	var pics []*Picture
	for _, src := range picsSrcs {
		pics = append(pics, &Picture{Source: src})
	}
	prod := &Product{
		Title: title, CategoryId: categoryId,
		Condition: Condition(condition), BuyingMode: BuyingMode(buyingMode), ListingTypeId: ListingTypeId(listingTypeId),
		Price:             price,
		AvailableQuantity: stock,
		Pictures:          pics,
	}
	err := prod.validate(false)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func NewExistantProduct(
	title, condition, buyingMode, listingTypeId string,
	categoryId CategoryId,
	price float64,
	stock *int,
	picsSrcs []string,
) (*Product, error) {
	var pics []*Picture
	for _, src := range picsSrcs {
		pics = append(pics, &Picture{Source: src})
	}
	prod := &Product{
		Title: title, CategoryId: categoryId,
		Condition: Condition(condition), BuyingMode: BuyingMode(buyingMode), ListingTypeId: ListingTypeId(listingTypeId),
		Price:             price,
		AvailableQuantity: stock,
		Pictures:          pics,
	}
	err := prod.validate(true)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (prod *Product) validate(exists bool) error {
	if len(prod.CategoryId) > 0 && len(prod.CategoryId) < 4 {
		return ErrInvalidCategoryId
	}

	if prod.Price == 0 {
		return ErrNilPrice
	}
	if exists {
		return nil
	}

	if prod.Title == "" {
		return ErrNilProductTitle
	}
	if len(prod.CategoryId) < 4 {
		return ErrNilCategoryId
	}
	if prod.AvailableQuantity == nil {
		return ErrNilStock
	}
	if err := prod.Condition.validate(); err != nil {
		return err
	}
	if err := prod.BuyingMode.validate(); err != nil {
		return err
	}
	if err := prod.ListingTypeId.validate(prod.site()); err != nil {
		return err
	}
	if prod.Pictures == nil {
		return ErrNilPictures
	}
	return nil
}

func (prod *Product) site() SiteId {
	if len(prod.CategoryId) < 3 {
		return ""
	}
	return SiteId(prod.CategoryId[0:2])
}

type ProductEdge struct {
	SellerId string      `json:"seller_id"`
	Query    interface{} `json:"query"`
	Paging   struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	} `json:"paging"`
	Results  []ProductId `json:"results"`
	ScrollId string      `json:"scroll_id"`
	Orders   []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"orders"`
	AvailableOrders []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"available_orders"`
}

func (ml *MeLi) FetchProducts() ([]*Product, error) {
	prodIds, err := ml.ScanAllProducts()
	if err != nil {
		return nil, err
	}
	chunkedProdIds := chunkProductIds(prodIds, 20)
	prodsChunksCh := make(chan []*Product)
	errCh := make(chan error)
	for _, chunk := range chunkedProdIds {
		chunk := chunk
		go func() {
			prods, err := ml.GetProducts(chunk)
			if err != nil {
				errCh <- err
				return
			}
			prodsChunksCh <- prods
		}()
	}
	var prods []*Product
	for i := 0; i < len(chunkedProdIds); i++ {
		select {
		case prodsChunk := <-prodsChunksCh:
			prods = append(prods, prodsChunk...)
		case err := <-errCh:
			return nil, err
		}
	}
	return prods, nil
}

func (ml *MeLi) GetProducts(ids []ProductId) ([]*Product, error) {
	if len(ids) > 20 {
		return nil, ErrInvalidMultigetQuantity
	}
	params, err := ml.paramsWithToken()
	if err != nil {
		params = nil // retrieve public product in case of not having credentials
	}

	params.Set("ids", csvProductIds(ids))

	URL, err := ml.RouteTo("/items/%v", params)
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
	var prods []*Product
	err = json.NewDecoder(resp.Body).Decode(&prods)
	if err != nil {
		return nil, err
	}
	return prods, nil
}

func (ml *MeLi) ScanAllProducts() ([]ProductId, error) {
	var scrollId string
	var prodIds []ProductId
	for {
		edge, err := ml.ScanProducts(scrollId)
		if err != nil {
			return nil, err
		}
		if len(edge.Results) == 0 {
			break
		}
		prodIds = append(prodIds, edge.Results...)
	}
	return prodIds, nil
}

func (ml *MeLi) ScanProducts(scrollId string) (*ProductEdge, error) {
	params, err := ml.paramsWithToken()
	if err != nil {
		return nil, err
	}

	params.Set("scroll_id", scrollId)
	params.Set("status", "active")
	params.Set("limit", "100")
	params.Set("search_type", "scan")
	URL, err := ml.RouteTo("/users/%v/items/search", params)
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
	edge := &ProductEdge{}
	err = json.NewDecoder(resp.Body).Decode(edge)
	if err != nil {
		return nil, err
	}
	return edge, nil
}

func (ml *MeLi) GetProduct(prodId ProductId) (*Product, error) {
	params, err := ml.paramsWithToken()
	if err != nil {
		params = nil // retrieve public product in case of not having credentials
	}
	URL, err := ml.RouteTo("/items/%v", params, prodId)
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

func (ml *MeLi) DeleteProduct(id ProductId) (*Product, error) {
	prod := &Product{Id: id}
	prod.Close()
	prod, err := ml.updateProduct(prod)
	if err != nil {
		return nil, err
	}
	time.Sleep(1 * time.Second)
	prod.Delete()
	return ml.updateProduct(prod)
}

func (prod *Product) Close() {
	prod.Status = "closed"
}

func (ml *MeLi) SetProduct(prod *Product) (newProd *Product, err error) {
	if prod == nil {
		return nil, ErrNilProduct
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
	URL, err := ml.RouteTo("/items/%v", params)
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
	URL, err := ml.RouteTo("/items/%v", params, prod.Id)

	if err != nil {
		return nil, err
	}
	prod.Id = "" // Unset id since its in the route
	if prod.Deleted {
		prod = &Product{Deleted: true}
	}
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
	if p.AvailableQuantity == nil {
		p.AvailableQuantity = &stock
		return
	}
	*p.AvailableQuantity += stock
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
		return ErrIncompatibleVar
	}
	p.Variants = append(p.Variants, v)
	return nil
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

func (prod *Product) removeVariant(varId VariantId) (v *Variant) {
	for i, pV := range prod.Variants {
		if pV.Id == varId {
			v = pV
			prod.rmVariantByIdx(i)
		}
	}
	return
}

func (prod *Product) rmVariantByIdx(i int) {
	prod.lock.Lock()
	lastIndex := len(prod.Variants) - 1
	prod.Variants[i] = prod.Variants[lastIndex]
	prod.Variants[lastIndex] = nil // Notices the GC to rm the last elem to avoid mem-leak
	prod.Variants = prod.Variants[:lastIndex]
	prod.lock.Unlock()
}

// Yes, it can be done wout wrapping out to do it O(n) but readability can be harm boi
func chunkProductIds(prodIds []ProductId, sz int) (chunkedProdIds [][]ProductId) {
	for i := 0; i < len(prodIds); i += sz {
		end := i + sz
		if end > len(prodIds) {
			end = len(prodIds)
		}
		chunkedProdIds = append(chunkedProdIds, prodIds[i:end])
	}
	return
}

func csvProductIds(ids []ProductId) (csv string) {
	for _, id := range ids {
		csv += string(id)
		csv += ","
	}
	return strings.TrimSuffix(csv, ",")
}
