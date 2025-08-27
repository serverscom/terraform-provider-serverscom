---
page_title: "Servers.com: serverscom_l2_segment_members"
---

# serverscom_l2_segment_members

Get information on L2 segment members for use in other resources. This data source provides all of the L2 segment members and their properties.

## Example Usage

Get the L2 segment members by L2 segment ID:

```hcl
data "serverscom_l2_segment_members" "example" {
  id = "L2Seg7zQnVb"
}

output "l2_segment_members_example" {
  value = data.serverscom_l2_segment_members.example.members
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the L2 segment to get members for.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the L2 segment.
* `members` - A list of L2 segment members. Each member has the following attributes:
  * `id` - The ID of the member.
  * `title` - The title of the member.
  * `mode` - The mode of the member connection.
  * `vlan` - The VLAN ID assigned to the member.
  * `status` - The status of the member.
  * `labels` - A map of labels assigned to the member.
  * `created_at` - The creation time of the member.
  * `updated_at` - The last update time of the member.
