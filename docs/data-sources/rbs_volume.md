---
page_title: "Servers.com: serverscom_rbs_volume"
---

# serverscom_rbs_volume

Get information on a Remote Block Storage (RBS) volume for use in other resources. This data source provides all of the RBS volume properties except the credentials, which are available in the `serverscom_rbs_volume_credentials` data source.

## Example Usage

Get the RBS volume by ID:

```hcl
data "serverscom_rbs_volume" "example" {
  id = "RBSvol7zQnVb"
}

output "rbs_volume_ip" {
  value = data.serverscom_rbs_volume.example.ip_address
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Required, string) The ID of the RBS volume.

## Attributes Reference

The following attributes are exported:

- `id` - (string) The ID of the RBS volume.
- `name` - (string) A name of the volume.
- `size` - (int) Volume size in GB.
- `status` - (string) Status of the volume. One of: `creating`, `pending`, `active`, `removing`.
- `labels` - (map) A map of labels assigned to the volume.
- `location_id` - (int) An ID of the volume location.
- `location_code` - (string) A location code of the volume. For example: `AMS1`.
- `ip_address` - (string) A private IP address of the volume (Volume IP address).
- `flavor_id` - (int) An ID of the volume flavor.
- `flavor_name` - (string) A name of the volume flavor.
- `iops` - (number) The maximum IOPS of the volume.
- `bandwidth` - (number) The maximum bandwidth of the volume in MB/s.
- `target_iqn` - (string) iSCSI Qualified Name (IQN) of the volume.
- `created_at` - (string) The creation time of the volume.
- `updated_at` - (string) The last update time of the volume.
