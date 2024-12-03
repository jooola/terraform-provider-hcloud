---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "hcloud_load_balancer_type Data Source - hcloud"
subcategory: ""
description: |-
  Provides details about a specific Hetzner Cloud Load Balancer Type.
  Use this resource to get detailed information about a specific Load Balancer Type.
---

# hcloud_load_balancer_type (Data Source)

Provides details about a specific Hetzner Cloud Load Balancer Type.

Use this resource to get detailed information about a specific Load Balancer Type.

## Example Usage

```terraform
data "hcloud_load_balancer_type" "by_id" {
  id = 1
}

data "hcloud_load_balancer_type" "by_name" {
  name = "lb11"
}

resource "hcloud_load_balancer" "main" {
  name               = "my-load-balancer"
  load_balancer_type = data.hcloud_load_balancer_type.name
  location           = "fsn1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (Number) ID of the Load Balancer Type.
- `name` (String) Name of the Load Balancer Type.

### Read-Only

- `description` (String) Description of the Load Balancer Type.
- `max_assigned_certificates` (Number) Maximum number of certificates that can be assigned for the Load Balancer of this type.
- `max_connections` (Number) Maximum number of simultaneous open connections for the Load Balancer of this type.
- `max_services` (Number) Maximum number of services for the Load Balancer of this type.
- `max_targets` (Number) Maximum number of targets for the Load Balancer of this type.