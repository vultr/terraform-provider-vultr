package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
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
		firewallGroup, meta, _, err := client.FirewallGroup.List(ctx, options)
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
	if err := d.Set("description", firewallGroupList[0].Description); err != nil {
		return diag.Errorf("unable to set firewall_group `description` read value: %v", err)
	}
	if err := d.Set("date_created", firewallGroupList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set firewall_group `date_created` read value: %v", err)
	}
	if err := d.Set("date_modified", firewallGroupList[0].DateModified); err != nil {
		return diag.Errorf("unable to set firewall_group `date_modified` read value: %v", err)
	}
	if err := d.Set("instance_count", firewallGroupList[0].InstanceCount); err != nil {
		return diag.Errorf("unable to set firewall_group `instance_count` read value: %v", err)
	}
	if err := d.Set("rule_count", firewallGroupList[0].RuleCount); err != nil {
		return diag.Errorf("unable to set firewall_group `rule_count` read value: %v", err)
	}
	if err := d.Set("max_rule_count", firewallGroupList[0].MaxRuleCount); err != nil {
		return diag.Errorf("unable to set firewall_group `max_rule_count` read value: %v", err)
	}
	return nil
}
