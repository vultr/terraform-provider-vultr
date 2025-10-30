package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrOrganizationPolicy() *schema.Resource {
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
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_system_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"document": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"statement": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"effect": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"action": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"resource": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"users": {
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

func dataSourceVultrOrganizationPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var policyList []govultr.OrganizationPolicy
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		list, meta, _, err := client.Organization.ListPolicies(ctx, options)
		if err != nil {
			return diag.Errorf("error getting organization policies : %v", err)
		}

		for i := range list {
			sm, err := structToMap(list[i])

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				policyList = append(policyList, list[i])
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(policyList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(policyList) < 1 {
		return diag.Errorf("no results were found")
	}

	groupList, _, _, err := client.Organization.ListPolicyGroups(ctx, policyList[0].ID, nil)
	if err != nil {
		return diag.Errorf("error getting organization policy groups : %v", err)
	}

	var groups []string
	for i := range groupList {
		groups = append(groups, groupList[i].ID)
	}

	userList, _, _, err := client.Organization.ListPolicyUsers(ctx, policyList[0].ID, nil)
	if err != nil {
		return diag.Errorf("error getting organization policy users : %v", err)
	}

	var users []string
	for i := range userList {
		users = append(users, userList[i].ID)
	}

	var statementFlat []map[string]interface{}
	for i := range policyList[0].Document.Statement {
		statementFlat = append(statementFlat, map[string]interface{}{
			"effect":   policyList[0].Document.Statement[i].Effect,
			"action":   policyList[0].Document.Statement[i].Action,
			"resource": policyList[0].Document.Statement[i].Resource,
		})
	}

	policyDocument := []map[string]interface{}{
		{
			"version":   policyList[0].Document.Version,
			"statement": statementFlat,
		},
	}

	if err := d.Set("name", policyList[0].Name); err != nil {
		return diag.Errorf("unable to set organization policy `name` read value: %v", err)
	}
	if err := d.Set("description", policyList[0].Description); err != nil {
		return diag.Errorf("unable to set organization policy `description` read value: %v", err)
	}
	if err := d.Set("is_system_policy", policyList[0].SystemPolicy); err != nil {
		return diag.Errorf("unable to set organization policy `is_system_policy` read value: %v", err)
	}
	if err := d.Set("document", policyDocument); err != nil {
		return diag.Errorf("unable to set organization policy `document` read value: %v", err)
	}
	if err := d.Set("users", users); err != nil {
		return diag.Errorf("unable to set organization policy `users` read value: %v", err)
	}
	if err := d.Set("groups", groups); err != nil {
		return diag.Errorf("unable to set organization policy `groups` read value: %v", err)
	}
	if err := d.Set("date_created", policyList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set organization `date_created` read value: %v", err)
	}

	return nil
}
