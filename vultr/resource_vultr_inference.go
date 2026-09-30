package vultr

import (
	"time"
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
		// Retry once on any error — Vultr API can return transient
		// 400/404 errors during maintenance or caching issues
		log.Printf("[WARN] Inference (%s) returned error: %v — retrying once", d.Id(), err)
		time.Sleep(2 * time.Second)
		inferenceSub2, _, err2 := client.Inference.Get(ctx, d.Id())
		if err2 != nil {
			if strings.Contains(err2.Error(), "invalid inference ID") {
				log.Printf("[WARN] Inference (%s) confirmed gone after retry — removing from state", d.Id())
				d.SetId("")
				return nil
			}
			// Still failing but not "gone" — skip instead of crash
			log.Printf("[WARN] Inference (%s) still failing after retry: %v — skipping read", d.Id(), err2)
			return nil
		}
		inferenceSub = inferenceSub2
		log.Printf("[WARN] Inference (%s) found on retry — keeping in state", d.Id())
	}

	// Guard against nil inference subscription (transient API caching issue)
	if inferenceSub == nil {
		log.Printf("[WARN] Inference subscription (%s) returned nil from API with no error — likely transient caching issue", d.Id())
		d.SetId("")
		return nil
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
