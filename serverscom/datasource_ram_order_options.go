package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomRamOrderOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomRamOrderOptionsRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server_model_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ram_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServerscomRamOrderOptionsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)

	options, err := client.Locations.RAMOptions(int64(locationID), int64(serverModelID)).Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving RAM order options: %s", err.Error())
	}

	optionList := make([]map[string]any, 0, len(options))
	for _, ram := range options {
		optionList = append(optionList, map[string]any{
			"ram":  ram.RAM,
			"type": ram.Type,
		})
	}

	id := fmt.Sprintf("ram_options-%d-%d", locationID, serverModelID)
	d.SetId(id)

	if err := d.Set("ram_options", optionList); err != nil {
		return fmt.Errorf("Error setting RAM order options: %s", err.Error())
	}

	return nil
}
