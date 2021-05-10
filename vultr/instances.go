package vultr

import (
	"context"
	"fmt"

	"github.com/vultr/govultr/v2"
)

func getPrivateNetworks(client *govultr.Client, instanceID string) ([]string, error) {
	options := &govultr.ListOptions{}
	var pn []string
	for {
		networks, meta, err := client.Instance.ListPrivateNetworks(context.Background(), instanceID, options)
		if err != nil {
			return nil, fmt.Errorf("error getting list of attached private networks : %v", err)
		}

		if len(networks) == 0 {
			break
		}

		for _, v := range networks {
			pn = append(pn, v.NetworkID)
		}

		if meta.Links.Next == "" {
			break
		}
		options.Cursor = meta.Links.Next
	}
	return pn, nil
}
