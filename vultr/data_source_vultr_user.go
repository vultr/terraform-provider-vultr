package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrUserRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"acl": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceVultrUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	options := &govultr.ListOptions{}
	userList := []govultr.User{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	for {
		users, meta, err := client.User.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting users: %v", err)
		}

		for _, u := range users {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(u)
			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				userList = append(userList, u)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(userList) > 1 {
		return diag.Errorf("your search returned too many results : %d. Please refine your search to be more specific", len(userList))
	}
	if len(userList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(userList[0].ID)
	d.Set("name", userList[0].Name)
	d.Set("email", userList[0].Email)
	d.Set("api_enabled", userList[0].APIEnabled)
	if err := d.Set("acl", userList[0].ACL); err != nil {
		return diag.Errorf("error setting `acl`: %#v", err)
	}
	return nil
}
