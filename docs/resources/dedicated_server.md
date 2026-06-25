---
page_title: "Servers.com: serverscom_dedicated_server"
---

# serverscom_dedicated_server

Provides a Servers.com dedicated server resource. This can be used to create, modify, and delete servers. Servers also support [provisioning](https://www.terraform.io/docs/provisioners/index.html).

## Example Usage

Create a new dedicated server:

```hcl
resource "serverscom_dedicated_server" "node_1" {
  hostname             = "node-1"
  location             = "SJC1"
  server_model         = "Dell R440 / 2xIntel Xeon Silver-4114 / 32 GB RAM / 1x480 GB SSD"
  ram_size             = 32

  operating_system     = "Ubuntu 24.04-server x86_64"

  private_uplink       = "Private 10 Gbps with redundancy"
  public_uplink        = "Public 10 Gbps with redundancy"
  bandwidth            = "20002 Gb"

  ssh_key_fingerprints = [
    "cf:1d:09:ab:cb:47:97:3f:50:9a:f0:34:14:78:fa:1b"
  ]

  ipv6 = true

  slot {
    drive_model = "480 GB SSD SATA"
    position    = 0
  }

  layout {
    slot_positions = [0]

    partition {
      target = "/"
      size = 10240
      fill = false
      fs = "ext4"
    }

    partition {
      target = "/home"
      size = 1
      fill = true
      fs = "ext4"
    }

    partition {
      target = "swap"
      size = 4096
      fill = false
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `hostname` - (Required, string) Name of the dedicated server (according to RFC 1123 specification).
- `location` - (Required, string) Location code of the dedicated server. For example: `AMS1`, `SJC1`, etc.
- `server_model` - (Required, string) Name of the dedicated server model.
- `ram_size` - (Optional, int) Size of the RAM (GB).
- `operating_system` - (Optional, string) The dedicated server operating system name.
- `private_uplink` - (Required, string) The dedicated server private uplink name.
- `public_uplink` - (Optional, string) The dedicated server public uplink name.
- `bandwidth` - (Optional, string) The dedicated server public bandwidth name.
- `ssh_key_fingerprints` - (Optional, list) SSH key fingerprint.
- `private_ipv4_network_id` - (Optional, string) Private IPv4 network ID.
- `public_ipv4_network_id` - (Optional, string) Public IPv4 network ID.
- `user_data` - (Optional, string, sensitive) A string of the desired user data for the dedicated server. It is applied in-place and is also sent with every OS reinstall. Changing it no longer recreates the server. Because the raw value is needed to update it in place, it is stored in the Terraform state (mark it carefully / use a remote backend if it contains secrets).
- `reinstall_trigger` - (Optional, string) An opt-in marker that performs an OS reinstall. Defaults to `none`. Changing it to any other value reinstalls the operating system, which **destroys all data on disk**. Changing `operating_system` or `layout` requires also changing `reinstall_trigger` in the same plan, otherwise the plan fails. See [Reinstalling the OS](#reinstalling-the-os).
- `ipv6` - (Optional, bool) Is IPv6 enabled. Defaults to `false`.
- `slot` - (Optional, list) List of drive slots. Slots used in partioning have to be listed.
- `slot.0.position` - (Required, int) Slot position.
- `slot.0.drive_model` - (Optional, string) The name of drive model to place in the slot.
- `layout` - (Optional, list) List of layouts.
- `layout.0.slot_positions` - (Required, list) List of slots which should be used in the layout.
- `layout.0.raid` - (Optional, int) RAID level for the layout.
- `layout.0.partition` - (Required, list) List of partitions for the layout.
- `layout.0.partition.0.target` - (Required, string) Target/Mount point for the partition.
- `layout.0.partition.0.size` - (Required, int) Size of the partition (MB).
- `layout.0.partition.0.fill` - (Optional, bool) Autofill partition by all unused space. When set to `true`, the partition will use all remaining available space.
- `layout.0.partition.0.fs` - (Optional, string) Filesystem type for the partition.
- `labels` - (Optional, map) A map of labels assigned to the dedicated server.

## Attributes Reference

The following attributes are exported:

- `id` - (string) Unique identifier of the dedicated server.
- `hostname` - (string) Name of the dedicated server.
- `location` - (string) Location code of the dedicated server.
- `server_model` - (string) Name of the dedicated server model.
- `configuration` - (string) Configuration description of the dedicated server.
- `private_ipv4_address` - (string) Private IPv4 address.
- `public_ipv4_address` - (string) Public IPv4 address.
- `status` - (string) Status of the dedicated server.
- `operational_status` - (string) Operational status of the dedicated server. After a reinstall the provider waits for this to return to `normal`.
- `labels` - (map) A map of labels assigned to the dedicated server.

## Reinstalling the OS

A reinstall reprovisions the operating system and **destroys all data on disk**, so it never happens implicitly. It is gated by the `reinstall_trigger` argument:

- `reinstall_trigger` defaults to `none`. While it stays `none`, no reinstall is performed.
- Changing `reinstall_trigger` to any other value (e.g. `"1"`, `"2024-06-redeploy"`) performs a reinstall using the current `operating_system`, `layout`, `ssh_key_fingerprints` and `user_data`. This also lets you re-run a reinstall without changing anything else — just bump the value again.
- `operating_system` and `layout` can only be applied through a reinstall. Changing either of them **without** also changing `reinstall_trigger` fails the plan with an explanatory error. Change `reinstall_trigger` in the same plan to acknowledge and apply.
- `user_data` is applied in place (without a reinstall) and is additionally included in the reinstall payload.
- `ssh_key_fingerprints` is sent with the reinstall payload (so a reinstall provisions the freshly installed OS with the configured keys). Changing it on its own does not trigger a reinstall and only takes effect on the next reinstall.

First-time adoption is safe: simply adding `reinstall_trigger` to an existing resource (its default `none`) does not trigger a reinstall.

The provider waits for `operational_status` to become `normal` after a reinstall. It also waits for `normal` before applying any change, because the API rejects reinstall/update requests while a previous operation is still running. The wait is bounded by the `update` timeout (default `4h`):

```hcl
resource "serverscom_dedicated_server" "node_1" {
  # ...
  operating_system  = "Ubuntu 24.04-server x86_64"
  reinstall_trigger = "1"

  timeouts {
    update = "2h"
  }
}
```

## Import

Dedicated servers can be imported using the dedicated server `id`:

```bash
terraform import serverscom_dedicated_server.node_1 <id>
```
