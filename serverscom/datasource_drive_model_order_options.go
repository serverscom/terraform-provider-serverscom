package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomDriveModelOrderOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomDriveModelOrderOptionsRead,

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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Return drives containing this value in their name.",
						},
						"media_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter by media type: HDD or SSD.",
						},
						"interface": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter by interface: SATA1, SATA2, SATA3, SAS, NVMe-PCIe.",
						},
					},
				},
			},
			"drive_models": {
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
				},
			},
		},
	}
}

func dataSourceServerscomDriveModelOrderOptionsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	locationID := d.Get("location_id").(int)
	serverModelID := d.Get("server_model_id").(int)

	collection := client.Locations.DriveModelOptions(int64(locationID), int64(serverModelID))

	id := fmt.Sprintf("drive_models-%d-%d", locationID, serverModelID)

	if v, ok := d.GetOk("filter"); ok {
		filter := v.([]any)[0].(map[string]any)

		if searchPattern, ok := filter["search_pattern"]; ok && searchPattern.(string) != "" {
			collection = collection.SetParam("search_pattern", searchPattern.(string))
		}
		if mediaType, ok := filter["media_type"]; ok && mediaType.(string) != "" {
			collection = collection.SetParam("media_type", mediaType.(string))
		}
		if iface, ok := filter["interface"]; ok && iface.(string) != "" {
			collection = collection.SetParam("interface", iface.(string))
		}

		hash, err := hashFilter(filter)
		if err != nil {
			return err
		}
		id = fmt.Sprintf("%s-%s", id, hash)
	}

	driveModels, err := collection.Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving drive model order options: %s", err.Error())
	}

	driveModelList := make([]map[string]any, 0, len(driveModels))
	for _, model := range driveModels {
		driveModelList = append(driveModelList, map[string]any{
			"id":          int(model.ID),
			"name":        model.Name,
			"capacity":    model.Capacity,
			"interface":   model.Interface,
			"form_factor": model.FormFactor,
			"media_type":  model.MediaType,
		})
	}

	d.SetId(id)
	if err := d.Set("drive_models", driveModelList); err != nil {
		return fmt.Errorf("Error setting drive model order options: %s", err.Error())
	}

	return nil
}
