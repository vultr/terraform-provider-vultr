package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
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
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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

	name := d.Get("name").(string)
	script := d.Get("script").(string)

	var scriptType string
	t, ok := d.GetOk("type")
	if ok {
		scriptType = t.(string)
	}

	s, err := client.StartupScript.Create(context.Background(), name, script, scriptType)
	if err != nil {
		return fmt.Errorf("Error creating startup script: %v", err)
	}

	d.SetId(s.ScriptID)
	log.Printf("[INFO] startup script ID: %s", d.Id())

	return resourceVultrStartupScriptRead(d, meta)
}

func resourceVultrStartupScriptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	scripts, err := client.StartupScript.List(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting startup scripts: %v", err)
	}

	var script *govultr.StartupScript
	for i := range scripts {
		if scripts[i].ScriptID == d.Id() {
			script = &scripts[i]
			break
		}
	}

	if script == nil {
		log.Printf("[WARN] Vultr startup script (%s) not found", d.Id())
		d.SetId("")
		return nil
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

	script := &govultr.StartupScript{
		ScriptID: d.Id(),
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Script:   d.Get("script").(string),
	}

	log.Printf("[INFO] Updating startup script: %s", d.Id())
	if err := client.StartupScript.Update(context.Background(), script); err != nil {
		return fmt.Errorf("Error updating startup script (%s): %v", d.Id(), err)
	}

	return resourceVultrStartupScriptRead(d, meta)
}

func resourceVultrStartupScriptDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting startup script: %s", d.Id())
	if err := client.StartupScript.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("Error destroying startup script (%s): %v", d.Id(), err)
	}

	return nil
}
