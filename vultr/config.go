package vultr

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/meta"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/vultr/govultr"
)

// Config is the configuration structure used to instantiate Vultr
type Config struct {
	APIKey     string
	RateLimit  int
	RetryLimit int
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

	userAgent := fmt.Sprintf("Terraform/%s", meta.SDKVersionString())

	transport := &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	}
	client := http.DefaultClient
	client.Transport = transport

	client.Transport = logging.NewTransport("Vultr", client.Transport)

	vultrClient := govultr.NewClient(client, c.APIKey)
	vultrClient.SetUserAgent(userAgent)

	if c.RateLimit != 0 {
		vultrClient.SetRateLimit(time.Duration(c.RateLimit) * time.Millisecond)
	}

	if c.RetryLimit != 0 {
		vultrClient.SetRetryLimit(c.RetryLimit)
	}

	return &Client{client: vultrClient}, nil
}
