package vultr

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Lookup changes on a TF field and convert schema.Set to []string
func tfChangeToSlices(fieldname string, d *schema.ResourceData) ([]string, []string) {
	oldVal, newVal := d.GetChange(fieldname)

	var oldSlice []string
	for _, v := range oldVal.(*schema.Set).List() {
		oldSlice = append(oldSlice, v.(string))
	}

	var newSlice []string
	for _, v := range newVal.(*schema.Set).List() {
		newSlice = append(newSlice, v.(string))
	}

	return oldSlice, newSlice
}
