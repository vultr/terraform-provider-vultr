# Tutorials

The tutorials in this directory give some insights in creating the 
server infrastructure using Vultr as cloud provider.

To run these tutorials, you will need to preconfigure your account as 
described in the [Read Me](../README.md)

## What will you need?
You will need to generate a Vultr API Key.
Some tutorial may require that you set up the API Key as bash
environment variable.

```
export VULTR_API_KEY="Your Vultr API Key"
```
In other tutorial, the API Key is configured in the file,
terraform.tfvars

## How to run?

Once you have exported your Vultr API Key you can run the following

```
terraform init
terraform plan -out temp.tfplan
terraform apply temp.tfplan
```

To delete all the resources that were created run
```
terraform destroy
```
