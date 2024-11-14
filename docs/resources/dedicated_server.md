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
- `user_data` - (Optional, string) A string of the desired user data for the dedicated server.
- `ipv6` - (Optional, bool) Is IPv6 enabled. Defaults to `false`.
- `slot` - (Optional, list) List of drive slots. Slots used in partioning have to be listed.
- `slot.0.position` - (Required, int) Slot position.
- `slot.0.drive_model_name` - (Optional, string) The name of drive model to place in the slot.
- `layout` - (Optional, list) List of layouts.
- `layout.0.slot_positions` - (Required, list) List of slots which should be used in the layout.
- `layout.0.raid` - (Optional, int) RAID level for the layout.
- `layout.0.partition` - (Required, list) List of partitions for the layout.
- `layout.0.partition.0.target` - (Required, string) Target/Mount point for the partition.
- `layout.0.partition.0.size` - (Required, int) Size of the partition (MB).
- `layout.0.partition.0.fill` - (Optional, bool) Autofill partition by all unused space. When set to `true`
