package serverscom

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomDNSDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerscomDNSDomainsRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_pattern":    {Type: schema.TypeString, Optional: true},
						"delegation_status": {Type: schema.TypeString, Optional: true},
						"label_selector":    {Type: schema.TypeString, Optional: true},
					},
				},
			},

			"dns_domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":    {Type: schema.TypeString, Computed: true},
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
				},
			},
		},
	}
}

func dataSourceServerscomDNSDomainsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	col := client.DNS.Collection()

	id := "dns-domains"

	if v, ok := d.GetOk("filter"); ok {
		filter := v.([]any)[0].(map[string]any)

		if sp, ok := filter["search_pattern"]; ok && sp.(string) != "" {
			col = col.SetParam("search_pattern", sp.(string))
		}
		if ds, ok := filter["delegation_status"]; ok && ds.(string) != "" {
			col = col.SetParam("delegation_status", ds.(string))
		}
		if ls, ok := filter["label_selector"]; ok && ls.(string) != "" {
			col = col.SetParam("label_selector", ls.(string))
		}

		hash, err := hashFilter(filter)
		if err != nil {
			return diag.FromErr(err)
		}
		id = fmt.Sprintf("dns-domains-%s", hash)
	}

	domains, err := col.Collect(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving dns domains: %s", err))
	}

	list := make([]map[string]any, 0, len(domains))
	for _, domain := range domains {
		m := map[string]any{
			"id":                domain.ID,
			"name":              domain.Name,
			"email":             domain.Email,
			"ttl":               domain.TTL,
			"labels":            domain.Labels,
			"delegation_status": string(domain.DelegationStatus),
			"created_at":        domain.Created.Format(time.RFC3339),
			"updated_at":        domain.Updated.Format(time.RFC3339),
		}
		if domain.UnpublishDate != nil {
			m["unpublish_date"] = domain.UnpublishDate.Format(time.RFC3339)
		}
		list = append(list, m)
	}

	d.SetId(id)
	if err := d.Set("dns_domains", list); err != nil {
		return diag.FromErr(fmt.Errorf("error setting dns domains: %s", err))
	}

	return nil
}
