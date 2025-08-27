package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomOperatingSystemOrderOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomOperatingSystemOrderOptionsRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server_model_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"operating_systems": {
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

func dataSourceServerscomOperatingSystemOrderOptionsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)

	collection := client.Locations.OperatingSystemOptions(int64(locationID), int64(serverModelID))

	operatingSystems, err := collection.Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving operating system order options: %s", err.Error())
	}

	osList := make([]map[string]any, 0, len(operatingSystems))
	for _, os := range operatingSystems {
		osList = append(osList, map[string]any{
			"id":          int(os.ID),
			"full_name":   os.FullName,
			"name":        os.Name,
			"version":     os.Version,
			"arch":        os.Arch,
			"filesystems": os.Filesystems,
		})
	}

	d.SetId(fmt.Sprintf("operating_systems-%d-%d", locationID, serverModelID))
	if err := d.Set("operating_systems", osList); err != nil {
		return fmt.Errorf("Error setting operating system order options: %s", err.Error())
	}

	return nil
}
