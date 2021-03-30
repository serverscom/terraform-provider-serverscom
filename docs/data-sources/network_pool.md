---
page_title: "Servers.com: serverscom_network_pool"
---

# network_pool

Get information on a Network Pool for use in other resources. This data source provides all of the Network Pool properties.

## Example Usage

Get the Network Pool by ID:

```hcl
data "serverscom_network_pool" "example" {
  id = "QeZ89zQb"
}

output "network_pool_example" {
  value = data.serverscom_network_pool.example.cidr
}
```

Get the Network Pool by CIDR:

```hcl
data "serverscom_network_pool" "example" {
  cidr = "10.0.0.0/20"
}
```
## Argument Reference

One of the following arguments must be provided:

* `id` - (Optional) The ID of the Network Pool.
* `cidr` - (Optional) The CIDR of the Network Pool.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the Network Pool.
* `cidr` - The CIDR of the Network Pool.
* `title` - The title of the Network Pool.
* `type` - Type of the Network Pool.
* `created_at` - The creation time of the Network Pool.
