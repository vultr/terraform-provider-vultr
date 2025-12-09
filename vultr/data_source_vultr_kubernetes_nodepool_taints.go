package vultr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrKubernetesNodePoolTaints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrKubernetesNodePoolTaintsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodepool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": dataSourceFiltersSchema(),
			"taints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVultrKubernetesNodePoolTaintsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)
	nodepoolID := d.Get("nodepool_id").(string)

	taints, _, err := client.Kubernetes.ListNodePoolTaints(ctx, clusterID, nodepoolID)
	if err != nil {
		return diag.Errorf("error getting nodepool taints: %v", err)
	}

	var filteredTaints []govultr.NodePoolTaint
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		f := buildVultrDataSourceFilter(filters.(*schema.Set))
		for _, taint := range taints {
			sm, err := structToMap(taint)
			if err != nil {
				return diag.FromErr(err)
			}
			if filterLoop(f, sm) {
				filteredTaints = append(filteredTaints, taint)
			}
		}
	} else {
		filteredTaints = taints
	}

	if err := d.Set("taints", flattenNodePoolTaints(filteredTaints)); err != nil {
		return diag.Errorf("error setting taints: %v", err)
	}

	d.SetId(fmt.Sprintf("%s-%s", clusterID, nodepoolID))

	return nil
}

func flattenNodePoolTaints(taints []govultr.NodePoolTaint) []map[string]interface{} {
	var result []map[string]interface{}

	for _, taint := range taints {
		t := map[string]interface{}{
			"id":     taint.ID,
			"key":    taint.Key,
			"value":  taint.Value,
			"effect": taint.Effect,
		}
		result = append(result, t)
	}

	return result
}
