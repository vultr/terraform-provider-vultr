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

func resourceVultrOrganizationRoleGroupAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationRoleGroupAttachmentCreate,
		ReadContext:   resourceVultrOrganizationRoleGroupAttachmentRead,
		DeleteContext: resourceVultrOrganizationRoleGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role_id": {
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

func resourceVultrOrganizationRoleGroupAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	roleID := d.Get("role_id").(string)
	groupID := d.Get("group_id").(string)

	log.Print("[INFO] Creating organization role group attachment")

	if _, _, err := client.Organization.AttachRoleGroup(ctx, roleID, groupID); err != nil {
		return diag.Errorf("error while creating organization role group attachment: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s", roleID, groupID))

	return resourceVultrOrganizationRoleGroupAttachmentRead(ctx, d, meta)
}

func resourceVultrOrganizationRoleGroupAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	roleID := d.Get("role_id").(string)
	groupID := d.Get("group_id").(string)

	options := &govultr.ListOptions{}
	found := false
	for {
		groups, meta, _, err := client.Organization.ListRoleGroups(ctx, roleID, options)
		if err != nil {
			return diag.Errorf("error getting organization role groups : %v", err)
		}

		for i := range groups {
			if groups[i].GroupID == groupID {
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
		tflog.Warn(ctx, fmt.Sprintf("Removing organization role group attachment (%s) because it is gone", d.Id()))
		d.SetId("")
		return nil
	}

	return nil
}

func resourceVultrOrganizationRoleGroupAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	roleID := d.Get("role_id").(string)
	groupID := d.Get("group_id").(string)

	log.Printf("[INFO] Deleting organization role group attachment (%s)", d.Id())
	if _, _, err := client.Organization.DetachRolePolicy(ctx, roleID, groupID); err != nil {
		return diag.Errorf("error deleting organization role group attachment %s : %v", d.Id(), err)
	}

	return nil
}
