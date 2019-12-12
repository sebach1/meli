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
	err := ml.Credentials.validateClient()
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

// curl -X POST https://api.mercadolibre.com/oauth/token?grant_type=authorization_code
// &client_id=$APP_ID&client_secret=$SECRET_KEY&code=$SERVER_GENERATED_AUTHORIZATION_CODE&redirect_uri=$REDIRECT_URI

func (ml *MeLi) SetCredentialsFromCode(code string, redirectURI string) error {
	err := ml.Credentials.validateServer()
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Set("code", code)
	params.Set("client_id", string(ml.Credentials.ApplicationId))
	params.Set("client_secret", string(ml.Credentials.Secret))
	params.Set("redirect_uri", redirectURI)
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

func (ml *MeLi) GetAuthURL(site SiteId) (string, error) {
	if ml.Credentials.ApplicationId == "" {
		return "", errNilApplicationId
	}
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", string(ml.Credentials.ApplicationId))
	return ml.AuthRouteTo("authorization", params, site)
}
