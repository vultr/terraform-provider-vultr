output "ip" {
  value = vultr_instance.my_instance.main_ip
}

output "fwg_id" {
  value = vultr_firewall_group.fwg.id
}