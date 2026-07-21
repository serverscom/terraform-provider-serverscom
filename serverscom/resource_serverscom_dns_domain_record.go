package serverscom

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

var dnsRecordTypes = []string{"A", "AAAA", "CAA", "CNAME", "MX", "NS", "TXT", "SRV", "ALIAS"}

func resourceServerscomDNSDomainRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerscomDNSDomainRecordCreate,
		ReadContext:   resourceServerscomDNSDomainRecordRead,
		UpdateContext: resourceServerscomDNSDomainRecordUpdate,
		DeleteContext: resourceServerscomDNSDomainRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServerscomDNSDomainRecordImport,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"dns_domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(dnsRecordTypes, false),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"data": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 604800),
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"created_at": {Type: schema.TypeString, Computed: true},
			"updated_at": {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceServerscomDNSDomainRecordCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	domainID := d.Get("dns_domain_id").(string)

	input := scgo.DNSRecordCreateInput{
		Name:     d.Get("name").(string),
		Type:     scgo.DNSRecordType(d.Get("type").(string)),
		Data:     d.Get("data").(string),
		TTL:      optionalConfigInt(d, "ttl"),
		Priority: optionalConfigInt(d, "priority"),
	}

	record, err := client.DNS.CreateRecord(ctx, domainID, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(record.ID)

	return resourceServerscomDNSDomainRecordRead(ctx, d, meta)
}

func resourceServerscomDNSDomainRecordRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	domainID := d.Get("dns_domain_id").(string)

	record, err := client.DNS.GetRecord(ctx, domainID, d.Id())
	if err != nil {
		if _, ok := err.(*scgo.NotFoundError); ok {
			log.Printf("[WARN] Serverscom dns record (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

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

func resourceServerscomDNSDomainRecordUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	domainID := d.Get("dns_domain_id").(string)

	input := scgo.DNSRecordUpdateInput{}
	if d.HasChange("data") {
		input.Data = d.Get("data").(string)
	}
	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("ttl") {
		v := d.Get("ttl").(int)
		input.TTL = &v
	}
	if d.HasChange("priority") {
		v := d.Get("priority").(int)
		input.Priority = &v
	}

	if _, err := client.DNS.UpdateRecord(ctx, domainID, d.Id(), input); err != nil {
		return diag.FromErr(err)
	}

	return resourceServerscomDNSDomainRecordRead(ctx, d, meta)
}

func resourceServerscomDNSDomainRecordDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	domainID := d.Get("dns_domain_id").(string)

	err := client.DNS.DeleteRecord(ctx, domainID, d.Id())
	if err != nil {
		if _, ok := err.(*scgo.NotFoundError); ok {
			log.Printf("[WARN] Serverscom dns record (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceServerscomDNSDomainRecordImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q, expected format <dns_domain_id>/<record_id>", d.Id())
	}

	d.Set("dns_domain_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
