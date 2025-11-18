package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomRBSVolumes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerscomRBSVolumesRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location_id":    {Type: schema.TypeInt, Optional: true},
						"label_selector": {Type: schema.TypeString, Optional: true},
						"search_pattern": {Type: schema.TypeString, Optional: true},
					},
				},
			},

			"rbs_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":     {Type: schema.TypeString, Computed: true},
						"name":   {Type: schema.TypeString, Computed: true},
						"size":   {Type: schema.TypeInt, Computed: true},
						"status": {Type: schema.TypeString, Computed: true},
						"labels": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{Type: schema.TypeString}, Computed: true,
						},
						"location_id":   {Type: schema.TypeInt, Computed: true},
						"location_code": {Type: schema.TypeString, Computed: true},
						"ip_address":    {Type: schema.TypeString, Computed: true},
						"flavor_id":     {Type: schema.TypeInt, Computed: true},
						"flavor_name":   {Type: schema.TypeString, Computed: true},
						"iops":          {Type: schema.TypeFloat, Computed: true},
						"bandwidth":     {Type: schema.TypeFloat, Computed: true},
						"target_iqn":    {Type: schema.TypeString, Computed: true},
						"created_at":    {Type: schema.TypeString, Computed: true},
						"updated_at":    {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceServerscomRBSVolumesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	col := client.RemoteBlockStorageVolumes.Collection()

	id := "rbs-volumes"

	if v, ok := d.GetOk("filter"); ok {
		filter := v.([]any)[0].(map[string]any)

		if id, ok := filter["location_id"]; ok && id.(int) > 0 {
			locationID := strconv.Itoa(id.(int))
			col = col.SetParam("location_id", locationID)
		}
		if ls, ok := filter["label_selector"]; ok && ls.(string) != "" {
			col = col.SetParam("label_selector", ls.(string))
		}
		if sp, ok := filter["search_pattern"]; ok && sp.(string) != "" {
			col = col.SetParam("search_pattern", sp.(string))
		}

		hash, err := hashFilter(filter)
		if err != nil {
			return diag.FromErr(err)
		}
		id = fmt.Sprintf("rbs-volumes-%s", hash)
	}

	vols, err := col.Collect(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving RBS volumes: %s", err))
	}

	list := make([]map[string]any, 0, len(vols))
	for _, vol := range vols {
		m := map[string]any{
			"id":            vol.ID,
			"name":          vol.Name,
			"size":          vol.Size,
			"status":        vol.Status,
			"labels":        vol.Labels,
			"location_id":   vol.LocationID,
			"location_code": vol.LocationCode,
			"flavor_id":     vol.FlavorID,
			"flavor_name":   vol.FlavorName,
			"iops":          vol.IOPS,
			"bandwidth":     vol.Bandwidth,
			"created_at":    vol.CreatedAt.String(),
			"updated_at":    vol.UpdatedAt.String(),
		}

		if vol.IPAddress != nil {
			m["ip_address"] = *vol.IPAddress
		}
		if vol.TargetIQN != nil {
			m["target_iqn"] = *vol.TargetIQN
		}

		list = append(list, m)
	}

	d.SetId(id)
	if err := d.Set("rbs_volumes", list); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rbs volumes: %s", err.Error()))
	}

	return nil
}
