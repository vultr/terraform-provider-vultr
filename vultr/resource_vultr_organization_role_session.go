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

func resourceVultrOrganizationRoleSession() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationRoleSessionCreate,
		ReadContext:   resourceVultrOrganizationRoleSessionRead,
		DeleteContext: resourceVultrOrganizationRoleSessionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"session_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conditions_met": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"remaining_duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auth_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_assumed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_expires": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrOrganizationRoleSessionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	sessionReq := &govultr.OrganizationRoleSessionReq{
		UserID:      d.Get("user").(string),
		RoleID:      d.Get("role").(string),
		SessionName: d.Get("session_name").(string),
		Duration:    d.Get("duration").(int),
		Context: govultr.OrganizationRoleSessionReqContext{
			IPAddress: d.Get("ip_address").(string),
		},
	}

	log.Print("[INFO] Creating organization role session")

	session, _, err := client.Organization.CreateRoleSession(ctx, sessionReq)
	if err != nil {
		return diag.Errorf("error while creating organization role session: %s", err)
	}

	d.SetId(session.Token)

	return resourceVultrOrganizationRoleSessionRead(ctx, d, meta)
}

func resourceVultrOrganizationRoleSessionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	session, _, err := client.Organization.GetRoleSession(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Assumed Role Not Found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing organization role session (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting organization role session : %v", err)
	}

	if err := d.Set("session_name", session.SessionName); err != nil {
		return diag.Errorf("unable to set resource organization role session `session_name` read value: %v", err)
	}
	if err := d.Set("user", session.UserID); err != nil {
		return diag.Errorf("unable to set resource organization role session `user` read value: %v", err)
	}
	if err := d.Set("role", session.RoleID); err != nil {
		return diag.Errorf("unable to set resource organization role session `role` read value: %v", err)
	}
	if err := d.Set("remaining_duration", session.RemainingDuration); err != nil {
		return diag.Errorf("unable to set resource organization role session `remaining_duration` read value: %v", err)
	}
	if err := d.Set("source_ip", session.SourceIP); err != nil {
		return diag.Errorf("unable to set resource organization role session `source_ip` read value: %v", err)
	}
	if err := d.Set("auth_method", session.AuthMethod); err != nil {
		return diag.Errorf("unable to set resource organization role session `auth_method` read value: %v", err)
	}
	if err := d.Set("conditions_me", session.ConditionsMet); err != nil {
		return diag.Errorf("unable to set resource organization role session `conditions_met` read value: %v", err)
	}
	if err := d.Set("date_assumed", session.DateAssumed); err != nil {
		return diag.Errorf("unable to set resource organization role session `date_assumed` read value: %v", err)
	}
	if err := d.Set("date_expires", session.DateExpires); err != nil {
		return diag.Errorf("unable to set resource organization role session `date_assumed` read value: %v", err)
	}

	return nil
}

func resourceVultrOrganizationRoleSessionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting organization role session (%s)", d.Id())
	if err := client.Organization.RevokeRoleSession(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting organization role session %s : %v", d.Id(), err)
	}

	return nil
}
