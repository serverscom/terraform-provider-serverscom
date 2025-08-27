package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomLocationsRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This pattern is used to return resources containing the parameter value in its name.",
						},
					},
				},
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceServerscomLocationsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	collection := client.Locations.Collection()

	id := "locations"
	if v, ok := d.GetOk("filter"); ok {
		filter := v.([]any)[0].(map[string]any)

		if searchPattern, ok := filter["search_pattern"]; ok && searchPattern.(string) != "" {
			collection = collection.SetParam("search_pattern", searchPattern.(string))
		}
		hash, err := hashFilter(filter)
		if err != nil {
			return err
		}
		id = fmt.Sprintf("locations-%s", hash)
	}

	locations, err := collection.Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving locations: %s", err.Error())
	}

	locationList := make([]map[string]any, 0, len(locations))
	for _, location := range locations {
		locationMap := map[string]any{
			"id":                     int(location.ID),
			"name":                   location.Name,
			"status":                 location.Status,
			"code":                   location.Code,
			"supported_features":     location.SupportedFeatures,
			"l2_segments_enabled":    location.L2SegmentsEnabled,
			"private_racks_enabled":  location.PrivateRacksEnabled,
			"load_balancers_enabled": location.LoadBalancersEnabled,
		}
		locationList = append(locationList, locationMap)
	}

	d.SetId(id)
	if err := d.Set("locations", locationList); err != nil {
		return fmt.Errorf("Error setting locations: %s", err.Error())
	}

	return nil
}
