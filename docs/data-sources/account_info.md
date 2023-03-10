---
page_title: "omglol_account_info Data Source - omglol"
subcategory: ""
description: |-
  Retrieve account info.
---

# omglol_account_info (Data Source)

Retrieve account info.

## Example Usage

```terraform
data omglol_account_info this {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `created` (String) The RFC 3339 representation of the time that the account was created. This can be used in conjunction with the [formatdate](https://developer.hashicorp.com/terraform/language/functions/formatdate) function.
- `email` (String) The email address associated with the account.
- `id` (String) The ID of this resource.
- `name` (String) The name associated with the account.