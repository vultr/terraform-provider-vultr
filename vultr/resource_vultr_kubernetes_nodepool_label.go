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

func resourceVultrKubernetesNodePoolLabel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrKubernetesNodePoolLabelCreate,
		ReadContext:   resourceVultrKubernetesNodePoolLabelRead,
		UpdateContext: resourceVultrKubernetesNodePoolLabelUpdate,
		DeleteContext: resourceVultrKubernetesNodePoolLabelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// Expected format: "clusterID nodePoolID labelID"
				ids := strings.SplitN(d.Id(), " ", 3)
				if len(ids) != 3 || ids[0] == "" || ids[1] == "" || ids[2] == "" {
					return nil, fmt.Errorf("unexpected format of label import IDs (%s): expected 'clusterID nodePoolID labelID'", d.Id())
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
		},
	}
}

func resourceVultrKubernetesNodePoolLabelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)
	nodepoolID := d.Get("nodepool_id").(string)

	req := &govultr.NodePoolLabelReq{
		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
	}

	label, _, err := client.Kubernetes.CreateNodePoolLabel(ctx, clusterID, nodepoolID, req)
	if err != nil {
		return diag.Errorf("error creating nodepool label: %v", err)
	}

	d.SetId(label.ID)

	return resourceVultrKubernetesNodePoolLabelRead(ctx, d, meta)
}

func resourceVultrKubernetesNodePoolLabelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)
	nodepoolID := d.Get("nodepool_id").(string)

	label, _, err := client.Kubernetes.ReadNodePoolLabel(ctx, clusterID, nodepoolID, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid") || strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] Kubernetes NodePool Label (%v) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting nodepool label: %v", err)
	}

	if err := d.Set("key", label.Key); err != nil {
		return diag.Errorf("unable to set `key`: %v", err)
	}
	if err := d.Set("value", label.Value); err != nil {
		return diag.Errorf("unable to set `value`: %v", err)
	}

	return nil
}

func resourceVultrKubernetesNodePoolLabelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Labels are immutable except for value - if key changes, ForceNew will handle it
	// For value changes, we need to delete and recreate
	if d.HasChange("value") {
		client := meta.(*Client).govultrClient()
		clusterID := d.Get("cluster_id").(string)
		nodepoolID := d.Get("nodepool_id").(string)

		// Delete the old label
		if err := client.Kubernetes.DeleteNodePoolLabel(ctx, clusterID, nodepoolID, d.Id()); err != nil {
			return diag.Errorf("error deleting old label during update: %v", err)
		}

		// Create new label
		req := &govultr.NodePoolLabelReq{
			Key:   d.Get("key").(string),
			Value: d.Get("value").(string),
		}

		label, _, err := client.Kubernetes.CreateNodePoolLabel(ctx, clusterID, nodepoolID, req)
		if err != nil {
			return diag.Errorf("error creating new label during update: %v", err)
		}

		d.SetId(label.ID)
	}

	return resourceVultrKubernetesNodePoolLabelRead(ctx, d, meta)
}

func resourceVultrKubernetesNodePoolLabelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)
	nodepoolID := d.Get("nodepool_id").(string)

	if err := client.Kubernetes.DeleteNodePoolLabel(ctx, clusterID, nodepoolID, d.Id()); err != nil {
		return diag.Errorf("error deleting nodepool label %v: %v", d.Id(), err)
	}

	return nil
}
