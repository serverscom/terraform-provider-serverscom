---
page_title: "Servers.com: serverscom_uplink_model_order_option"
---

# serverscom_uplink_model_order_option

Get information on an uplink model available for a server model in a location. An uplink is a channel that carries your traffic in a private or public network. To list all uplink models, use the `serverscom_uplink_model_order_options` data source.

The `name` attribute is the value of the `private_uplink` and `public_uplink` arguments of the `serverscom_dedicated_server` resource. The `id` is the `uplink_model_id` input of the `serverscom_bandwidth_order_options` data source.

## Example Usage

```hcl
data "serverscom_uplink_model_order_option" "example" {
  location_id     = 1
  server_model_id = 100
  id              = 400
}

output "uplink_speed" {
  value = data.serverscom_uplink_model_order_option.example.speed
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `server_model_id` - (Required, int) The ID of the server model. Use the `serverscom_server_model_order_options` data source to find it.
- `id` - (Required, int) The ID of the uplink model. Use the `serverscom_uplink_model_order_options` data source to list IDs.

## Attributes Reference

The following attributes are exported:

- `id` - (int) The ID of the uplink model.
- `name` - (string) A human-readable name of the uplink model.
- `type` - (string) The uplink type: `private` or `public`.
- `speed` - (int) The maximum throughput of the uplink in Mbps.
- `redundancy` - (bool) Whether the uplink has redundancy.
