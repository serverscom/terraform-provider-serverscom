package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomSbmOperatingSystemOrderOption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomSbmOperatingSystemOrderOptionRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"sbm_flavor_model_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filesystems": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceServerscomSbmOperatingSystemOrderOptionRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	sbmFlavorModelID := d.Get("sbm_flavor_model_id").(int)
	osID := d.Get("id").(int)

	os, err := client.Locations.GetSBMOperatingSystemOption(ctx, int64(locationID), int64(sbmFlavorModelID), int64(osID))
	if err != nil {
		return fmt.Errorf("Error retrieving SBM operating system order option: %s", err.Error())
	}

	d.SetId(strconv.Itoa(int(os.ID)))
	d.Set("full_name", os.FullName)
	d.Set("name", os.Name)
	d.Set("version", os.Version)
	d.Set("arch", os.Arch)
	d.Set("filesystems", os.Filesystems)

	return nil
}
