---
page_title: "Servers.com Provider"
subcategory: ""
description: |-

---

# Servers.com Provider

The Servers.com Provider allows interacting with Servers.com services. The provider has to be set up properly before using, as shown in the Example Usage. Available resources and helpful guides are listed in the navigation bar.


## Example Usage

Terraform 0.13 and later:

```
terraform {
  required_providers {
    serverscom = {
      source = "serverscom/serverscom"
      version = "0.2.2"
    }
  }
}

provider "serverscom" {
  # Configuration options
}
```

Terraform 0.12 and earlier:

```
terraform
# Configure the Servers.com Provider
provider "serverscom" {
  token = "<your API token>"
  endpoint = "https://api.servers.com/v1"
}

# Create a dedicated server
resource "serverscom_dedicated_server" "node_1" {
  ...
}
```

## Schema

- `token` (Required, string) - A token to perform API requests for Servers.com services. It can be obtained in the [Customer Portal](https://portal.servers.com/#/profile/api-tokens).
- `endpoint` (Optional, string) - The Servers.com API endpoint. In most cases, the default one is used: `https://api.servers.com/v1`.
