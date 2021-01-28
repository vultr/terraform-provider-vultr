##################################################################################
#  Note: 
##################################################################################
# VARIABLES
##################################################################################
variable "vultr_api_key"   {}
variable "region"          {}
variable "public_key_file" {}
variable "startup_script"  {}


##################################################################################
# PROVIDERS
##################################################################################

provider "vultr" {
	 api_key = var.vultr_api_key
}

##################################################################################
# DATA
##################################################################################
#data "aws_vpc" "default" {
#   default = true
#}


##################################################################################
# RESOURCES
##################################################################################

resource "vultr_ssh_key" "bastion" {
	 name    = "vultr"
   ssh_key = file(var.public_key_file)	 
}

resource "vultr_startup_script" "bastion" {
	name     = "startup"
	script   = base64encode( file(var.startup_script) )
}

resource "vultr_instance" "bastion" {
	 plan        = "vc2-1c-1gb"
	 region      = var.region
	 os_id       = 387
	 ssh_key_ids = [vultr_ssh_key.bastion.id]
	 script_id   = vultr_startup_script.bastion.id
	 enable_ipv6      = false
	 activation_email = false

   connection {
       type     = "ssh"
       host     = self.main_ip
       user     = "root"
       password = self.default_password
    }   
    provisioner "remote-exec" {
       inline = [
          "echo ${self.main_ip} > /var/tmp/externalip",
          "/usr/bin/apt-get update "
       ]
    }   
}

resource "vultr_instance_ipv4" "bastion_ipv4" {
	 instance_id = vultr_instance.bastion.id
}

### OUTPUT ###
output "instance_main_ip" {
	 value = vultr_instance.bastion.main_ip
}
output "instance_defaultpassword" {
	 value = vultr_instance.bastion.default_password
}
