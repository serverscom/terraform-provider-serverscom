---
page_title: "Servers.com: serverscom_dns_domain"
---

# serverscom_dns_domain

Get information on a DNS domain by its ID.

## Example Usage

```hcl
data "serverscom_dns_domain" "example" {
  id = "abc123"
}

output "dns_domain_name" {
  value = data.serverscom_dns_domain.example.name
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Required) The ID of the DNS domain.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the DNS domain.
- `name` - Domain FQDN.
- `email` - SOA responsible person (RNAME).
- `ttl` - Default TTL in seconds.
- `labels` - A map of labels assigned to the domain.
- `delegation_status` - Delegation status. One of: `delegated`, `verified`, `undelegated`.
- `unpublish_date` - When the domain becomes unavailable.
- `created_at` - The creation time of the domain.
- `updated_at` - The last update time of the domain.
