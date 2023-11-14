package vultr

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/vultr/govultr/v3"
	"golang.org/x/oauth2"
)

const terraformSDKPath string = "github.com/hashicorp/terraform-plugin-sdk/v2"

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

// terraformSDKVersion looks up the module version of the Terraform SDK for use
// in the User Agent client string
func terraformSDKVersion() string {
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return "0.0.0"
	}

	for _, module := range i.Deps {
		if module.Path == terraformSDKPath {
			return module.Version
		}
	}

	return "0.0.0"
}

// Client configures govultr and returns an initialized client
func (c *Config) Client() (*Client, error) {
	userAgent := fmt.Sprintf("Terraform/%s", terraformSDKVersion())
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.APIKey,
	})

	client := oauth2.NewClient(context.Background(), tokenSrc)
	client.Transport = logging.NewSubsystemLoggingHTTPTransport("Vultr", client.Transport)

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
