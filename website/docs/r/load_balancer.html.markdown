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
  region = "ewr"
  label     = "vultr-load-balancer"
  balancing_algorithm = "roundrobin"
  nodes = 3

  forwarding_rules {
    frontend_protocol = "http"
    frontend_port = 82
    backend_protocol = "http"
    backend_port = 81
  }

  health_check {
    path = "/test"
    port = 8080
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

* `region` - (Required) The region your load balancer is deployed in.
* `nodes` - (Optional) The number of nodes to add to the load balancer (1-99). Must be an odd number. Default value is 1.
* `label` - (Optional) The load balancer's label.
* `balancing_algorithm` - (Optional) The balancing algorithm for your load balancer. Options are `roundrobin` or `leastconn`. Default value is `roundrobin`
* `proxy_protocol` - (Optional) Boolean value that indicates if Proxy Protocol is enabled.
* `cookie_name` - (Optional) Name for your given sticky session.
* `ssl_redirect` - (Optional) Boolean value that indicates if HTTP calls will be redirected to HTTPS.
* `http_version` - (Optional) Integer value that indicates if HTTP/2 or HTTP/3 is enabled. Allowed values 2 or 3.
* `attached_instances` - (Optional) Array of instances that are currently attached to the load balancer.
* `auto_ssl_domain` - (Optional) The auto SSL domain configuration for a load balancer. This can be a root domain (example.com) or include a subdomain (sub.example.com).
* `private_network` (Optional) (Deprecated: use `vpc` instead) A private network ID that the load balancer should be attached to.
* `vpc` (Optional)- A VPC ID that the load balancer should be attached to.

`ssl` (Optional) block supports the following elements:
* `private_key` - (Required) The SSL certificates private key.
* `certificate` - (Required) The SSL Certificate.
* `chain` - (Optional) The SSL certificate chain.

`health_check` (Optional) block supports the following elements:
* `protocol` - (Optional) The protocol used to traffic requests to the load balancer. Possible values are `http`, or `tcp`. Default value is `http`.
* `path` - (Optional) The path on the attached instances that the load balancer should check against. Default value is `/`
* `port` - (Optional) The assigned port (integer) on the attached instances that the load balancer should check against. Default value is `80`.
* `check_interval` - (Optional) Time in seconds to perform health check. Default value is 15.
* `response_timeout` - (Optional) Time in seconds to wait for a health check response. Default value is 5.
* `unhealthy_threshold` - (Optional) Number of failed attempts encountered before failover. Default value is 5.
* `healthy_threshold` - (Optional)  Number of failed attempts encountered before failover. Default value is 5. 

`forwarding_rules` (Required) is a list of items the following elements:
* `frontend_protocol` - (Required) Protocol on load balancer side. Possible values: "http", "https", "tcp".
* `frontend_port` - (Required) Port on load balancer side.
* `backend_protocol` - (Required) Protocol on instance side. Possible values: "http", "https", "tcp".
* `backend_port` - (Required) Port on instance side.

`firewall_rules` (Optional) is a list of items the following elements:
* `frontend_port` - (Required) Port on load balancer side.
* `ip_type` - (Required) The type of ip this rule is - may be either v4 or v6.
* `source` - (Required) IP address with subnet that is allowed through the firewall. You may also pass in `cloudflare` which will allow only CloudFlares IP range.

`global_regions` (Optional) is a list of items the following elements:
* `region_id` - (Required) The three letter region code which should have a child load balancer
* `vpc_id` - (Optional) The VPC ID to attach to the child load balancer

## Attributes Reference

The following attributes are exported:
* `id` - The load balancer ID.
* `region` - The region your load balancer is deployed in.
* `label` - The load balancer's label.
* `nodes` - The number of nodes for the load balancer.
* `balancing_algorithm` - The balancing algorithm for your load balancer.
* `proxy_protocol` - Boolean value that indicates if Proxy Protocol is enabled.
* `cookie_name` - Name for your given sticky session.
* `ssl_redirect` - Boolean value that indicates if HTTP calls will be redirected to HTTPS.
* `http_version` - Integer value that indicates if HTTP/2 or HTTP/3 is enabled.
* `has_ssl` - Boolean value that indicates if SSL is enabled.
* `auto_ssl_domain` - The auto SSL domain configuration for a load balancer.
* `attached_instances` - Array of instances that are currently attached to the load balancer.
* `status` - Current status for the load balancer
* `ipv4` - IPv4 address for your load balancer.
* `ipv6` - IPv6 address for your load balancer.
* `private_network` - (Deprecated: use `vpc` instead) Defines the private network the load balancer is attached to.
* `vpc` - Defines the VPC the load balancer is attached to.

`health_check` exports the following elements:
* `protocol` - The protocol used to traffic requests to the load balancer. Possible values are "http", "https", or "tcp".
* `path` - The path on the attached instances that the load balancer should check against.
* `port` - The assigned port (integer) on the attached instances that the load balancer should check against.
* `check_interval` - Time in seconds to perform health check. Default value is 15.
* `response_timeout` - Time in seconds to wait for a health check response. Default value is 5.
* `unhealthy_threshold` - Number of failed attempts encountered before failover. Default value is 5.
* `healthy_threshold` -  Number of failed attempts encountered before failover. Default value is 5. 

`forwarding_rules` is a list of items with the following elements:
* `frontend_protocol` - Protocol on load balancer side. Possible values: "http", "https", "tcp".
* `frontend_port` - Port on load balancer side.
* `backend_protocol` - Protocol on instance side. Possible values: "http", "https", "tcp".
* `target_port` - Port on instance side.

`firewall_rules` is a list of items with the following elements:
* `frontend_port` - Port on load balancer side.
* `ip_type` - The type of ip this rule is - may be either v4 or v6.
* `source` - IP address with subnet that is allowed through the firewall. You may also pass in `cloudflare` which will allow only CloudFlares IP range.

`global_regions` is a list of items with the following elements:
* `region_id` - The three letter region code of the child load balancer 
* `vpc_id` - The VPC ID attached to the child load balancer


## Import

Load Balancers can be imported using the load balancer `ID`, e.g.

```
terraform import vultr_load_balancer.lb-name b6a859c5-b299-49dd-8888-b1abbc517d08
```
