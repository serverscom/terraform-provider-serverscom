package serverscom

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomDNSDomainRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerscomDNSDomainRecordRead,

		Schema: map[string]*schema.Schema{
			"dns_domain_id": {Type: schema.TypeString, Required: true},
			"id":            {Type: schema.TypeString, Required: true},

			"name":       {Type: schema.TypeString, Computed: true},
			"type":       {Type: schema.TypeString, Computed: true},
			"data":       {Type: schema.TypeString, Computed: true},
			"ttl":        {Type: schema.TypeInt, Computed: true},
			"priority":   {Type: schema.TypeInt, Computed: true},
			"created_at": {Type: schema.TypeString, Computed: true},
			"updated_at": {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceServerscomDNSDomainRecordRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	domainID := d.Get("dns_domain_id").(string)
	recordID := d.Get("id").(string)

	record, err := client.DNS.GetRecord(ctx, domainID, recordID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving dns record: %s", err))
	}

	d.SetId(record.ID)
	d.Set("dns_domain_id", record.DomainID)
	d.Set("name", record.Name)
	d.Set("type", string(record.Type))

	if record.Data != nil {
		d.Set("data", *record.Data)
	}
	if record.TTL != nil {
		d.Set("ttl", *record.TTL)
	}
	if record.Priority != nil {
		d.Set("priority", *record.Priority)
	}

	d.Set("created_at", record.Created.Format(time.RFC3339))
	d.Set("updated_at", record.Updated.Format(time.RFC3339))

	return nil
}
