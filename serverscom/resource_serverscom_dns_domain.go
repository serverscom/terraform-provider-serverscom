package serverscom

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func resourceServerscomDNSDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerscomDNSDomainCreate,
		ReadContext:   resourceServerscomDNSDomainRead,
		UpdateContext: resourceServerscomDNSDomainUpdate,
		DeleteContext: resourceServerscomDNSDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 604800),
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"delegation_status": {Type: schema.TypeString, Computed: true},
			"unpublish_date":    {Type: schema.TypeString, Computed: true},
			"created_at":        {Type: schema.TypeString, Computed: true},
			"updated_at":        {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceServerscomDNSDomainCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	input := scgo.DNSDomainCreateInput{
		Name:  d.Get("name").(string),
		Email: d.Get("email").(string),
	}
	if v, ok := d.GetOk("ttl"); ok {
		input.TTL = v.(int)
	}
	if labelsRaw, ok := d.GetOk("labels"); ok {
		input.Labels = toStringMap(labelsRaw.(map[string]any))
	}

	domain, err := client.DNS.CreateDomain(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain.ID)

	return resourceServerscomDNSDomainRead(ctx, d, meta)
}

func resourceServerscomDNSDomainRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	domain, err := client.DNS.GetDomain(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*scgo.NotFoundError); ok {
			log.Printf("[WARN] Serverscom dns domain (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", domain.Name)
	d.Set("email", domain.Email)
	d.Set("ttl", domain.TTL)
	d.Set("labels", domain.Labels)
	d.Set("delegation_status", string(domain.DelegationStatus))

	if domain.UnpublishDate != nil {
		d.Set("unpublish_date", domain.UnpublishDate.Format(time.RFC3339))
	} else {
		d.Set("unpublish_date", nil)
	}
	d.Set("created_at", domain.Created.Format(time.RFC3339))
	d.Set("updated_at", domain.Updated.Format(time.RFC3339))

	return nil
}

func resourceServerscomDNSDomainUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	// Only labels can be updated in place. The API's domain update endpoint
	// (DNSDomainUpdateInput) does not accept email or ttl changes.
	input := scgo.DNSDomainUpdateInput{Labels: map[string]string{}}
	if labelsRaw, ok := d.GetOk("labels"); ok {
		input.Labels = toStringMap(labelsRaw.(map[string]any))
	}

	if _, err := client.DNS.UpdateDomain(ctx, d.Id(), input); err != nil {
		return diag.FromErr(err)
	}

	return resourceServerscomDNSDomainRead(ctx, d, meta)
}

func resourceServerscomDNSDomainDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	err := client.DNS.DeleteDomain(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*scgo.NotFoundError); ok {
			log.Printf("[WARN] Serverscom dns domain (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
