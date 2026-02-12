package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomRBSFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerscomRBSFlavorsRead,

		Schema: map[string]*schema.Schema{
			"location_id": {Type: schema.TypeInt, Required: true},

			"rbs_flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":               {Type: schema.TypeInt, Computed: true},
						"name":             {Type: schema.TypeString, Computed: true},
						"iops_per_gb":      {Type: schema.TypeFloat, Computed: true},
						"bandwidth_per_gb": {Type: schema.TypeFloat, Computed: true},
						"min_size_gb":      {Type: schema.TypeInt, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceServerscomRBSFlavorsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	locationID := d.Get("location_id").(int)

	collection := client.Locations.RemoteBlockStorageFlavors(int64(locationID))

	flavors, err := collection.Collect(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving RBS flavor order options: %s", err.Error()))
	}

	list := make([]map[string]any, 0, len(flavors))
	for _, f := range flavors {
		list = append(list, map[string]any{
			"id":               f.ID,
			"name":             f.Name,
			"iops_per_gb":      f.IOPSPerGB,
			"bandwidth_per_gb": f.BandwidthPerGB,
			"min_size_gb":      f.MinSizeGB,
		})
	}

	id := fmt.Sprintf("rbs-flavors-%d", locationID)
	d.SetId(id)

	if err := d.Set("rbs_flavors", list); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting RBS flavor order options: %s", err.Error()))
	}

	return nil
}
