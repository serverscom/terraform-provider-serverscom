---
page_title: "Servers.com: serverscom_dns_domain_record"
---

# serverscom_dns_domain_record

Get information on a single DNS record by domain ID and record ID.

## Example Usage

```hcl
data "serverscom_dns_domain_record" "www" {
  dns_domain_id = serverscom_dns_domain.example.id
  id            = "rec123"
}
```

## Argument Reference

The following arguments are supported:

- `dns_domain_id` - (Required) The ID of the parent DNS domain.
- `id` - (Required) The ID of the DNS record.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the DNS record.
- `dns_domain_id` - The ID of the parent DNS domain.
- `name` - Record FQDN.
- `type` - Record type.
- `data` - Record data.
- `ttl` - TTL in seconds.
- `priority` - Priority for `MX` and `SRV` records.
- `created_at` - The creation time of the record.
- `updated_at` - The last update time of the record.
