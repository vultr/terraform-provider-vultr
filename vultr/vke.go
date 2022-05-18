package vultr

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func nodePoolSchema(isNodePool bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"label": {
			Type:     schema.TypeString,
			Required: true,
		},
		"plan": {
			Type:     schema.TypeString,
			Required: true,
		},
		"node_quantity": {
			Type:         schema.TypeInt,
			ValidateFunc: validation.IntAtLeast(1),
			Required:     true,
		},
		"auto_scaler": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"min_nodes": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"max_nodes": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		//computed fields
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"date_created": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"date_updated": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nodes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"date_created": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"label": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	if isNodePool {
		s["cluster_id"] = &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
		}
		s["tag"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		}

	} else {
		// Make tags unmodifiable for the vultr_kubernetes resource
		// This lets us know which node pool was part of the vultr_kubernetes resource
		s["tag"] = &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		}
	}

	return s
}
