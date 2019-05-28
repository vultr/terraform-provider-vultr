package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
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
				Type:     schema.TypeString,
				Sensitive: true,
				Required: true,
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

	name := d.Get("name").(string)
	email := d.Get("email").(string)
	password := d.Get("password").(string)

	// optional param
	apiEnabled := d.Get("api_enabled")

	api := ""
	if apiEnabled.(bool) == true {
		api = "yes"
	} else {
		api = "no"
	}

	acl, aclOK := d.GetOk("acl")
	a := acl.([]interface{})
	aclMap := []string{}
	if aclOK {
		for _, v := range a {
			aclMap = append(aclMap, v.(string))
		}
	}

	user, err := client.User.Create(context.Background(), email, name, password, api, aclMap)

	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	d.SetId(user.UserID)
	d.Set("api_key", user.APIKey)

	return resourceVultrUsersRead(d, meta)
}

func resourceVultrUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	users, err := client.User.GetList(context.Background())

	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	counter := 0
	flag := false
	for _, v := range users {
		if d.Id() == v.UserID {
			flag = true
			break
		}
		counter++
	}

	if !flag {
		log.Printf("[WARN] Removing user (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	enabled := false
	if users[counter].APIEnabled == "yes" {
		enabled = true
	}

	d.Set("name", users[counter].Name)
	d.Set("email", users[counter].Email)
	d.Set("api_enabled", enabled)
	d.Set("acl", users[counter].ACL)
	return nil
}

func resourceVultrUsersUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	user := &govultr.User{
		UserID:   d.Id(),
		Name:     d.Get("name").(string),
		Email:    d.Get("email").(string),
		Password: d.Get("password").(string),
	}

	if d.Get("api_enabled").(bool) == true {
		user.APIEnabled = "yes"
	} else {
		user.APIEnabled = "no"
	}

	acl, aclOK := d.GetOk("acl")
	a := acl.([]interface{})
	aclMap := []string{}
	if aclOK {
		for _, v := range a {
			aclMap = append(aclMap, v.(string))
		}
	}

	user.ACL = aclMap

	err := client.User.Update(context.Background(), user)
	if err != nil {
		return fmt.Errorf("Error updating user %s : %v", d.Id(), err)
	}

	return resourceVultrUsersRead(d, meta)
}

func resourceVultrUsersDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Destroying User %s", d.Id())

	err := client.User.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting user %s : %v", d.Id(), err)
	}
	return nil
}
