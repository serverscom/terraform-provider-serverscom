---
page_title: "Servers.com: serverscom_sbm_server"
---

# serverscom_sbm_server

Provides a resource of the Servers.com Scalable Bare Metal server. This can be used to create, modify, and delete SBM servers. Learn more about the service [in the knowledge base](https://www.servers.com/support/knowledge/scalable-bare-metal/).

## Example Usage

Create a new SBM server:

```hcl
resource "serverscom_sbm_server" "node_01" {
  hostname             = "node-01"
  location             = "AMS1"
  flavor               = "SBM-01"
  operating_system     = "Ubuntu 24.04-server x86_64"
  ssh_key_fingerprints = [
    "cf:1d:09:ab:cb:47:97:3f:50:9a:f0:34:14:78:fa:1b"
  ]
} 
```

## Argument Reference

The following arguments are supported:

- `hostname` - (Required, string) A name of the SBM server.
- `location` - (Required, string) A location code of the SBM server. For example: `AMS1`, `SJC1`, etc.
- `flavor` - (Required, string) A flavor of an SBM server.
- `operating_system` - (Required, string) A name of an operating system.
- `ssh_key_fingerprints` - (Optional, list) An SSH key fingerprint.
- `user_data` - (Optional, string) A user data string for the SBM server.
- `private_ipv4_network_id` - (Optional, string) An ID of a private IPv4 network.
- `private_ipv4_address` - (Optional, string) A private IPv4 address for the SBM server.
- `public_ipv4_network_id` - (Optional, string) An ID of a public IPv4 network.
- `public_ipv4_address` - (Optional, string) A public IPv4 address for the SBM server.
- `labels` - (Optional, map) A map of labels assigned to the SBM server.

## Attributes Reference

The following attributes are exported:

- `id` - (string) Unique identifier of the SBM server.
- `hostname` - (string) A name of the SBM server.
- `location` - (string) A location code of the SBM server.
- `flavor` - (string) A flavor of an SBM server.
- `operating_system` - (string) A name of an operating system.
- `private_ipv4_address` - (string) A private IPv4 address for the SBM server.
- `public_ipv4_address` - (string) A public IPv4 address for the SBM server.
- `status` - (string) Status of the SBM server.
- `labels` - (map) A map of labels assigned to the SBM server.

## Import

SBM servers can be imported using the SBM server `id`:

```bash
terraform import serverscom_sbm_server.node_01 <id>
```
