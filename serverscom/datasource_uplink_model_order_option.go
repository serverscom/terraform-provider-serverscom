package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomUplinkModelOrderOption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomUplinkModelOrderOptionRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server_model_id": {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"speed": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"redundancy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceServerscomUplinkModelOrderOptionRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)
	uplinkModelID := d.Get("id").(int)

	uplink, err := client.Locations.GetUplinkOption(ctx, int64(locationID), int64(serverModelID), int64(uplinkModelID))
	if err != nil {
		return fmt.Errorf("Error retrieving uplink model order option: %s", err.Error())
	}

	d.SetId(strconv.Itoa(int(uplink.ID)))
	d.Set("name", uplink.Name)
	d.Set("type", uplink.Type)
	d.Set("speed", uplink.Speed)
	d.Set("redundancy", uplink.Redundancy)

	return nil
}
