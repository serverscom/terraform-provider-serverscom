package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomLocationRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supported_features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"l2_segments_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"private_racks_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"load_balancers_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceServerscomLocationRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("id").(int)

	location, err := client.Locations.GetLocation(ctx, int64(locationID))
	if err != nil {
		return fmt.Errorf("Error retrieving location: %s", err.Error())
	}

	d.SetId(strconv.Itoa(int(location.ID)))
	d.Set("name", location.Name)
	d.Set("status", location.Status)
	d.Set("code", location.Code)
	d.Set("supported_features", location.SupportedFeatures)
	d.Set("l2_segments_enabled", location.L2SegmentsEnabled)
	d.Set("private_racks_enabled", location.PrivateRacksEnabled)
	d.Set("load_balancers_enabled", location.LoadBalancersEnabled)

	return nil
}
