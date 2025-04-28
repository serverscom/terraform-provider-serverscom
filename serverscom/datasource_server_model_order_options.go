package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomServerModelOrderOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomServerModelOrderOptionsRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This pattern is used to return resources containing the parameter value in its name.",
						},
						"has_raid_controller": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Specify by true or false value, if only servers with RAID controller should be taken to an output.",
						},
					},
				},
			},
			"server_models": {
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
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"has_raid_controller": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"raid_controller_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"drive_slots_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServerscomServerModelOrderOptionsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)

	collection := client.Locations.ServerModelOptions(int64(locationID))

	id := fmt.Sprintf("server_models-%d", locationID)
	if v, ok := d.GetOk("filter"); ok {
		filter := v.([]any)[0].(map[string]any)

		if searchPattern, ok := filter["search_pattern"]; ok && searchPattern.(string) != "" {
			collection = collection.SetParam("search_pattern", searchPattern.(string))
		}
		if hasRaidController, ok := filter["has_raid_controller"]; ok {
			collection = collection.SetParam("has_raid_controller", fmt.Sprintf("%v", hasRaidController.(bool)))
		}

		hash, err := hashFilter(filter)
		if err != nil {
			return err
		}
		id = fmt.Sprintf("%s-%s", id, hash)
	}

	options, err := collection.Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving server model order options: %s", err.Error())
	}

	optionList := make([]map[string]any, 0, len(options))
	for _, option := range options {
		optionMap := map[string]any{
			"id":                   int(option.ID),
			"name":                 option.Name,
			"cpu_name":             option.CPUName,
			"cpu_count":            option.CPUCount,
			"cpu_cores_count":      option.CPUCoresCount,
			"cpu_frequency":        option.CPUFrequency,
			"ram":                  option.RAM,
			"ram_type":             option.RAMType,
			"max_ram":              option.MaxRAM,
			"has_raid_controller":  option.HasRAIDController,
			"raid_controller_name": option.RAIDControllerName,
			"drive_slots_count":    option.DriveSlotsCount,
		}
		optionList = append(optionList, optionMap)
	}

	d.SetId(id)
	if err := d.Set("server_models", optionList); err != nil {
		return fmt.Errorf("Error setting server model order options: %s", err.Error())
	}

	return nil
}
