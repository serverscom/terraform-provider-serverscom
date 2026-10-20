---
page_title: "Servers.com: serverscom_sbm_operating_system_order_option"
---

# serverscom_sbm_operating_system_order_option

Get information on an operating system available for a Scalable Bare Metal (SBM) flavor in a location. Windows operating systems are not available for SBM servers. To list all operating systems, use the `serverscom_sbm_operating_system_order_options` data source.

The `operating_system` argument of the `serverscom_sbm_server` resource is matched, case-insensitively, against the OS name, version, and architecture joined with spaces, for example `Ubuntu 24.04-server x86_64`.

## Example Usage

```hcl
data "serverscom_sbm_operating_system_order_option" "example" {
  location_id         = 1
  sbm_flavor_model_id = 10
  id                  = 300
}

output "sbm_os_name" {
  value = data.serverscom_sbm_operating_system_order_option.example.full_name
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `sbm_flavor_model_id` - (Required, int) The ID of the SBM flavor. Use the `serverscom_sbm_flavor_order_options` data source to find it.
- `id` - (Required, int) The ID of the operating system. Use the `serverscom_sbm_operating_system_order_options` data source to list IDs.

## Attributes Reference

The following attributes are exported:

- `id` - (int) The ID of the operating system.
- `full_name` - (string) The full name of the operating system: name, version, and architecture.
- `name` - (string) The name of the operating system.
- `version` - (string) The version of the operating system.
- `arch` - (string) The architecture of the operating system.
- `filesystems` - (list of strings) File systems supported by the operating system.
