---
page_title: "Servers.com: serverscom_locations"
---

# serverscom_locations

Lists Servers.com locations (data centers) with optional filtering. Use it to find a location ID by its code or name.

## Example Usage

List all locations:

```hcl
data "serverscom_locations" "all" {}

output "location_codes" {
  value = [for l in data.serverscom_locations.all.locations : l.code]
}
```

Find the ID of a location by its code:

```hcl
locals {
  ams1_id = one([
    for l in data.serverscom_locations.all.locations : l.id if l.code == "AMS1"
  ])
}
```

Filter locations by name and check which ones support iPXE Boot:

```hcl
data "serverscom_locations" "amsterdam" {
  filter {
    search_pattern = "Amsterdam"
  }
}

output "ipxe_locations" {
  value = [
    for l in data.serverscom_locations.amsterdam.locations : l.code
    if contains(l.supported_features, "public_ipxe_boot")
  ]
}
```

## Argument Reference

The following arguments are supported:

- `filter` - (Optional, block) A block to filter the returned locations. Supports:
  - `search_pattern` - (Optional, string) Substring match on the location name.

## Attributes Reference

The following attributes are exported:

- `locations` - A list of locations. Each element contains:
  - `id` - The ID of the location.
  - `name` - A human-readable name of the location.
  - `code` - A technical identifier of the location. For example: `AMS1`.
  - `status` - The status of the location.
  - `supported_features` - A list of features supported in the location. See the `serverscom_location` data source for the possible values.
  - `l2_segments_enabled` - Whether L2 Segments are available in the location.
  - `private_racks_enabled` - Whether Private Racks are available in the location.
  - `load_balancers_enabled` - Whether the Load Balancer service is available in the location.
