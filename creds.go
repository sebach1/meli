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

func (c *creds) validateClient() error {
	if err := c.validateServer(); err != nil {
		return err
	}
	if c.Refresh == "" {
		return errNilRefreshToken
	}
	if c.Access == "" {
		return errNilAccessToken
	}
	return nil
}

func (c *creds) validateServer() error {
	if c == nil {
		return errNilCredentials
	}
	if c.ApplicationId == "" {
		return errNilApplicationId
	}
	if c.Secret == "" {
		return errNilSecret
	}
	return nil
}
