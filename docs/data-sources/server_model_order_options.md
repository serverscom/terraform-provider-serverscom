---
page_title: "Servers.com: serverscom_server_model_order_options"
---

# serverscom_server_model_order_options

Lists dedicated server models (chassis) available for ordering in a location, with optional filtering. Use it to find the server model name for the `serverscom_dedicated_server` resource or the server model ID for other order options data sources.

## Example Usage

List all server models in a location:

```hcl
data "serverscom_server_model_order_options" "all" {
  location_id = 1
}

output "server_models" {
  value = {
    for m in data.serverscom_server_model_order_options.all.server_models :
    m.name => m.id
  }
}
```

List only models with a RAID controller:

```hcl
data "serverscom_server_model_order_options" "raid" {
  location_id = 1

  filter {
    search_pattern      = "Dell"
    has_raid_controller = true
  }
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `filter` - (Optional, block) A block to filter the returned models. Supports:
  - `search_pattern` - (Optional, string) Substring match on the server model name.
  - `has_raid_controller` - (Optional, bool) Return only models with (`true`) or without (`false`) a RAID controller.

~> **Note:** Terraform cannot tell an omitted boolean from `false`. When a `filter` block is present, `has_raid_controller` is always sent to the API, so set it explicitly to the value you need.

## Attributes Reference

The following attributes are exported:

- `server_models` - A list of server models. Each element contains:
  - `id` - The ID of the server model.
  - `name` - The name of the server model.
  - `cpu_name` - The name of the CPU model.
  - `cpu_count` - The number of CPUs.
  - `cpu_cores_count` - The number of physical cores per CPU.
  - `cpu_frequency` - The CPU frequency in MHz.
  - `ram` - The default RAM size in GB.
  - `ram_type` - The DDR type of RAM.
  - `max_ram` - The maximum RAM size in GB.
  - `has_raid_controller` - Whether the server model has a RAID controller.
  - `raid_controller_name` - The name of the RAID controller.
  - `drive_slots_count` - The number of drive slots.

The list does not include drive slot details.
