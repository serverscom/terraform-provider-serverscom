package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomBandwidthOrderOption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomBandwidthOrderOptionRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server_model_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"uplink_model_id": {
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
			"commit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceServerscomBandwidthOrderOptionRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)
	uplinkModelID := d.Get("uplink_model_id").(int)
	bandwidthID := d.Get("id").(int)

	bw, err := client.Locations.GetBandwidthOption(ctx, int64(locationID), int64(serverModelID), int64(uplinkModelID), int64(bandwidthID))
	if err != nil {
		return fmt.Errorf("Error retrieving bandwidth order option: %s", err.Error())
	}

	d.SetId(strconv.Itoa(int(bw.ID)))
	d.Set("name", bw.Name)
	d.Set("type", bw.Type)

	if bw.Commit != nil {
		d.Set("commit", int(*bw.Commit))
	} else {
		d.Set("commit", 0)
	}

	return nil
}
