---
page_title: "Servers.com: serverscom_subnetwork"
---

# serverscom_subnetwork

Provides a Servers.com Subnetwork resource to manage Subnetworks for dedicated server networks provisioning.

## Example Usage

Create a new Subnetwork by mask:

```hcl
resource "serverscom_subnetwork" "private_network" {
  network_pool_id = "QeZ89zQb"
  title = "Custom private_network"
  mask = 29
}
```

Create a new Subnetwork by CIDR:

```hcl
resource "serverscom_subnetwork" "private_network" {
  network_pool_id = "QeZ89zQb"
  title = "Custom private_network"
  cidr = "10.0.0.0/29"
}
```

## Argument Reference

The following arguments are supported:

- `network_pool_id` - (Required, string) ID of the Network Pool.
- `title` - (Optional, string) Title of the Subnetwork.
- `cidr` - (Optional, string) CIDR of the Subnetwork.
- `mask` - (Optional, int) Mask of the Subnetwork.

## Attributes Reference

The following attributes are exported:

- `id` - (string) Unique identifier of the Subnetwork.
- `title` - (Optional, string) Title of the Subnetwork.
- `cidr` - (string) CIDR of the Subnetwork.
- `mask` - (int) Mask of the Subnetwork.
- `network_pool_id` - (string) Network Pool ID of the subnetwork.

## Import

Subnetworks can be imported using the Subnetwork `id`:

```bash
terraform import serverscom_subnetwork.private_network <id>
```

