---
layout: "vultr"
page_title: "Vultr: vultr_load_balancer"
sidebar_current: "docs-vultr-resource-load-balancer"
description: |-
  Get information about a Vultr Load Balancer.
---

# vultr_load_balancer

Get information about a Vultr load balancer.

## Example Usage

Create a new load balancer:

```hcl
resource "vultr_load_balancer" "lb" {
  region_id = 1
  label     = "terraform lb example"
  balancing_algorithm = "roundrobin"

  forwarding_rules {
    frontend_protocol = "http"
    frontend_port = 82
    backend_protocol = "http"
    backend_port = 81
  }

  health_check {
    path = "/test"
    port = "8080"
    protocol = "http"
    response_timeout = 1
    unhealthy_threshold =2 
    check_interval = 3
    healthy_threshold =4
  }
}
```

## Argument Reference

The follow arguments are supported:

* `region_id` - (Required) The region your load balancer is deployed in.
* `forwarding_rules` - (Required) List of forwarding rules for a load balancer. The configuration of a `forwarding_rules` is listened below.
* `label` - (Optional) The load balancers label.
* `balancing_algorithm` - (Optional) The balancing algorithm for your load balancer. Options are `roundrobin` or `leastconn`
* `proxy_protocol` - (Optional) Boolean value that indicates if Proxy Protocol is enabled.
* `cookie_name` - (Optional) Name for your given sticky session.
* `ssl_redirect` - (Optional) Boolean value that indicates if HTTP calls will be redirected to HTTPS.
* `attached_instances` - (Optional) Array of instances that are currently attached to the load balancer.
* `health_check` - (Optional) A block that defines the way load balancers should check for health. The configuration of a `health_check` is listed below.
* `ssl` - (Optional) A block that supplies your ssl configuration to be used with HTTPS. The configuration of a `ssl` is listed below.

`health_check` supports the following

* `protocol` - (Optional) The protocol used to traffic requests to the load balancer. Possible values are `http`, or `tcp`. Default value is `http`.
* `path` - (Optional) The path on the attached instances that the load balancer should check against. Default value is `/`
* `port` - (Optional) The assigned port (integer) on the attached instances that the load balancer should check against. Default value is `80`.
* `check_interval` - (Optional) Time in seconds to perform health check. Default value is 15.
* `response_timeout` - (Optional) Time in seconds to wait for a health check response. Default value is 5.
* `unhealthy_threshold` - (Optional) Number of failed attempts encountered before failover. Default value is 5.
* `healthy_threshold` - (Optional)  Number of failed attempts encountered before failover. Default value is 5. 

`forwarding_rules` supports the following

* `frontend_protocol` - (Required) Protocol on load balancer side. Possible values: "http", "https", "tcp".
* `frontend_port` - (Required) Port on load balancer side.
* `backend_protocol` - (Required) Protocol on instance side. Possible values: "http", "https", "tcp".
* `target_port` - (Required) Port on instance side.

`ssl` supports the following

* `private_key` - (Required) The SSL certificates private key.
* `certificate` - (Required) The SSL Certificate.
* `chain` - (Optional) The SSL certificate chain.

## Attributes Reference

The following attributes are exported:
* `id` - The load balancer ID.
* `region_id` - The region your load balancer is deployed in.
* `label` - The load balancers label.
* `balancing_algorithm` - The balancing algorithm for your load balancer.
* `proxy_protocol` - Boolean value that indicates if Proxy Protocol is enabled.
* `cookie_name` - Name for your given sticky session.
* `ssl_redirect` - Boolean value that indicates if HTTP calls will be redirected to HTTPS.
* `has_ssl` - Boolean value that indicates if SSL is enabled.
* `attached_instances` - Array of instances that are currently attached to the load balancer.
* `status` - Current status for the load balancer
* `ipv4` - IPv4 address for your load balancer.
* `ipv6` - IPv6 address for your load balancer.
* `health_check` - Defines the way load balancers should check for health. 
* `forwarding_rules` - Defines the forwarding rules for a load balancer.
