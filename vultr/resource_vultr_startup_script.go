package vultr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrStartupScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrStartupScriptCreate,
		Read:   resourceVultrStartupScriptRead,
		Update: resourceVultrStartupScriptUpdate,
		Delete: resourceVultrStartupScriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"script": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsBase64,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "boot",
				ValidateFunc: validation.StringInSlice([]string{"boot", "pxe"}, false),
				ForceNew:     true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrStartupScriptCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	scriptReq := &govultr.StartupScriptReq{
		Name:   d.Get("name").(string),
		Script: d.Get("script").(string),
		Type:   d.Get("type").(string),
	}

	s, err := client.StartupScript.Create(context.Background(), scriptReq)
	if err != nil {
		return fmt.Errorf("Error creating startup script: %v", err)
	}

	d.SetId(s.ID)
	log.Printf("[INFO] startup script ID: %s", d.Id())

	return resourceVultrStartupScriptRead(d, meta)
}

func resourceVultrStartupScriptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	script, err := client.StartupScript.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting startup script: %v", err)
	}

	d.Set("name", script.Name)
	d.Set("script", script.Script)
	d.Set("type", script.Type)
	d.Set("date_created", script.DateCreated)
	d.Set("date_modified", script.DateModified)

	return nil
}

func resourceVultrStartupScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	if d.HasChange("name") || d.HasChange("type") || d.HasChange("script") {
		scriptReq := &govultr.StartupScriptReq{
			Name:   d.Get("name").(string),
			Type:   d.Get("type").(string),
			Script: d.Get("script").(string),
		}

		log.Printf("[INFO] Updating startup script: %s", d.Id())
		if err := client.StartupScript.Update(context.Background(), d.Id(), scriptReq); err != nil {
			return fmt.Errorf("Error updating startup script (%s): %v", d.Id(), err)
		}
	}

	return resourceVultrStartupScriptRead(d, meta)
}

func resourceVultrStartupScriptDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting startup script: %s", d.Id())
	if err := client.StartupScript.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error destroying startup script (%s): %v", d.Id(), err)
	}

	return nil
}
