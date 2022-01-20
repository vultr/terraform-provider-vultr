package vultr

import (
	"context"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVultrAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrAccountRead,
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
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"pending_charges": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"last_payment_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_payment_amount": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	account, err := client.Account.Get(ctx)

	if err != nil {
		return diag.Errorf("error getting account info: %v", err)
	}

	d.SetId("account")
	d.Set("name", account.Name)
	d.Set("email", account.Email)
	d.Set("balance", math.Round(float64(account.Balance)*100)/100)
	d.Set("pending_charges", math.Round(float64(account.PendingCharges)*100)/100)
	d.Set("last_payment_date", account.LastPaymentDate)
	d.Set("last_payment_amount", math.Round(float64(account.LastPaymentAmount)*100)/100)
	if err := d.Set("acl", account.ACL); err != nil {
		return diag.Errorf("error setting `acls`: %#v", err)
	}
	return nil
}
