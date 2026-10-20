---
page_title: "Servers.com: serverscom_rbs_flavor_order_option"
---

# serverscom_rbs_flavor_order_option

Get information on a Remote Block Storage (RBS) flavor in a location. A flavor defines the IOPS and bandwidth per GB of volume size and the minimum volume size.

To list all flavors available in a location, use the `serverscom_rbs_flavor_order_options` data source.

## Example Usage

Get the RBS flavor by ID:

```hcl
data "serverscom_rbs_flavor_order_option" "example" {
  location_id = 1
  id          = 2
}

output "rbs_flavor_min_size" {
  value = data.serverscom_rbs_flavor_order_option.example.min_size_gb
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `id` - (Required, int) The ID of the RBS flavor. Use the `serverscom_rbs_flavor_order_options` data source to list IDs.

## Attributes Reference

The following attributes are exported:

- `id` - (int) The ID of the RBS flavor. Pass it as `flavor_id` to the `serverscom_rbs_volume` resource.
- `name` - (string) The flavor name.
- `iops_per_gb` - (number) The Input/Output Operations Per Second (IOPS) per 1 GB of volume size.
- `bandwidth_per_gb` - (number) Bandwidth per 1 GB of volume size (MB/s).
- `min_size_gb` - (int) The minimum volume size in GB.
