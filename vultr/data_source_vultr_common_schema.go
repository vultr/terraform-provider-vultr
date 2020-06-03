package vultr

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"strings"
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
	err = json.Unmarshal(a, &structMap)

	newMap := make(map[string]interface{})
	for k, v := range structMap {
		switch v.(type) {
		case string:
			newMap[strings.ToLower(k)] = v.(string)
		case bool:
			newMap[strings.ToLower(k)] = strconv.FormatBool(v.(bool))
		case int:
			newMap[strings.ToLower(k)] = strconv.FormatInt(int64(v.(int)), 10)
		case float64:
			newMap[strings.ToLower(k)] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
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

func valuesLoop(values []string, i interface{}) bool {
	for _, v := range values {
		if v == i {
			return true
		}
	}
	return false
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
