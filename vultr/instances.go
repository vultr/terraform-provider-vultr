package vultr

import (
	"context"
	"fmt"

	"github.com/vultr/govultr/v3"
)

func getVPCs(client *govultr.Client, instanceID string) ([]string, error) {
	options := &govultr.ListOptions{}
	var vpcs []string
	for {
		vpcInfo, meta,_, err := client.Instance.ListVPCInfo(context.Background(), instanceID, options)
		if err != nil {
			return nil, fmt.Errorf("error getting list of attached VPCs: %v", err)
		}

		if len(vpcInfo) == 0 {
			break
		}

		for _, v := range vpcInfo {
			vpcs = append(vpcs, v.ID)
		}

		if meta.Links.Next == "" {
			break
		}
		options.Cursor = meta.Links.Next
	}
	return vpcs, nil
}
