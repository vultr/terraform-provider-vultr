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

func resourceVultrOrganizationRoleTrust() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationRoleTrustCreate,
		ReadContext:   resourceVultrOrganizationRoleTrustRead,
		UpdateContext: resourceVultrOrganizationRoleTrustUpdate,
		DeleteContext: resourceVultrOrganizationRoleTrustDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_range": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hour_start": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"hour_end": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"date_expires": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrOrganizationRoleTrustCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	var ipRanges []string
	if ips, ipsOK := d.GetOk("ip_range"); ipsOK {
		ipsVal := ips.(*schema.Set).List()
		for i := range ipsVal {
			ipRanges = append(ipRanges, ipsVal[i].(string))
		}
	}

	trustReq := &govultr.OrganizationRoleTrustReq{
		UserID:  d.Get("user").(string),
		GroupID: d.Get("group").(string),
		RoleID:  d.Get("role").(string),
		Type:    d.Get("type").(string),
		Conditions: govultr.OrganizationRoleTrustCondition{
			TimeOfDay: govultr.OrganizationRoleTrustConditionTime{
				Start: d.Get("hour_start").(int),
				End:   d.Get("hour_end").(int),
			},
			IPRanges: ipRanges,
		},
		DateExpires: d.Get("date_expires").(string),
	}

	log.Print("[INFO] Creating organization role trust")

	trust, _, err := client.Organization.CreateRoleTrust(ctx, trustReq)
	if err != nil {
		return diag.Errorf("error while creating organization role trust: %s", err)
	}

	d.SetId(trust.ID)

	return resourceVultrOrganizationRoleTrustRead(ctx, d, meta)
}

func resourceVultrOrganizationRoleTrustRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	trust, _, err := client.Organization.GetRoleTrust(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Role trust not found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing organization role trust (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting organization role trust : %v", err)
	}

	if err := d.Set("type", trust.Type); err != nil {
		return diag.Errorf("unable to set resource organization role trust `type` read value: %v", err)
	}
	if err := d.Set("hour_start", trust.Conditions.TimeOfDay.Start); err != nil {
		return diag.Errorf("unable to set resource organization role trust `hour_start` read value: %v", err)
	}
	if err := d.Set("hour_end", trust.Conditions.TimeOfDay.End); err != nil {
		return diag.Errorf("unable to set resource organization role trust `hour_end` read value: %v", err)
	}
	if err := d.Set("ip_range", trust.Conditions.IPRanges); err != nil {
		return diag.Errorf("unable to set resource organization role trust `ip_range` read value: %v", err)
	}
	if err := d.Set("date_created", trust.DateCreated); err != nil {
		return diag.Errorf("unable to set resource organization role trust `date_created` read value: %v", err)
	}

	return nil
}

func resourceVultrOrganizationRoleTrustUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updating organization role trust (%s)", d.Id())

	var ips []string
	ipRangeVal := d.Get("ip_range").(*schema.Set).List()
	if len(ipRangeVal) != 0 {
		for i := range ipRangeVal {
			ips = append(ips, ipRangeVal[i].(string))
		}
	}

	req := &govultr.OrganizationRoleTrustReq{
		Type: d.Get("type").(string),
		Conditions: govultr.OrganizationRoleTrustCondition{
			TimeOfDay: govultr.OrganizationRoleTrustConditionTime{
				Start: d.Get("hour_start").(int),
				End:   d.Get("hour_end").(int),
			},
			IPRanges: ips,
		},
	}

	if _, _, err := client.Organization.UpdateRoleTrust(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating organization role trust %s : %v", d.Id(), err)
	}

	return resourceVultrOrganizationRoleTrustRead(ctx, d, meta)
}

func resourceVultrOrganizationRoleTrustDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting organization role trust (%s)", d.Id())
	if err := client.Organization.DeleteRoleTrust(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting organization role trust %s : %v", d.Id(), err)
	}

	return nil
}
