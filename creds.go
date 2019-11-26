package meli

type token string

type applicationId string
type refreshToken token
type accessToken token

type creds struct {
	Access        accessToken
	Refresh       refreshToken
	ApplicationId applicationId
	Secret        token
}

func (c *creds) validate() error {
	if c.ApplicationId == "" {
		return errNilApplicationId
	}
	if c.Secret == "" {
		return errNilSecret
	}
	if c.Refresh == "" {
		return errNilRefreshToken
	}
	if c.Access == "" {
		return errNilAccessToken
	}
	return nil
}
