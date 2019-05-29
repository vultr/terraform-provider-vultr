package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrSSHKeyRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrSSHKeyRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	sshKeys, err := client.SSHKey.List(context.Background())

	if err != nil {
		return fmt.Errorf("error getting SSH keys: %v", err)
	}

	sshKeyList := []govultr.SSHKey{}

	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	for _, ssh := range sshKeys {
		sm, err := structToMap(ssh)

		if err != nil {
			return err
		}

		if filterLoop(f, sm) {
			sshKeyList = append(sshKeyList, ssh)
		}
	}

	if len(sshKeyList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(sshKeyList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(sshKeyList[0].SSHKeyID)
	d.Set("name", sshKeyList[0].Name)
	d.Set("ssh_key", sshKeyList[0].Key)
	d.Set("date_created", sshKeyList[0].DateCreated)
	return nil
}
