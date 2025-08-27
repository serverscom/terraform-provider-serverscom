---
page_title: "Servers.com: serverscom_sbm_server"
---

# serverscom_sbm_server

Get information on a Server Building Manager (SBM) server for use in other resources. This data source provides all of the SBM server properties.

## Example Usage

Get the SBM server by ID:

```hcl
data "serverscom_sbm_server" "example" {
  id = "SBM7zQnVb"
}

output "sbm_server_example" {
  value = data.serverscom_sbm_server.example.title
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the SBM server.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SBM server.
* `rack_id` - The rack ID where the server is located.
* `title` - The title/hostname of the SBM server.
* `location_id` - The location ID of the SBM server.
* `location_code` - The location code of the SBM server (e.g., "AMS1", "SJC1").
* `status` - The status of the SBM server.
* `operational_status` - The operational status of the SBM server.
* `power_status` - The power status of the SBM server.
* `configuration` - The configuration description of the SBM server.
* `private_ipv4_address` - The private IPv4 address of the SBM server.
* `public_ipv4_address` - The public IPv4 address of the SBM server.
* `lease_start_at` - The lease start date and time.
* `scheduled_release_at` - The scheduled release date and time.
* `type` - The type of the SBM server.
* `oob_ipv4_address` - The out-of-band IPv4 address.
* `ram_size` - The RAM size in GB.
* `server_model_id` - The server model ID.
* `server_model_name` - The server model name.
* `bandwidth_id` - The bandwidth ID.
* `bandwidth_name` - The bandwidth name.
* `private_uplink_id` - The private uplink ID.
* `private_uplink_name` - The private uplink name.
* `public_uplink_id` - The public uplink ID.
* `public_uplink_name` - The public uplink name.
* `operating_system_id` - The operating system ID.
* `operating_system_full_name` - The full name of the operating system.
* `labels` - A map of labels assigned to the SBM server.
* `created_at` - The creation time of the SBM server.
* `updated_at` - The last update time of the SBM server.
