---
page_title: "Servers.com: serverscom_rbs_volume"
---

# serverscom_rbs_volume

Provides a Servers.com Remote Block Storage (RBS) volume resource. This can be used to create, resize, rename, label, and delete RBS volumes. Learn more about the service [in the knowledge base](https://www.servers.com/docs/rbs/concepts).

RBS volumes are exposed to your hosts over iSCSI. The provider creates and manages the volume itself; connecting it to a Linux, Windows, or Kubernetes host is done on the host side, see the Connecting a volume section below.

## Example Usage

Create a new RBS volume. Use the the `serverscom_rbs_flavor_order_options` data source to find the flavor ID available in a location:

```hcl
data "serverscom_locations" "all" {}

data "serverscom_rbs_flavor_order_options" "ams1" {
  location_id = local.location_id
}

locals {
  location_id = one([for l in data.serverscom_locations.all.locations : l.id if l.code == "AMS1"])
}

resource "serverscom_rbs_volume" "data" {
  name        = "data-volume"
  size        = 100
  location_id = local.location_id
  flavor_id   = data.serverscom_rbs_flavor_order_options.ams1.rbs_flavors[0].id

  labels = {
    env = "production"
  }
}
```

Increase the size of an existing volume by raising `size` and applying the configuration again. The volume goes through the `pending` status and the provider waits until it is `active` again.

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) A name of the volume. From 1 to 255 characters.
- `size` - (Required, int) Volume size in GB. The size can only be increased; decreasing it is rejected by the API. The minimum size depends on the flavor (see `min_size_gb` in the `serverscom_rbs_flavor_order_options` data source).
- `location_id` - (Required, int) An ID of the location where the volume is created. Use the the `serverscom_locations` data source to find it. Changing this forces a new resource to be created.
- `flavor_id` - (Required, int) An ID of the RBS flavor. A flavor defines IOPS and bandwidth per GB of volume size. Use the the `serverscom_rbs_flavor_order_options` data source to list flavors available in the location.
- `labels` - (Optional, map) A map of labels assigned to the volume.
- `timeouts` - (Optional, block) Customizes how long the provider waits for an operation:
  - `create` - (Optional, string) Time to wait for the volume to become `active`. Defaults to `1h`.
  - `delete` - (Optional, string) Time to wait for the volume removal. Defaults to `30m`.

~> **Note:** Changing `flavor_id` of an existing volume is not supported. The change is not applied, the volume keeps its flavor, and Terraform keeps showing the difference in the plan. To use another flavor, create a new volume with the desired flavor and migrate your data.

~> **Note:** Service limits apply: a volume can be mounted to one server at a time, the maximum volume size is 2 TB (up to 10 TB upon request), up to 99 volumes and 10 TB of total volume size per account. The billing period is a calendar month; a volume canceled before the period is over is charged for the full month. See [Limitations](https://www.servers.com/docs/rbs/limitations) and [Billing](https://www.servers.com/docs/rbs/billing).

## Attributes Reference

The following attributes are exported:

- `id` - (string) Unique identifier of the volume.
- `name` - (string) A name of the volume.
- `size` - (int) Volume size in GB.
- `location_id` - (int) An ID of the volume location.
- `location_code` - (string) A location code of the volume. For example: `AMS1`.
- `flavor_id` - (int) An ID of the volume flavor.
- `flavor_name` - (string) A name of the volume flavor.
- `labels` - (map) A map of labels assigned to the volume.
- `status` - (string) Status of the volume. One of: `creating`, `pending`, `active`, `removing`. `pending` means that an operation, for example a resize, is in progress.
- `ip_address` - (string) A private IP address of the volume (Volume IP address). Use it as the iSCSI portal address.
- `target_iqn` - (string) iSCSI Qualified Name (IQN) of the volume. Use it as the iSCSI target name.
- `iops` - (number) The maximum IOPS of the volume, calculated from the flavor and the volume size.
- `bandwidth` - (number) The maximum bandwidth of the volume in MB/s, calculated from the flavor and the volume size.
- `created_at` - (string) The creation time of the volume (RFC 3339).
- `updated_at` - (string) The last update time of the volume (RFC 3339).

## Connecting a volume

To attach the volume to a host you need the volume IP address, the target IQN, and the CHAP credentials. The IP address and the IQN are exported by this resource; the credentials are available through the the `serverscom_rbs_volume_credentials` data source:

```hcl
data "serverscom_rbs_volume_credentials" "data" {
  volume_id = serverscom_rbs_volume.data.id
}

output "iscsi_target" {
  value = {
    portal = "${serverscom_rbs_volume.data.ip_address}:3260"
    iqn    = serverscom_rbs_volume.data.target_iqn
  }
}
```

The `iscsiadm`, Windows, and Kubernetes steps are described in the [Getting started](https://www.servers.com/docs/rbs/getting-started) article of the knowledge base. For best performance, use 4 KB block sizes.

Before deleting a volume, disconnect it from the host (unmount, then log out from the iSCSI target) to avoid data loss.

## Import

RBS volumes can be imported using the volume `id`:

```bash
terraform import serverscom_rbs_volume.data <id>
```
