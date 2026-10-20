---
page_title: "Servers.com: serverscom_sbm_operating_system_order_options"
---

# serverscom_sbm_operating_system_order_options

Lists operating systems available for a Scalable Bare Metal (SBM) flavor in a location. Windows operating systems are not available for SBM servers.

## Example Usage

```hcl
data "serverscom_sbm_operating_system_order_options" "example" {
  location_id         = 1
  sbm_flavor_model_id = 10
}

output "sbm_operating_systems" {
  value = [for os in data.serverscom_sbm_operating_system_order_options.example.sbm_operating_systems : os.full_name]
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `sbm_flavor_model_id` - (Required, int) The ID of the SBM flavor. Use the `serverscom_sbm_flavor_order_options` data source to find it.

## Attributes Reference

The following attributes are exported:

- `sbm_operating_systems` - A list of operating systems. Each element contains:
  - `id` - The ID of the operating system.
  - `full_name` - The full name of the operating system: name, version, and architecture.
  - `name` - The name of the operating system.
  - `version` - The version of the operating system.
  - `arch` - The architecture of the operating system.
  - `filesystems` - File systems supported by the operating system.
