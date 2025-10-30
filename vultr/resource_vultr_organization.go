package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationCreate,
		ReadContext:   resourceVultrOrganizationRead,
		UpdateContext: resourceVultrOrganizationUpdate,
		DeleteContext: resourceVultrOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	orgReq := &govultr.OrganizationReq{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	log.Print("[INFO] Creating organization")

	org, _, err := client.Organization.CreateOrganization(ctx, orgReq)
	if err != nil {
		return diag.Errorf("error while creating organization : %s", err)
	}

	d.SetId(org.ID)

	return resourceVultrOrganizationRead(ctx, d, meta)
}

func resourceVultrOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	org, _, err := client.Organization.GetOrganization(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Organization not found.") {
			tflog.Warn(ctx, fmt.Sprintf("Removing organization (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting organization : %v", err)
	}

	if err := d.Set("name", org.Name); err != nil {
		return diag.Errorf("unable to set resource organization `name` read value: %v", err)
	}
	if err := d.Set("date_created", org.DateCreated); err != nil {
		return diag.Errorf("unable to set resource organization `date_created` read value: %v", err)
	}
	if err := d.Set("type", org.Type); err != nil {
		return diag.Errorf("unable to set resource organization `type` read value: %v", err)
	}

	return nil
}

func resourceVultrOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updating organization (%s)", d.Id())

	req := &govultr.OrganizationReq{Name: d.Get("name").(string)}
	if _, _, err := client.Organization.UpdateOrganization(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating organization %s: %v", d.Id(), err)
	}

	return resourceVultrOrganizationRead(ctx, d, meta)
}

func resourceVultrOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting organization (%s)", d.Id())
	if err := client.Organization.DeleteOrganization(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting organization %s: %v", d.Id(), err)
	}

	return nil
}
