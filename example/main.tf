terraform {
  required_providers {
    vultr = {
      source = "vultr/vultr"
      version = "2.16.3"
    }
  }
}

provider "vultr" {
  # In your .bashrc you need to set
  # export VULTR_API_KEY="Your Vultr API Key"
}

resource "vultr_instance" "my_instance" {
  plan                   = var.one_cpu_one_gb_ram
  region                 = var.vultr_seattle
  app_id                 = var.docker_centos
  label                  = "terraform example"
  enable_ipv6            = true
  backups                = "enabled"
  enable_private_network = true
  activation_email       = false
  ddos_protection        = true
  tags                   = ["my_tag"]
  firewall_group_id      = vultr_firewall_group.fwg.id
  backups_type {
    type = "daily"
  }
}

resource "vultr_firewall_group" "fwg" {
  description = "docker-fwg"
}

resource "vultr_firewall_rule" "tcp" {
  firewall_group_id = vultr_firewall_group.fwg.id
  protocol          = "udp"
  subnet            = vultr_instance.my_instance.main_ip
  subnet_size       = 32
  port              = "8080"
  ip_type           = "v4"
}

resource "vultr_dns_domain" "my_domain" {
  domain    = "tf-domain.com"
  ip        = vultr_instance.my_instance.main_ip
}

resource "vultr_dns_record" "a-record" {
  data   = vultr_instance.my_instance.main_ip
  domain = vultr_dns_domain.my_domain.id
  name   = "www"
  type   = "A"
  ttl    = 3600
}

resource "vultr_load_balancer" "lb" {
  region              = "ewr"
  label               = "terraform lb example"
  balancing_algorithm = "roundrobin"

  forwarding_rules {
    frontend_protocol = "http"
    frontend_port     = 80
    backend_protocol  = "http"
    backend_port      = 80
  }

  health_check {
    protocol            = "http"
    port                = 80
    path                = "/health"
    check_interval      = 15
    response_timeout    = 5
    unhealthy_threshold = 5
    healthy_threshold   = 5
  }

}
