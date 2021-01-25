package vultr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceVultrDNSDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrDNSDomainRead,
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

func dataSourceVultrDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	domain, err := client.Domain.Get(context.Background(), d.Get("domain").(string))
	if err != nil {
		return fmt.Errorf("error getting dns domains: %v", err)
	}

	d.SetId(domain.Domain)
	d.Set("domain", domain.Domain)
	d.Set("date_created", domain.DateCreated)
	return nil
}
