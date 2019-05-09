package vultr

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vultr/govultr"
)

// Config is the configuration structure used to instantiate Vultr
type Config struct {
	APIKey string
}

// Client wraps govultr
type Client struct {
	client *govultr.Client
}

func (c *Client) govultrClient() *govultr.Client {
	return c.client
}

// Client configures govultr and returns an initialized client
func (c *Config) Client() (*Client, error) {

	userAgent := fmt.Sprintf("Terraform/%s", terraform.VersionString())

	transport := &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	}
	client := http.DefaultClient
	client.Transport = transport

	client.Transport = logging.NewTransport("Vultr", client.Transport)

	vultrClient := govultr.NewClient(client, c.APIKey)
	vultrClient.SetUserAgent(userAgent)
	vultrClient.SetRateLimit(300)

	return &Client{client: vultrClient}, nil
}
