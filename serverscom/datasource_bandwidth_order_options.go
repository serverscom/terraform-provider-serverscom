package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomBandwidthOrderOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomBandwidthOrderOptionsRead,

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
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This parameter filters output according to the bandwidth option's type.",
						},
					},
				},
			},
			"bandwidth_options": {
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
						"commit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServerscomBandwidthOrderOptionsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)
	uplinkModelID := d.Get("uplink_model_id").(int)

	collection := client.Locations.BandwidthOptions(int64(locationID), int64(serverModelID), int64(uplinkModelID))

	id := fmt.Sprintf("bandwidth_options-%d-%d-%d", locationID, serverModelID, uplinkModelID)
	if v, ok := d.GetOk("filter"); ok {
		filter := v.([]any)[0].(map[string]any)

		if bwType, ok := filter["type"]; ok && bwType.(string) != "" {
			collection = collection.SetParam("type", bwType.(string))
		}

		hash, err := hashFilter(filter)
		if err != nil {
			return err
		}
		id = fmt.Sprintf("%s-%s", id, hash)
	}

	bandwidthOptions, err := collection.Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving bandwidth order options: %s", err.Error())
	}

	optionList := make([]map[string]any, 0, len(bandwidthOptions))
	for _, bw := range bandwidthOptions {
		commit := 0
		if bw.Commit != nil {
			commit = int(*bw.Commit)
		}
		optionList = append(optionList, map[string]any{
			"id":     int(bw.ID),
			"name":   bw.Name,
			"type":   bw.Type,
			"commit": commit,
		})
	}

	d.SetId(id)
	if err := d.Set("bandwidth_options", optionList); err != nil {
		return fmt.Errorf("Error setting bandwidth order options: %s", err.Error())
	}

	return nil
}
