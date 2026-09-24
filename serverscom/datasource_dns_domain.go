package serverscom

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomDNSDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerscomDNSDomainRead,

		Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Required: true},

			"name":  {Type: schema.TypeString, Computed: true},
			"email": {Type: schema.TypeString, Computed: true},
			"ttl":   {Type: schema.TypeInt, Computed: true},
			"labels": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{Type: schema.TypeString}, Computed: true,
			},
			"delegation_status": {Type: schema.TypeString, Computed: true},
			"unpublish_date":    {Type: schema.TypeString, Computed: true},
			"created_at":        {Type: schema.TypeString, Computed: true},
			"updated_at":        {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceServerscomDNSDomainRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	id := d.Get("id").(string)

	domain, err := client.DNS.GetDomain(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving dns domain: %s", err))
	}

	d.SetId(domain.ID)
	d.Set("name", domain.Name)
	d.Set("email", domain.Email)
	d.Set("ttl", domain.TTL)
	d.Set("labels", domain.Labels)
	d.Set("delegation_status", string(domain.DelegationStatus))

	if domain.UnpublishDate != nil {
		d.Set("unpublish_date", domain.UnpublishDate.Format(time.RFC3339))
	}
	d.Set("created_at", domain.Created.Format(time.RFC3339))
	d.Set("updated_at", domain.Updated.Format(time.RFC3339))

	return nil
}
