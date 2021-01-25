package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
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
				Type:     schema.TypeList,
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
	a := acl.([]interface{})
	var aclMap []string
	if aclOK {
		for _, v := range a {
			aclMap = append(aclMap, v.(string))
		}

		userReq.ACL = aclMap
	}

	user, err := client.User.Create(context.Background(), userReq)
	if err != nil {
		return diag.Errorf("error creating user: %v", err)
	}

	d.SetId(user.ID)
	d.Set("api_key", user.APIKey)

	return resourceVultrUsersRead(ctx, d, meta)
}

func resourceVultrUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	user, err := client.User.Get(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error getting user: %v", err)
	}

	d.Set("name", user.Name)
	d.Set("email", user.Email)
	d.Set("api_enabled", user.APIEnabled)
	if err := d.Set("acl", user.ACL); err != nil {
		return diag.Errorf("error setting `acl`: %#v", err)
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
	a := acl.([]interface{})
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
