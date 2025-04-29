package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomSbmFlavorOrderOption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomSbmFlavorOrderOptionRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_cores_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_frequency": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"drives_configuration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_uplink_model_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_uplink_model_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_uplink_model_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_uplink_model_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServerscomSbmFlavorOrderOptionRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	sbmFlavorID := d.Get("id").(int)

	flavor, err := client.Locations.GetSBMFlavorOption(ctx, int64(locationID), int64(sbmFlavorID))
	if err != nil {
		return fmt.Errorf("Error retrieving SBM flavor order option: %s", err.Error())
	}

	d.SetId(strconv.Itoa(int(flavor.ID)))
	d.Set("name", flavor.Name)
	d.Set("cpu_name", flavor.CPUName)
	d.Set("cpu_count", flavor.CPUCount)
	d.Set("cpu_cores_count", flavor.CPUCoresCount)
	d.Set("cpu_frequency", flavor.CPUFrequency)
	d.Set("ram_size", flavor.RAMSize)
	d.Set("drives_configuration", flavor.DrivesConfiguration)
	d.Set("public_uplink_model_id", flavor.PublicUplinkModelID)
	d.Set("public_uplink_model_name", flavor.PublicUplinkModelName)
	d.Set("private_uplink_model_id", flavor.PrivateUplinkModelID)
	d.Set("private_uplink_model_name", flavor.PrivateUplinkModelName)
	d.Set("bandwidth_id", flavor.BandwidthID)
	d.Set("bandwidth_name", flavor.BandwidthName)

	return nil
}
