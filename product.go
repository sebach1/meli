package meli

import "time"

type Product struct {
	Id              string `json:"id"`
	SiteId          string `json:"site_id"`
	Title           string `json:"title"`
	Status          string `json:"status"`
	SellerId        int    `json:"seller_id"`
	CategoryId      string `json:"category_id"`
	OfficialStoreId int    `json:"official_store_id"`

	Price      int    `json:"price"`
	BasePrice  int    `json:"base_price"`
	CurrencyId string `json:"currency_id"`

	AvailableQuantity int `json:"available_quantity"`
	InitialQuantity   int `json:"initial_quantity"`
	SoldQuantity      int `json:"sold_quantity"`

	BuyingMode      string `json:"buying_mode"`
	ListingTypeId   string `json:"listing_type_id"`
	Condition       string `json:"condition"`
	Permalink       string `json:"permalink"`
	Thumbnail       string `json:"thumbnail"`
	SecureThumbnail string `json:"secure_thumbnail"`

	Shipping      *Shipping
	SellerAddress *SellerAddress `json:"seller_address"`
	SaleTerms     []*SaleTerm    `json:"sale_terms"`
	Attributes    []*Attribute
	Pictures      []*Picture `json:"pictures"`
	Variants      []*Variant `json:"variations"`

	Descriptions []struct {
		Id string `json:"id"`
	} `json:"descriptions"`

	StartTime   time.Time `json:"start_time"`
	StopTime    time.Time `json:"stop_time"`
	DateCreated time.Time `json:"date_created"`
	LastUpdated time.Time `json:"last_updated"`

	Warranty string `json:"warranty"`

	AcceptsMercadopago        bool        `json:"accepts_mercadopago"`
	InternationalDeliveryMode string      `json:"international_delivery_mode"`
	SellerContact             interface{} `json:"seller_contact"`
	Geolocation               struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"geolocation"`
	ListingSource   string   `json:"listing_source"`
	Tags            []string `json:"tags"`
	DomainId        string   `json:"domain_id"`
	AutomaticRelist bool     `json:"automatic_relist"`

	Health                       float64       `json:"health"`
	CatalogListing               bool          `json:"catalog_listing"`
	Location                     struct{}      `json:"location"`
	VideoId                      string        `json:"video_id"`
	SubStatus                    []interface{} `json:"sub_status"`
	CatalogProductId             string        `json:"catalog_product_id"`
	ParentItemId                 string        `json:"parent_item_id"`
	DifferentialPricing          interface{}   `json:"differential_pricing"`
	DealIds                      []interface{} `json:"deal_ids"`
	CoverageAreas                []interface{} `json:"coverage_areas"`
	Warnings                     []interface{} `json:"warnings"`
	OriginalPrice                float64       `json:"original_price"`
	NonMercadoPagoPaymentMethods []interface{} `json:"non_mercado_pago_payment_methods"`
	Subtitle                     string        `json:"subtitle"`

	Deleted bool `json:"deleted"`
}
