package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/mod/semver"
)

func dataSourceVultrKubernetesVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrKubernetesVersionRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"available": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"upgrades": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"latest": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrKubernetesVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Client).govultrClient()

	versions, _, err := client.Kubernetes.GetVersions(ctx)
	if err != nil {
		return diag.Errorf("unable to retrieve kubernetes versions: %v", err)
	}

	filter, filterOk := d.GetOk("filter")

	if filterOk && !semver.IsValid(filter.(string)) {
		warn := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "filter is not a valid semver string and cannot be compared",
		}

		diags = append(diags, warn)
	}

	versionList := []string{}
	upgradeList := []string{}
	for i := range versions.Versions {

		if filterOk {
			compResult := semver.MajorMinor(versions.Versions[i])
			compFilter := semver.MajorMinor(filter.(string))
			compare := semver.Compare(compResult, compFilter)
			if compare < 0 {
				continue
			}

			if compare > 0 {
				upgradeList = append(upgradeList, versions.Versions[i])
			}
		}

		versionList = append(versionList, versions.Versions[i])
	}

	semver.Sort(versionList)

	d.SetId("-") // no natural ID but still let terraform know this exists

	if len(versionList) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "no versions available, check your filter",
		})

		return diags
	}

	if err := d.Set("available", versionList); err != nil {
		return diag.Errorf("unable to set kubernetes version `versions` read value: %v", err)
	}
	if err := d.Set("upgrades", upgradeList); err != nil {
		return diag.Errorf("unable to set kubernetes version `upgrades` read value: %v", err)
	}
	if err := d.Set("latest", versionList[len(versionList)-1]); err != nil {
		return diag.Errorf("unable to set kubernetes version `latest` read value: %v", err)
	}

	return diags
}
