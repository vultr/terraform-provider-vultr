package vultr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVultrAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrAccountRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acl": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"balance": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pending_charges": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_payment_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_payment_amount": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	account, err := client.Account.Get(context.Background())

	if err != nil {
		return fmt.Errorf("error getting account info: %v", err)
	}

	d.SetId("account")
	d.Set("name", account.Name)
	d.Set("email", account.Email)
	d.Set("balance", account.Balance)
	d.Set("pending_charges", account.PendingCharges)
	d.Set("last_payment_date", account.LastPaymentDate)
	d.Set("last_payment_amount", account.LastPaymentAmount)
	if err := d.Set("acl", account.ACL); err != nil {
		return fmt.Errorf("error setting `acls`: %#v", err)
	}
	return nil
}
