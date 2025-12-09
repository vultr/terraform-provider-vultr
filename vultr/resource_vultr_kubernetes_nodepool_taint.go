package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrKubernetesNodePoolTaint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrKubernetesNodePoolTaintCreate,
		ReadContext:   resourceVultrKubernetesNodePoolTaintRead,
		UpdateContext: resourceVultrKubernetesNodePoolTaintUpdate,
		DeleteContext: resourceVultrKubernetesNodePoolTaintDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// Expected format: "clusterID nodePoolID taintID"
				ids := strings.SplitN(d.Id(), " ", 3)
				if len(ids) != 3 || ids[0] == "" || ids[1] == "" || ids[2] == "" {
					return nil, fmt.Errorf("unexpected format of taint import IDs (%s): expected 'clusterID nodePoolID taintID'", d.Id())
				}

				d.SetId(ids[2])
				if err := d.Set("cluster_id", ids[0]); err != nil {
					return nil, fmt.Errorf("unable to set cluster_id: %v", err)
				}
				if err := d.Set("nodepool_id", ids[1]); err != nil {
					return nil, fmt.Errorf("unable to set nodepool_id: %v", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nodepool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"effect": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validEffects := []string{"NoSchedule", "PreferNoSchedule", "NoExecute"}
					for _, effect := range validEffects {
						if v == effect {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validEffects, v))
					return
				},
			},
		},
	}
}

func resourceVultrKubernetesNodePoolTaintCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)
	nodepoolID := d.Get("nodepool_id").(string)

	req := &govultr.NodePoolTaintReq{
		Key:    d.Get("key").(string),
		Value:  d.Get("value").(string),
		Effect: d.Get("effect").(string),
	}

	taint, _, err := client.Kubernetes.CreateNodePoolTaint(ctx, clusterID, nodepoolID, req)
	if err != nil {
		return diag.Errorf("error creating nodepool taint: %v", err)
	}

	d.SetId(taint.ID)

	return resourceVultrKubernetesNodePoolTaintRead(ctx, d, meta)
}

func resourceVultrKubernetesNodePoolTaintRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)
	nodepoolID := d.Get("nodepool_id").(string)

	taint, _, err := client.Kubernetes.ReadNodePoolTaint(ctx, clusterID, nodepoolID, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid") || strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] Kubernetes NodePool Taint (%v) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting nodepool taint: %v", err)
	}

	if err := d.Set("key", taint.Key); err != nil {
		return diag.Errorf("unable to set `key`: %v", err)
	}
	if err := d.Set("value", taint.Value); err != nil {
		return diag.Errorf("unable to set `value`: %v", err)
	}
	if err := d.Set("effect", taint.Effect); err != nil {
		return diag.Errorf("unable to set `effect`: %v", err)
	}

	return nil
}

func resourceVultrKubernetesNodePoolTaintUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Taints are immutable except for value and effect - if key changes, ForceNew will handle it
	// For value/effect changes, we need to delete and recreate
	if d.HasChange("value") || d.HasChange("effect") {
		client := meta.(*Client).govultrClient()
		clusterID := d.Get("cluster_id").(string)
		nodepoolID := d.Get("nodepool_id").(string)

		// Delete the old taint
		if err := client.Kubernetes.DeleteNodePoolTaint(ctx, clusterID, nodepoolID, d.Id()); err != nil {
			return diag.Errorf("error deleting old taint during update: %v", err)
		}

		// Create new taint
		req := &govultr.NodePoolTaintReq{
			Key:    d.Get("key").(string),
			Value:  d.Get("value").(string),
			Effect: d.Get("effect").(string),
		}

		taint, _, err := client.Kubernetes.CreateNodePoolTaint(ctx, clusterID, nodepoolID, req)
		if err != nil {
			return diag.Errorf("error creating new taint during update: %v", err)
		}

		d.SetId(taint.ID)
	}

	return resourceVultrKubernetesNodePoolTaintRead(ctx, d, meta)
}

func resourceVultrKubernetesNodePoolTaintDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)
	nodepoolID := d.Get("nodepool_id").(string)

	if err := client.Kubernetes.DeleteNodePoolTaint(ctx, clusterID, nodepoolID, d.Id()); err != nil {
		return diag.Errorf("error deleting nodepool taint %v: %v", d.Id(), err)
	}

	return nil
}
