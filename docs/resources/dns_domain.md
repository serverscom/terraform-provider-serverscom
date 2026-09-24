---
page_title: "Servers.com: serverscom_dns_domain"
---

# serverscom_dns_domain

Provides a Servers.com DNS domain (zone) resource.

## Example Usage

```hcl
resource "serverscom_dns_domain" "example" {
  name  = "example.com"
  email = "admin@example.com"
  ttl   = 3600

  labels = {
    env = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) Domain FQDN, lowercase, max 255 characters. Changing this forces a new resource to be created.
- `email` - (Required, string) SOA responsible person (RNAME).
- `ttl` - (Optional, number) Default TTL in seconds (0-604800).
- `labels` - (Optional, map) A map of labels assigned to the domain.

~> **Note:** The SOA record is system-managed and cannot be created or edited directly. Only `labels` can be updated in place. The DNS domain API does not accept `email` or `ttl` updates on an existing domain, so changing either of these after creation will not be applied by the provider. PTR domains are not supported.

## Attributes Reference

The following attributes are exported:

- `id` - (string) Unique identifier of the domain.
- `name` - (string) Domain FQDN.
- `email` - (string) SOA responsible person (RNAME).
- `ttl` - (number) Default TTL in seconds.
- `labels` - (map) A map of labels assigned to the domain.
- `delegation_status` - (string) Delegation status. One of: `delegated`, `verified`, `undelegated`.
- `unpublish_date` - (string) When the domain becomes unavailable.
- `created_at` - (string) The creation time of the domain.
- `updated_at` - (string) The last update time of the domain.

## Import

DNS domains can be imported using the domain `id`:

```bash
terraform import serverscom_dns_domain.example <id>
```
