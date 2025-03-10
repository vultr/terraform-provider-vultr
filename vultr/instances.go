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
		vpcInfo, meta, _, err := client.Instance.ListVPCInfo(context.Background(), instanceID, options)
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

func getVPC2s(client *govultr.Client, instanceID string) ([]string, error) {
	options := &govultr.ListOptions{}
	var vpcs []string
	for {
		vpcInfo, meta, _, err := client.Instance.ListVPC2Info(context.Background(), instanceID, options) //nolint:staticcheck
		if err != nil {
			return nil, fmt.Errorf("error getting list of attached VPCs 2.0: %v", err)
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

func getBareMetalServerVPC2s(client *govultr.Client, serverID string) ([]string, error) {
	var vpcs []string

	vpcInfo, _, err := client.BareMetalServer.ListVPC2Info(context.Background(), serverID) //nolint:staticcheck
	if err != nil {
		return nil, fmt.Errorf("error getting list of attached VPCs 2.0: %v", err)
	}

	for _, v := range vpcInfo {
		vpcs = append(vpcs, v.ID)
	}

	return vpcs, nil
}
