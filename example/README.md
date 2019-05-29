# Vultr VPS, DNS Domain/Record, and Firewall Group/Rule example

This example launches a single Vultr VPS and configures a DNS Domain + Record along with a firewall group and rule that are then mapped to the VPS.


To run this example, you will need to preconfigure your Vultr provider as described in the [Read Me](../README.md)

## What will you need?
You will need to export your Vultr API Key as a environmental variable

```
export VULTR_API_KEY="Your Vultr API Key"
```

## How to run?

Once you have exported your Vultr API Key you can run the following

```
terraform plan
terraform apply
```

To delete all the resources that were created run
```
terraform destroy
```