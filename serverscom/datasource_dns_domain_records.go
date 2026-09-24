package serverscom

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomDNSDomainRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerscomDNSDomainRecordsRead,

		Schema: map[string]*schema.Schema{
			"dns_domain_id": {Type: schema.TypeString, Required: true},

			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {Type: schema.TypeString, Optional: true},
						"name": {Type: schema.TypeString, Optional: true},
					},
				},
			},

			"dns_domain_records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":         {Type: schema.TypeString, Computed: true},
						"name":       {Type: schema.TypeString, Computed: true},
						"type":       {Type: schema.TypeString, Computed: true},
						"data":       {Type: schema.TypeString, Computed: true},
						"ttl":        {Type: schema.TypeInt, Computed: true},
						"priority":   {Type: schema.TypeInt, Computed: true},
						"created_at": {Type: schema.TypeString, Computed: true},
						"updated_at": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceServerscomDNSDomainRecordsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	domainID := d.Get("dns_domain_id").(string)

	col := client.DNS.Records(domainID)

	id := fmt.Sprintf("dns-domain-records-%s", domainID)

	if v, ok := d.GetOk("filter"); ok {
		filter := v.([]any)[0].(map[string]any)

		if t, ok := filter["type"]; ok && t.(string) != "" {
			col = col.SetParam("type", t.(string))
		}
		if n, ok := filter["name"]; ok && n.(string) != "" {
			col = col.SetParam("name", n.(string))
		}

		hash, err := hashFilter(filter)
		if err != nil {
			return diag.FromErr(err)
		}
		id = fmt.Sprintf("dns-domain-records-%s-%s", domainID, hash)
	}

	records, err := col.Collect(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving dns records: %s", err))
	}

	list := make([]map[string]any, 0, len(records))
	for _, record := range records {
		m := map[string]any{
			"id":         record.ID,
			"name":       record.Name,
			"type":       string(record.Type),
			"created_at": record.Created.Format(time.RFC3339),
			"updated_at": record.Updated.Format(time.RFC3339),
		}
		if record.Data != nil {
			m["data"] = *record.Data
		}
		if record.TTL != nil {
			m["ttl"] = *record.TTL
		}
		if record.Priority != nil {
			m["priority"] = *record.Priority
		}
		list = append(list, m)
	}

	d.SetId(id)
	if err := d.Set("dns_domain_records", list); err != nil {
		return diag.FromErr(fmt.Errorf("error setting dns records: %s", err))
	}

	return nil
}
