---
page_title: "Servers.com: serverscom_rbs_flavor_order_options"
---

# serverscom_rbs_flavor_order_options

Lists Remote Block Storage (RBS) flavors available in a location. Use the flavor `id` as the `flavor_id` argument of the `serverscom_rbs_volume` resource.

## Example Usage

```hcl
data "serverscom_rbs_flavor_order_options" "example" {
  location_id = 1
}

output "rbs_flavors" {
  value = {
    for f in data.serverscom_rbs_flavor_order_options.example.rbs_flavors :
    f.name => f.id
  }
}
```

Pick a flavor by name and use it for a volume:

```hcl
locals {
  rbs_flavor = one([
    for f in data.serverscom_rbs_flavor_order_options.example.rbs_flavors : f
    if f.name == "<flavor name>"
  ])
}

resource "serverscom_rbs_volume" "data" {
  name        = "data-volume"
  size        = max(100, local.rbs_flavor.min_size_gb)
  location_id = 1
  flavor_id   = local.rbs_flavor.id
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.

## Attributes Reference

The following attributes are exported:

- `rbs_flavors` - A list of RBS flavors. Each element contains:
  - `id` - The ID of the flavor.
  - `name` - The flavor name.
  - `iops_per_gb` - The IOPS per 1 GB of volume size.
  - `bandwidth_per_gb` - Bandwidth per 1 GB of volume size (MB/s).
  - `min_size_gb` - The minimum volume size in GB.
