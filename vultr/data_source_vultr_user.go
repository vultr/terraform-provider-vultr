package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrUserRead,
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

func dataSourceVultrUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	options := &govultr.ListOptions{}
	userList := []govultr.User{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	for {
		users, meta, err := client.User.List(context.Background(), options)
		if err != nil {
			return fmt.Errorf("error getting users: %v", err)
		}

		for _, u := range users {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(u)
			if err != nil {
				return err
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
		return fmt.Errorf("your search returned too many results : %d. Please refine your search to be more specific", len(userList))
	}
	if len(userList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(userList[0].ID)
	d.Set("name", userList[0].Name)
	d.Set("email", userList[0].Email)
	d.Set("api_enabled", userList[0].APIEnabled)
	if err := d.Set("acl", userList[0].ACL); err != nil {
		return fmt.Errorf("error setting `acl`: %#v", err)
	}
	return nil
}
