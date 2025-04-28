package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomSbmOperatingSystemOrderOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomSbmOperatingSystemOrderOptionsRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"sbm_flavor_model_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"sbm_operating_systems": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceServerscomSbmOperatingSystemOrderOptionsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	sbmFlavorModelID := d.Get("sbm_flavor_model_id").(int)

	collection := client.Locations.SBMOperatingSystemOptions(int64(locationID), int64(sbmFlavorModelID))

	options, err := collection.Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving SBM operating system order options: %s", err.Error())
	}

	optionList := make([]map[string]any, 0, len(options))
	for _, os := range options {
		optionList = append(optionList, map[string]any{
			"id":          int(os.ID),
			"full_name":   os.FullName,
			"name":        os.Name,
			"version":     os.Version,
			"arch":        os.Arch,
			"filesystems": os.Filesystems,
		})
	}

	id := fmt.Sprintf("sbm_operating_systems-%d-%d", locationID, sbmFlavorModelID)
	d.SetId(id)

	if err := d.Set("sbm_operating_systems", optionList); err != nil {
		return fmt.Errorf("Error setting SBM operating system order options: %s", err.Error())
	}

	return nil
}
