---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "omglol_account_settings Resource - omglol"
subcategory: ""
description: |-
  Manage omg.lol account settings through terraform. This resource will update the exisiting account settings, it cannot be imported or destroyed. Specifying more than one of this resource will have unpredictable results.
---

# omglol_account_settings (Resource)

Manage omg.lol account settings through terraform. This resource will update the exisiting account settings, it cannot be imported or destroyed. Specifying more than one of this resource will have unpredictable results.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `communication` (String) Commuinication preferences. Valid values are `email_ok` and `email_not_ok`
- `date_format` (String) Date preferences. Valid values are: `iso_8601` for *YYYY-MM-DD*, `dmy` for *DD-MM-YYYY*, and `mdy` for *MM-DD-YYYY*.

### Read-Only

- `id` (String) The ID of this resource.
- `last_updated` (String)

