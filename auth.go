package meli

import (
	"encoding/json"
	"net/url"
)

type authBody struct {
	AccessToken  accessToken  `json:"access_token,omitempty"`
	TokenType    string       `json:"token_type,omitempty"`
	ExpiresIn    int          `json:"expires_in,omitempty"`
	Scope        string       `json:"scope,omitempty"`
	UserId       int          `json:"user_id,omitempty"`
	RefreshToken refreshToken `json:"refresh_token,omitempty"`
}

func (ml *MeLi) RefreshToken() error {
	err := ml.Credentials.validate()
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", string(ml.Credentials.Refresh))
	params.Set("client_id", string(ml.Credentials.ApplicationId))
	params.Set("client_secret", string(ml.Credentials.Secret))
	URL, err := ml.RouteTo("/oauth/token", params)
	if err != nil {
		return err
	}
	resp, err := ml.Post(URL, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return errFromReader(resp.Body)
	}
	body := &authBody{}
	err = json.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		return err
	}

	if body.AccessToken == "" || body.RefreshToken == "" {
		return errRemoteInconsistency
	}
	ml.Credentials.Access = body.AccessToken
	ml.Credentials.Refresh = body.RefreshToken
	return nil
}
