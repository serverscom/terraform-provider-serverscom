---
page_title: "Servers.com: serverscom_rbs_volumes"
---

# serverscom_rbs_volumes

Lists Remote Block Storage (RBS) volumes with optional filtering.

## Example Usage

List all volumes:

```hcl
data "serverscom_rbs_volumes" "all" {}
```

List volumes in a location with a specific label:

```hcl
data "serverscom_rbs_volumes" "production" {
  filter {
    location_id    = 1
    label_selector = "env=production"
  }
}

output "production_volume_ids" {
  value = [for v in data.serverscom_rbs_volumes.production.rbs_volumes : v.id]
}
```

## Argument Reference

The following arguments are supported:

- `filter` - (Optional, block) A block to filter the returned volumes. Only one `filter` block is used. Supports:
  - `location_id` - (Optional, int) Return only volumes from the location with this ID.
  - `label_selector` - (Optional, string) Label selector, e.g. `env=production`.
  - `search_pattern` - (Optional, string) Substring match on the volume name.

## Attributes Reference

The following attributes are exported:

- `rbs_volumes` - A list of RBS volumes. Each element contains:
  - `id` - The ID of the volume.
  - `name` - A name of the volume.
  - `size` - Volume size in GB.
  - `status` - Status of the volume.
  - `labels` - A map of labels assigned to the volume.
  - `location_id` - An ID of the volume location.
  - `location_code` - A location code of the volume.
  - `ip_address` - A private IP address of the volume.
  - `flavor_id` - An ID of the volume flavor.
  - `flavor_name` - A name of the volume flavor.
  - `iops` - The maximum IOPS of the volume.
  - `bandwidth` - The maximum bandwidth of the volume in MB/s.
  - `target_iqn` - iSCSI Qualified Name (IQN) of the volume.
  - `created_at` - The creation time of the volume.
  - `updated_at` - The last update time of the volume.
