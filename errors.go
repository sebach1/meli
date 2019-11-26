package meli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	errNonexistantPath = errors.New("the given PATH does NOT EXISTS")

	errNilApplicationId = errors.New("the APPLICATION Id is nil")
	errNilAccessToken   = errors.New("the ACCESS TOKEN is NIL")
	errNilRefreshToken  = errors.New("the REFRESH TOKEN is NIL")
	errNilSecret        = errors.New("the SECRET is NIL")

	errNilProduct = errors.New("the given PRODUCT is NIL")

	errRemoteInconsistency = errors.New("the SERVER had an inconsistency while performing a request (status code != real behavior)")
)

type Error struct {
	Message     string `json:"message,omitempty"`
	ResponseErr string `json:"error,omitempty"`
	Status      int    `json:"status,omitempty"`
}

func (body *Error) Error() string {
	return fmt.Sprintf("%s: %s", body.ResponseErr, body.Message)
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
