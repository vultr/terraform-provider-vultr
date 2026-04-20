package vultr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiError struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

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

func checkIsMissing(e error, missingMsg string) (bool, error) {
	apiError := apiError{}
	if err := json.Unmarshal([]byte(e.Error()), &apiError); err != nil {
		return false, fmt.Errorf("unable to unmarshal api response: %v", err)
	}

	if apiError.Status == http.StatusNotFound {
		return true, nil
	}

	if missingMsg != "" && strings.Contains(apiError.Message, missingMsg) {
		return true, nil
	}

	return false, nil
}

// IgnoreCase implement a DiffSupressFunc to ignore case
func IgnoreCase(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
