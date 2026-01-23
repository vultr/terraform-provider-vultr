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

func resourceVultrOrganizationPolicyGroupAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationPolicyGroupAttachmentCreate,
		ReadContext:   resourceVultrOrganizationPolicyGroupAttachmentRead,
		DeleteContext: resourceVultrOrganizationPolicyGroupAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVultrOrganizationPolicyGroupAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Print("[INFO] Creating organization policy group attachment")

	policyID := d.Get("policy_id").(string)
	groupID := d.Get("group_id").(string)

	if err := client.Organization.AttachPolicyGroup(ctx, policyID, groupID); err != nil {
		return diag.Errorf("error creating organization attachment of policy %s to group %s : %v", policyID, groupID, err)
	}

	d.SetId(fmt.Sprintf("%s_%s", policyID, groupID))

	return resourceVultrOrganizationPolicyGroupAttachmentRead(ctx, d, meta)
}

func resourceVultrOrganizationPolicyGroupAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	policyID := d.Get("policy_id").(string)
	groupID := d.Get("group_id").(string)

	options := &govultr.ListOptions{}
	found := false
	for {
		groups, meta, _, err := client.Organization.ListPolicyGroups(ctx, policyID, options)
		if err != nil {
			return diag.Errorf("error getting organization policy groups : %v", err)
		}

		for i := range groups {
			if groups[i].ID == groupID {
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
		tflog.Warn(ctx, fmt.Sprintf("Removing organization policy group attachment (%s) because it is gone", d.Id()))
		d.SetId("")
		return nil
	}

	return nil
}

func resourceVultrOrganizationPolicyGroupAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Print("[INFO] Deleting organization policy group attachment")

	policyID := d.Get("policy_id").(string)
	groupID := d.Get("group_id").(string)

	log.Printf("[INFO] Deleting organization policy group attachment (%s)", d.Id())
	if err := client.Organization.DetachPolicyGroup(ctx, policyID, groupID); err != nil {
		return diag.Errorf("error deleting organization policy group attachment %s : %v", d.Id(), err)
	}

	return nil
}
