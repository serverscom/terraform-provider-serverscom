---
page_title: "Servers.com: serverscom_cloud_computing_instance"
---

# serverscom_cloud_computing_instance

Get information on a cloud computing instance for use in other resources. This data source provides all of the cloud computing instance properties.

## Example Usage

Get the cloud computing instance by ID:

```hcl
data "serverscom_cloud_computing_instance" "example" {
  id = "CCInst7zQnVb"
}

output "cloud_instance_example" {
  value = data.serverscom_cloud_computing_instance.example.name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the cloud computing instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the cloud computing instance.
* `openstack_uuid` - The OpenStack UUID of the instance.
* `name` - The name of the cloud computing instance.
* `status` - The status of the cloud computing instance.
* `region_id` - The region ID where the instance is located.
* `region_code` - The region code where the instance is located.
* `flavor_id` - The flavor ID of the instance.
* `flavor_name` - The flavor name of the instance.
* `image_id` - The image ID used to create the instance.
* `image_name` - The image name used to create the instance.
* `public_ipv4_address` - The public IPv4 address of the instance.
* `private_ipv4_address` - The private IPv4 address of the instance.
* `local_ipv4_address` - The local IPv4 address of the instance.
* `public_ipv6_address` - The public IPv6 address of the instance.
* `ipv6_enabled` - Whether IPv6 is enabled for the instance.
* `gpn_enabled` - Whether GPN (Global Private Network) is enabled for the instance.
* `backup_copies` - The number of backup copies for the instance.
* `public_port_blocked` - Whether the public port is blocked for the instance.
* `labels` - A map of labels assigned to the cloud computing instance.
* `created_at` - The creation time of the cloud computing instance.
* `updated_at` - The last update time of the cloud computing instance.
