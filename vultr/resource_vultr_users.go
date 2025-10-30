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

func resourceVultrUsers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrUsersCreate,
		ReadContext:   resourceVultrUsersRead,
		UpdateContext: resourceVultrUsersUpdate,
		DeleteContext: resourceVultrUsersDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"api_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"acl": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrUsersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	test := d.Get("api_enabled").(bool)
	userReq := &govultr.UserReq{
		Email:      d.Get("email").(string),
		Name:       d.Get("name").(string),
		Password:   d.Get("password").(string),
		APIEnabled: &test,
	}

	acl, aclOK := d.GetOk("acl")
	a := acl.(*schema.Set).List()
	var aclMap []string
	if aclOK {
		for _, v := range a {
			aclMap = append(aclMap, v.(string))
		}

		userReq.ACL = aclMap
	}

	user, _, err := client.User.Create(context.Background(), userReq)
	if err != nil {
		return diag.Errorf("error creating user: %v", err)
	}

	d.SetId(user.ID)
	if err := d.Set("api_key", user.APIKey); err != nil {
		return diag.Errorf("unable to set resource user `api_key` create value: %v", err)
	}

	if groups, groupsOK := d.GetOk("groups"); groupsOK {
		groupList := groups.(*schema.Set).List()
		for i := range groupList {
			addReq := &govultr.OrganizationGroupMemberReq{UserID: d.Id()}
			if err := client.Organization.AddGroupMember(ctx, groupList[i].(string), addReq); err != nil {
				log.Printf("[ERROR] error adding user %s to organization group %s : %v", d.Id(), groupList[i], err)
			}
		}
	}

	return resourceVultrUsersRead(ctx, d, meta)
}

func resourceVultrUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	user, _, err := client.User.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid user") {
			tflog.Warn(ctx, fmt.Sprintf("Removing user (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting user: %v", err)
	}

	groups, _, _, err := client.Organization.ListUserGroups(ctx, d.Id(), nil)
	if err != nil {
		return diag.Errorf("error getting user groups : %v", err)
	}

	var groupList []string
	for i := range groups {
		groupList = append(groupList, groups[i].ID)
	}
	if err := d.Set("groups", groupList); err != nil {
		return diag.Errorf("unable to set resource user `groups` read value: %v", err)
	}

	if err := d.Set("name", user.Name); err != nil {
		return diag.Errorf("unable to set resource user `name` read value: %v", err)
	}
	if err := d.Set("email", user.Email); err != nil {
		return diag.Errorf("unable to set resource user `email` read value: %v", err)
	}
	if err := d.Set("api_enabled", user.APIEnabled); err != nil {
		return diag.Errorf("unable to set resource user `api_enabled` read value: %v", err)
	}
	if err := d.Set("acl", user.ACL); err != nil {
		return diag.Errorf("unable to set resource user `acl` read value: %v", err)
	}

	return nil
}

func resourceVultrUsersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	userReq := &govultr.UserReq{}

	if d.HasChange("email") {
		userReq.Email = d.Get("email").(string)
	}

	if d.HasChange("name") {
		userReq.Name = d.Get("name").(string)
	}

	if d.HasChange("password") {
		userReq.Password = d.Get("password").(string)
	}

	if d.HasChange("api_enabled") {
		api := d.Get("api_enabled").(bool)
		userReq.APIEnabled = &api
	}

	acl, aclOK := d.GetOk("acl")
	a := acl.(*schema.Set).List()
	var aclMap []string
	if aclOK {
		for _, v := range a {
			aclMap = append(aclMap, v.(string))
		}
		userReq.ACL = aclMap
	}

	err := client.User.Update(context.Background(), d.Id(), userReq)
	if err != nil {
		return diag.Errorf("Error updating user %s : %v", d.Id(), err)
	}

	if d.HasChange("groups") {
		log.Printf("[INFO] Updating user groups")

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
				if err := client.Organization.RemoveGroupMember(ctx, removeIDs[i], d.Id()); err != nil {
					return diag.Errorf("error removing user %s from organization group %s : %v", d.Id(), removeIDs[i], err)
				}
			}
		}

		if len(addIDs) > 0 {
			for i := range addIDs {
				addReq := &govultr.OrganizationGroupMemberReq{UserID: d.Id()}
				if err := client.Organization.AddGroupMember(ctx, addIDs[i], addReq); err != nil {
					return diag.Errorf("error adding user %s to organization group %s : %v", d.Id(), addIDs[i], err)
				}
			}
		}
	}

	return resourceVultrUsersRead(ctx, d, meta)
}

func resourceVultrUsersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting User %s", d.Id())

	err := client.User.Delete(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error deleting user %s : %v", d.Id(), err)
	}
	return nil
}
