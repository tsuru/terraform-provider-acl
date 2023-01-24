---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "acl_destination_rule Resource - terraform-provider-acl"
subcategory: ""
description: |-
  
---

# acl_destination_rule (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `instance` (String) ACL Instance Name

### Optional

- `app` (String)
- `dns` (String)
- `ip` (String)
- `pool` (String)
- `port` (Block List) (see [below for nested schema](#nestedblock--port))
- `rpaas` (Block List, Max: 1) (see [below for nested schema](#nestedblock--rpaas))
- `service_name` (String) ACL Service Name

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--port"></a>
### Nested Schema for `port`

Required:

- `number` (Number)
- `protocol` (String)


<a id="nestedblock--rpaas"></a>
### Nested Schema for `rpaas`

Optional:

- `instance` (String)
- `service_name` (String)

