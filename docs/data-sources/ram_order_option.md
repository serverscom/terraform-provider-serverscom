---
page_title: "Servers.com: serverscom_ram_order_option"
---

# serverscom_ram_order_option

Checks that a RAM size and type is available for a server model in a location. The data source finds the option matching both `ram` and `type` and fails with an error if there is none, so it can be used to validate a `ram_size` before ordering a server. To list all RAM options, use the `serverscom_ram_order_options` data source.

## Example Usage

```hcl
data "serverscom_ram_order_option" "example" {
  location_id     = 1
  server_model_id = 100
  ram             = 64
  type            = "DDR4"
}

resource "serverscom_dedicated_server" "node_1" {
  # ...
  ram_size = data.serverscom_ram_order_option.example.ram
}
```

## Argument Reference

The following arguments are supported:

- `location_id` - (Required, int) The ID of the location. Use the `serverscom_locations` data source to find it.
- `server_model_id` - (Required, int) The ID of the server model. Use the `serverscom_server_model_order_options` data source to find it.
- `ram` - (Required, int) RAM size in GB to match.
- `type` - (Required, string) RAM type to match, for example `DDR4`.

## Attributes Reference

The following attributes are exported:

- `id` - (string) A generated identifier of the matched option.
- `ram` - (int) The RAM size in GB.
- `type` - (string) The DDR type of RAM.
