package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomUplinkModelOrderOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomUplinkModelOrderOptionsRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server_model_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redundancy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The true value indicates an uplink having redundancy.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This type defines uplinks with private or public connection.",
						},
					},
				},
			},
			"uplink_models": {
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
				},
			},
		},
	}
}

func dataSourceServerscomUplinkModelOrderOptionsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)

	collection := client.Locations.UplinkOptions(int64(locationID), int64(serverModelID))

	id := fmt.Sprintf("uplink_models-%d-%d", locationID, serverModelID)
	if v, ok := d.GetOk("filter"); ok {
		filter := v.([]any)[0].(map[string]any)

		if redundancy, ok := filter["redundancy"]; ok {
			collection = collection.SetParam("redundancy", strconv.FormatBool(redundancy.(bool)))
		}
		if uplinkType, ok := filter["type"]; ok && uplinkType.(string) != "" {
			collection = collection.SetParam("type", uplinkType.(string))
		}

		hash, err := hashFilter(filter)
		if err != nil {
			return err
		}
		id = fmt.Sprintf("%s-%s", id, hash)
	}

	uplinkModels, err := collection.Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving uplink model order options: %s", err.Error())
	}

	uplinkList := make([]map[string]any, 0, len(uplinkModels))
	for _, uplink := range uplinkModels {
		uplinkList = append(uplinkList, map[string]any{
			"id":         int(uplink.ID),
			"name":       uplink.Name,
			"type":       uplink.Type,
			"speed":      uplink.Speed,
			"redundancy": uplink.Redundancy,
		})
	}

	d.SetId(id)
	if err := d.Set("uplink_models", uplinkList); err != nil {
		return fmt.Errorf("Error setting uplink model order options: %s", err.Error())
	}

	return nil
}
