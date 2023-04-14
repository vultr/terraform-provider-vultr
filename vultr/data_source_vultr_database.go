package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrDatabaseRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"plan_ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"plan_vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"plan_replicas": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_dow": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_backup": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trusted_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"mysql_sql_modes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"mysql_require_primary_key": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mysql_slow_query_log": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mysql_long_query_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"redis_eviction_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var databaseList []govultr.Database
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.DBListOptions{}
	databases, _, _, err := client.Database.List(ctx, options)
	if err != nil {
		return diag.Errorf("error getting databases: %v", err)
	}

	for s := range databases {
		// we need convert the a struct INTO a map so we can easily manipulate the data here
		sm, err := structToMap(databases[s])

		if err != nil {
			return diag.FromErr(err)
		}

		if filterLoop(f, sm) {
			databaseList = append(databaseList, databases[s])
		}
	}

	if len(databaseList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(databaseList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(databaseList[0].ID)
	if err := d.Set("date_created", databaseList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set database `date_created` read value: %v", err)
	}

	return nil
}
