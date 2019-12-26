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
		return ErrNilRefreshToken
	}
	if c.Access == "" {
		return ErrNilAccessToken
	}
	return nil
}

func (c *creds) validateServer() error {
	if c == nil {
		return ErrNilCredentials
	}
	if c.ApplicationId == "" {
		return ErrNilApplicationId
	}
	if c.Secret == "" {
		return ErrNilSecret
	}
	return nil
}
