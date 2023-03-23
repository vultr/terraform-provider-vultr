package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrStartupScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrStartupScriptCreate,
		ReadContext:   resourceVultrStartupScriptRead,
		UpdateContext: resourceVultrStartupScriptUpdate,
		DeleteContext: resourceVultrStartupScriptDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceVultrStartupScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	scriptReq := &govultr.StartupScriptReq{
		Name:   d.Get("name").(string),
		Script: d.Get("script").(string),
		Type:   d.Get("type").(string),
	}

	s,_, err := client.StartupScript.Create(ctx, scriptReq)
	if err != nil {
		return diag.Errorf("Error creating startup script: %v", err)
	}

	d.SetId(s.ID)
	log.Printf("[INFO] startup script ID: %s", d.Id())

	return resourceVultrStartupScriptRead(ctx, d, meta)
}

func resourceVultrStartupScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	script,_, err := client.StartupScript.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid startup script ID") {
			tflog.Warn(ctx, fmt.Sprintf("Removing startup script (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting startup script: %v", err)
	}

	if err := d.Set("name", script.Name); err != nil {
		return diag.Errorf("unable to set resource startup_script `name` read value: %v", err)
	}
	if err := d.Set("script", script.Script); err != nil {
		return diag.Errorf("unable to set resource startup_script `script` read value: %v", err)
	}
	if err := d.Set("type", script.Type); err != nil {
		return diag.Errorf("unable to set resource startup_script `type` read value: %v", err)
	}
	if err := d.Set("date_created", script.DateCreated); err != nil {
		return diag.Errorf("unable to set resource startup_script `date_created` read value: %v", err)
	}
	if err := d.Set("date_modified", script.DateModified); err != nil {
		return diag.Errorf("unable to set resource startup_script `date_modified` read value: %v", err)
	}

	return nil
}

func resourceVultrStartupScriptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	if d.HasChange("name") || d.HasChange("type") || d.HasChange("script") {
		scriptReq := &govultr.StartupScriptReq{
			Name:   d.Get("name").(string),
			Type:   d.Get("type").(string),
			Script: d.Get("script").(string),
		}

		log.Printf("[INFO] Updating startup script: %s", d.Id())
		if err := client.StartupScript.Update(ctx, d.Id(), scriptReq); err != nil {
			return diag.Errorf("Error updating startup script (%s): %v", d.Id(), err)
		}
	}

	return resourceVultrStartupScriptRead(ctx, d, meta)
}

func resourceVultrStartupScriptDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting startup script: %s", d.Id())
	if err := client.StartupScript.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying startup script (%s): %v", d.Id(), err)
	}

	return nil
}
