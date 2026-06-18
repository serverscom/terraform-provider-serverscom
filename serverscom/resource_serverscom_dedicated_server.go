package serverscom

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	serverscomDedicatedServerDefaultCreateTimeout = 24 * time.Hour
	serverscomDedicatedServerDefaultUpdateTimeout = 4 * time.Hour
	serverscomDedicatedServerDefaultDeleteTimeout = 1 * time.Hour

	// dedicatedServerReinstallTriggers are the fields that can only be applied
	// through an OS reinstall (which destroys all data on disk). Changing any of
	// them requires bumping reinstall_trigger to acknowledge the reinstall.
	// ssh_key_fingerprints and user_data are intentionally NOT here: they are
	// applied in-place via the API (and also passed into the reinstall payload).
	dedicatedServerReinstallTriggers = []string{"operating_system", "layout"}
)

func resourceServerscomDedicatedServer() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceServerscomDedicatedServerRead,
		UpdateContext: resourceServerscomDedicatedServerUpdate,
		DeleteContext: resourceServerscomDedicatedServerDelete,
		CreateContext: resourceServerscomDedicatedServerCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(serverscomDedicatedServerDefaultCreateTimeout),
			Update: schema.DefaultTimeout(serverscomDedicatedServerDefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(serverscomDedicatedServerDefaultDeleteTimeout),
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"location": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareStrings,
				ValidateFunc:     validation.NoZeroValues,
			},
			"server_model": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareStrings,
				ValidateFunc:     validation.NoZeroValues,
			},
			"ram_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"operating_system": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareStrings,
			},
			"public_uplink": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_uplink": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareStrings,
				ValidateFunc:     validation.NoZeroValues,
			},
			"bandwidth": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareStrings,
			},
			"slot": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"position": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"drive_model": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: compareStrings,
						},
					},
				},
			},
			"layout": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slot_positions": {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"raid": {
							Type:         schema.TypeInt,
							Default:      0,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{0, 1, 5, 6, 10, 50, 60}),
						},
						"partition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": {
										Type:     schema.TypeString,
										Required: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"fill": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"fs": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"ssh_key_fingerprints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"user_data": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},
			"configuration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ipv4_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"public_ipv4_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operational_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reinstall_trigger": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  reinstallTriggerNone,
				Description: "Changing this to any value other than \"" + reinstallTriggerNone + "\" performs an OS reinstall, " +
					"which DESTROYS ALL DATA on disk. A reinstall is required to apply changes to operating_system, " +
					"layout or ssh_key_fingerprints; changing one of those without also changing reinstall_trigger fails the plan.",
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},

		CustomizeDiff: resourceServerscomDedicatedServerCustomizeDiff,
	}
}

// dedicatedServerOperationalStatusNormal is the operational_status a dedicated
// server reports when it is not undergoing any operation (e.g. a reinstall).
const dedicatedServerOperationalStatusNormal = "normal"

// reinstallTriggerNone is the sentinel value of reinstall_trigger that means
// "no reinstall". A reinstall is performed only when reinstall_trigger changes
// to some other value.
const reinstallTriggerNone = "none"

// dedicatedServerWillReinstall reports whether the planned change to
// reinstall_trigger should perform an OS reinstall. A reinstall happens only
// when the value actually changes to something other than the "none" sentinel,
// so first-time adoption (the field appearing with its default) is a no-op.
func dedicatedServerWillReinstall(oldTrigger, newTrigger string) bool {
	return oldTrigger != newTrigger && newTrigger != reinstallTriggerNone
}

func resourceServerscomDedicatedServerCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	// On create there is no reinstall: the OS is installed by the create flow.
	if d.Id() == "" {
		return nil
	}

	var changedTriggers []string
	for _, field := range dedicatedServerReinstallTriggers {
		if d.HasChange(field) {
			changedTriggers = append(changedTriggers, field)
		}
	}

	oldTrigger, newTrigger := d.GetChange("reinstall_trigger")
	willReinstall := dedicatedServerWillReinstall(oldTrigger.(string), newTrigger.(string))

	if len(changedTriggers) > 0 && !willReinstall {
		return fmt.Errorf(
			"changing %v requires an OS reinstall, which DESTROYS ALL DATA on disk; "+
				"set reinstall_trigger to a new value (other than %q) to acknowledge and perform the reinstall, "+
				"or revert these fields",
			changedTriggers, reinstallTriggerNone,
		)
	}

	return nil
}

func resourceServerscomDedicatedServerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	dedicatedServer, err := client.Hosts.GetDedicatedServer(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom dedicated server (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.Errorf("error retrieving dedicated server: %s", err)
		}
	}

	if dedicatedServer.ScheduledRelease != nil {
		log.Printf("[WARN] Serverscom dedicated server (%s) marked as scheduled to release", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("hostname", dedicatedServer.Title)
	d.Set("configuration", dedicatedServer.Configuration)
	d.Set("private_ipv4_address", dedicatedServer.PrivateIPv4Address)
	d.Set("public_ipv4_address", dedicatedServer.PublicIPv4Address)
	d.Set("status", dedicatedServer.Status)
	d.Set("operational_status", dedicatedServer.OperationalStatus)
	d.Set("server_model", dedicatedServer.ConfigurationDetails.ServerModelName)
	d.Set("public_uplink", dedicatedServer.ConfigurationDetails.PublicUplinkName)
	d.Set("private_uplink", dedicatedServer.ConfigurationDetails.PrivateUplinkName)
	d.Set("bandwidth", dedicatedServer.ConfigurationDetails.BandwidthName)
	d.Set("operating_system", dedicatedServer.ConfigurationDetails.OperatingSystemFullName)
	d.Set("ram_size", dedicatedServer.ConfigurationDetails.RAMSize)
	d.Set("location", dedicatedServer.LocationCode)
	d.Set("labels", dedicatedServer.Labels)

	if dedicatedServer.Status != "active" {
		return nil
	}

	slots, err := client.Hosts.DedicatedServerDriveSlots(d.Id()).Collect(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	driveSlots := getDriveSlots(slots)

	sort.SliceStable(driveSlots, func(i, j int) bool {
		posI := driveSlots[i]["position"].(int)
		posJ := driveSlots[j]["position"].(int)

		return posI < posJ
	})

	d.Set("slot", driveSlots)

	if dedicatedServer.PublicIPv4Address != nil {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": *dedicatedServer.PublicIPv4Address,
		})
	}

	return nil
}

func resourceServerscomDedicatedServerUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	// A reinstall or update request is rejected while the server is still busy
	// with a previous operation, so wait until it settles to "normal" first.
	if err := waitForDedicatedServerOperationalNormal(ctx, d, meta); err != nil {
		return diag.Errorf("error waiting for dedicated server (%s) to become ready for update: %s", d.Id(), err)
	}

	oldTrigger, newTrigger := d.GetChange("reinstall_trigger")
	willReinstall := dedicatedServerWillReinstall(oldTrigger.(string), newTrigger.(string))

	if willReinstall {
		reinstallInput, err := buildDedicatedServerReinstallInput(ctx, d)
		if err != nil {
			return diag.FromErr(err)
		}

		reinstalledServer, err := client.Hosts.ReinstallOperatingSystemForDedicatedServer(ctx, d.Id(), *reinstallInput)
		if err != nil {
			return diag.FromErr(err)
		}

		// Rely on the reinstall response: it returns the server already in its
		// in-progress operational status (e.g. "installation"). Anchor the wait
		// on that status, then wait for it to return to "normal".
		pending := []string{}
		if s := reinstalledServer.OperationalStatus; s != "" && s != dedicatedServerOperationalStatusNormal {
			pending = append(pending, s)
		}

		if _, err := waitForDedicatedServerAttribute(ctx, d, dedicatedServerOperationalStatusNormal, pending, "operational_status", meta, schema.TimeoutUpdate); err != nil {
			return diag.Errorf("error waiting for dedicated server (%s) reinstall to complete: %s", d.Id(), err)
		}
	}

	input := scgo.DedicatedServerUpdateInput{}

	hasChanges := false
	if d.HasChange("labels") {
		hasChanges = true
		if labelsRaw, ok := d.GetOk("labels"); ok {
			labels := labelsRaw.(map[string]any)
			stringLabels := make(map[string]string)
			for k, v := range labels {
				stringLabels[k] = v.(string)
			}
			input.Labels = stringLabels
		}
	}

	if d.HasChange("hostname") {
		hasChanges = true
		if title, ok := d.GetOk("hostname"); ok {
			input.Title = title.(string)
		}
	}

	// user_data is applied in-place via PATCH; unlike SSH keys it cannot be part
	// of the reinstall payload (the dedicated reinstall API has no user_data field).
	if d.HasChange("user_data") {
		hasChanges = true
		userData := d.Get("user_data").(string)
		input.UserData = &userData
	}

	if hasChanges {
		if _, err := client.Hosts.UpdateDedicatedServer(ctx, d.Id(), input); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceServerscomDedicatedServerRead(ctx, d, meta)
}

// buildDedicatedServerReinstallInput assembles the reinstall payload from the
// current (post-change) resource configuration. Only the fields the reinstall
// API accepts are included: hostname, operating system, drive layout and SSH
// keys. SSH keys and user_data are not managed independently for now — keys are
// simply passed into the reinstall as-is, and user_data stays ForceNew.
func buildDedicatedServerReinstallInput(ctx context.Context, d *schema.ResourceData) (*scgo.OperatingSystemReinstallInput, error) {
	location, err := getLocation(ctx, d.Get("location").(string))
	if err != nil {
		return nil, err
	}

	serverModel, err := getServerModel(ctx, location.ID, d.Get("server_model").(string))
	if err != nil {
		return nil, err
	}

	input := &scgo.OperatingSystemReinstallInput{
		Hostname: d.Get("hostname").(string),
	}

	if operatingSystemName, ok := d.GetOk("operating_system"); ok {
		operatingSystem, err := getOperatingSystem(ctx, location.ID, serverModel.ID, operatingSystemName.(string))
		if err != nil {
			return nil, err
		}
		input.OperatingSystemID = &operatingSystem.ID
	}

	layouts := getLayouts(d)
	reinstallLayouts := make([]scgo.OperatingSystemReinstallLayoutInput, 0, len(layouts))
	for _, layout := range layouts {
		reinstallLayout := scgo.OperatingSystemReinstallLayoutInput{
			SlotPositions: layout.SlotPositions,
			Raid:          layout.Raid,
		}

		for _, partition := range layout.Partitions {
			reinstallLayout.Partitions = append(
				reinstallLayout.Partitions,
				scgo.OperatingSystemReinstallPartitionInput{
					Target: partition.Target,
					Size:   partition.Size,
					Fs:     partition.Fs,
					Fill:   partition.Fill,
				},
			)
		}

		reinstallLayouts = append(reinstallLayouts, reinstallLayout)
	}
	input.Drives = scgo.OperatingSystemReinstallDrivesInput{Layout: reinstallLayouts}

	if val, ok := d.GetOk("ssh_key_fingerprints"); ok {
		input.SSHKeyFingerprints = expandedStringList(val.([]any))
	}

	if userData, ok := d.GetOk("user_data"); ok {
		userDataValue := userData.(string)
		input.UserData = &userDataValue
	}

	return input, nil
}

func resourceServerscomDedicatedServerDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	dedicatedServer, err := client.Hosts.GetDedicatedServer(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom dedicated server (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.Errorf("error retrieving dedicated server: %s", err.Error())
		}
	}

	if dedicatedServer.ScheduledRelease != nil {
		log.Printf("[WARN] Serverscom dedicated server (%s) already scheduled to release", d.Id())
		d.SetId("")
		return nil
	}

	if dedicatedServer.Status == "pending" || dedicatedServer.Status == "init" {
		_, err = waitForDedicatedServerAttribute(ctx, d, "active", []string{"init", "pending"}, "status", meta, schema.TimeoutDelete)
		if err != nil {
			return diag.Errorf("error waiting for dedicated server (%s) to become ready: %s", d.Id(), err)
		}
	}

	if _, err := client.Hosts.ScheduleReleaseForDedicatedServer(ctx, d.Id(), scgo.ScheduleReleaseInput{}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServerscomDedicatedServerCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		location        *scgo.Location
		serverModel     *scgo.ServerModelOption
		operatingSystem *scgo.OperatingSystemOption
		publicUplink    *scgo.UplinkOption
		bandwidth       *scgo.BandwidthOption
		privateUplink   *scgo.UplinkOption
		slots           []scgo.DedicatedServerSlotInput
		layouts         []scgo.DedicatedServerLayoutInput

		publicIpv4NetworkId  *string
		privateIpv4NetworkId *string

		err error
	)

	input := &DedicatedServerCreateInput{}

	if id, ok := d.GetOk("public_ipv4_network_id"); ok {
		publicIpv4NetworkIdValue := id.(string)
		publicIpv4NetworkId = &publicIpv4NetworkIdValue
	}

	if id, ok := d.GetOk("private_ipv4_network_id"); ok {
		privateIpv4NetworkIdValue := id.(string)
		privateIpv4NetworkId = &privateIpv4NetworkIdValue
	}

	hostname := d.Get("hostname").(string)
	input.Hosts = []scgo.DedicatedServerHostInput{
		{
			Hostname:             hostname,
			PublicIPv4NetworkID:  publicIpv4NetworkId,
			PrivateIPv4NetworkID: privateIpv4NetworkId,
		},
	}
	if labelsRaw, ok := d.GetOk("labels"); ok {
		labels := labelsRaw.(map[string]any)
		stringLabels := make(map[string]string)
		for k, v := range labels {
			stringLabels[k] = v.(string)
		}
		input.Hosts[0].Labels = stringLabels
	}

	location, err = getLocation(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	input.LocationID = location.ID

	serverModel, err = getServerModel(ctx, location.ID, d.Get("server_model").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	input.ServerModelID = serverModel.ID

	if ramSize, ok := d.GetOk("ram_size"); ok {
		input.RAMSize = ramSize.(int)
	} else {
		input.RAMSize = serverModel.RAM
	}

	if operatingSystemName, ok := d.GetOk("operating_system"); ok {
		operatingSystem, err = getOperatingSystem(ctx, location.ID, serverModel.ID, operatingSystemName.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		input.OperatingSystemID = &operatingSystem.ID
	}

	input.UplinkModels = scgo.DedicatedServerUplinkModelsInput{}

	if publicUplinkName, ok := d.GetOk("public_uplink"); ok {
		publicUplink, err = getUplink(ctx, location.ID, serverModel.ID, publicUplinkName.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		input.UplinkModels.Public = &scgo.DedicatedServerPublicUplinkInput{}
		input.UplinkModels.Public.ID = publicUplink.ID
	}

	if bandwidthName, ok := d.GetOk("bandwidth"); ok && publicUplink != nil {
		bandwidth, err = getBandwidth(ctx, location.ID, serverModel.ID, publicUplink.ID, bandwidthName.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		input.UplinkModels.Public.BandwidthModelID = bandwidth.ID
	} else if !ok && publicUplink != nil {
		return diag.Errorf("bandwidth must be specified, when public uplink is present")
	}

	privateUplink, err = getUplink(ctx, location.ID, serverModel.ID, d.Get("private_uplink").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	input.UplinkModels.Private.ID = privateUplink.ID

	slots, err = getSlots(ctx, d, location.ID, serverModel.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: Populate slots from model when len(slots) is zero
	err = verifySlots(slots)
	if err != nil {
		return diag.FromErr(err)
	}

	input.Drives.Slots = slots

	layouts = getLayouts(d)

	input.Drives.Layout = layouts

	if val, ok := d.GetOk("ssh_key_fingerprints"); ok {
		input.SSHKeyFingerprints = expandedStringList(val.([]any))
	}

	if ipv6, ok := d.GetOk("ipv6"); ok {
		input.IPv6 = ipv6.(bool)
	}

	if userData, ok := d.GetOk("user_data"); ok {
		userDataValue := userData.(string)
		input.UserData = &userDataValue
	}

	resultChan, err := serverCollector.AddRequest(ctx, "dedicated", input)
	if err != nil {
		return diag.FromErr(err)
	}

	// waiting for result from collector
	result := <-resultChan
	if result.Error != nil {
		return diag.FromErr(result.Error)
	}

	if result.Servers.Count() == 0 {
		return diag.Errorf("invalid dedicated servers count returned by api")
	}

	// find corresponding server by title matching hostname
	id := result.Servers.GetIdByHostname(hostname)
	if id == "" {
		return diag.Errorf("can't find the server with title '%s' in api response", hostname)
	}

	d.SetId(id)

	_, err = waitForDedicatedServerAttribute(ctx, d, "active", []string{"init", "pending"}, "status", meta, schema.TimeoutCreate)
	if err != nil {
		return diag.Errorf("error waiting for dedicated server (%s) to become ready: %s", d.Id(), err)
	}

	return nil
}

func getDriveSlots(slots []scgo.HostDriveSlot) []map[string]any {
	driveSlots := make([]map[string]any, 0)

	for _, slot := range slots {
		var currentSlot = make(map[string]any)

		currentSlot["position"] = slot.Position

		if slot.DriveModel != nil {
			currentSlot["drive_model"] = slot.DriveModel.Name
		} else {
			continue
		}

		driveSlots = append(driveSlots, currentSlot)
	}

	return driveSlots
}

func getLocation(ctx context.Context, code string) (*scgo.Location, error) {
	locations, err := cache.Locations(ctx)
	if err != nil {
		return nil, err
	}

	for _, loc := range locations {
		if normalizeString(loc.Code) == normalizeString(code) {
			return &loc, nil
		}
	}

	return nil, fmt.Errorf("can't find location by: %s", code)
}

func getServerModel(ctx context.Context, locationID int64, name string) (*scgo.ServerModelOption, error) {
	serverModels, err := cache.ServerModels(ctx, locationID)
	if err != nil {
		return nil, err
	}

	for _, sm := range serverModels {
		if normalizeString(sm.Name) == normalizeString(name) {
			return &sm, nil
		}
	}

	return nil, fmt.Errorf("can't find server model by: %s", name)
}

func getDriveModel(ctx context.Context, locationID int64, serverModelID int64, name string) (*scgo.DriveModel, error) {
	driveModels, err := cache.DriveModels(ctx, locationID, serverModelID)
	if err != nil {
		return nil, err
	}

	for _, dm := range driveModels {
		if normalizeString(dm.Name) == normalizeString(name) {
			return &dm, nil
		}
	}

	return nil, fmt.Errorf("can't find drive model by: %s", name)
}

func getOperatingSystem(ctx context.Context, locationID int64, serverModelID int64, name string) (*scgo.OperatingSystemOption, error) {
	operatingSystems, err := cache.OperatingSystems(ctx, locationID, serverModelID)
	if err != nil {
		return nil, err
	}

	for _, os := range operatingSystems {
		fullName := fmt.Sprintf("%s %s %s", os.Name, os.Version, os.Arch)

		if normalizeString(fullName) == normalizeString(name) {
			return &os, nil
		}
	}

	return nil, fmt.Errorf("can't find operating system by: %s", name)
}

func getUplink(ctx context.Context, locationID int64, serverModelID int64, name string) (*scgo.UplinkOption, error) {
	uplinks, err := cache.Uplinks(ctx, locationID, serverModelID)
	if err != nil {
		return nil, err
	}

	for _, uplink := range uplinks {
		if normalizeString(name) == normalizeString(uplink.Name) {
			return &uplink, nil
		}
	}

	return nil, fmt.Errorf("can't find uplink by: %s", name)
}

func getBandwidth(ctx context.Context, locationID int64, serverModelID int64, uplinkModelID int64, name string) (*scgo.BandwidthOption, error) {
	bandwidthList, err := cache.Bandwidth(ctx, locationID, serverModelID, uplinkModelID)
	if err != nil {
		return nil, err
	}

	for _, bandwidth := range bandwidthList {
		if normalizeString(name) == normalizeString(bandwidth.Name) {
			return &bandwidth, nil
		}
	}

	return nil, fmt.Errorf("can't find bandwidth by: %s", name)
}

func getSlots(ctx context.Context, d *schema.ResourceData, locationID int64, serverModelID int64) ([]scgo.DedicatedServerSlotInput, error) {
	var slotsInput []scgo.DedicatedServerSlotInput

	if slotsList, ok := d.GetOk("slot"); ok {
		for _, slotSchema := range slotsList.([]any) {
			slot := slotSchema.(map[string]any)

			var driveModelID *int64

			if value, ok := slot["drive_model"]; ok && len(value.(string)) != 0 {
				driveModel, err := getDriveModel(ctx, locationID, serverModelID, value.(string))
				if err != nil {
					return nil, err
				}

				driveModelID = &driveModel.ID
			}

			slotsInput = append(
				slotsInput,
				scgo.DedicatedServerSlotInput{
					Position:     slot["position"].(int),
					DriveModelID: driveModelID,
				},
			)

		}
	}

	return slotsInput, nil
}

func verifySlots(slots []scgo.DedicatedServerSlotInput) error {
	var hasDriveInFirstSlot bool

	for _, slot := range slots {
		if slot.Position == 0 && slot.DriveModelID != nil {
			hasDriveInFirstSlot = true
		}
	}

	if len(slots) == 0 {
		return fmt.Errorf("at least one slot must be specified")
	} else if !hasDriveInFirstSlot {
		return fmt.Errorf("slot with position 0 must be filled")
	} else {
		return nil
	}
}

func getLayouts(d *schema.ResourceData) []scgo.DedicatedServerLayoutInput {
	var layoutInput []scgo.DedicatedServerLayoutInput

	if layoutsList, ok := d.GetOk("layout"); ok {
		for _, layoutSchema := range layoutsList.([]any) {
			layout := layoutSchema.(map[string]any)

			currentLayout := scgo.DedicatedServerLayoutInput{}
			currentLayout.SlotPositions = expandIntList(layout["slot_positions"].([]any))

			if len(currentLayout.SlotPositions) > 1 {
				raidLevel := layout["raid"].(int)
				currentLayout.Raid = &raidLevel
			}

			currentLayout.Partitions = []scgo.DedicatedServerLayoutPartitionInput{}

			partitionsList := layout["partition"].([]any)

			for _, partitionSchema := range partitionsList {
				partition := partitionSchema.(map[string]any)

				currentPartition := scgo.DedicatedServerLayoutPartitionInput{}
				currentPartition.Target = partition["target"].(string)
				currentPartition.Size = partition["size"].(int)

				currentFs := partition["fs"].(string)
				if currentFs != "" {
					currentPartition.Fs = &currentFs
				}

				currentPartition.Fill = partition["fill"].(bool)

				currentLayout.Partitions = append(
					currentLayout.Partitions,
					currentPartition,
				)
			}

			layoutInput = append(
				layoutInput,
				currentLayout,
			)
		}
	}

	return layoutInput
}

// waitForDedicatedServerOperationalNormal blocks until the server's
// operational_status is "normal", i.e. no operation (such as a reinstall) is in
// progress. It checks the current status first to avoid the poll delay when the
// server is already idle. The API rejects reinstall/update requests while an
// operation is still running, so callers must settle to "normal" beforehand.
func waitForDedicatedServerOperationalNormal(ctx context.Context, d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)

	dedicatedServer, err := client.Hosts.GetDedicatedServer(ctx, d.Id())
	if err != nil {
		return err
	}

	if dedicatedServer.OperationalStatus == dedicatedServerOperationalStatusNormal {
		return nil
	}

	_, err = waitForDedicatedServerAttribute(ctx, d, dedicatedServerOperationalStatusNormal, []string{}, "operational_status", meta, schema.TimeoutUpdate)
	return err
}

func waitForDedicatedServerAttribute(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta any, timeoutKey string) (any, error) {
	log.Printf(
		"[INFO] Waiting for dedicated server (%s) to have %s of %s",
		d.Id(), attribute, target,
	)

	stateConf := &retry.StateChangeConf{
		Pending:      pending,
		Target:       []string{target},
		Refresh:      newDedicatedServerStateRefreshFunc(ctx, d, attribute, meta),
		Timeout:      d.Timeout(timeoutKey),
		PollInterval: 1 * time.Minute,
		Delay:        1 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newDedicatedServerStateRefreshFunc(ctx context.Context, d *schema.ResourceData, attribute string, meta any) retry.StateRefreshFunc {
	return func() (any, string, error) {
		diags := resourceServerscomDedicatedServerRead(ctx, d, meta)
		if diags.HasError() {
			return nil, "", errors.New(diags[0].Summary)
		}

		// See if we can access our attribute
		if attr, ok := d.GetOk(attribute); ok {
			switch attr := attr.(type) {
			case bool:
				return d, strconv.FormatBool(attr), nil
			default:
				return d, attr.(string), nil
			}
		}

		return nil, "", nil
	}
}

// DedicatedServerCreateInput implements ServerCreateInput interface for dedicated servers
type DedicatedServerCreateInput struct {
	scgo.DedicatedServerCreateInput
}

// GetHosts returns hosts from server input
func (d *DedicatedServerCreateInput) GetHosts() []any {
	hosts := make([]any, len(d.Hosts))
	for i, h := range d.Hosts {
		hosts[i] = h
	}
	return hosts
}

// SetHosts sets hosts for server create input
func (d *DedicatedServerCreateInput) SetHosts(hosts []any) {
	if hosts == nil {
		d.Hosts = nil
		return
	}
	dedicatedHosts := make([]scgo.DedicatedServerHostInput, len(hosts))
	for i, h := range hosts {
		dedicatedHosts[i] = h.(scgo.DedicatedServerHostInput)
	}
	d.Hosts = dedicatedHosts
}

// SBMServerResponse implements ServersResponse interface
type DedicatedServerResponse struct {
	servers []scgo.DedicatedServer
}

// GetIdByHostname returns server id if hostname match
func (r *DedicatedServerResponse) GetIdByHostname(h string) string {
	for _, s := range r.servers {
		if s.Title == h {
			return s.ID
		}
	}
	return ""
}

// Count returns amount of servers in response
func (r *DedicatedServerResponse) Count() int {
	return len(r.servers)
}
