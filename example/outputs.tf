output "ip" {
  value = "${vultr_server.my_server.main_ip}"
}

output "fwg_id" {
  value = "${vultr_firewall_group.fwg.id}"
}