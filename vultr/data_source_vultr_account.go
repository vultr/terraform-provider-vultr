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
			"balance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pending_charges": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_payment_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_payment_amount": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	account, err := client.Account.GetInfo(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting account info: %v", err)
	}

	d.SetId("account")
	d.Set("balance", account.Balance)
	d.Set("pending_charges", account.PendingCharges)
	d.Set("last_payment_date", account.LastPaymentDate)
	d.Set("last_payment_amount", account.LastPaymentAmount)
	return nil
}
