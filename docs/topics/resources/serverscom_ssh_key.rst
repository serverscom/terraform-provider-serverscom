.. _resource_serverscom_ssh_key:

SSH Key
=======

Provides a Servers.com SSH key resource to manage SSH keys for dedicated server/cloud computing instance access.

Example
*******

Create a new SSH key

.. sourcecode:: terraform

  resource "serverscom_ssh_key" "default" {
    name = "Terraform Example"
    public_key = "${file("~/.ssh/id_rsa.pub")}"
  }


Argument Reference
******************

The following arguments are supported:

- `name` - (Required, string) Name of the SSH key.
- `public_key` - (Required, string) The public key. If this is a file, it can be read using the file interpolation function

Attributes Reference
********************

The following attributes are exported:

- `id` - (int) The unique ID of the key.
- `name` - (string) The name of the SSH key
- `public_key` - (string) The text of the public key
- `fingerprint` - (string) The fingerprint of the SSH key


Import
******

SSH keys can be imported using the SSH key `fingeprint`:

.. sourcecode:: bash

        terraform import serverscom_ssh_key.default <fingerprint>

.. vi: textwidth=78
