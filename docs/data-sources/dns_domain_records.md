---
page_title: "Servers.com: serverscom_dns_domain_records"
---

# serverscom_dns_domain_records

Lists all records for a given DNS domain with optional filtering.

## Example Usage

```hcl
data "serverscom_dns_domain_records" "all" {
  dns_domain_id = serverscom_dns_domain.example.id
}

data "serverscom_dns_domain_records" "a_records" {
  dns_domain_id = serverscom_dns_domain.example.id

  filter {
    type = "A"
  }
}
```

## Argument Reference

The following arguments are supported:

- `dns_domain_id` - (Required) The ID of the parent DNS domain.
- `filter` - (Optional) A block to filter the returned records. Supports:
  - `type` - (Optional) Filter by record type.
  - `name` - (Optional) Filter by record name.

## Attributes Reference

The following attributes are exported:

- `dns_domain_records` - A list of DNS records. Each element contains:
  - `id` - The ID of the DNS record.
  - `name` - Record FQDN.
  - `type` - Record type.
  - `data` - Record data.
  - `ttl` - TTL in seconds.
  - `priority` - Priority for `MX` and `SRV` records.
  - `created_at` - The creation time of the record.
  - `updated_at` - The last update time of the record.
