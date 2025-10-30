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

func resourceVultrOrganizationGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationGroupCreate,
		ReadContext:   resourceVultrOrganizationGroupRead,
		UpdateContext: resourceVultrOrganizationGroupUpdate,
		DeleteContext: resourceVultrOrganizationGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"roles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policies": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrOrganizationGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	grpReq := &govultr.OrganizationGroupReq{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Print("[INFO] Creating organization group")

	grp, _, err := client.Organization.CreateGroup(ctx, grpReq)
	if err != nil {
		return diag.Errorf("error while creating organization group : %s", err)
	}

	// if members, membersOK := d.GetOk("members"); membersOK {
	// 	memList := members.(*schema.Set).List()
	// 	for i := range memList {
	// 		addReq := &govultr.OrganizationGroupMemberReq{UserID: memList[i].(string)}
	// 		if err := client.Organization.AddGroupMember(ctx, grp.ID, addReq); err != nil {
	// 			log.Print("[ERROR] error adding organization group %s to member %s : %v", grp.ID, memList[i], err)
	// 		}
	// 	}
	// }

	// if roles, rolesOK := d.GetOk("roles"); rolesOK {
	// 	roleList := roles.(*schema.Set).List()
	// 	for i := range roleList {
	// 		if _, _, err := client.Organization.AttachRoleGroup(ctx, roleList[i].(string), grp.ID); err != nil {
	// 			log.Print("[ERROR] error attaching organization group %s to role %s : %v", grp.ID, roleList[i], err)
	// 		}
	// 	}
	// }

	// if policies, policiesOK := d.GetOk("policies"); policiesOK {
	// 	polList := policies.(*schema.Set).List()
	// 	for i := range polList {
	// 		if err := client.Organization.AttachPolicyGroup(ctx, polList[i].(string), grp.ID); err != nil {
	// 			log.Print("[ERROR] error attaching organization group %s to policy %s : %v", grp.ID, polList[i], err)
	// 		}
	// 	}
	// }

	d.SetId(grp.ID)

	return resourceVultrOrganizationGroupRead(ctx, d, meta)
}

func resourceVultrOrganizationGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	grp, _, err := client.Organization.GetGroup(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Group not found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing organization group (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting organization group : %v", err)
	}

	var users []string
	for i := range grp.Members {
		users = append(users, grp.Members[i].ID)
	}

	polList, _, _, err := client.Organization.ListGroupPolicies(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error getting organization group policies : %v", err)
	}

	var pols []string
	for i := range polList.All {
		pols = append(pols, polList.All[i].ID)
	}

	roleList, _, _, err := client.Organization.ListGroupRoles(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error getting organization group roles : %v", err)
	}

	var roles []string
	for i := range roleList {
		roles = append(roles, roleList[i].ID)
	}

	if err := d.Set("name", grp.Name); err != nil {
		return diag.Errorf("unable to set resource organization group `name` read value: %v", err)
	}
	if err := d.Set("description", grp.Description); err != nil {
		return diag.Errorf("unable to set resource organization group `description` read value: %v", err)
	}
	if err := d.Set("users", users); err != nil {
		return diag.Errorf("unable to set resource organization group `users` read value: %v", err)
	}
	if err := d.Set("policies", pols); err != nil {
		return diag.Errorf("unable to set resource organization group `policies` read value: %v", err)
	}
	if err := d.Set("roles", roles); err != nil {
		return diag.Errorf("unable to set resource organization group `roles` read value: %v", err)
	}
	if err := d.Set("date_created", grp.DateCreated); err != nil {
		return diag.Errorf("unable to set resource organization group `date_created` read value: %v", err)
	}

	return nil
}

func resourceVultrOrganizationGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updating organization group (%s)", d.Id())

	req := &govultr.OrganizationGroupReq{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	if _, _, err := client.Organization.UpdateGroup(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating organization group %s : %v", d.Id(), err)
	}

	// if d.HasChange("members") {
	// 	log.Printf("[INFO] Updating organization group members")

	// 	oldMem, newMem := d.GetChange("members")
	// 	oldMemList := oldMem.(*schema.Set).List()
	// 	newMemList := newMem.(*schema.Set).List()

	// 	var oldIDs, newIDs []string
	// 	for i := range oldMemList {
	// 		oldIDs = append(oldIDs, oldMemList[i].(string))
	// 	}

	// 	for i := range newMemList {
	// 		newIDs = append(newIDs, newMemList[i].(string))
	// 	}

	// 	removeIDs := diffSlice(newIDs, oldIDs)
	// 	addIDs := diffSlice(oldIDs, newIDs)

	// 	if len(removeIDs) > 0 {
	// 		for j := range removeIDs {
	// 			if err := client.Organization.RemoveGroupMember(ctx, d.Id(), removeIDs[j]); err != nil {
	// 				return diag.Errorf("error removing organization group %s from member %s : %v", d.Id(), removeIDs[j], err)
	// 			}
	// 		}
	// 	}

	// 	if len(addIDs) > 0 {
	// 		for k := range addIDs {
	// 			addReq := &govultr.OrganizationGroupMemberReq{UserID: addIDs[k]}
	// 			if err := client.Organization.AddGroupMember(ctx, d.Id(), addReq); err != nil {
	// 				return diag.Errorf("error adding organization group %s to member %s : %v", d.Id(), addIDs[k], err)
	// 			}
	// 		}
	// 	}
	// }

	// if d.HasChange("roles") {
	// 	log.Printf("[INFO] Updating organization group roles")

	// 	oldRoles, newRoles := d.GetChange("roles")
	// 	oldRoleList := oldRoles.(*schema.Set).List()
	// 	newRoleList := newRoles.(*schema.Set).List()

	// 	var oldIDs, newIDs []string
	// 	for i := range oldRoleList {
	// 		oldIDs = append(oldIDs, oldRoleList[i].(string))
	// 	}

	// 	for i := range newRoleList {
	// 		newIDs = append(newIDs, newRoleList[i].(string))
	// 	}

	// 	removeIDs := diffSlice(newIDs, oldIDs)
	// 	addIDs := diffSlice(oldIDs, newIDs)

	// 	if len(removeIDs) > 0 {
	// 		for j := range removeIDs {
	// 			if err := client.Organization.DetachRoleGroup(ctx, removeIDs[j], d.Id()); err != nil {
	// 				return diag.Errorf("error detaching organization group %s from role %s : %v", d.Id(), removeIDs[j], err)
	// 			}
	// 		}
	// 	}

	// 	if len(addIDs) > 0 {
	// 		for k := range addIDs {
	// 			if _, _, err := client.Organization.AttachRoleGroup(ctx, addIDs[k], d.Id()); err != nil {
	// 				return diag.Errorf("error attaching organization group %s to role %s : %v", d.Id(), addIDs[k], err)
	// 			}
	// 		}
	// 	}
	// }

	// if d.HasChange("policies") {
	// 	log.Printf("[INFO] Updating organization group policies")

	// 	oldPolicies, newPolicies := d.GetChange("policies")
	// 	oldPolicyList := oldPolicies.(*schema.Set).List()
	// 	newPolicyList := newPolicies.(*schema.Set).List()

	// 	var oldIDs, newIDs []string
	// 	for i := range oldPolicyList {
	// 		oldIDs = append(oldIDs, oldPolicyList[i].(string))
	// 	}

	// 	for i := range newPolicyList {
	// 		newIDs = append(newIDs, newPolicyList[i].(string))
	// 	}

	// 	removeIDs := diffSlice(newIDs, oldIDs)
	// 	addIDs := diffSlice(oldIDs, newIDs)

	// 	if len(removeIDs) > 0 {
	// 		for i := range removeIDs {
	// 			if err := client.Organization.DetachPolicyGroup(ctx, removeIDs[i], d.Id()); err != nil {
	// 				return diag.Errorf("error detaching organization group %s from policy %s : %v", d.Id(), removeIDs[i], err)
	// 			}
	// 		}
	// 	}

	// 	if len(addIDs) > 0 {
	// 		for i := range addIDs {
	// 			if err := client.Organization.AttachPolicyGroup(ctx, addIDs[i], d.Id()); err != nil {
	// 				return diag.Errorf("error attaching organization group %s to policy %s : %v", d.Id(), addIDs[i], err)
	// 			}
	// 		}
	// 	}
	// }

	return resourceVultrOrganizationRead(ctx, d, meta)
}

func resourceVultrOrganizationGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting organization group (%s)", d.Id())
	if err := client.Organization.DeleteGroup(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting organization group %s : %v", d.Id(), err)
	}

	return nil
}
