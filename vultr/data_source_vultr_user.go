package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
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
				Type:     schema.TypeString,
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

	users, err := client.User.List(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting applications: %v", err)
	}

	userList := []govultr.User{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))

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

	if len(userList) > 1 {
		return fmt.Errorf("your search returned too many results : %d. Please refine your search to be more specific", len(userList))
	}

	if len(userList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(userList[0].UserID)
	d.Set("name", userList[0].Name)
	d.Set("email", userList[0].Email)
	d.Set("api_enabled", userList[0].APIEnabled)
	if err := d.Set("acl", userList[0].ACL); err != nil {
		return fmt.Errorf("Error setting `acl`: %#v", err)
	}
	return nil
}
