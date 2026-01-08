package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrOrganizationGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrOrganizationGroupRead,
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
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"roles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceVultrOrganizationGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var groupList []govultr.OrganizationGroup
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{PerPage: 10}
	for {
		list, meta, _, err := client.Organization.ListGroups(ctx, options)
		if err != nil {
			return diag.Errorf("error getting organization groups : %v", err)
		}

		for i := range list {
			sm, err := structToMap(list[i])

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				groupList = append(groupList, list[i])
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(groupList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(groupList) < 1 {
		return diag.Errorf("no results were found")
	}

	var users []string
	for i := range groupList[0].Members {
		users = append(users, groupList[0].Members[i].ID)
	}

	polList, _, _, err := client.Organization.ListGroupPolicies(ctx, groupList[0].ID)
	if err != nil {
		return diag.Errorf("error getting organization group policies : %v", err)
	}

	var pols []string
	for i := range polList.All {
		pols = append(pols, polList.All[i].ID)
	}

	roleList, _, _, err := client.Organization.ListGroupRoles(ctx, groupList[0].ID)
	if err != nil {
		return diag.Errorf("error getting organization group roles : %v", err)
	}

	var roles []string
	allRoles := roleList.All
	for i := range allRoles {
		roles = append(roles, allRoles[i].ID)
	}

	d.SetId(groupList[0].ID)
	if err := d.Set("name", groupList[0].Name); err != nil {
		return diag.Errorf("unable to set organization group `name` read value: %v", err)
	}
	if err := d.Set("description", groupList[0].Description); err != nil {
		return diag.Errorf("unable to set organization group `description` read value: %v", err)
	}
	if err := d.Set("users", users); err != nil {
		return diag.Errorf("unable to set organization group `users` read value: %v", err)
	}
	if err := d.Set("policies", pols); err != nil {
		return diag.Errorf("unable to set organization group `policies` read value: %v", err)
	}
	if err := d.Set("roles", roles); err != nil {
		return diag.Errorf("unable to set organization group `roles` read value: %v", err)
	}
	if err := d.Set("date_created", groupList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set organization `date_created` read value: %v", err)
	}

	return nil
}
