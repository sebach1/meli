package meli

type SaleTerm struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Values []*Value `json:"values"`
}

type SellerAddress struct {
	City struct {
		Id   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"city,omitempty"`
	State struct {
		Id   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"state,omitempty"`
	Country struct {
		Id   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"country,omitempty"`
	SearchLocation struct {
		Neighbourhood struct {
			Id   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"neighborhood,omitempty"` // :@
		City struct {
			Id   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"city,omitempty"`
		State struct {
			Id   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"state,omitempty"`
	} `json:"search_location,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Id        int     `json:"id,omitempty"`
}

type Shipping struct {
	Mode         string        `json:"mode,omitempty"`
	Methods      []interface{} `json:"methods,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	Dimensions   interface{}   `json:"dimensions,omitempty"`
	LocalPickUp  bool          `json:"local_pick_up,omitempty"`
	FreeShipping bool          `json:"free_shipping,omitempty"`
	LogisticType string        `json:"logistic_type,omitempty"`
	StorePickUp  bool          `json:"store_pick_up,omitempty"`
}
