---
page_title: "Servers.com: serverscom_dns_domain_record"
---

# serverscom_dns_domain_record

Provides a Servers.com DNS record resource to manage a single record within a DNS domain.

## Example Usage

```hcl
resource "serverscom_dns_domain_record" "www" {
  dns_domain_id = serverscom_dns_domain.example.id
  name          = "www.example.com"
  type          = "A"
  data          = "192.0.2.1"
  ttl           = 300
}

resource "serverscom_dns_domain_record" "mail" {
  dns_domain_id = serverscom_dns_domain.example.id
  name          = "example.com"
  type          = "MX"
  data          = "mail.example.com"
  priority      = 10
  ttl           = 3600
}

resource "serverscom_dns_domain_record" "spf" {
  dns_domain_id = serverscom_dns_domain.example.id
  name          = "example.com"
  type          = "TXT"
  data          = "v=spf1 include:example.com ~all"
}

resource "serverscom_dns_domain_record" "caa" {
  dns_domain_id = serverscom_dns_domain.example.id
  name          = "example.com"
  type          = "CAA"
  data          = "0 issue letsencrypt.org"
}

resource "serverscom_dns_domain_record" "srv" {
  dns_domain_id = serverscom_dns_domain.example.id
  name          = "_sip._tcp.example.com"
  type          = "SRV"
  data          = "10 20 5060 sip.example.com"
}
```

## Argument Reference

The following arguments are supported:

- `dns_domain_id` - (Required, string) Parent domain ID. Changing this forces a new resource to be created.
- `type` - (Required, string) Record type. One of: `A`, `AAAA`, `CAA`, `CNAME`, `MX`, `NS`, `TXT`, `SRV`, `ALIAS`. Changing this forces a new resource to be created.
- `name` - (Required, string) Record FQDN, must be within the parent domain.
- `data` - (Required, string) Record data. Format depends on `type`.
- `ttl` - (Optional, number) TTL in seconds (0-604800).
- `priority` - (Optional, number) Priority for `MX` and `SRV` records.

~> **Note:** `SOA` and `PTR` records are system-managed (read-only) and not supported by this resource.

## Attributes Reference

The following attributes are exported:

- `id` - (string) Unique identifier of the record.
- `dns_domain_id` - (string) Parent domain ID.
- `type` - (string) Record type.
- `name` - (string) Record FQDN.
- `data` - (string) Record data.
- `ttl` - (number) TTL in seconds.
- `priority` - (number) Priority for `MX` and `SRV` records.
- `created_at` - (string) The creation time of the record.
- `updated_at` - (string) The last update time of the record.

## Import

DNS records can be imported using a composite id of the parent domain ID and the record ID:

```bash
terraform import serverscom_dns_domain_record.example <dns_domain_id>/<record_id>
```
