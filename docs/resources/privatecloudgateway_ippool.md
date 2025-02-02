---
page_title: "spectrocloud_privatecloudgateway_ippool Resource - terraform-provider-spectrocloud"
subcategory: ""
description: |-
  
---

# Resource `spectrocloud_privatecloudgateway_ippool`





## Schema

### Required

- **gateway** (String)
- **name** (String)
- **network_type** (String)
- **prefix** (Number)
- **private_cloud_gateway_id** (String)

### Optional

- **id** (String) The ID of this resource.
- **ip_end_range** (String)
- **ip_start_range** (String)
- **nameserver_addresses** (Set of String)
- **nameserver_search_suffix** (Set of String)
- **restrict_to_single_cluster** (Boolean)
- **subnet_cidr** (String)
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)


