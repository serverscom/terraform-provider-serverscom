---
page_title: "Servers.com: serverscom_ssh_key"
---

# serverscom_ssh_key

Provides a Servers.com SSH key resource to manage SSH keys for dedicated server/cloud computing instance access.

## Example Usage

Create a new SSH key

```hcl
resource "serverscom_ssh_key" "default" {
  name = "Terraform Example"
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) Name of the SSH key.
- `public_key` - (Required, string) Public key. If this is a file, it can be read using the file interpolation function.

## Attributes Reference

The following attributes are exported:

- `id` - (int) Unique identifier of the SSH key.
- `name` - (string) Name of the SSH key.
- `public_key` - (string) Public part of the SSH key.
- `fingerprint` - (string) Fingerprint of the SSH key.

## Import

SSH keys can be imported using the SSH key `fingeprint`:

```bash
terraform import serverscom_ssh_key.default <fingerprint>
```

