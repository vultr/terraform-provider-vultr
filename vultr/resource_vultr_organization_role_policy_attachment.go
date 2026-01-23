package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrOrganizationRolePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationRolePolicyAttachmentCreate,
		ReadContext:   resourceVultrOrganizationRolePolicyAttachmentRead,
		DeleteContext: resourceVultrOrganizationRolePolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVultrOrganizationRolePolicyAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	roleID := d.Get("role_id").(string)
	policyID := d.Get("policy_id").(string)

	log.Print("[INFO] Creating organization role policy attachment")

	_, _, err := client.Organization.AttachRolePolicy(
		ctx,
		roleID,
		&govultr.OrganizationRolePolicyReq{PolicyID: policyID},
	)
	if err != nil {
		return diag.Errorf("error while creating organization role policy attachment: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s", roleID, policyID))

	return resourceVultrOrganizationRolePolicyAttachmentRead(ctx, d, meta)
}

func resourceVultrOrganizationRolePolicyAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	roleID := d.Get("role_id").(string)
	policyID := d.Get("policy_id").(string)

	options := &govultr.ListOptions{}
	found := false
	for {
		policies, meta, _, err := client.Organization.ListRolePolicies(ctx, roleID, options)
		if err != nil {
			return diag.Errorf("error getting organization role policies : %v", err)
		}

		for i := range policies {
			if policies[i].PolicyID == policyID {
				found = true
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if !found {
		tflog.Warn(ctx, fmt.Sprintf("Removing organization role policy attachment (%s) because it is gone", d.Id()))
		d.SetId("")
		return nil
	}

	return nil
}

func resourceVultrOrganizationRolePolicyAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	roleID := d.Get("role_id").(string)
	policyID := d.Get("policy_id").(string)

	log.Printf("[INFO] Deleting organization role policy attachment (%s)", d.Id())
	if _, _, err := client.Organization.DetachRolePolicy(ctx, roleID, policyID); err != nil {
		return diag.Errorf("error deleting organization role policy attachment %s : %v", d.Id(), err)
	}

	return nil
}
