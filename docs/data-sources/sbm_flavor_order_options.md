---
page_title: "Servers.com: serverscom_sbm_flavor_order_options"
---

# serverscom_sbm_flavor_order_options

Lists Scalable Bare Metal (SBM) flavors available in a location, with optional filtering. Use it to find the name for the `flavor` argument of the `serverscom_sbm_server` resource.

## Example Usage

```hcl
data "serverscom_sbm_flavor_order_options" "example" {
  location_id = 1
}

output "sbm_flavors" {
  value = {
    for f in data.serverscom_sbm_flavor_order_options.example.sbm_flavors :
    f.name => f.id
  }
}
```

Search flavors by name and include unavailable ones:

```hcl
data "serverscom_sbm_flavor_order_options" "all" {
  location_id = 1

  filter {
    search_pattern = "SBM"
    show_all       = true
  }
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `filter` - (Optional, block) A block to filter the returned flavors. Supports:
  - `search_pattern` - (Optional, string) Substring match on the flavor name.
  - `show_all` - (Optional, bool) By default, only available flavors are returned. Set to `true` to include unavailable flavors.

## Attributes Reference

The following attributes are exported:

- `sbm_flavors` - A list of SBM flavors. Each element contains:
  - `id` - The ID of the SBM flavor.
  - `name` - A human-readable name of the SBM flavor.
  - `cpu_name` - The name of the CPU.
  - `cpu_count` - The number of CPUs per server.
  - `cpu_cores_count` - The number of cores per CPU.
  - `cpu_frequency` - The CPU frequency.
  - `ram_size` - The RAM size in GB.
  - `drives_configuration` - The number of drives and their parameters.
  - `public_uplink_model_id` - The ID of the public uplink model.
  - `public_uplink_model_name` - The name of the public uplink model.
  - `private_uplink_model_id` - The ID of the private uplink model.
  - `private_uplink_model_name` - The name of the private uplink model.
  - `bandwidth_id` - The ID of the bandwidth option.
  - `bandwidth_name` - The name of the bandwidth option.
