package meli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	errNilApplicationId = errors.New("the APPLICATION Id is nil")
	errNilAccessToken   = errors.New("the ACCESS TOKEN is NIL")
	errNilRefreshToken  = errors.New("the REFRESH TOKEN is NIL")
	errNilSecret        = errors.New("the SECRET is NIL")

	errNilProductId    = errors.New("the given PRODUCT ID is NIL")
	errNilProduct      = errors.New("the given PRODUCT TITLE is NIL")
	errNilPictures     = errors.New("the given PRODUCT PICTURES are NIL")
	errNilStock        = errors.New("the given PRODUCT STOCK is NIL")
	errNilPrice        = errors.New("the given PRODUCT PRICE is NIL")
	errNilProductTitle = errors.New("the given PRODUCT is NIL")
	errNilVariant      = errors.New("the given VARIANT is NIL")
	errVariantNotFound = errors.New("the given VARIANT does NOT EXISTS")
	// errNilCategory     = errors.New("the given CATEGORY is NIL")
	errNilCategoryId   = errors.New("the given CATEGORY ID is NIL")
	errNilCombinations = errors.New("the given ATTR COMBINATIONS are NIL")

	errInvalidBuyingMode = errors.New("the BUYING MODE is invalid")
	errInvalidCondition  = errors.New("the CONDITION is invalid")

	errNilVarStock    = errors.New("the VARIANT wanted to be created has NIL STOCK")
	errNilVarPrice    = errors.New("the VARIANT wanted to be created has NIL PRICE")
	errNilVarPictures = errors.New("the VARIANT wanted to be created has NIL PICTURES")

	errIncompatibleVar = errors.New("the given VARIANT is INCOMPATIBLE")

	errRemoteInconsistency = errors.New("the SERVER had an inconsistency while performing a request (status code != real behaviour)")
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
		body.Message = errRemoteInconsistency.Error()
	}
	return body
}
