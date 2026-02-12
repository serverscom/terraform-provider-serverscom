package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomRBSVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerscomRBSVolumeRead,

		Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Required: true},

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
	}
}

func dataSourceServerscomRBSVolumeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	id := d.Get("id").(string)

	vol, err := client.RemoteBlockStorageVolumes.Get(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving RBS volume: %s", err))
	}

	d.SetId(vol.ID)
	d.Set("name", vol.Name)
	d.Set("size", vol.Size)
	d.Set("status", vol.Status)
	d.Set("labels", vol.Labels)
	d.Set("location_id", vol.LocationID)
	d.Set("location_code", vol.LocationCode)

	if vol.IPAddress != nil {
		d.Set("ip_address", *vol.IPAddress)
	}

	d.Set("flavor_id", vol.FlavorID)
	d.Set("flavor_name", vol.FlavorName)
	d.Set("iops", vol.IOPS)
	d.Set("bandwidth", vol.Bandwidth)

	if vol.TargetIQN != nil {
		d.Set("target_iqn", *vol.TargetIQN)
	}

	d.Set("created_at", vol.CreatedAt.String())
	d.Set("updated_at", vol.UpdatedAt.String())

	return nil
}
