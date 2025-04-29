package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomL2Segment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomL2SegmentRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"location_group_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServerscomL2SegmentRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	L2SegmentID := d.Get("id").(string)

	L2Segment, err := client.L2Segments.Get(ctx, L2SegmentID)
	if err != nil {
		return fmt.Errorf("Error retrieving L2 segment: %s", err.Error())
	}

	d.SetId(L2Segment.ID)
	d.Set("name", L2Segment.Name)
	d.Set("type", L2Segment.Type)
	d.Set("status", L2Segment.Status)
	d.Set("location_group_id", L2Segment.LocationGroupID)
	d.Set("location_group_code", L2Segment.LocationGroupCode)
	d.Set("labels", L2Segment.Labels)
	d.Set("created_at", L2Segment.Created)
	d.Set("updated_at", L2Segment.Updated)

	return nil
}
