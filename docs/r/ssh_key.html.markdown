---
layout: "vultr"
page_title: "Vultr: vultr_ssh_key"
sidebar_current: "docs-vultr-resource-ssh-key"
description: |-
  Provides a Vultr SSH key resource. This can be used to create, read, modify, and delete SSH keys.
---

# vultr_ssh_key

Provides a Vultr SSH key resource. This can be used to create, read, modify, and delete SSH keys.

## Example Usage

Create an SSH key
```hcl
resource "vultr_ssh_key" "my_ssh_key" {
  name = "my-ssh-key"
  ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyVGaw1PuEl98f4/7Kq3O9ZIvDw2OFOSXAFVqilSFNkHlefm1iMtPeqsIBp2t9cbGUf55xNDULz/bD/4BCV43yZ5lh0cUYuXALg9NI29ui7PEGReXjSpNwUD6ceN/78YOK41KAcecq+SS0bJ4b4amKZIJG3JWmDKljtv1dmSBCrTmEAQaOorxqGGBYmZS7NQumRe4lav5r6wOs8OACMANE1ejkeZsGFzJFNqvr5DuHdDL5FAudW23me3BDmrM9ifUzzjl1Jwku3bnRaCcjaxH8oTumt1a00mWci/1qUlaVFft085yvVq7KZbF2OPPbl+erDW91+EZ2FgEi+v1/CSJ5 your_username@hostname"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name/label of the SSH key.
* `ssh_key` - (Required) The public SSH key.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SSH key.
* `name` - The name/label of the SSH key.
* `ssh_key` - The public SSH key.
* `date_created` - The date the SSH key was added to your Vultr account.

## Import

SSH keys can be imported using the SSH key `SSHKEYID`, e.g.

```
terraform import vultr_ssh_key.my_key 541b4960f23bd
```