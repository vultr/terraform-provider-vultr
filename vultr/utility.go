package vultr

import (
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Lookup changes on a TF field and convert schema.Set to []string
func tfChangeToSlices(fieldname string, d *schema.ResourceData) ([]string, []string) { //nolint:unparam
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
	return strings.EqualFold(old, new)
}

// suppressIPDiff returns true when old and new are the same IP address
// just written differently (leading zeros, mixed case, etc).
// Handles both v4 and v6 transparently.
func suppressIPDiff(_, old, new string, _ *schema.ResourceData) bool {
	oldIP := net.ParseIP(old)
	newIP := net.ParseIP(new)

	// if either side doesn't parse, fall through to normal string compare
	if oldIP == nil || newIP == nil {
		return old == new
	}

	return oldIP.Equal(newIP)
}

// canonicalizeIP parses and re-renders an IP address to its canonical
// string form. For v6 this strips leading zeros and lowercases per
// RFC 5952, matching what the Vultr API returns.
// If it can't parse (bad input or empty string), just pass through.
func canonicalizeIP(val interface{}) string {
	raw := val.(string)
	ip := net.ParseIP(raw)
	if ip == nil {
		return raw
	}
	return ip.String()
}
