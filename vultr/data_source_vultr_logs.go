package vultr

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrLogsRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_status_code": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_parameters": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"continue_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"returned_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unreturned_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrLogsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	logs, logsMeta, _, err := client.Logs.List(ctx, govultr.LogsOptions{
		StartTime:    d.Get("start_time").(string),
		EndTime:      d.Get("end_time").(string),
		LogLevel:     d.Get("log_level").(string),
		ResourceType: d.Get("resource_type").(string),
		ResourceID:   d.Get("resource_id").(string),
	})
	if err != nil {
		return diag.Errorf("error getting logs : %v", err)
	}

	if len(logs) < 1 {
		return diag.Errorf("no logs were found")
	}

	ts := time.Now()
	d.SetId(strconv.Itoa(int(ts.UnixMilli())))
	if err := d.Set("continue_time", logsMeta.ContinueTime); err != nil {
		return diag.Errorf("unable to set logs `continue_time` read value: %v", err)
	}
	if err := d.Set("returned_count", logsMeta.ReturnedCount); err != nil {
		return diag.Errorf("unable to set logs `returned_count` read value: %v", err)
	}
	if err := d.Set("unreturned_count", logsMeta.UnreturnedCount); err != nil {
		return diag.Errorf("unable to set logs `unreturned_count` read value: %v", err)
	}
	if err := d.Set("total_count", logsMeta.TotalCount); err != nil {
		return diag.Errorf("unable to set logs `total_count` read value: %v", err)
	}

	allLogs := []map[string]interface{}{}
	for i := range logs {
		aLog := map[string]interface{}{
			"level":            logs[i].Level,
			"resource_id":      logs[i].ResourceID,
			"resource_type":    logs[i].ResourceType,
			"message":          logs[i].Message,
			"timestamp":        logs[i].Timestamp,
			"user_id":          logs[i].Metadata.UserID,
			"user_name":        logs[i].Metadata.UserName,
			"ip_address":       logs[i].Metadata.IPAddress,
			"http_status_code": logs[i].Metadata.HTTPStatusCode,
			"method":           logs[i].Metadata.Method,
			"request_path":     logs[i].Metadata.RequestPath,
			"request_body":     logs[i].Metadata.RequestBody,
			"query_parameters": logs[i].Metadata.QueryParameters,
		}

		allLogs = append(allLogs, aLog)
	}

	if err := d.Set("results", allLogs); err != nil {
		return diag.Errorf("unable to set logs `results` read value: %v", err)
	}

	return nil
}
