---
layout: "serverscom"
page_title: "Servers.com: serverscom_ssh_key"
sidebar_current: "docs-serverscom-resource-ssh-key"
description: |-
  Provides a Servers.com SSH key resource to manage SSH keys for dedicated server/cloud computing instance access.
---

# serverscom_ssh_key

Provides a Servers.com SSH key resource to manage SSH keys for dedicated server/cloud computing instance access.

## Example Usage

```hcl
# Create a new SSH key
resource "serverscom_ssh_key" "default" {
  name = "Terraform Example"
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) Name of the SSH key.
- `public_key` - (Required, string) The public key. If this is a file, it can be read using the file interpolation function

## Attributes Reference

The following attributes are exported:

- `id` - (int) The unique ID of the key.
- `name` - (string) The name of the SSH key
- `public_key` - (string) The text of the public key
- `fingerprint` - (string) The fingerprint of the SSH key
