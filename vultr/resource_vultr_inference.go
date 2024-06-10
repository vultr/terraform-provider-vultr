package vultr

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrInference() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrInferenceCreate,
		ReadContext:   resourceVultrInferenceRead,
		UpdateContext: resourceVultrInferenceUpdate,
		DeleteContext: resourceVultrInferenceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"usage": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     schema.TypeInt,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
		},
	}
}

func resourceVultrInferenceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.InferenceCreateUpdateReq{
		Label: d.Get("label").(string),
	}

	log.Printf("[INFO] Creating inference subscription")
	inferenceSub, _, err := client.Inference.Create(ctx, req)
	if err != nil {
		return diag.Errorf("error creating inference subscription: %v", err)
	}

	d.SetId(inferenceSub.ID)

	return resourceVultrInferenceRead(ctx, d, meta)
}

func resourceVultrInferenceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	inferenceSub, _, err := client.Inference.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "invalid inference ID") {
			log.Printf("[WARN] Removing inference subscription (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting inference subscription (%s): %v", d.Id(), err)
	}

	if err := d.Set("date_created", inferenceSub.DateCreated); err != nil {
		return diag.Errorf("unable to set resource inference `date_created` read value: %v", err)
	}

	if err := d.Set("label", inferenceSub.Label); err != nil {
		return diag.Errorf("unable to set resource inference `label` read value: %v", err)
	}

	if err := d.Set("api_key", inferenceSub.APIKey); err != nil {
		return diag.Errorf("unable to set resource inference `api_key` read value: %v", err)
	}

	// Grab usage
	usage, _, err := client.Inference.GetUsage(ctx, d.Id())
	if err == nil {
		if err := d.Set("usage", flattenInferenceUsage(usage)); err != nil {
			return diag.Errorf("unable to set resource inference `usage` read value: %v", err)
		}
	}

	return nil
}
func resourceVultrInferenceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.InferenceCreateUpdateReq{
		Label: d.Get("label").(string),
	}

	if _, _, err := client.Inference.Update(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating inference subscription %s : %s", d.Id(), err.Error())
	}

	return resourceVultrInferenceRead(ctx, d, meta)
}

func resourceVultrInferenceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting inference subscription (%s)", d.Id())

	if err := client.Inference.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying inference subscription %s : %v", d.Id(), err)
	}

	return nil
}
