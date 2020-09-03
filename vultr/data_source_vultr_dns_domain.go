package vultr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceVultrDnsDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrDnsDomainRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the domain",
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

	domain, err := client.Domain.Get(context.Background(), d.Get("domain").(string))
	if err != nil {
		return fmt.Errorf("error getting dns domains: %v", err)
	}

	d.SetId(domain.Domain)
	d.Set("date_created", domain.DateCreated)
	return nil
}
