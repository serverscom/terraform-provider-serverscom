---
page_title: "Servers.com: serverscom_cloud_computing_instance"
---

# serverscom_cloud_computing_instance

Provides an Servers.com cloud computing instance resource. This can be used to create, modify, and delete cloud computing instances. Cloud computing instances also support [provisioning](https://www.terraform.io/docs/provisioners/index.html).

## Example Usage

Create a new cloud computing instance

```hcl
resource "serverscom_cloud_computing_instance" "instance_1" {
  name = "instance-1"
  region = "NL01"
  image = "Ubuntu 20.04-server (64 bit)"

  flavor = "SSD.50"

  gpn_enabled = true
  ipv6_enabled = true
  backup_copies = 5

  ssh_key_fingerprint = "cf:1d:09:ab:cb:47:97:3f:50:9a:f0:34:14:78:fa:1b"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) Name of the cloud instance (according to RFC 1123 specification).
- `region` - (Required, string) Cloud computing region code.
- `image` - (Required, string) Name of the image.
- `flavor` - (Required, string) Name of the flavor.
- `gpn_enabled` - (Optional, bool) Is GPN network enabled. Defaults to `false`.
- `ipv6_enabled` - (Optional, bool) Is IPv6 enabled. Defaults to `false`.
- `backup_copies` - (Optional, int) Count of backup copies. Defaults to `0`.
- `ssh_key_fingerprint` - (Optional, string) SSH key fingerprint.
- `user_data` - (Optional, string) A string of the desired user data for the cloud computing instance.
- `labels` - (Optional, map) A map of labels assigned to the cloud computing instance.

## Attributes Reference

The following attributes are exported:

- `id` - (string) Unique identifier of the cloud computing instance.
- `name` - (string) Name of the cloud instance (according to RFC 1123 specification).
- `region` - (string) Cloud computing region code.
- `image` - (string) Name of the image.
- `flavor` - (string) Name of the flavor.
- `gpn_enabled` - (bool) Is GPN network enabled. Defaults to `false`.
- `ipv6_enabled` - (bool) Is IPv6 enabled. Defaults to `false`.
- `backup_copies` - (int) Count of backup copies. Defaults to `0`.
- `status` - (string) Status of the cloud computing instance.
- `private_ipv4_address` - (string) Private IPv4 address.
- `public_ipv4_address` - (string) Public IPv4 address.
- `public_ipv6_address` - (string) Public IPv6 address.
- `openstack_uuid` - (string) OpenStack unique identifier (UUID) of the cloud computing instance.
- `labels` - (map) A map of labels assigned to the cloud computing instance.

## Import

Cloud computing instances can be imported using the cloud computing
instance `id`:

```bash
terraform import serverscom_cloud_computing_instance.instance_1 <id>
```
