---
page_title: "spectrocloud_cluster_eks Resource - terraform-provider-spectrocloud"
subcategory: ""
description: |-
  
---

# Resource `spectrocloud_cluster_eks`



## Example Usage

```terraform
data "spectrocloud_cloudaccount_aws" "account" {
  # id = <uid>
  name = var.cluster_cloud_account_name
}

data "spectrocloud_cluster_profile" "profile" {
  # id = <uid>
  name = var.cluster_cluster_profile_name
}

data "spectrocloud_backup_storage_location" "bsl" {
  name = var.backup_storage_location_name
}

resource "spectrocloud_cluster_eks" "cluster" {
  name             = var.cluster_name
  tags             = ["dev", "department:devops", "owner:bob"]
  cloud_account_id = data.spectrocloud_cloudaccount_aws.account.id

  cloud_config {
    ssh_key_name = "default"
    region       = "us-west-2"
  }

  cluster_profile {
    id = data.spectrocloud_cluster_profile.profile.id

    # To override or specify values for a cluster:

    # pack {
    #   name   = "spectro-byo-manifest"
    #   tag    = "1.0.x"
    #   values = <<-EOT
    #     manifests:
    #       byo-manifest:
    #         contents: |
    #           # Add manifests here
    #           apiVersion: v1
    #           kind: Namespace
    #           metadata:
    #             labels:
    #               app: wordpress
    #               app2: wordpress2
    #             name: wordpress
    #   EOT
    # }
  }

  backup_policy {
    schedule                  = "0 0 * * SUN"
    backup_location_id        = data.spectrocloud_backup_storage_location.bsl.id
    prefix                    = "prod-backup"
    expiry_in_hour            = 7200
    include_disks             = true
    include_cluster_resources = true
  }

  scan_policy {
    configuration_scan_schedule = "0 0 * * SUN"
    penetration_scan_schedule   = "0 0 * * SUN"
    conformance_scan_schedule   = "0 0 1 * *"
  }

  machine_pool {
    control_plane           = true
    control_plane_as_worker = true
    name                    = "master-pool"
    count                   = 1
    instance_type           = "t3.large"
    disk_size_gb            = 62
    az_subnets = {
      "us-west-2a" = "subnet-0d4978ddbff16c"
    }
  }

  machine_pool {
    name          = "worker-basic"
    count         = 1
    instance_type = "t3.large"
    az_subnets = {
      "us-west-2a" = "subnet-0d4978ddbff16c"
    }
  }

}
```

## Schema

### Required

- **cloud_account_id** (String)
- **cloud_config** (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--cloud_config))
- **machine_pool** (Block List, Min: 1) (see [below for nested schema](#nestedblock--machine_pool))
- **name** (String)

### Optional

- **backup_policy** (Block List, Max: 1) (see [below for nested schema](#nestedblock--backup_policy))
- **cluster_profile** (Block List) (see [below for nested schema](#nestedblock--cluster_profile))
- **cluster_profile_id** (String, Deprecated)
- **fargate_profile** (Block List) (see [below for nested schema](#nestedblock--fargate_profile))
- **id** (String) The ID of this resource.
- **pack** (Block List) (see [below for nested schema](#nestedblock--pack))
- **scan_policy** (Block List, Max: 1) (see [below for nested schema](#nestedblock--scan_policy))
- **tags** (Set of String)
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-only

- **cloud_config_id** (String)
- **kubeconfig** (String)

<a id="nestedblock--cloud_config"></a>
### Nested Schema for `cloud_config`

Required:

- **region** (String)

Optional:

- **az_subnets** (Map of String)
- **azs** (List of String)
- **endpoint_access** (String)
- **public_access_cidrs** (Set of String)
- **ssh_key_name** (String)
- **vpc_id** (String)


<a id="nestedblock--machine_pool"></a>
### Nested Schema for `machine_pool`

Required:

- **count** (Number)
- **disk_size_gb** (Number)
- **instance_type** (String)
- **name** (String)

Optional:

- **az_subnets** (Map of String)
- **azs** (List of String)


<a id="nestedblock--backup_policy"></a>
### Nested Schema for `backup_policy`

Required:

- **backup_location_id** (String)
- **expiry_in_hour** (Number)
- **prefix** (String)
- **schedule** (String)

Optional:

- **include_cluster_resources** (Boolean)
- **include_disks** (Boolean)
- **namespaces** (Set of String)


<a id="nestedblock--cluster_profile"></a>
### Nested Schema for `cluster_profile`

Required:

- **id** (String) The ID of this resource.

Optional:

- **pack** (Block List) (see [below for nested schema](#nestedblock--cluster_profile--pack))

<a id="nestedblock--cluster_profile--pack"></a>
### Nested Schema for `cluster_profile.pack`

Required:

- **name** (String)
- **tag** (String)
- **values** (String)



<a id="nestedblock--fargate_profile"></a>
### Nested Schema for `fargate_profile`

Required:

- **name** (String)
- **selector** (Block List, Min: 1) (see [below for nested schema](#nestedblock--fargate_profile--selector))

Optional:

- **additional_tags** (Map of String)
- **subnets** (List of String)

<a id="nestedblock--fargate_profile--selector"></a>
### Nested Schema for `fargate_profile.selector`

Required:

- **namespace** (String)

Optional:

- **labels** (Map of String)



<a id="nestedblock--pack"></a>
### Nested Schema for `pack`

Required:

- **name** (String)
- **tag** (String)
- **values** (String)


<a id="nestedblock--scan_policy"></a>
### Nested Schema for `scan_policy`

Required:

- **configuration_scan_schedule** (String)
- **conformance_scan_schedule** (String)
- **penetration_scan_schedule** (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)


