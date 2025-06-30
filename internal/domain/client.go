package domain

// dynamic registrationはとりあえずサポートしない https://openid.net/specs/openid-connect-registration-1_0.html
type Client struct {
	clientID     string
	clientName   string
	secret       string
	redirectURIs []string
}

func (c *Client) ClientID() string      { return c.clientID }
func (c *Client) ClientName() string    { return c.clientName }
func (c *Client) Secret() string        { return c.secret }
func (c *Client) RedirectURI() []string { return c.redirectURIs }
