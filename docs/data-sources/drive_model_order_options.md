---
page_title: "Servers.com: serverscom_drive_model_order_options"
---

# serverscom_drive_model_order_options

Lists drive models available for a server model in a location, with optional filtering. Use it to find the drive model name for the `slot` block of the `serverscom_dedicated_server` resource.

## Example Usage

```hcl
data "serverscom_drive_model_order_options" "ssd" {
  location_id     = 1
  server_model_id = 100

  filter {
    media_type = "SSD"
  }
}

output "ssd_drives" {
  value = [for d in data.serverscom_drive_model_order_options.ssd.drive_models : d.name]
}
```

Use a drive model in a server order:

```hcl
resource "serverscom_dedicated_server" "node_1" {
  # ...
  slot {
    position    = 0
    drive_model = data.serverscom_drive_model_order_options.ssd.drive_models[0].name
  }
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `server_model_id` - (Required, int) The ID of the server model. Use the `serverscom_server_model_order_options` data source to find it.
- `filter` - (Optional, block) A block to filter the returned drive models. Supports:
  - `search_pattern` - (Optional, string) Substring match on the drive model name.
  - `media_type` - (Optional, string) Media type: `HDD` or `SSD`.
  - `interface` - (Optional, string) Drive interface, for example `SATA1`, `SATA2`, `SATA3`, `SAS`, `NVMe-PCIe`.

## Attributes Reference

The following attributes are exported:

- `drive_models` - A list of drive models. Each element contains:
  - `id` - The ID of the drive model.
  - `name` - The name of the drive model.
  - `capacity` - The capacity of the drive in GB.
  - `interface` - The interface type of the drive.
  - `form_factor` - The form factor of the drive.
  - `media_type` - The media type: `HDD` or `SSD`.
