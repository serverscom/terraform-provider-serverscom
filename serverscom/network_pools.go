package serverscom

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func networkPoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "id of the Network Pool",
		},
		"title": {
			Type:        schema.TypeString,
			Description: "title of the Network Pool",
		},
		"type": {
			Type:        schema.TypeString,
			Description: "the type of the Network Pool",
		},
		"cidr": {
			Type:        schema.TypeString,
			Description: "CIDR of the Network Pool",
		},
		"location_ids": {
			Type:        schema.TypeSet,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Description: "list of location id's of the Network Pool",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "the creation time of the Network Pool",
		},
	}
}

func flattenServerscomNetworkPool(rawNetworkPool, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	networkPool := rawNetworkPool.(*scgo.NetworkPool)

	flattenNetworkPool := map[string]interface{}{
		"id":         networkPool.ID,
		"title":      networkPool.Title,
		"type":       networkPool.Type,
		"cidr":       networkPool.CIDR,
		"created_at": networkPool.Created.String(),
	}

	locationIds := make([]int, len(networkPool.LocationIDs))
	for i, locationId := range networkPool.LocationIDs {
		locationIds[i] = locationId
	}

	flattenNetworkPool["location_ids"] = locationIds

	return flattenNetworkPool, nil
}
