package serverscom

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	serverscomL2SegmentDefaultCreateTimeout = 30 * time.Minute
	serverscomL2SegmentDefaultDeleteTimeout = 1 * time.Hour
	serverscomL2SegmentDefaultUpdateTimeout = 15 * time.Minute
)

func resourceServerscomL2Segment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServerscomL2SegmentRead,
		Update: resourceServerscomL2SegmentUpdate,
		Delete: resourceServerscomL2SegmentDelete,
		Create: resourceServerscomL2SegmentCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(serverscomL2SegmentDefaultCreateTimeout),
			Delete: schema.DefaultTimeout(serverscomL2SegmentDefaultDeleteTimeout),
			Update: schema.DefaultTimeout(serverscomL2SegmentDefaultUpdateTimeout),
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
			},
			"location_group": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.NoZeroValues,
				DiffSuppressFunc: compareStrings,
			},
			"member": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"native", "trunk"}, false),
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan": {
							Type:     schema.TypeString,
							Computed: true,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceServerscomL2SegmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)

	ctx := context.TODO()

	l2Segment, err := client.L2Segments.Get(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom l2 segment (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving l2 segment: %s", err)
		}
	}

	if l2Segment.Status == "removing" {
		log.Printf("[WARN] Serverscom l2 segment (%s) in removing status", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", l2Segment.Name)
	d.Set("location_group", l2Segment.LocationGroupCode)
	d.Set("type", l2Segment.Type)
	d.Set("status", l2Segment.Status)
	d.Set("created_at", l2Segment.Created.String())
	d.Set("updated_at", l2Segment.Updated.String())
	d.Set("labels", l2Segment.Labels)

	if l2Segment.Status != "active" {
		return nil
	}

	members, err := client.L2Segments.Members(d.Id()).Collect(ctx)
	if err != nil {
		return err
	}

	l2Members := getMembers(members)

	d.Set("member", l2Members)

	return nil
}

func resourceServerscomL2SegmentUpdate(d *schema.ResourceData, meta interface{}) error {
	input := scgo.L2SegmentUpdateInput{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		input.Name = &name
	}

	if d.HasChange("member") {
		input.Members = getSchemaMembers(d)
	}

	if d.HasChange("labels") {
		if labelsRaw, ok := d.GetOk("labels"); ok {
			labels := labelsRaw.(map[string]interface{})
			stringLabels := make(map[string]string)
			for k, v := range labels {
				stringLabels[k] = v.(string)
			}
			input.Labels = stringLabels
		}
	}

	if d.HasChanges("name", "member") {
		client := meta.(*scgo.Client)

		ctx := context.TODO()

		if _, err := waitForL2SegmentAttribute(d, "active", []string{"pending"}, "status", meta, schema.TimeoutUpdate); err != nil {
			return err
		}

		if _, err := client.L2Segments.Update(ctx, d.Id(), input); err != nil {
			return err
		}

		if _, err := waitForL2SegmentAttribute(d, "active", []string{"pending"}, "status", meta, schema.TimeoutUpdate); err != nil {
			return err
		}

		return resourceServerscomL2SegmentRead(d, meta)
	}

	return nil
}

func resourceServerscomL2SegmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)

	ctx := context.TODO()

	l2Segment, err := client.L2Segments.Get(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom l2 segment (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving l2 segment: %s", err)
		}
	}

	if l2Segment.Status == "removing" {
		log.Printf("[WARN] Serverscom l2 segment (%s) in removing status", d.Id())
		d.SetId("")
		return nil
	}

	if l2Segment.Status == "pending" {
		if _, err := waitForL2SegmentAttribute(d, "active", []string{"pending"}, "status", meta, schema.TimeoutDelete); err != nil {
			return err
		}
	}

	return client.L2Segments.Delete(ctx, d.Id())
}

func resourceServerscomL2SegmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)

	ctx := context.TODO()

	name := d.Get("name").(string)

	input := scgo.L2SegmentCreateInput{}
	input.Name = &name
	input.Type = d.Get("type").(string)

	locationGroup, err := getLocationGroup(input.Type, d.Get("location_group").(string))
	if err != nil {
		return err
	}

	input.LocationGroupID = locationGroup.ID
	input.Members = getSchemaMembers(d)

	if labelsRaw, ok := d.GetOk("labels"); ok {
		labels := labelsRaw.(map[string]interface{})
		stringLabels := make(map[string]string)
		for k, v := range labels {
			stringLabels[k] = v.(string)
		}
		input.Labels = stringLabels
	}

	l2Segment, err := client.L2Segments.Create(ctx, input)
	if err != nil {
		return err
	}

	d.SetId(l2Segment.ID)

	if _, err := waitForL2SegmentAttribute(d, "active", []string{"pending"}, "status", meta, schema.TimeoutCreate); err != nil {
		return err
	}

	return resourceServerscomL2SegmentRead(d, meta)
}

func getMembers(members []scgo.L2Member) []map[string]interface{} {
	l2Members := make([]map[string]interface{}, 0)

	for _, member := range members {
		var currentMember = make(map[string]interface{})

		currentMember["id"] = member.ID
		currentMember["mode"] = member.Mode
		currentMember["vlan"] = member.Vlan
		currentMember["status"] = member.Status
		currentMember["created_at"] = member.Created.String()
		currentMember["updated_at"] = member.Updated.String()

		l2Members = append(l2Members, currentMember)
	}

	return l2Members
}

func waitForL2SegmentAttribute(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, timeoutKey string) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for l2 segment (%s) to have %s of %s",
		d.Id(), attribute, target,
	)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newL2SegmentStateRefreshFunc(d, attribute, meta),
		Timeout:    d.Timeout(timeoutKey),
		Delay:      1 * time.Minute,
		MinTimeout: 15 * time.Second,
	}

	return stateConf.WaitForState()
}

func newL2SegmentStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	return func() (interface{}, string, error) {
		err := resourceServerscomL2SegmentRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		// See if we can access our attribute
		if attr, ok := d.GetOkExists(attribute); ok {
			l2Segment, err := client.L2Segments.Get(ctx, d.Id())

			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving l2 segment: %s", err)
			}

			switch attr.(type) {
			case bool:
				return &l2Segment, strconv.FormatBool(attr.(bool)), nil
			default:
				return &l2Segment, attr.(string), nil
			}
		}

		return nil, "", nil
	}
}

func getLocationGroup(groupType string, groupCode string) (*scgo.L2LocationGroup, error) {
	locationGroups, err := cache.LocationGroups()
	if err != nil {
		return nil, err
	}

	for _, locationGroup := range locationGroups {
		if locationGroup.GroupType == groupType && normalizeString(locationGroup.Code) == normalizeString(groupCode) {
			return &locationGroup, nil
		}
	}

	return nil, fmt.Errorf("Can't find location group by: %s (%s)", groupCode, groupType)
}

func getSchemaMembers(d *schema.ResourceData) []scgo.L2SegmentMemberInput {
	var members []scgo.L2SegmentMemberInput

	for _, memberMap := range d.Get("member").(*schema.Set).List() {
		member := memberMap.(map[string]interface{})

		currentMember := scgo.L2SegmentMemberInput{}
		currentMember.ID = member["id"].(string)
		currentMember.Mode = member["mode"].(string)

		members = append(members, currentMember)
	}

	return members
}
