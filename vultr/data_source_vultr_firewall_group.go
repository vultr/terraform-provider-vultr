package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrFirewallGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrFirewallGroupRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrFirewallGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	firewallGroupList := []govultr.FirewallGroup{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		firewallGroup, meta, err := client.FirewallGroup.List(context.Background(), options)
		if err != nil {
			return diag.Errorf("error getting firewall group: %v", err)
		}

		for _, fw := range firewallGroup {
			sm, err := structToMap(fw)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				firewallGroupList = append(firewallGroupList, fw)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(firewallGroupList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(firewallGroupList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(firewallGroupList[0].ID)
	d.Set("description", firewallGroupList[0].Description)
	d.Set("date_created", firewallGroupList[0].DateCreated)
	d.Set("date_modified", firewallGroupList[0].DateModified)
	d.Set("instance_count", firewallGroupList[0].InstanceCount)
	d.Set("rule_count", firewallGroupList[0].RuleCount)
	d.Set("max_rule_count", firewallGroupList[0].MaxRuleCount)
	return nil
}
