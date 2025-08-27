package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomSbmFlavorOrderOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomSbmFlavorOrderOptionsRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
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
						"show_all": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If true, all flavors including unavailable ones will be shown.",
						},
					},
				},
			},
			"sbm_flavors": {
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
						"cpu_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_cores_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_frequency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"drives_configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_uplink_model_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"public_uplink_model_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_uplink_model_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_uplink_model_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServerscomSbmFlavorOrderOptionsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)

	collection := client.Locations.SBMFlavorOptions(int64(locationID))

	id := fmt.Sprintf("sbm_flavors-%d", locationID)

	if v, ok := d.GetOk("filter"); ok {
		filter := v.([]any)[0].(map[string]any)

		if searchPattern, ok := filter["search_pattern"]; ok && searchPattern.(string) != "" {
			collection = collection.SetParam("search_pattern", searchPattern.(string))
		}
		if showAll, ok := filter["show_all"]; ok {
			collection = collection.SetParam("show_all", strconv.FormatBool(showAll.(bool)))
		}

		hash, err := hashFilter(filter)
		if err != nil {
			return err
		}
		id = fmt.Sprintf("%s-%s", id, hash)
	}

	options, err := collection.Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving SBM flavor order options: %s", err.Error())
	}

	optionList := make([]map[string]any, 0, len(options))
	for _, flavor := range options {
		optionList = append(optionList, map[string]any{
			"id":                        int(flavor.ID),
			"name":                      flavor.Name,
			"cpu_name":                  flavor.CPUName,
			"cpu_count":                 flavor.CPUCount,
			"cpu_cores_count":           flavor.CPUCoresCount,
			"cpu_frequency":             flavor.CPUFrequency,
			"ram_size":                  flavor.RAMSize,
			"drives_configuration":      flavor.DrivesConfiguration,
			"public_uplink_model_id":    flavor.PublicUplinkModelID,
			"public_uplink_model_name":  flavor.PublicUplinkModelName,
			"private_uplink_model_id":   flavor.PrivateUplinkModelID,
			"private_uplink_model_name": flavor.PrivateUplinkModelName,
			"bandwidth_id":              flavor.BandwidthID,
			"bandwidth_name":            flavor.BandwidthName,
		})
	}

	d.SetId(id)

	if err := d.Set("sbm_flavors", optionList); err != nil {
		return fmt.Errorf("Error setting SBM flavor order options: %s", err.Error())
	}

	return nil
}
