variable "vultr_seattle" {
  description = "Vultr Seattle Region"
  default = "4"
}

variable "os_type" {
  description = "Application"  
  default = 186
}

variable "docker_centos" {
  description = "Docker on CentOS 7"
  default = 17
}

variable "one_cpu_one_gb_ram" {
  description = "1024 MB RAM,25 GB SSD,1.00 TB BW"
  default = 201
}
