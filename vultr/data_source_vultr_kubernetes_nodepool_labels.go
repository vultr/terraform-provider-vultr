package vultr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrKubernetesNodePoolLabels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrKubernetesNodePoolLabelsRead,
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
			"labels": {
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
					},
				},
			},
		},
	}
}

func dataSourceVultrKubernetesNodePoolLabelsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)
	nodepoolID := d.Get("nodepool_id").(string)

	labels, _, err := client.Kubernetes.ListNodePoolLabels(ctx, clusterID, nodepoolID)
	if err != nil {
		return diag.Errorf("error getting nodepool labels: %v", err)
	}

	var filteredLabels []govultr.NodePoolLabel
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		f := buildVultrDataSourceFilter(filters.(*schema.Set))
		for _, label := range labels {
			sm, err := structToMap(label)
			if err != nil {
				return diag.FromErr(err)
			}
			if filterLoop(f, sm) {
				filteredLabels = append(filteredLabels, label)
			}
		}
	} else {
		filteredLabels = labels
	}

	if err := d.Set("labels", flattenNodePoolLabels(filteredLabels)); err != nil {
		return diag.Errorf("error setting labels: %v", err)
	}

	d.SetId(fmt.Sprintf("%s-%s", clusterID, nodepoolID))

	return nil
}

func flattenNodePoolLabels(labels []govultr.NodePoolLabel) []map[string]interface{} {
	var result []map[string]interface{}

	for _, label := range labels {
		l := map[string]interface{}{
			"id":    label.ID,
			"key":   label.Key,
			"value": label.Value,
		}
		result = append(result, l)
	}

	return result
}
