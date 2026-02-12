package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomRBSFlavor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerscomRBSFlavorRead,
		Schema: map[string]*schema.Schema{
			"location_id":      {Type: schema.TypeInt, Required: true},
			"id":               {Type: schema.TypeInt, Required: true},
			"name":             {Type: schema.TypeString, Computed: true},
			"iops_per_gb":      {Type: schema.TypeFloat, Computed: true},
			"bandwidth_per_gb": {Type: schema.TypeFloat, Computed: true},
			"min_size_gb":      {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceServerscomRBSFlavorRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	locationID := int64(d.Get("location_id").(int))
	flavorID := int64(d.Get("flavor_id").(int))

	flavor, err := client.Locations.GetRemoteBlockStorageFlavor(ctx, locationID, flavorID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving RBS flavor order option: %s", err.Error()))
	}

	d.SetId(strconv.Itoa(int(flavor.ID)))
	d.Set("name", flavor.Name)
	d.Set("iops_per_gb", flavor.IOPSPerGB)
	d.Set("bandwidth_per_gb", flavor.BandwidthPerGB)
	d.Set("min_size_gb", flavor.MinSizeGB)

	return nil
}
