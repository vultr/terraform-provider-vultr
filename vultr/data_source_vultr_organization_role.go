package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrOrganizationRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrOrganizationPolicyRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_session_duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policies": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrOrganizationRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var roleList []govultr.OrganizationRole
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{PerPage: 10}
	for {
		list, meta, _, err := client.Organization.ListRoles(ctx, options)
		if err != nil {
			return diag.Errorf("error getting organization roles : %v", err)
		}

		for i := range list {
			sm, err := structToMap(list[i])

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				roleList = append(roleList, list[i])
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(roleList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(roleList) < 1 {
		return diag.Errorf("no results were found")
	}

	policyList, _, _, err := client.Organization.ListRolePolicies(ctx, roleList[0].ID, nil)
	if err != nil {
		return diag.Errorf("error getting organization role policies : %v", err)
	}

	var policies []string
	for i := range policyList {
		policies = append(policies, policyList[i].ID)
	}

	if err := d.Set("name", roleList[0].Name); err != nil {
		return diag.Errorf("unable to set organization role `name` read value: %v", err)
	}
	if err := d.Set("description", roleList[0].Description); err != nil {
		return diag.Errorf("unable to set organization role `description` read value: %v", err)
	}
	if err := d.Set("max_session_duration", roleList[0].MaxSessionDuration); err != nil {
		return diag.Errorf("unable to set organization role `max_session_duration` read value: %v", err)
	}
	if err := d.Set("policies", policies); err != nil {
		return diag.Errorf("unable to set organization role `policies` read value: %v", err)
	}
	if err := d.Set("date_created", roleList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set organization `date_created` read value: %v", err)
	}

	return nil
}
