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

	Deleted bool `json:"deleted"`
	// Irrelevant now
	// AcceptsMercadopago        bool        `json:"accepts_mercadopago"`
	// InternationalDeliveryMode string      `json:"international_delivery_mode"`
	// SellerContact             interface{} `json:"seller_contact"`
	// Geolocation               struct {
	// 	Latitude  float64 `json:"latitude"`
	// 	Longitude float64 `json:"longitude"`
	// } `json:"geolocation"`
	// ListingSource   string   `json:"listing_source"`
	// Tags            []string `json:"tags"`
	// DomainId        string   `json:"domain_id"`
	// AutomaticRelist bool     `json:"automatic_relist"`

	// Health         float64 `json:"health"`
	// CatalogListing bool    `json:"catalog_listing"`
	// Location      struct { } `json:"location"`
	// VideoId           interface{} `json:"video_id"`
	// SubStatus           []interface{} `json:"sub_status"`
	// CatalogProductId    interface{}   `json:"catalog_product_id"`
	// ParentItemId        interface{}   `json:"parent_item_id"`
	// DifferentialPricing interface{}   `json:"differential_pricing"`
	// DealIds             []interface{} `json:"deal_ids"`
	// CoverageAreas []interface{} `json:"coverage_areas"`
	// Warnings      []interface{} `json:"warnings"`
	// OriginalPrice     interface{} `json:"original_price"`
	// NonMercadoPagoPaymentMethods []interface{} `json:"non_mercado_pago_payment_methods"`
	// Subtitle          interface{} `json:"subtitle"`
}

type Picture struct {
	Id        string `json:"id"`
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Size      string `json:"size"`
	MaxSize   string `json:"max_size"`
	Quality   string `json:"quality"`
}

type Attribute struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	Values             []*Value `json:"values"`
	AttributeGroupId   string   `json:"attribute_group_id"`
	AttributeGroupName string   `json:"attribute_group_name"`
}

type SaleTerm struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Values []*Value `json:"values"`
}

type SellerAddress struct {
	City struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"city"`
	State struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"state"`
	Country struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"country"`
	SearchLocation struct {
		Neighborhood struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"neighborhood"`
		City struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"city"`
		State struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"state"`
	} `json:"search_location"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Id        int     `json:"id"`
}

type Shipping struct {
	Mode         string        `json:"mode"`
	Methods      []interface{} `json:"methods"`
	Tags         []string      `json:"tags"`
	Dimensions   interface{}   `json:"dimensions"`
	LocalPickUp  bool          `json:"local_pick_up"`
	FreeShipping bool          `json:"free_shipping"`
	LogisticType string        `json:"logistic_type"`
	StorePickUp  bool          `json:"store_pick_up"`
}

type Variant struct {
	Id    int64 `json:"id"`
	Price int   `json:"price"`

	AvailableQuantity int `json:"available_quantity"`
	SoldQuantity      int `json:"sold_quantity"`

	AttributeCombinations []*AttributeCombination `json:"attribute_combinations"`
	SaleTerms             []*SaleTerm             `json:"sale_terms"`
	PictureIds            []string                `json:"picture_ids"`
	CatalogProductId      interface{}             `json:"catalog_product_id"`
}

type Value struct {
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Struct interface{} `json:"struct"`
}

type AttributeCombination struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Values []*Value `json:"values"`
}

type WebHook struct {
	Resource      string `json:"resource"`
	UserId        int    `json:"user_id"`
	Topic         string `json:"topic"`
	ApplicationId int64  `json:"application_id"`
	Attempts      int    `json:"attempts"`

	Sent     time.Time `json:"sent"`
	Received time.Time `json:"received"`
}
