.. _resource_serverscom_l2_segment:

L2 Segment
==========

Provides an Servers.com l2 segment resource. This can be used to create, modify, and l2 segments.

Example
*******

Create a new L2 segment

.. sourcecode:: terraform

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

Argument Reference
******************

The following arguments are supported:

- `name` - (Optional, string) Name of the l2 segment.
- `type` - (Required, string) Type of the l2 segment.
- `location_group` - (Required, string) Location group code.
- `member` - (Required, list) List of the l2 segment members.
- `member.0.id` - (Required, int) ID of the dedicated server.
- `member.0.mode` - (Required, string) Membership mode of the dedicated server.

Attributes Reference
********************

The following attributes are exported:

- `id` - (string) Unique ID of the l2 segment.
- `name` - (string) Name of the l2 segment.
- `type` - (string) Type of the l2 segment.
- `location_group` - (string) Location group code.
- `member` - (list) List of the l2 segment members.
- `member.0.id` - (int) ID of the dedicated server.
- `member.0.mode` - (string) Membership mode of the dedicated server.
- `member.0.status` - (string) Status of the membership.
- `member.0.vlan` - (int) VLAN number of the member.
- `member.0.created_at` - (string) Member created at.
- `member.0.updated_at` - (string) Member updated at.
- `status` - (string) Status of the l2 segment.
- `created_at` - (string) L2 segment created at.
- `updated_at` - (string) L2 segment updated at.

Import
******

L2 Segments can be imported using the l2 segment `id`:

.. sourcecode:: bash

  terraform import serverscom_l2_segment.segment_1 <id>

.. vi: textwidth=78
