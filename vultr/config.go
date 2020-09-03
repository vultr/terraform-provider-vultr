package vultr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/meta"
	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"
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
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.APIKey,
	})

	client := oauth2.NewClient(context.Background(), tokenSrc)
	client.Transport = logging.NewTransport("Vultr", client.Transport)

	vultrClient := govultr.NewClient(client)
	vultrClient.SetUserAgent(userAgent)

	if c.RateLimit != 0 {
		vultrClient.SetRateLimit(time.Duration(c.RateLimit) * time.Millisecond)
	}

	if c.RetryLimit != 0 {
		vultrClient.SetRetryLimit(c.RetryLimit)
	}

	return &Client{client: vultrClient}, nil
}
