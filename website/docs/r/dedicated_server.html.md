---
layout: "serverscom"
page_title: "Servers.com: serverscom_dedicated_server"
sidebar_current: "docs-serverscom-resource-dedicated-server"
description: |-
  Provides an Servers.com dedicated server resource. This can be used to create, modify, and delete servers. Servers also support provisioning.
---

# serverscom_dedicated_server

Provides an Servers.com dedicated server resource. This can be used to create, modify, and delete servers. Servers also support [provisioning](https://www.terraform.io/docs/provisioners/index.html).

## Example Usage

```hcl
# Create a new dedicated server running ubuntu
resource "serverscom_dedicated_server" "node_1" {
  hostname             = "node-1"
  location             = "SJC1"
  server_model         = "Dell R440 / 2xIntel Xeon Silver-4114 / 32 GB RAM / 1x480 GB SSD"
  ram_size             = 32

  operating_system     = "Ubuntu 16.04-server x86_64"

  private_uplink       = "Private 10 Gbps with redundancy"
  public_uplink        = "Public 10 Gbps with redundancy"
  bandwidth            = "200000 GB"

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

- `hostname` - (Required, string) Name of the dedicated server to create (must be a valid hostname as per RFC 1123).
- `location` - (Required, string) The location code to create the dedicated server in. `AMS1`, `SJC1`, etc.
- `server_model` - (Required, string) Name of the server model the dedicated server is created from.
- `ram_size` - (Optional, int) Size of the RAM (GB).
- `operating_system` - (Optional, string) The dedicated server operating system name.
- `private_uplink` - (Required, string) The dedicated server private uplink name.
- `public_uplink` - (Optional, string) The dedicated server public uplink name.
- `bandwidth` - (Optional, string) The dedicated server public bandwidth name.
- `ssh_key_fingerprints.0` - (Optional, string) SSH key fingerprint.
- `ipv6` - (Optional, bool) Is IPv6 enabled. Defaults to `false`.
- `slots` - (Optional, list) List of drive slots.
- `slots.0.position` - (Required, int) The slot position.
- `slots.0.drive_model_name` - (Optional, string) The name of drive model to place in the slot.
- `layout` - (Optional, list) List of layouts.
- `layout.0.slot_positions` - (Required, list) List of slots which should be used by the layout.
- `layout.0.raid` - (Optional, int) RAID level for the layout.
- `layout.0.partition` - (Required, list) List of partitions for the layout.
- `layout.0.partition.0.target` - (Required, string) Target/Mountoint for the partition.
- `layout.0.partition.0.size` - (Required, int) Size of partition (MB).
- `layout.0.partition.0.fill` - (Optional, bool) Autofill partition by all unused space. When set to `true`, the size will be ignored.
- `layout.0.partition.0.fs` - (Optional, string) Filesystem for the partition.

## Attributes Reference

The following attributes are exported:

- `id` - (string) Unique ID of the dedicated server.
- `hostname` - (string) Name of the dedicated server.
- `location` - (string) The location code.
- `server_model` - (string) Name of the server model.
- `ram_size` - (int) Size of the RAM (GB).
- `operating_system` - (string) The dedicated server operating system name.
- `private_uplink` - (string) The dedicated server private uplink name.
- `public_uplink` - (string) The dedicated server public uplink name.
- `bandwidth` - (string) The dedicated server public bandwidth name.
- `status` - (string) The status of the dedicated server.
- `private_ipv4_address` - (string) The private IPv4 address.
- `public_ipv4_address` - (string) The public IPv4 address.
- `configuration` - (string) The current configuration name of the dedicated server.
- `slots` - (list) List of drive slots in the dedicated server.
- `slots.0.position` - (int) The slot position.
- `slots.0.drive_model_name` - (string) The name of drive model.

## Import

Dedicated servers can be imported using the dedicated server `id`:

```
terraform import serverscom_dedicated_server.node_1 <id>
```
