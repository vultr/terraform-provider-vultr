package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVultrKubernetesVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrKubernetesVersionsRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func dataSourceVultrKubernetesVersionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	versions, _, err := client.Kubernetes.GetVersions(ctx)
	if err != nil {
		return diag.Errorf("error getting kubernetes versions: %v", err)
	}

	if err := d.Set("versions", versions); err != nil {
		return diag.Errorf("unable to set kubernetes `versions` read value: %v", err)
	}

	return nil
}
