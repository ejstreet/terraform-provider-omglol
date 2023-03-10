---
page_title: "omglol_purls Data Source - omglol"
subcategory: ""
description: |-
  List all PURLs for a given omg.lol address.
---

# omglol_purls (Data Source)

List all PURLs for a given omg.lol address.

## Example Usage

```terraform
data omglol_purls example {
  address = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `address` (String) The omg.lol address to read the purls from.

### Read-Only

- `purls` (Attributes List) A list of all the PURLs for the given address. (see [below for nested schema](#nestedatt--purls))

<a id="nestedatt--purls"></a>
### Nested Schema for `purls`

Read-Only:

- `counter` (Number) The number of time a PURL has been accessed.
- `id` (String) Unique ID of the PURL. Can be used for imports.
- `listed` (Boolean) Returns `true` if listed on your `address`.url.lol page.
- `name` (String) The name of the PURL. The name field is how you will access your designated URL.
- `url` (String) The url that is pointed to.