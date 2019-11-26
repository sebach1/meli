package meli

import "time"

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
		Neighbourhood struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"neighborhood"` // :@
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
