package vultr

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Lookup changes on a TF field and convert schema.Set to []string
func tfChangeToSlices(fieldname string, d *schema.ResourceData) ([]string, []string) {
	oldVal, newVal := d.GetChange(fieldname)

	oldSlice := []string{}
	for _, v := range oldVal.(*schema.Set).List() {
		oldSlice = append(oldSlice, v.(string))
	}

	newSlice := []string{}
	for _, v := range newVal.(*schema.Set).List() {
		newSlice = append(newSlice, v.(string))
	}

	return oldSlice, newSlice
}

// Compare two slices and return elements that are in x but not in y
func diffSlice(x, y []string) []string {
	var diff []string

	b := map[string]string{}
	for i := range x {
		b[x[i]] = ""
	}

	for i := range y {
		if _, ok := b[y[i]]; !ok {
			diff = append(diff, y[i])
		}
	}

	return diff
}

// IgnoreCase implement a DiffSupressFunc to ignore case
func IgnoreCase(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}
