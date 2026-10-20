---
page_title: "Servers.com: serverscom_operating_system_order_option"
---

# serverscom_operating_system_order_option

Get information on an operating system available for a dedicated server model in a location. For Scalable Bare Metal (SBM) servers, use the `serverscom_sbm_operating_system_order_option` data source. To list all operating systems, use the `serverscom_operating_system_order_options` data source.

The `operating_system` argument of the `serverscom_dedicated_server` resource is matched, case-insensitively, against the OS name, version, and architecture joined with spaces, for example `Ubuntu 24.04-server x86_64`. The `filesystems` attribute lists the values accepted by the `fs` argument of a partition.

## Example Usage

```hcl
data "serverscom_operating_system_order_option" "example" {
  location_id     = 1
  server_model_id = 100
  id              = 300
}

output "os_filesystems" {
  value = data.serverscom_operating_system_order_option.example.filesystems
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `server_model_id` - (Required, int) The ID of the server model. Use the `serverscom_server_model_order_options` data source to find it.
- `id` - (Required, int) The ID of the operating system. Use the `serverscom_operating_system_order_options` data source to list IDs.

## Attributes Reference

The following attributes are exported:

- `id` - (int) The ID of the operating system.
- `full_name` - (string) The full name of the operating system: name, version, and architecture.
- `name` - (string) The name of the operating system.
- `version` - (string) The version of the operating system.
- `arch` - (string) The architecture of the operating system.
- `filesystems` - (list of strings) File systems supported by the operating system.
