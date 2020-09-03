package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrUsers() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrUsersCreate,
		Read:   resourceVultrUsersRead,
		Update: resourceVultrUsersUpdate,
		Delete: resourceVultrUsersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceVultrUsersCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	userReq := &govultr.UserReq{
		Email:    d.Get("email").(string),
		Name:     d.Get("name").(string),
		Password: d.Get("password").(string),
		//todo fix govultr
		//APIEnabled:  d.Get("apiEnabled").(string),
	}

	//todo fix govultr
	//acl, aclOK := d.GetOk("acl")
	//a := acl.([]interface{})
	//aclMap := []string{}
	//if aclOK {
	//	for _, v := range a {
	//		aclMap = append(aclMap, v.(string))
	//	}
	//
	//	userReq.ACL = aclMap
	//}

	user, err := client.User.Create(context.Background(), userReq)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	d.SetId(user.ID)
	d.Set("api_key", user.APIKey)

	return resourceVultrUsersRead(d, meta)
}

func resourceVultrUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	user, err := client.User.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	d.Set("name", user.Name)
	d.Set("email", user.Email)
	d.Set("api_enabled", user.APIEnabled)
	if err := d.Set("acl", user.ACL); err != nil {
		return fmt.Errorf("error setting `acl`: %#v", err)
	}

	return nil
}

func resourceVultrUsersUpdate(d *schema.ResourceData, meta interface{}) error {
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
		userReq.APIEnabled = d.Get("api_enabled").(string)
	}

	acl, aclOK := d.GetOk("acl")
	a := acl.([]interface{})
	aclMap := []string{}
	if aclOK {
		for _, v := range a {
			aclMap = append(aclMap, v.(string))
		}
		userReq.ACL = aclMap
	}

	err := client.User.Update(context.Background(), d.Id(), userReq)
	if err != nil {
		return fmt.Errorf("Error updating user %s : %v", d.Id(), err)
	}

	return resourceVultrUsersRead(d, meta)
}

func resourceVultrUsersDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting User %s", d.Id())

	err := client.User.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error deleting user %s : %v", d.Id(), err)
	}
	return nil
}
