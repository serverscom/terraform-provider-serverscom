package serverscom

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomServerModelOrderOption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomServerModelOrderOptionRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
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
			"drive_slots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"position": {
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
						"drive_model_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hot_swappable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServerscomServerModelOrderOptionRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)

	option, err := client.Locations.GetServerModelOption(ctx, int64(locationID), int64(serverModelID))
	if err != nil {
		return fmt.Errorf("Error retrieving server model order option: %s", err.Error())
	}

	d.SetId(strconv.Itoa(int(option.ID)))
	d.Set("name", option.Name)
	d.Set("cpu_name", option.CPUName)
	d.Set("cpu_count", option.CPUCount)
	d.Set("cpu_cores_count", option.CPUCoresCount)
	d.Set("cpu_frequency", option.CPUFrequency)
	d.Set("ram", option.RAM)
	d.Set("ram_type", option.RAMType)
	d.Set("max_ram", option.MaxRAM)
	d.Set("has_raid_controller", option.HasRAIDController)
	d.Set("raid_controller_name", option.RAIDControllerName)
	d.Set("drive_slots_count", option.DriveSlotsCount)

	var driveSlots []map[string]interface{}
	for _, slot := range option.DriveSlots {
		driveSlot := map[string]interface{}{
			"position":       slot.Position,
			"interface":      slot.Interface,
			"form_factor":    slot.FormFactor,
			"drive_model_id": slot.DriveModelID,
			"hot_swappable":  slot.HotSwappable,
		}
		driveSlots = append(driveSlots, driveSlot)
	}
	if err := d.Set("drive_slots", driveSlots); err != nil {
		return fmt.Errorf("Error setting drive_slots: %s", err.Error())
	}

	return nil
}
