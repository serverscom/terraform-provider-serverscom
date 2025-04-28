package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomRamOrderOption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomRamOrderOptionRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server_model_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ram": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "RAM size in GB to match",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "RAM type to match (e.g., DDR4)",
			},
		},
	}
}

func dataSourceServerscomRamOrderOptionRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)
	id := fmt.Sprintf("drive_models-%d-%d", locationID, serverModelID)

	targetRAM := d.Get("ram").(int)
	targetType := d.Get("type").(string)

	options, err := client.Locations.RAMOptions(int64(locationID), int64(serverModelID)).Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving RAM order options: %s", err.Error())
	}

	for _, ram := range options {
		if ram.RAM == targetRAM && ram.Type == targetType {
			id = fmt.Sprintf("%s-%d-%s", id, ram.RAM, ram.Type)
			d.SetId(id)
			d.Set("ram", ram.RAM)
			d.Set("type", ram.Type)
			return nil
		}
	}

	return fmt.Errorf("No RAM option found matching ram=%d and type=%s", targetRAM, targetType)
}
