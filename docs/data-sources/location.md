---
page_title: "Servers.com: serverscom_location"
---

# serverscom_location

Get information on a Servers.com location (data center), including the features and services it supports. Location IDs are used by the order options data sources and by resources such as the `serverscom_rbs_volume` resource. To find an ID by a location code or name, use the `serverscom_locations` data source.

## Example Usage

Get the location by ID:

```hcl
data "serverscom_location" "example" {
  id = 1
}

output "location_code" {
  value = data.serverscom_location.example.code
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Required, int) The ID of the location.

## Attributes Reference

The following attributes are exported:

- `id` - (int) The ID of the location.
- `name` - (string) A human-readable name of the location.
- `code` - (string) A technical identifier of the location. For example: `AMS1`, `SJC1`. This is the value used for the `location` argument of server resources.
- `status` - (string) The status of the location: being configured, active, or disabled.
- `supported_features` - (list of strings) Features supported in the location: `disaggregated_public_ports`, `disaggregated_private_ports`, `no_public_ip_address`, `no_private_ip`, `host_rescue_mode`, `oob_public_access`, `no_public_network` (deprecated), `private_ipxe_boot`, `public_ipxe_boot`.
- `l2_segments_enabled` - (bool) Whether L2 Segments are available in the location.
- `private_racks_enabled` - (bool) Whether Private Racks are available in the location.
- `load_balancers_enabled` - (bool) Whether the Load Balancer service is available in the location.
