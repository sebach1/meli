package meli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	ErrNilCredentials   = errors.New("the CREDENTIALS are NIL")
	ErrNilApplicationId = errors.New("the APPLICATION ID is NIL")
	ErrNilAccessToken   = errors.New("the ACCESS TOKEN is NIL")
	ErrNilRefreshToken  = errors.New("the REFRESH TOKEN is NIL")
	ErrNilSecret        = errors.New("the SECRET is NIL")

	ErrNilProductId         = errors.New("the given PRODUCT ID is NIL")
	ErrNilProduct           = errors.New("the given PRODUCT TITLE is NIL")
	ErrNilPictures          = errors.New("the given PRODUCT PICTURES are NIL")
	ErrNilStock             = errors.New("the given PRODUCT STOCK is NIL")
	ErrNilPrice             = errors.New("the given PRODUCT PRICE is NIL")
	ErrNilProductTitle      = errors.New("the given PRODUCT is NIL")
	ErrNilVariant           = errors.New("the given VARIANT is NIL")
	ErrVariantNotFound      = errors.New("the given VARIANT does NOT EXISTS")
	ErrNilCategory          = errors.New("the given CATEGORY is NIL")
	ErrInvalidListingTypeId = errors.New("the given LISTING TYPE ID is INVALID")

	ErrInvalidCategoryId = errors.New("the given CATEGORY ID is INVALID")
	ErrNilCategoryId     = errors.New("the given CATEGORY ID is NIL")
	ErrNilCombinations   = errors.New("the given ATTR COMBINATIONS are NIL")

	ErrInvalidBuyingMode = errors.New("the BUYING MODE is invalid")
	ErrInvalidCondition  = errors.New("the CONDITION is invalid")

	ErrNilVarStock    = errors.New("the VARIANT wanted to be created has NIL STOCK")
	ErrNilVarPrice    = errors.New("the VARIANT wanted to be created has NIL PRICE")
	ErrNilVarPictures = errors.New("the VARIANT wanted to be created has NIL PICTURES")

	ErrIncompatibleVar = errors.New("the given VARIANT is INCOMPATIBLE")

	ErrRemoteInconsistency = errors.New("the SERVER had an inconsistency while performing a request (status code != real behaviour)")

	ErrInvalidMultigetQuantity = errors.New("invalid quantity of elements for multiget request type")
)

type Error struct {
	Message     string      `json:"message,omitempty"`
	ResponseErr string      `json:"error,omitempty"`
	Status      int         `json:"status,omitempty"`
	Cause       []*errCause `json:"cause,omitempty"`
}

type errCause struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (svErr *Error) Error() string {
	if len(svErr.Cause) == 0 {
		return fmt.Sprintf("%s: %s", svErr.ResponseErr, svErr.Message)
	}
	var strErr string
	for _, cause := range svErr.Cause {
		strErr += fmt.Sprintf("%s: %s", cause.Code, cause.Message)
		strErr += "; "
	}
	return strErr
}

func errFromReader(stream io.Reader) error {
	body := &Error{}
	err := json.NewDecoder(stream).Decode(body)
	if err != nil {
		return err
	}
	if body.ResponseErr == "" && body.Message == "" {
		body.ResponseErr = "remote_inconsistency"
		body.Message = ErrRemoteInconsistency.Error()
	}
	return body
}
