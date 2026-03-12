---
layout: "vultr"
page_title: "Vultr: vultr_logs"
sidebar_current: "docs-vultr-datasource-logs"
description: |-
  Retrieve Vultr API logs
---

# vultr_logs

Retrieve Vultr API logs

## Example Usage

Define a unique log set by supplying a name/ID and at least the required arguments:

```hcl
data "vultr_logs" "my_logs" {
  start_time = "2026-02-26T00:00:00Z"
  end_time = "2026-02-27T00:00:00Z"
  log_level = "debug"
  resource_type = "kubernetes"
  resource_id = "6164bc1f-a8bf-40a2-b3e2-b8a8b436a9fd"
}
```

## Argument Reference

The following arguments are supported:

* `start_time` - A UTC timestamp for the start of the time period from which to return logs. 

Logs with a timestamp equal to, or after `start_time` are included in the response.
Expected Format: `yyyy-mm-ddThh:mm:ssZ` (i.e. `2025-06-26T00:00:00Z`)
`start_time` and `end_time` may not be more than 30 days and 1 hour apart

* `end_time` - A UTC timestamp for the end of the time period from which to return logs. 

Only logs with a timestamp before the `end_time` are included in the response. 
Expected Format: `yyyy-mm-ddThh:mm:ssZ` (i.e. `2025-06-26T00:00:00Z`)
`start_time` and `end_time` may not be more than 30 days and 1 hour apart

* `log_level` - (Optional) Level with which to filter the logs by (must be one of `info`, `debug`, `warning`, `error`, or `critical`). 
* `resource_type` - (Optional) Filter the logs by the type of a resource (i.e. `instances`, `kubernetes`, `bare-metals`).
* `resource_id` - (Optional) Filter the logs by the UUID of a specific resource.

## Attributes Reference

The following attributes are exported:

* `results` - The list of log files returned in the query. Made up of:

  `level` - The log level of the result.
  `resource_id` - The UUID for the resource that was interacted with. Only set if the logged interaction relates to a specific resource with a UUID.
  `resource_type` - The type of resource that was interacted with.
  `message` - A message relating to the event that is being logged.
  `timestamp` - A UTC timestamp of the time at which the log was generated. 
  `user_id` - The UUID for the user who triggered the event that is being logged.
  `user_name` - The email address of a user who is logging in (if applicable).
  `ip_address` - The IP address from which the request that generated the log originated.
  `http_status_code` - The status code returned for and API request (if applicable).
  `method` - The HTTP request method of the API request being logged (if applicable).
  `request_path` - The URI path of the API request being logged (if applicable).
  `request_body` - The request body provided for the API request being logged (if applicable).
  `query_parameters` - The query string provided for the API request being logged (if applicable).

* `continue_time` - In the event that there are more logs found for a specified time period that can be returned, this field will be set with a UTC timestamp of where the logs were left off.
* `returned_count` - The number of log records that were returned. There is a maximum limit of 5,000 logs returned by any request.
* `unreturned_count` - The number of log records from the specified time period that were not returned due to the maximum return limit of 5,000 logs.
* `total_count` - The total number of records that were found for the specified time period.
