---
page_title: "Servers.com: serverscom_l2_segment"
---

# serverscom_l2_segment

Get information on an L2 segment for use in other resources. This data source provides all of the L2 segment properties.

## Example Usage

Get the L2 segment by ID:

```hcl
data "serverscom_l2_segment" "example" {
  id = "L2Seg7zQnVb"
}

output "l2_segment_example" {
  value = data.serverscom_l2_segment.example.name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the L2 segment.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the L2 segment.
* `name` - The name of the L2 segment.
* `type` - The type of the L2 segment.
* `status` - The status of the L2 segment.
* `location_group_id` - The location group ID where the L2 segment is available.
* `location_group_code` - The location group code where the L2 segment is available.
* `labels` - A map of labels assigned to the L2 segment.
* `created_at` - The creation time of the L2 segment.
* `updated_at` - The last update time of the L2 segment.
