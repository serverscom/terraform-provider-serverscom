---
page_title: "Servers.com: serverscom_uplink_model_order_options"
---

# serverscom_uplink_model_order_options

Lists uplink models available for a server model in a location, with optional filtering. Use it to find the names for the `private_uplink` and `public_uplink` arguments of the `serverscom_dedicated_server` resource, and the uplink IDs for the `serverscom_bandwidth_order_options` data source.

## Example Usage

```hcl
data "serverscom_uplink_model_order_options" "public" {
  location_id     = 1
  server_model_id = 100

  filter {
    type       = "public"
    redundancy = true
  }
}

output "public_uplinks" {
  value = [for u in data.serverscom_uplink_model_order_options.public.uplink_models : u.name]
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `server_model_id` - (Required, int) The ID of the server model. Use the `serverscom_server_model_order_options` data source to find it.
- `filter` - (Optional, block) A block to filter the returned uplink models. Supports:
  - `type` - (Optional, string) The uplink type: `private` or `public`.
  - `redundancy` - (Optional, bool) Return only uplinks with (`true`) or without (`false`) redundancy.

~> **Note:** Terraform cannot tell an omitted boolean from `false`. When a `filter` block is present, `redundancy` is always sent to the API, so set it explicitly to the value you need.

## Attributes Reference

The following attributes are exported:

- `uplink_models` - A list of uplink models. Each element contains:
  - `id` - The ID of the uplink model.
  - `name` - A human-readable name of the uplink model.
  - `type` - The uplink type: `private` or `public`.
  - `speed` - The maximum throughput of the uplink in Mbps.
  - `redundancy` - Whether the uplink has redundancy.
