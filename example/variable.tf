variable "vultr_seattle" {
  description = "Vultr Seattle Region"
  default = "sea"
}

variable "docker_centos" {
  description = "Docker on CentOS 7"
  default = 17
}

variable "one_cpu_two_gb_ram" {
  description = "2048 MB RAM,25 GB SSD,1.00 TB BW"
  default = "vc2-1c-2gb"
}
