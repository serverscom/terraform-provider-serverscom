package serverscom

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	rbsDefaultCreateTimeout = 1 * time.Hour
	rbsDefaultDeleteTimeout = 30 * time.Minute
)

func resourceServerscomRBSVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerscomRBSVolumeCreate,
		ReadContext:   resourceServerscomRBSVolumeRead,
		UpdateContext: resourceServerscomRBSVolumeUpdate,
		DeleteContext: resourceServerscomRBSVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(rbsDefaultCreateTimeout),
			Delete: schema.DefaultTimeout(rbsDefaultDeleteTimeout),
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"status":        {Type: schema.TypeString, Computed: true},
			"ip_address":    {Type: schema.TypeString, Computed: true},
			"target_iqn":    {Type: schema.TypeString, Computed: true},
			"location_code": {Type: schema.TypeString, Computed: true},
			"flavor_name":   {Type: schema.TypeString, Computed: true},
			"iops":          {Type: schema.TypeFloat, Computed: true},
			"bandwidth":     {Type: schema.TypeFloat, Computed: true},
			"created_at":    {Type: schema.TypeString, Computed: true},
			"updated_at":    {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceServerscomRBSVolumeCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	input := scgo.RemoteBlockStorageVolumeCreateInput{
		Name:       d.Get("name").(string),
		Size:       int64(d.Get("size").(int)),
		LocationID: d.Get("location_id").(int),
		FlavorID:   d.Get("flavor_id").(int),
	}
	if labelsRaw, ok := d.GetOk("labels"); ok {
		input.Labels = toStringMap(labelsRaw.(map[string]interface{}))
	}

	vol, err := client.RemoteBlockStorageVolumes.Create(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vol.ID)

	if vol.Status == "creating" || vol.Status == "pending" {
		pending := []string{"creating", "pending"}
		target := []string{"active"}
		_, err := waitForRBSVolumeAttribute(ctx, d, target, pending, "status", meta, "create")
		if err != nil {
			return diag.FromErr(fmt.Errorf("waiting for volume to become active: %w", err))
		}
	}

	return resourceServerscomRBSVolumeRead(ctx, d, meta)
}

func resourceServerscomRBSVolumeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	vol, err := client.RemoteBlockStorageVolumes.Get(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*scgo.NotFoundError); ok {
			log.Printf("[WARN] Serverscom rbs volume (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", vol.Name)
	d.Set("size", int(vol.Size))
	d.Set("location_id", vol.LocationID)
	d.Set("flavor_id", vol.FlavorID)
	d.Set("labels", vol.Labels)

	d.Set("status", vol.Status)
	if vol.IPAddress != nil {
		d.Set("ip_address", *vol.IPAddress)
	} else {
		d.Set("ip_address", nil)
	}
	if vol.TargetIQN != nil {
		d.Set("target_iqn", *vol.TargetIQN)
	} else {
		d.Set("target_iqn", nil)
	}
	d.Set("location_code", vol.LocationCode)
	d.Set("flavor_name", vol.FlavorName)
	d.Set("iops", vol.IOPS)
	d.Set("bandwidth", vol.Bandwidth)
	d.Set("created_at", vol.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", vol.UpdatedAt.Format(time.RFC3339))

	return nil
}

func resourceServerscomRBSVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	input := scgo.RemoteBlockStorageVolumeUpdateInput{}
	changesSize := false

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("labels") {
		if labelsRaw, ok := d.GetOk("labels"); ok {
			input.Labels = toStringMap(labelsRaw.(map[string]interface{}))
		} else {
			input.Labels = map[string]string{}
		}
	}
	if d.HasChange("size") {
		newSize := int64(d.Get("size").(int))
		input.Size = newSize
		changesSize = true
	}

	vol, err := client.RemoteBlockStorageVolumes.Update(ctx, d.Id(), input)
	if err != nil {
		return diag.FromErr(err)
	}

	if changesSize && (vol.Status == "pending") {
		pending := []string{"pending"}
		target := []string{"active"}
		_, err := waitForRBSVolumeAttribute(ctx, d, target, pending, "status", meta, "update")
		if err != nil {
			return diag.FromErr(fmt.Errorf("waiting for volume resize to complete: %w", err))
		}
	}

	return resourceServerscomRBSVolumeRead(ctx, d, meta)
}

func resourceServerscomRBSVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scgo.Client)

	err := client.RemoteBlockStorageVolumes.Delete(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*scgo.NotFoundError); ok {
			log.Printf("[WARN] Serverscom rbs volume (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	pending := []string{"pending"}
	target := []string{"removing", "removed"} // removed - state when read returns 404
	_, err = waitForRBSVolumeAttribute(ctx, d, target, pending, "status", meta, "delete")
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error waiting for rbs volume (%s) to become deleted: %s", d.Id(), err))
	}
	d.SetId("")
	return nil
}

func toStringMap(m map[string]interface{}) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v.(string)
	}
	return out
}

func waitForRBSVolumeAttribute(ctx context.Context, d *schema.ResourceData, target []string, pending []string, attribute string, meta any, timeoutKey string) (any, error) {
	client := meta.(*scgo.Client)

	log.Printf("[INFO] Waiting for rbs volume (%s) attribute %s -> %s", d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:      pending,
		Target:       target,
		Refresh:      newRBSVolumeStateRefreshFunc(ctx, d, attribute, client),
		Timeout:      d.Timeout(timeoutKey),
		PollInterval: 15 * time.Second,
		Delay:        15 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newRBSVolumeStateRefreshFunc(ctx context.Context, d *schema.ResourceData, attribute string, meta interface{}) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		diags := resourceServerscomRBSVolumeRead(ctx, d, meta)

		if diags.HasError() {
			return nil, "", errors.New(diags[0].Summary)
		}

		if d.Id() == "" {
			return nil, "removed", nil
		}

		if attr, ok := d.GetOk(attribute); ok {
			switch v := attr.(type) {
			case bool:
				return d, strconv.FormatBool(v), nil
			case string:
				return d, v, nil
			default:
				return d, fmt.Sprintf("%v", v), nil
			}
		}

		return nil, "", nil
	}
}
