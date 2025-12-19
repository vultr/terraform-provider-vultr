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

func resourceVultrOrganizationRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationRoleCreate,
		ReadContext:   resourceVultrOrganizationRoleRead,
		UpdateContext: resourceVultrOrganizationRoleUpdate,
		DeleteContext: resourceVultrOrganizationRoleDelete,
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"max_session_duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrOrganizationRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	roleReq := &govultr.OrganizationRoleReq{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Type:               d.Get("type").(string),
		MaxSessionDuration: d.Get("max_session_duration").(int),
	}

	log.Print("[INFO] Creating organization role")

	role, _, err := client.Organization.CreateRole(ctx, roleReq)
	if err != nil {
		return diag.Errorf("error while creating organization role : %s", err)
	}

	d.SetId(role.ID)

	if policies, policiesOK := d.GetOk("policies"); policiesOK {
		polList := policies.(*schema.Set).List()
		for i := range polList {
			polReq := &govultr.OrganizationRolePolicyReq{PolicyID: polList[i].(string)}
			if _, _, err := client.Organization.AttachRolePolicy(ctx, d.Id(), polReq); err != nil {
				return diag.Errorf("error attaching organization role %s to policy %s : %v", role.ID, polList[i], err)
			}
		}
	}

	if groups, groupsOK := d.GetOk("groups"); groupsOK {
		groupsList := groups.(*schema.Set).List()
		for i := range groupsList {
			if _, _, err := client.Organization.AttachRoleGroup(ctx, d.Id(), groupsList[i].(string)); err != nil {
				return diag.Errorf("error attaching organization role %s to group %s : %v", role.ID, groupsList[i], err)
			}
		}
	}

	return resourceVultrOrganizationRoleRead(ctx, d, meta)
}

func resourceVultrOrganizationRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	role, _, err := client.Organization.GetRole(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Role not found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing organization role (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting organization role : %v", err)
	}

	var pol []string
	for i := range role.Policies {
		pol = append(pol, role.Policies[i].ID)
	}

	if err := d.Set("name", role.Name); err != nil {
		return diag.Errorf("unable to set resource organization role `name` read value: %v", err)
	}
	if err := d.Set("description", role.Description); err != nil {
		return diag.Errorf("unable to set resource organization role `description` read value: %v", err)
	}
	if err := d.Set("type", role.Type); err != nil {
		return diag.Errorf("unable to set resource organization role `type` read value: %v", err)
	}
	if err := d.Set("max_session_duration", role.MaxSessionDuration); err != nil {
		return diag.Errorf("unable to set resource organization role `max_session_duration` read value: %v", err)
	}
	if err := d.Set("policies", pol); err != nil {
		return diag.Errorf("unable to set resource organization role `policies` read value: %v", err)
	}
	if err := d.Set("date_created", role.DateCreated); err != nil {
		return diag.Errorf("unable to set resource organization role `date_created` read value: %v", err)
	}

	groups, _, _, err := client.Organization.ListRoleGroups(ctx, d.Id(), nil)
	if err != nil {
		return diag.Errorf("error getting role groups : %v", err)
	}

	var groupsList []string
	for i := range groups {
		groupsList = append(groupsList, groups[i].GroupID)
	}
	if err := d.Set("groups", groupsList); err != nil {
		return diag.Errorf("unable to set resource role `groups` read value: %v", err)
	}

	return nil
}

func resourceVultrOrganizationRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updating organization role (%s)", d.Id())

	req := &govultr.OrganizationRoleReq{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Type:               d.Get("type").(string),
		MaxSessionDuration: d.Get("max_session_duration").(int),
	}
	if _, _, err := client.Organization.UpdateRole(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating organization role %s : %v", d.Id(), err)
	}

	if d.HasChange("policies") {
		log.Printf("[INFO] Updating organization role policies")

		oldPolicies, newPolicies := d.GetChange("policies")
		oldPolicyList := oldPolicies.(*schema.Set).List()
		newPolicyList := newPolicies.(*schema.Set).List()

		var oldIDs, newIDs []string
		for i := range oldPolicyList {
			oldIDs = append(oldIDs, oldPolicyList[i].(string))
		}

		for i := range newPolicyList {
			newIDs = append(newIDs, newPolicyList[i].(string))
		}

		removeIDs := diffSlice(newIDs, oldIDs)
		addIDs := diffSlice(oldIDs, newIDs)

		if len(removeIDs) > 0 {
			for i := range removeIDs {
				if _, _, err := client.Organization.DetachRolePolicy(ctx, d.Id(), removeIDs[i]); err != nil {
					return diag.Errorf("error detaching organization role %s policy %s : %v", d.Id(), removeIDs[i], err)
				}
			}
		}

		if len(addIDs) > 0 {
			for i := range addIDs {
				addReq := &govultr.OrganizationRolePolicyReq{PolicyID: addIDs[i]}
				if _, _, err := client.Organization.AttachRolePolicy(ctx, d.Id(), addReq); err != nil {
					return diag.Errorf("error attaching organization role %s policy %s : %v", d.Id(), addIDs[i], err)
				}
			}
		}
	}

	if d.HasChange("groups") {
		log.Printf("[INFO] Updating role groups")

		oldGroups, newGroups := d.GetChange("groups")
		oldGroupsList := oldGroups.(*schema.Set).List()
		newGroupsList := newGroups.(*schema.Set).List()

		var oldIDs, newIDs []string
		for i := range oldGroupsList {
			oldIDs = append(oldIDs, oldGroupsList[i].(string))
		}

		for i := range newGroupsList {
			newIDs = append(newIDs, newGroupsList[i].(string))
		}

		removeIDs := diffSlice(newIDs, oldIDs)
		addIDs := diffSlice(oldIDs, newIDs)

		if len(removeIDs) > 0 {
			for i := range removeIDs {
				if err := client.Organization.DetachRoleGroup(ctx, d.Id(), removeIDs[i]); err != nil {
					return diag.Errorf("error detaching organization role %s group %s : %v", d.Id(), removeIDs[i], err)
				}
			}
		}

		if len(addIDs) > 0 {
			for i := range addIDs {
				if _, _, err := client.Organization.AttachRoleGroup(ctx, d.Id(), addIDs[i]); err != nil {
					return diag.Errorf("error attaching organization role %s group %s : %v", d.Id(), addIDs[i], err)
				}
			}
		}
	}

	return resourceVultrOrganizationRead(ctx, d, meta)
}

func resourceVultrOrganizationRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting organization role (%s)", d.Id())
	if err := client.Organization.DeleteRole(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting organization role %s : %v", d.Id(), err)
	}

	return nil
}
