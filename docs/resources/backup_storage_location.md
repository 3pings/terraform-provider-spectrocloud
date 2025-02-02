---
page_title: "spectrocloud_backup_storage_location Resource - terraform-provider-spectrocloud"
subcategory: ""
description: |-
  
---

# Resource `spectrocloud_backup_storage_location`



## Example Usage

```terraform
resource "spectrocloud_backup_storage_location" "bsl1" {
  name        = "dev-backup-s3"
  is_default  = false
  region      = "us-east-2"
  bucket_name = "dev-backup"
  s3 {
    credential_type     = var.credential_type
    access_key          = var.aws_access_key
    secret_key          = var.aws_secret_key
    s3_force_path_style = false

    #s3_url             = "http://10.90.78.23"
  }
}

resource "spectrocloud_backup_storage_location" "bsl2" {
  name        = "prod-backup-s3"
  is_default  = false
  region      = "us-east-2"
  bucket_name = "prod-backup"
  s3 {
    credential_type     = var.credential_type
    arn                 = var.arn
    external_id         = var.external_id
    s3_force_path_style = false
    #s3_url             = "http://10.90.78.23"
  }
}
```

## Schema

### Required

- **bucket_name** (String)
- **is_default** (Boolean)
- **name** (String)
- **region** (String)
- **s3** (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--s3))

### Optional

- **ca_cert** (String)
- **id** (String) The ID of this resource.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

<a id="nestedblock--s3"></a>
### Nested Schema for `s3`

Required:

- **credential_type** (String)

Optional:

- **access_key** (String)
- **arn** (String)
- **external_id** (String)
- **s3_force_path_style** (Boolean)
- **s3_url** (String)
- **secret_key** (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)


