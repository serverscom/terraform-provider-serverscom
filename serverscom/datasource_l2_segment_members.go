package serverscom

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomL2SegmentMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomL2SegmentMembersRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the L2 segment to get members for",
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
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
				},
			},
		},
	}
}

func dataSourceServerscomL2SegmentMembersRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	L2SegmentID := d.Get("id").(string)

	members, err := client.L2Segments.Members(L2SegmentID).Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving L2 members: %s", err.Error())
	}

	d.SetId(L2SegmentID)

	membersList := make([]map[string]any, len(members))
	for i, member := range members {
		memberMap := map[string]any{
			"id":         member.ID,
			"title":      member.Title,
			"mode":       member.Mode,
			"vlan":       member.Vlan,
			"status":     member.Status,
			"labels":     member.Labels,
			"created_at": member.Created.Format(time.RFC3339),
			"updated_at": member.Updated.Format(time.RFC3339),
		}
		membersList[i] = memberMap
	}

	if err := d.Set("members", membersList); err != nil {
		return fmt.Errorf("Error setting members: %s", err.Error())
	}

	return nil
}
