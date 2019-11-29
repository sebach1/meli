package meli

import (
	"strings"
	"time"
)

type WebHook struct {
	Resource      string `json:"resource,omitempty"`
	UserId        int    `json:"user_id,omitempty"`
	Topic         string `json:"topic,omitempty"`
	ApplicationId int64  `json:"application_id,omitempty"`
	Attempts      int    `json:"attempts,omitempty"`

	Sent     time.Time `json:"sent,omitempty"`
	Received time.Time `json:"received,omitempty"`
}

func (ml *MeLi) ProcessProductWebHook(wh *WebHook) (*Product, error) {
	return ml.GetProduct(ProductId(wh.ResourceID()))
}

func (wh *WebHook) ResourceID() string {
	if wh.Resource == "" {
		return ""
	}
	idx := strings.LastIndex(wh.Resource, "/")
	return wh.Resource[0:idx]
}
