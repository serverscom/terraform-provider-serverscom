---
page_title: "Servers.com: serverscom_server_model_order_option"
---

# serverscom_server_model_order_option

Get information on a dedicated server model (chassis). The data source exposes CPU, RAM, RAID controller, and drive slot attributes. To list the server models available in a location, use the `serverscom_server_model_order_options` data source.

The `name` attribute is the value of the `server_model` argument of the `serverscom_dedicated_server` resource. The `id` is the `server_model_id` input of the drive model, operating system, RAM, uplink model, and bandwidth order options data sources.

## Example Usage

```hcl
data "serverscom_server_model_order_option" "example" {
  id              = 100
  server_model_id = 100
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Required, int) The ID of the server model.
- `server_model_id` - (Required, int) The ID of the server model.

## Attributes Reference

The following attributes are exported:

- `id` - (int) The ID of the server model.
- `name` - (string) The name of the server model.
- `cpu_name` - (string) The name of the CPU model.
- `cpu_count` - (int) The number of CPUs.
- `cpu_cores_count` - (int) The number of physical cores per CPU.
- `cpu_frequency` - (int) The CPU frequency in MHz.
- `ram` - (int) The default RAM size in GB.
- `ram_type` - (string) The DDR type of RAM.
- `max_ram` - (int) The maximum RAM size available for the server model in GB.
- `has_raid_controller` - (bool) Whether the server model has a RAID controller.
- `raid_controller_name` - (string) The name of the RAID controller. Empty if there is no RAID controller.
- `drive_slots_count` - (int) The number of drive slots.
- `drive_slots` - (list) Drive slots of the server model. Each element contains:
  - `position` - (int) The slot position, starting from `0`. This is the `position` of the `slot` block in the `serverscom_dedicated_server` resource.
  - `interface` - (string) The slot interface.
  - `form_factor` - (string) The slot form factor.
  - `drive_model_id` - (int) The ID of the drive model in the slot.
  - `hot_swappable` - (bool) Whether the drive in the slot is hot swappable.
