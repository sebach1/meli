package meli

type Picture struct {
	Id        string `json:"id,omitempty"`
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secure_url,omitempty"`
	Size      string `json:"size,omitempty"`
	MaxSize   string `json:"max_size,omitempty"`
	Quality   string `json:"quality,omitempty"`
}
