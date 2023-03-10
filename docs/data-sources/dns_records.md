---
page_title: "omglol_dns_records Data Source - omglol"
subcategory: ""
description: |-
  List all DNS records for a given omg.lol address.
---

# omglol_dns_records (Data Source)

List all DNS records for a given omg.lol address.

## Example Usage

```terraform
data omglol_dns_records example {
  address = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `address` (String) The omg.lol address to read the records from.

### Read-Only

- `id` (String) The ID of this resource.
- `records` (Attributes List) A list of all the DNS records for the given address. (see [below for nested schema](#nestedatt--records))

<a id="nestedatt--records"></a>
### Nested Schema for `records`

Read-Only:

- `created_at` (String)
- `data` (String) The data entered into the record.
- `fqdn` (String) The fully qualified domain name of the record. Made by combining DNS name, address, and omg.lol top-level.
- `id` (Number)
- `name` (String) The prefix attached before the address. `@` represents the apex.
- `priority` (Number) The priority of the record. Only applies to MX records.
- `ttl` (Number) The Time-To-Live (TTL) of the record.
- `type` (String) The record type.
- `updated_at` (String)