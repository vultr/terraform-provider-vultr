package vultr

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type filter struct {
	name   string
	values []string
}

func buildVultrDataSourceFilter(set *schema.Set) []filter {
	var filters []filter

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var values []string
		for _, value := range m["values"].([]interface{}) {
			values = append(values, value.(string))
		}
		filters = append(filters, filter{
			name:   m["name"].(string),
			values: values,
		})
	}

	return filters
}

func structToMap(data interface{}) (map[string]interface{}, error) {
	var structMap map[string]interface{}

	a, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(a, &structMap); err != nil {
		return nil, err
	}

	newMap := make(map[string]interface{})
	for k, v := range structMap {
		switch val := v.(type) {
		case string:
			newMap[strings.ToLower(k)] = val
		case bool:
			newMap[strings.ToLower(k)] = val
		case int:
			newMap[strings.ToLower(k)] = strconv.FormatInt(int64(val), 10)
		case float64:
			newMap[strings.ToLower(k)] = strconv.FormatFloat(val, 'f', -1, 64)
		default:
			newMap[strings.ToLower(k)] = v
		}
	}

	return newMap, nil
}

func filterLoop(f []filter, m map[string]interface{}) bool {
	for _, filter := range f {
		if !valuesLoop(filter.values, m[filter.name]) {
			return false
		}
	}
	return true
}

func valuesLoop(values []string, actual interface{}) bool {
	switch a := actual.(type) {
	case []interface{}:
		// It's an array of strings, so something like: location
		var found bool
		for _, i := range values {
			found = false
			for _, j := range a {
				if i == j {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	default:
		// It's a string, so something like: ram, type, vcpu_count
		for _, i := range values {
			if actual == i {
				return true
			}
		}
		return false
	}
}

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}
