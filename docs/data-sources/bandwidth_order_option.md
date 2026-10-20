---
page_title: "Servers.com: serverscom_bandwidth_order_option"
---

# serverscom_bandwidth_order_option

Get information on a bandwidth option (a traffic billing policy) available for an uplink model in a location. To list all bandwidth options, use the `serverscom_bandwidth_order_options` data source.

The `name` attribute is the value of the `bandwidth` argument of the `serverscom_dedicated_server` resource. Bandwidth is required when a public uplink is used.

## Example Usage

```hcl
data "serverscom_bandwidth_order_option" "example" {
  location_id     = 1
  server_model_id = 100
  uplink_model_id = 400
  id              = 500
}

output "bandwidth_name" {
  value = data.serverscom_bandwidth_order_option.example.name
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `server_model_id` - (Required, int) The ID of the server model. Use the `serverscom_server_model_order_options` data source to find it.
- `uplink_model_id` - (Required, int) The ID of the public uplink model. Use the `serverscom_uplink_model_order_options` data source to find it.
- `id` - (Required, int) The ID of the bandwidth option. Use the `serverscom_bandwidth_order_options` data source to list IDs.

## Attributes Reference

The following attributes are exported:

- `id` - (int) The ID of the bandwidth option.
- `name` - (string) A human-readable name of the bandwidth option.
- `type` - (string) The type of the bandwidth option: `bytes` (a quantity of included traffic), `bandwidth` (throughput), or `unmetered` (no limits on traffic).
- `commit` - (int) The number of units included in the commit. It is `0` for an unmetered option.
