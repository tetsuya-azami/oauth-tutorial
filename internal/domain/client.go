package domain

type ClientType int

const (
	ConfidentialClient ClientType = iota
	PublicClient
)

type ClientID string

// dynamic registrationはとりあえずサポートしない https://openid.net/specs/openid-connect-registration-1_0.html
type Client struct {
	clientID     ClientID
	clientName   string
	clientType   ClientType
	secret       string
	redirectURIs []string
}

func ReconstructClient(clientID ClientID, clientName string, clientType ClientType, secret string, redirectURIs []string) *Client {
	return &Client{
		clientID:     clientID,
		clientName:   clientName,
		clientType:   clientType,
		secret:       secret,
		redirectURIs: redirectURIs,
	}
}

func (c *Client) ContainsRedirectURI(redirectURI string) bool {
	for _, uri := range c.redirectURIs {
		if uri == redirectURI {
			return true
		}
	}

	return false
}

func (c *Client) ClientID() ClientID     { return c.clientID }
func (c *Client) ClientName() string     { return c.clientName }
func (c *Client) ClientType() ClientType { return c.clientType }
func (c *Client) Secret() string         { return c.secret }
func (c *Client) RedirectURI() []string  { return c.redirectURIs }
