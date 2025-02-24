---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_backup_users"
description: |-
  Provides information about users with the PaaS Backup User project privileges.
---

# Data Source: aws_paas_backup_users

Provides information about users with the PaaS Backup User project privileges.

## Example Usage

```terraform
data "aws_paas_backup_users" "selected" {
  active_only = true
}

output "backup-user-logins" {
  value = data.aws_paas_backup_users.selected.users[*].login
}
```

## Argument Reference

The following arguments are supported:

* `active_only` - (Optional) Indicates whether to filter only active users.
    * _Default value_: `false`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The region.
* `users` - List of users. Each user has the following structure:
    * `email` - The user email.
    * `enabled` - Indicates whether the user is active.
    * `id` - The ID of the user.
    * `login` - The user login.
    * `name` - The user name.
