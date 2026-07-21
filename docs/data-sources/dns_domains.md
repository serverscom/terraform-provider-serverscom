---
page_title: "Servers.com: serverscom_dns_domains"
---

# serverscom_dns_domains

Lists DNS domains with optional filtering.

## Example Usage

```hcl
data "serverscom_dns_domains" "production" {
  filter {
    label_selector    = "env=production"
    delegation_status = "delegated"
  }
}
```

## Argument Reference

The following arguments are supported:

- `filter` - (Optional) A block to filter the returned domains. Supports:
  - `search_pattern` - (Optional) Substring match on the domain name.
  - `delegation_status` - (Optional) One of: `delegated`, `verified`, `undelegated`.
  - `label_selector` - (Optional) Label selector, e.g. `env=production`.

## Attributes Reference

The following attributes are exported:

- `dns_domains` - A list of DNS domains. Each element contains:
  - `id` - The ID of the DNS domain.
  - `name` - Domain FQDN.
  - `email` - SOA responsible person (RNAME).
  - `ttl` - Default TTL in seconds.
  - `labels` - A map of labels assigned to the domain.
  - `delegation_status` - Delegation status.
  - `unpublish_date` - When the domain becomes unavailable.
  - `created_at` - The creation time of the domain.
  - `updated_at` - The last update time of the domain.
