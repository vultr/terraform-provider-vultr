package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/vultr/govultr"
)

func dataSourceVultrDnsDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrDnsDomainRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrDnsDomainRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).govultrClient()

	domain := d.Get("domain").(string)

	dnsDomains, err := client.DNSDomain.GetList(context.Background())

	if err != nil {
		return fmt.Errorf("error getting dns domains: %v", err)
	}

	dnsList := []govultr.DNSDomain{}

	for _, d := range dnsDomains {

		if d.Domain == domain {
			dnsList = append(dnsList, d)

		}
	}

	if len(dnsList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(dnsList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(dnsDomains[0].Domain)
	d.Set("date_created", dnsDomains[0].DateCreated)
	return nil
}
