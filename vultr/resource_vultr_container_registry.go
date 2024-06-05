package vultr

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-z0-9]+$`),
					"must be lowercase and contain only alpha-numeric characters",
				),
			},
			"plan": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"start_up",
						"business",
						"premium",
						"enterprise",
					},
					false,
				),
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"root_user": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	crReq := &govultr.ContainerRegistryReq{
		Name:   d.Get("name").(string),
		Region: d.Get("region").(string),
		Public: d.Get("public").(bool),
		Plan:   d.Get("plan").(string),
	}

	log.Print("[INFO] Creating container registry")
	cr, _, err := client.ContainerRegistry.Create(ctx, crReq)
	if err != nil {
		return diag.Errorf("error creating container registry: %v", err)
	}

	d.SetId(cr.ID)
	log.Printf("[INFO] Created container registry with ID: %s", d.Id())

	return resourceVultrContainerRegistryRead(ctx, d, meta)
}

func resourceVultrContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	cr, _, err := client.ContainerRegistry.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid container registry ID") {
			log.Printf("[WARN] Container registry (%s) not found and will be removed", d.Id())
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
	if err := d.Set("storage", flattenCRStorage(cr)); err != nil {
		return diag.Errorf("unable to set resource container registry `storage` read value: %v", err)
	}
	if err := d.Set("root_user", flattenCRRootUser(cr)); err != nil {
		return diag.Errorf("unable to set resource container registry `root_user` read value: %v", err)
	}
	if err := d.Set("public", cr.Public); err != nil {
		return diag.Errorf("unable to set resource container `public` read value: %v", err)
	}
	if err := d.Set("date_created", cr.DateCreated); err != nil {
		return diag.Errorf("unable to set resource container `date_created` read value: %v", err)
	}

	return nil
}

func resourceVultrContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	vcr := &govultr.ContainerRegistryUpdateReq{}
	log.Printf("[INFO] Updating container registry: %s", d.Id())

	if d.HasChange("plan") {
		log.Print("[INFO] Updating `plan`")
		vcr.Plan = govultr.StringToStringPtr(d.Get("plan").(string))
	}

	if d.HasChange("public") {
		log.Print("[INFO] Updating `public`")
		vcr.Public = govultr.BoolToBoolPtr(d.Get("public").(bool))
	}

	cr, _, err := client.ContainerRegistry.Update(ctx, d.Id(), vcr)
	if err != nil {
		return diag.Errorf("error updating container registry: %v", err)
	}

	d.SetId(cr.ID)

	return resourceVultrContainerRegistryRead(ctx, d, meta)
}

func resourceVultrContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting container registry: %s", d.Id())

	err := client.ContainerRegistry.Delete(ctx, d.Id())

	if err != nil {
		return diag.Errorf("error deleting container registry (%s): %v", d.Id(), err)
	}

	return nil
}

func flattenCRStorage(cr *govultr.ContainerRegistry) map[string]interface{} {
	return map[string]interface{}{
		"allowed": fmt.Sprintf("%.2f GB", cr.Storage.Allowed.GigaBytes),
		"used":    fmt.Sprintf("%.2f GB", cr.Storage.Used.GigaBytes),
	}
}

func flattenCRRootUser(cr *govultr.ContainerRegistry) map[string]interface{} {
	return map[string]interface{}{
		"id":            strconv.Itoa(cr.RootUser.ID),
		"username":      cr.RootUser.UserName,
		"password":      cr.RootUser.Password,
		"root":          strconv.FormatBool(cr.RootUser.Root),
		"date_created":  cr.RootUser.DateCreated,
		"date_modified": cr.RootUser.DateModified,
	}
}
