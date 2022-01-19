package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrSSHKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrSSHKeyRead,
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

func dataSourceVultrSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	sshKeyList := []govultr.SSHKey{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		sshKeys, meta, err := client.SSHKey.List(context.Background(), options)
		if err != nil {
			return diag.Errorf("error getting SSH keys: %v", err)
		}

		for _, ssh := range sshKeys {
			sm, err := structToMap(ssh)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				sshKeyList = append(sshKeyList, ssh)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}
	if len(sshKeyList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(sshKeyList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(sshKeyList[0].ID)
	d.Set("name", sshKeyList[0].Name)
	d.Set("ssh_key", sshKeyList[0].SSHKey)
	d.Set("date_created", sshKeyList[0].DateCreated)
	return nil
}
