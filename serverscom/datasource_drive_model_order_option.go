package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomDriveModelOrderOption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomDriveModelOrderOptionRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server_model_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"drive_model_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"interface": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"form_factor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"media_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServerscomDriveModelOrderOptionRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)
	driveModelID := d.Get("drive_model_id").(int)

	model, err := client.Locations.GetDriveModelOption(ctx, int64(locationID), int64(serverModelID), int64(driveModelID))
	if err != nil {
		return fmt.Errorf("Error retrieving drive model order option: %s", err.Error())
	}

	d.SetId(strconv.Itoa(int(model.ID)))
	d.Set("name", model.Name)
	d.Set("capacity", model.Capacity)
	d.Set("interface", model.Interface)
	d.Set("form_factor", model.FormFactor)
	d.Set("media_type", model.MediaType)

	return nil
}
