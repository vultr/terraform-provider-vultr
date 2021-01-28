# Vultr VPS as bastion host

This example launches a single Vultr VPS to be used as bastion host. 
Features of this tutorial are as follows,
  + the startup script creates a non-root user id, foobar, to be 
    configured later as admin user (hint: add it to sudo list)

  + make use of the terraform.tfvars

  + make use of the  provisioner "remote-exec"

  + the initial root password will be shown after the execution of
    `terraform apply` 

## What will you need?
You will need to include your Vultr API Key in the setting as
prescribed in the terraform.tfvars 

