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

func resourceVultrOrganizationPolicyUserAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationPolicyUserAttachmentCreate,
		ReadContext:   resourceVultrOrganizationPolicyUserAttachmentRead,
		DeleteContext: resourceVultrOrganizationPolicyUserAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVultrOrganizationPolicyUserAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Print("[INFO] Creating organization policy user attachment")

	policyID := d.Get("policy_id").(string)
	userID := d.Get("user_id").(string)

	if err := client.Organization.AttachPolicyUser(ctx, policyID, userID); err != nil {
		return diag.Errorf("error creating organization attachment of policy %s to user %s : %v", policyID, userID, err)
	}

	d.SetId(fmt.Sprintf("%s_%s", policyID, userID))

	return resourceVultrOrganizationPolicyUserAttachmentRead(ctx, d, meta)
}

func resourceVultrOrganizationPolicyUserAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	policyID := d.Get("policy_id").(string)
	userID := d.Get("user_id").(string)

	options := &govultr.ListOptions{}
	found := false
	for {
		users, meta, _, err := client.Organization.ListPolicyUsers(ctx, policyID, options)
		if err != nil {
			return diag.Errorf("error getting organization policy users: %v", err)
		}

		for i := range users {
			if users[i].ID == userID {
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
		tflog.Warn(ctx, fmt.Sprintf("Removing organization policy user attachment (%s) because it is gone", d.Id()))
		d.SetId("")
		return nil
	}

	return nil
}

func resourceVultrOrganizationPolicyUserAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Print("[INFO] Deleting organization policy user attachment")

	policyID := d.Get("policy_id").(string)
	userID := d.Get("user_id").(string)

	log.Printf("[INFO] Deleting organization policy user attachment (%s)", d.Id())
	if err := client.Organization.DetachPolicyUser(ctx, policyID, userID); err != nil {
		return diag.Errorf("error deleting organization policy user attachment %s : %v", d.Id(), err)
	}

	return nil
}
