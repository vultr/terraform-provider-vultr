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

func resourceVultrOrganizationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOrganizationPolicyCreate,
		ReadContext:   resourceVultrOrganizationPolicyRead,
		UpdateContext: resourceVultrOrganizationPolicyUpdate,
		DeleteContext: resourceVultrOrganizationPolicyDelete,
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
			"document": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"statement": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"effect": {
										Type:     schema.TypeString,
										Required: true,
									},
									"actions": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"resources": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Default:  nil,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Default:  nil,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_system_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrOrganizationPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	documentObj := d.Get("document").(*schema.Set).List()
	documentVal := documentObj[0].(map[string]interface{})
	statements := documentVal["statement"].([]interface{})

	var statementReq []govultr.OrganizationPolicyStatement
	for i := range statements {
		statementObj := statements[i].(map[string]interface{})

		actionObj := statementObj["actions"].([]interface{})
		var actionList []string
		for n := range actionObj {
			actionList = append(actionList, actionObj[n].(string))
		}

		resourceObj := statementObj["resources"].([]interface{})
		var resourceList []string
		for n := range resourceObj {
			resourceList = append(resourceList, resourceObj[n].(string))
		}

		statementReq = append(statementReq, govultr.OrganizationPolicyStatement{
			Effect:   statementObj["effect"].(string),
			Action:   actionList,
			Resource: resourceList,
		})
	}

	policyReq := &govultr.OrganizationPolicyReq{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PolicyDocument: govultr.OrganizationPolicyDocument{
			Version:   documentVal["version"].(string),
			Statement: statementReq,
		},
	}

	log.Print("[INFO] Creating organization policy")

	policy, _, err := client.Organization.CreatePolicy(ctx, policyReq)
	if err != nil {
		return diag.Errorf("error while creating organization policy : %s", err)
	}

	if groups, groupsOK := d.GetOk("groups"); groupsOK {
		groupList := groups.(*schema.Set).List()
		for i := range groupList {
			if err := client.Organization.AttachPolicyGroup(ctx, policy.ID, groupList[i].(string)); err != nil {
				return diag.Errorf("error attaching organization policy %s to group %s : %v", policy.ID, groupList[i], err)
			}
		}
	}

	if users, usersOK := d.GetOk("users"); usersOK {
		userList := users.(*schema.Set).List()
		for i := range userList {
			if err := client.Organization.AttachPolicyUser(ctx, policy.ID, userList[i].(string)); err != nil {
				return diag.Errorf("error attaching organization policy %s to user %s : %v", policy.ID, userList[i], err)
			}
		}
	}

	d.SetId(policy.ID)

	return resourceVultrOrganizationPolicyRead(ctx, d, meta)
}

func resourceVultrOrganizationPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	policy, _, err := client.Organization.GetPolicy(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Policy not found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing organization policy (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting organization policy : %v", err)
	}

	var statementFlat []map[string]interface{}
	for i := range policy.Document.Statement {
		statementFlat = append(statementFlat, map[string]interface{}{
			"effect":    policy.Document.Statement[i].Effect,
			"actions":   policy.Document.Statement[i].Action,
			"resources": policy.Document.Statement[i].Resource,
		})
	}

	policyDocument := []map[string]interface{}{
		{
			"version":   policy.Document.Version,
			"statement": statementFlat,
		},
	}

	users, _, _, err := client.Organization.ListPolicyUsers(ctx, d.Id(), nil)
	if err != nil {
		return diag.Errorf("error getting organization policy users : %v", err)
	}

	var userIDs []string
	if len(users) != 0 {
		for i := range users {
			userIDs = append(userIDs, users[i].ID)
		}
	}

	groups, _, _, err := client.Organization.ListPolicyGroups(ctx, d.Id(), nil)
	if err != nil {
		return diag.Errorf("error getting organization policy groups : %v", err)
	}

	var groupIDs []string
	if len(groups) != 0 {
		for i := range groups {
			groupIDs = append(groupIDs, groups[i].ID)
		}
	}

	if err := d.Set("name", policy.Name); err != nil {
		return diag.Errorf("unable to set resource organization policy `name` read value: %v", err)
	}
	if err := d.Set("description", policy.Description); err != nil {
		return diag.Errorf("unable to set resource organization policy `description` read value: %v", err)
	}
	if err := d.Set("document", policyDocument); err != nil {
		return diag.Errorf("unable to set resource organization policy `document` read value: %v", err)
	}
	if err := d.Set("is_system_policy", policy.SystemPolicy); err != nil {
		return diag.Errorf("unable to set resource organization policy `is_system_policy` read value: %v", err)
	}
	if err := d.Set("users", userIDs); err != nil {
		return diag.Errorf("unable to set resource organization policy `users` read value: %v", err)
	}
	if err := d.Set("groups", groupIDs); err != nil {
		return diag.Errorf("unable to set resource organization policy `groups` read value: %v", err)
	}
	if err := d.Set("date_created", policy.DateCreated); err != nil {
		return diag.Errorf("unable to set resource organization policy `date_created` read value: %v", err)
	}

	return nil
}

func resourceVultrOrganizationPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Updating organization policy (%s)", d.Id())

	documentObj := d.Get("document").(*schema.Set).List()
	documentVal := documentObj[0].(map[string]interface{})
	statements := documentVal["statement"].([]interface{})

	var statementReq []govultr.OrganizationPolicyStatement
	for i := range statements {
		statementObj := statements[i].(map[string]interface{})

		actionObj := statementObj["actions"].([]interface{})
		var actionList []string
		for n := range actionObj {
			actionList = append(actionList, actionObj[n].(string))
		}

		resourceObj := statementObj["resources"].([]interface{})
		var resourceList []string
		for n := range resourceObj {
			resourceList = append(resourceList, resourceObj[n].(string))
		}

		statementReq = append(statementReq, govultr.OrganizationPolicyStatement{
			Effect:   statementObj["effect"].(string),
			Action:   actionList,
			Resource: resourceList,
		})
	}

	policyReq := &govultr.OrganizationPolicyReq{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PolicyDocument: govultr.OrganizationPolicyDocument{
			Version:   documentVal["version"].(string),
			Statement: statementReq,
		},
	}

	if _, _, err := client.Organization.UpdatePolicy(ctx, d.Id(), policyReq); err != nil {
		return diag.Errorf("error while updating organization policy : %s", err)
	}

	if d.HasChange("groups") {
		log.Printf("[INFO] Updating organization policy groups")

		oldGroups, newGroups := d.GetChange("groups")
		oldGroupList := oldGroups.(*schema.Set).List()
		newGroupList := newGroups.(*schema.Set).List()

		var oldIDs, newIDs []string
		for i := range oldGroupList {
			oldIDs = append(oldIDs, oldGroupList[i].(string))
		}

		for i := range newGroupList {
			newIDs = append(newIDs, newGroupList[i].(string))
		}

		removeIDs := diffSlice(newIDs, oldIDs)
		addIDs := diffSlice(oldIDs, newIDs)

		if len(removeIDs) > 0 {
			for i := range removeIDs {
				if err := client.Organization.DetachPolicyGroup(ctx, d.Id(), removeIDs[i]); err != nil {
					return diag.Errorf("error detaching organization group %s from policy %s : %v", d.Id(), removeIDs[i], err)
				}
			}
		}

		if len(addIDs) > 0 {
			for i := range addIDs {
				if err := client.Organization.AttachPolicyGroup(ctx, d.Id(), addIDs[i]); err != nil {
					return diag.Errorf("error attaching organization group %s to policy %s : %v", d.Id(), addIDs[i], err)
				}
			}
		}
	}

	if d.HasChange("users") {
		log.Printf("[INFO] Updating organization policy users")

		oldUsers, newUsers := d.GetChange("users")
		oldUserList := oldUsers.(*schema.Set).List()
		newUserList := newUsers.(*schema.Set).List()

		var oldIDs, newIDs []string
		for i := range oldUserList {
			oldIDs = append(oldIDs, oldUserList[i].(string))
		}

		for i := range newUserList {
			newIDs = append(newIDs, newUserList[i].(string))
		}

		removeIDs := diffSlice(newIDs, oldIDs)
		addIDs := diffSlice(oldIDs, newIDs)

		if len(removeIDs) > 0 {
			for i := range removeIDs {
				if err := client.Organization.DetachPolicyUser(ctx, d.Id(), removeIDs[i]); err != nil {
					return diag.Errorf("error detaching organization user %s from policy %s : %v", d.Id(), removeIDs[i], err)
				}
			}
		}

		if len(addIDs) > 0 {
			for i := range addIDs {
				if err := client.Organization.AttachPolicyUser(ctx, d.Id(), addIDs[i]); err != nil {
					return diag.Errorf("error attaching organization user %s to policy %s : %v", d.Id(), addIDs[i], err)
				}
			}
		}
	}

	return resourceVultrOrganizationPolicyRead(ctx, d, meta)
}

func resourceVultrOrganizationPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting organization policy (%s)", d.Id())
	if err := client.Organization.DeletePolicy(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting organization policy %s : %v", d.Id(), err)
	}

	return nil
}
