package vultr

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func nodePoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:     schema.TypeString,
			Required: true,
		},
		"plan": {
			Type:     schema.TypeString,
			Required: true,
		},
		"node_quantity": {
			Type:     schema.TypeInt,
			Required: true,
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
}
