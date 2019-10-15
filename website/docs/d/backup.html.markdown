---
layout: "vultr"
page_title: "Vultr: vultr_backup"
sidebar_current: "docs-vultr-datasource-backup"
description: |-
  Get information about a Vultr backup.
---

# vultr_backup

Get information about a Vultr backup. This data source provides a list of backups which contain the description, size, status, and the creation date for your Vultr backup.

## Example Usage

Get the information for a backup by `description`:

```hcl
data "vultr_backup" "my_backup" {
  filter {
    name = "description"
    values = ["my-backup-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding backups.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `BACKUPID` - The ID of the backup
* `description` - The description of the backup.
* `size` - The size of the backup in bytes.
* `status` - The status of the backup.
* `date_created` - The date the backup was added to your Vultr account.