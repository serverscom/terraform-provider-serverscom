---
page_title: "Servers.com: serverscom_l2_segment"
---

# serverscom_l2_segment

Provides an Servers.com l2 segment resource. This can be used to create, modify, and delete L2 segments.

## Example Usage

Create a new L2 segment

```hcl
resource "serverscom_l2_segment" "segment_1" {
  name = "l2-segment-1"
  type = "private"
  location_group = "SJC1"

  member {
    id = "QBeXDWey"
    mode = "native"
  }

  member {
    id = "4QbYEKbz"
    mode = "native"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional, string) Name of the L2 segment.
- `type` - (Required, string) Type of the L2 segment.
- `location_group` - (Required, string) Location group code.
- `member` - (Required, list) List of the L2 segment members.
- `member.0.id` - (Required, string) ID of the dedicated server.
- `member.0.mode` - (Required, string) Membership mode of the dedicated server.
- `labels` - (Optional, map) A map of labels assigned to the L2 segment.

## Attributes Reference

The following attributes are exported:

- `id` - (string) Unique identifier of the L2 segment.
- `name` - (string) Name of the L2 segment.
- `type` - (string) Type of the L2 segment.
- `location_group` - (string) Location group code.
- `member` - (list) List of the L2 segment members.
- `member.0.id` - (string) ID of the dedicated server.
- `member.0.mode` - (string) Membership mode of the dedicated server.
- `member.0.status` - (string) Status of the membership.
- `member.0.vlan` - (string) VLAN number of the member.
- `member.0.created_at` - (string) Member created at.
- `member.0.updated_at` - (string) Member updated at.
- `status` - (string) Status of the L2 segment.
- `created_at` - (string) L2 segment created at.
- `updated_at` - (string) L2 segment updated at.
- `labels` - (map) A map of labels assigned to the L2 segment.

## Import

L2 Segments can be imported using the l2 segment `id`:

```bash
terraform import serverscom_l2_segment.segment_1 <id>
```
