---
page_title: "Servers.com: serverscom_sbm_flavor_order_option"
---

# serverscom_sbm_flavor_order_option

Get information on a Scalable Bare Metal (SBM) flavor available in a location. An SBM flavor is a fixed configuration of CPU, RAM, drives, uplinks, and bandwidth. To list all flavors, use the `serverscom_sbm_flavor_order_options` data source.

The `name` attribute is the value of the `flavor` argument of the `serverscom_sbm_server` resource. The `id` is the `sbm_flavor_model_id` input of the `serverscom_sbm_operating_system_order_options` data source.

## Example Usage

```hcl
data "serverscom_sbm_flavor_order_option" "example" {
  location_id = 1
  id          = 10
}

output "sbm_flavor_name" {
  value = data.serverscom_sbm_flavor_order_option.example.name
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `id` - (Required, int) The ID of the SBM flavor. Use the `serverscom_sbm_flavor_order_options` data source to list IDs.

## Attributes Reference

The following attributes are exported:

- `id` - (int) The ID of the SBM flavor.
- `name` - (string) A human-readable name of the SBM flavor.
- `cpu_name` - (string) The name of the CPU.
- `cpu_count` - (int) The number of CPUs per server.
- `cpu_cores_count` - (int) The number of cores per CPU.
- `cpu_frequency` - (string) The CPU frequency.
- `ram_size` - (int) The RAM size in GB.
- `drives_configuration` - (string) The number of drives and their parameters.
- `public_uplink_model_id` - (int) The ID of the public uplink model.
- `public_uplink_model_name` - (string) The name of the public uplink model.
- `private_uplink_model_id` - (int) The ID of the private uplink model.
- `private_uplink_model_name` - (string) The name of the private uplink model.
- `bandwidth_id` - (int) The ID of the bandwidth option.
- `bandwidth_name` - (string) The name of the bandwidth option.
