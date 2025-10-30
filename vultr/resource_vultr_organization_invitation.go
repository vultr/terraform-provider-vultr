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

func resourceVultrOrganizationInvitation() *schema.Resource {
	return &schema.Resource{CreateContext: resourceVultrOrganizationInvitationCreate,
		ReadContext: resourceVultrOrganizationInvitationRead,
		// UpdateContext: resourceVultrOrganizationInvitationUpdate,
		DeleteContext: resourceVultrOrganizationInvitationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// TODO: roles aren't supported in the API
			// "roles": {
			// 	Type:     schema.TypeSet,
			// 	Required: true,
			// 	ForceNew: true,
			// 	Elem:     &schema.Schema{Type: schema.TypeString},
			// },
			"api_enabled": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"email_registered": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_responded": {
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

func resourceVultrOrganizationInvitationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	// var roles []string
	// if roleObj, roleOK := d.GetOk("roles"); roleOK {
	// 	roleVals := roleObj.(*schema.Set).List()
	// 	if len(roleVals) != 0 {
	// 		for i := range roleVals {
	// 			roles = append(roles, roleVals[i].(string))
	// 		}
	// 	}
	// }

	invReq := &govultr.OrganizationInvitationReq{
		Email: d.Get("email").(string),
		Permissions: govultr.OrganizationInvitationPermission{
			APIEnabled: d.Get("api_enabled").(bool),
			// Roles:      roles,
		},
	}

	log.Print("[INFO] Creating organization invitation")

	inv, _, err := client.Organization.CreateInvitation(ctx, invReq)
	if err != nil {
		return diag.Errorf("error while creating organization invitation : %s", err)
	}

	d.SetId(inv.ID)

	return resourceVultrOrganizationInvitationRead(ctx, d, meta)
}

func resourceVultrOrganizationInvitationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	inv, _, err := client.Organization.GetInvitation(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invite not found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing organization invitation (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting organization invitation : %v", err)
	}

	if err := d.Set("email_registered", inv.EmailRegistered); err != nil {
		return diag.Errorf("unable to set resource organization invitation `email_registered` read value: %v", err)
	}
	if err := d.Set("status", inv.Status); err != nil {
		return diag.Errorf("unable to set resource organization invitation `status` read value: %v", err)
	}
	if err := d.Set("date_created", inv.DateCreated); err != nil {
		return diag.Errorf("unable to set resource organization invitation `date_created` read value: %v", err)
	}
	if err := d.Set("date_responded", inv.DateResponded); err != nil {
		return diag.Errorf("unable to set resource organization invitation `date_responded` read value: %v", err)
	}
	if err := d.Set("date_expires", inv.DateExpiration); err != nil {
		return diag.Errorf("unable to set resource organization invitation `date_expires` read value: %v", err)
	}

	return nil
}

func resourceVultrOrganizationInvitationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceVultrOrganizationInvitationRead(ctx, d, meta)
}

func resourceVultrOrganizationInvitationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
