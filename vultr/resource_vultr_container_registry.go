package vultr

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrContainerRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrContainerRegistryCreate,
		ReadContext:   resourceVultrContainerRegistryRead,
		UpdateContext: resourceVultrContainerRegistryUpdate,
		DeleteContext: resourceVultrContainerRegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"root_user": {
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceVultrContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	crReq := &govultr.ContainerRegistryReq{
		Name:   d.Get("name").(string),
		Region: d.Get("region").(string),
		Public: d.Get("public").(bool),
		Plan:   d.Get("plan").(string),
	}

	cr, _, err := client.ContainerRegistry.Create(ctx, crReq)
	if err != nil {
		return diag.Errorf("error creating container registry: %v", err)
	}

	d.SetId(cr.ID)
	log.Printf("[INFO] Container Registry ID: %s", d.Id())

	return resourceVultrContainerRegistryRead(ctx, d, meta)
}

func resourceVultrContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	cr, _, err := client.ContainerRegistry.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid container registry ID") {
			log.Printf("[WARN] Vultr container registry (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting container registry: %v", err)
	}

	if err := d.Set("name", cr.Name); err != nil {
		return diag.Errorf("unable to set resource container registry `name` read value: %v", err)
	}
	if err := d.Set("urn", cr.URN); err != nil {
		return diag.Errorf("unable to set resource container registry `urn` read value: %v", err)
	}

	if err := d.Set("storage", cr.Storage); err != nil {
		return diag.Errorf("unable to set resource container registry `storage` read value: %v", err)
	}
	if err := d.Set("root_user", cr.RootUser); err != nil {
		return diag.Errorf("unable to set resource container registry `root_user` read value: %v", err)
	}
	if err := d.Set("public", cr.Public); err != nil {
		return diag.Errorf("unable to set resource container `public` read value: %v", err)
	}

	return nil
}

func resourceVultrContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vcr := &govultr.ContainerRegistryUpdateReq{}

	if d.HasChange("plan") {
		vcr.Plan = govultr.StringToStringPtr(d.Get("plan").(string))
	}

	if d.HasChange("public") {
		vcr.Public = govultr.BoolToBoolPtr(d.Get("public").(bool))
	}

	log.Printf("[INFO] Updating container registry: %s", d.Id())

	cr, _, err := client.ContainerRegistry.Update(ctx, d.Id(), vcr)
	if err != nil {
		return diag.Errorf("error updating container registry: %v", err)
	}

	d.SetId(cr.ID)
	log.Printf("[INFO] Container Registry ID: %s", d.Id())

	return resourceVultrContainerRegistryRead(ctx, d, meta)
}

func resourceVultrContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting container: %s", d.Id())

	err := client.ContainerRegistry.Delete(ctx, d.Id())

	if err != nil {
		return diag.Errorf("error destroying Container Registry (%s): %v", d.Id(), err)
	}

	return nil
}
