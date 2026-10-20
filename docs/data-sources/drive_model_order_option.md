---
page_title: "Servers.com: serverscom_drive_model_order_option"
---

# serverscom_drive_model_order_option

Get information on a drive model available for a server model in a location. To list all drive models, use the `serverscom_drive_model_order_options` data source.

The `name` attribute is the value of the `drive_model` argument of the `slot` block in the `serverscom_dedicated_server` resource.

## Example Usage

```hcl
data "serverscom_drive_model_order_option" "example" {
  location_id     = 1
  server_model_id = 100
  id              = 200
}

output "drive_model_name" {
  value = data.serverscom_drive_model_order_option.example.name
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `server_model_id` - (Required, int) The ID of the server model. Use the `serverscom_server_model_order_options` data source to find it.
- `id` - (Required, int) The ID of the drive model. Use the `serverscom_drive_model_order_options` data source to list IDs.

## Attributes Reference

The following attributes are exported:

- `id` - (int) The ID of the drive model.
- `name` - (string) The name of the drive model.
- `capacity` - (int) The capacity of the drive in GB.
- `interface` - (string) The interface type of the drive.
- `form_factor` - (string) The form factor of the drive.
- `media_type` - (string) The media type: `HDD` or `SSD`.
