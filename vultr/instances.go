package vultr

import (
	"context"
	"fmt"

	"github.com/vultr/govultr/v2"
)

func getVPCs(client *govultr.Client, instanceID string) ([]string, error) {
	options := &govultr.ListOptions{}
	var vpcs []string
	for {
		vpcInfo, meta, err := client.Instance.ListVPCInfo(context.Background(), instanceID, options)
		if err != nil {
			return nil, fmt.Errorf("error getting list of attached VPCs: %v", err)
		}

		if len(vpcInfo) == 0 {
			break
		}

		for _, v := range vpcInfo {
			log.Printf("OUT vpc ID: %s\n", v.ID)
			log.Printf("OUT vpc MAC: %s\n", v.MacAddress)
			log.Printf("OUT vpc ip: %s\n", v.IPAddress)
			vpcs = append(vpcs, v.ID)
		}

		if meta.Links.Next == "" {
			break
		}
		options.Cursor = meta.Links.Next
	}
	return vpcs, nil
}
