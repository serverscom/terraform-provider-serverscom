---
page_title: "Servers.com: serverscom_operating_system_order_options"
---

# serverscom_operating_system_order_options

Lists operating systems available for a dedicated server model in a location. For Scalable Bare Metal (SBM) servers, use the `serverscom_sbm_operating_system_order_options` data source.

## Example Usage

List operating systems and pick the value for the `operating_system` argument of a server:

```hcl
data "serverscom_operating_system_order_options" "example" {
  location_id     = 1
  server_model_id = 100
}

output "operating_systems" {
  value = [for os in data.serverscom_operating_system_order_options.example.operating_systems : os.full_name]
}
```

Find Ubuntu images:

```hcl
output "ubuntu" {
  value = [
    for os in data.serverscom_operating_system_order_options.example.operating_systems :
    os.full_name if os.name == "Ubuntu"
  ]
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `server_model_id` - (Required, int) The ID of the server model. Use the `serverscom_server_model_order_options` data source to find it.

## Attributes Reference

The following attributes are exported:

- `operating_systems` - A list of operating systems. Each element contains:
  - `id` - The ID of the operating system.
  - `full_name` - The full name of the operating system: name, version, and architecture.
  - `name` - The name of the operating system.
  - `version` - The version of the operating system.
  - `arch` - The architecture of the operating system.
  - `filesystems` - File systems supported by the operating system.
