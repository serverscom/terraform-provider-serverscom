package serverscom

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	serverscomDedicatedServerDefaultCreateTimeout = 24 * time.Hour
	serverscomDedicatedServerDefaultDeleteTimeout = 1 * time.Hour
)

func resourceServerscomDedicatedServer() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServerscomDedicatedServerRead,
		Update: resourceServerscomDedicatedServerUpdate,
		Delete: resourceServerscomDedicatedServerDelete,
		Create: resourceServerscomDedicatedServerCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(serverscomDedicatedServerDefaultCreateTimeout),
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
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				StateFunc:    HashStringStateFunc(),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new != "" && old == d.Get("user_data")
				},
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
		},
	}
}

func resourceServerscomDedicatedServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)

	ctx := context.TODO()

	dedicatedServer, err := client.Hosts.GetDedicatedServer(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom dedicated server (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving dedicated server: %s", err)
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
	d.Set("server_model", dedicatedServer.ConfigurationDetails.ServerModelName)
	d.Set("public_uplink", dedicatedServer.ConfigurationDetails.PublicUplinkName)
	d.Set("private_uplink", dedicatedServer.ConfigurationDetails.PrivateUplinkName)
	d.Set("bandwidth", dedicatedServer.ConfigurationDetails.BandwidthName)
	d.Set("operating_system", dedicatedServer.ConfigurationDetails.OperatingSystemFullName)
	d.Set("ram_size", dedicatedServer.ConfigurationDetails.RAMSize)
	d.Set("location", dedicatedServer.LocationCode)

	if dedicatedServer.Status != "active" {
		return nil
	}

	slots, err := client.Hosts.DedicatedServerDriveSlots(d.Id()).Collect(ctx)
	if err != nil {
		return err
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

func resourceServerscomDedicatedServerUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceServerscomDedicatedServerRead(d, meta)
}

func resourceServerscomDedicatedServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	dedicatedServer, err := client.Hosts.GetDedicatedServer(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom dedicated server (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving dedicated server: %s", err.Error())
		}
	}

	if dedicatedServer.ScheduledRelease != nil {
		log.Printf("[WARN] Serverscom dedicated server (%s) already scheduled to release", d.Id())
		d.SetId("")
		return nil
	}

	if dedicatedServer.Status == "pending" || dedicatedServer.Status == "init" {
		_, err = waitForDedicatedServerAttribute(d, "active", []string{"init", "pending"}, "status", meta, schema.TimeoutDelete)
		if err != nil {
			return fmt.Errorf("Error waiting for dedicated server (%s) to become ready: %s", d.Id(), err)
		}
	}

	if _, err := client.Hosts.ScheduleReleaseForDedicatedServer(ctx, d.Id()); err != nil {
		return err
	}

	return nil
}

func resourceServerscomDedicatedServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)

	var location *scgo.Location
	var serverModel *scgo.ServerModelOption
	var operatingSystem *scgo.OperatingSystemOption
	var publicUplink *scgo.UplinkOption
	var bandwidth *scgo.BandwidthOption
	var privateUplink *scgo.UplinkOption
	var slots []scgo.DedicatedServerSlotInput
	var layouts []scgo.DedicatedServerLayoutInput
	var dedicatedServers []scgo.DedicatedServer

	var publicIpv4NetworkId *string
	var privateIpv4NetworkId *string

	var err error

	input := scgo.DedicatedServerCreateInput{}

	if id, ok := d.GetOk("public_ipv4_network_id"); ok {
		publicIpv4NetworkIdValue := id.(string)
		publicIpv4NetworkId = &publicIpv4NetworkIdValue
	}

	if id, ok := d.GetOk("private_ipv4_network_id"); ok {
		privateIpv4NetworkIdValue := id.(string)
		privateIpv4NetworkId = &privateIpv4NetworkIdValue
	}

	input.Hosts = []scgo.DedicatedServerHostInput{
		{
			Hostname:             d.Get("hostname").(string),
			PublicIPv4NetworkID:  publicIpv4NetworkId,
			PrivateIPv4NetworkID: privateIpv4NetworkId,
		},
	}

	location, err = getLocation(d.Get("location").(string))
	if err != nil {
		return err
	}

	input.LocationID = location.ID

	serverModel, err = getServerModel(location.ID, d.Get("server_model").(string))
	if err != nil {
		return err
	}

	input.ServerModelID = serverModel.ID

	if ramSize, ok := d.GetOk("ram_size"); ok {
		input.RAMSize = ramSize.(int)
	} else {
		input.RAMSize = serverModel.RAM
	}

	if operatingSystemName, ok := d.GetOk("operating_system"); ok {
		operatingSystem, err = getOperatingSystem(location.ID, serverModel.ID, operatingSystemName.(string))
		if err != nil {
			return err
		}

		input.OperatingSystemID = &operatingSystem.ID
	}

	input.UplinkModels = scgo.DedicatedServerUplinkModelsInput{}

	if publicUplinkName, ok := d.GetOk("public_uplink"); ok {
		publicUplink, err = getUplink(location.ID, serverModel.ID, publicUplinkName.(string))
		if err != nil {
			return err
		}

		input.UplinkModels.Public = &scgo.DedicatedServerPublicUplinkInput{}
		input.UplinkModels.Public.ID = publicUplink.ID
	}

	if bandwidthName, ok := d.GetOk("bandwidth"); ok && publicUplink != nil {
		bandwidth, err = getBandwidth(location.ID, serverModel.ID, publicUplink.ID, bandwidthName.(string))
		if err != nil {
			return err
		}

		input.UplinkModels.Public.BandwidthModelID = bandwidth.ID
	} else if !ok && publicUplink != nil {
		return fmt.Errorf("bandwidth must be specified, when public uplink is present")
	}

	privateUplink, err = getUplink(location.ID, serverModel.ID, d.Get("private_uplink").(string))
	if err != nil {
		return err
	}

	input.UplinkModels.Private.ID = privateUplink.ID

	slots, err = getSlots(d, location.ID, serverModel.ID)
	if err != nil {
		return err
	}

	// TODO: Populate slots from model when len(slots) is zero
	err = verifySlots(slots)
	if err != nil {
		return err
	}

	input.Drives.Slots = slots

	layouts = getLayouts(d)

	input.Drives.Layout = layouts

	if val, ok := d.GetOk("ssh_key_fingerprints"); ok {
		input.SSHKeyFingerprints = expandedStringList(val.([]interface{}))
	}

	if ipv6, ok := d.GetOk("ipv6"); ok {
		input.IPv6 = ipv6.(bool)
	}

	if userData, ok := d.GetOk("user_data"); ok {
		userDataValue := userData.(string)
		input.UserData = &userDataValue
	}

	ctx := context.TODO()

	dedicatedServers, err = client.Hosts.CreateDedicatedServers(ctx, input)
	if err != nil {
		return err
	}

	if len(dedicatedServers) == 0 {
		return fmt.Errorf("Invalid dedicated servers count returned by api")
	}

	dedicatedServer := dedicatedServers[0]

	d.SetId(dedicatedServer.ID)

	_, err = waitForDedicatedServerAttribute(d, "active", []string{"init", "pending"}, "status", meta, schema.TimeoutCreate)
	if err != nil {
		return fmt.Errorf("Error waiting for dedicated server (%s) to become ready: %s", d.Id(), err)
	}

	return resourceServerscomDedicatedServerRead(d, meta)
}

func getDriveSlots(slots []scgo.HostDriveSlot) []map[string]interface{} {
	driveSlots := make([]map[string]interface{}, 0)

	for _, slot := range slots {
		var currentSlot = make(map[string]interface{})

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

func getLocation(code string) (*scgo.Location, error) {
	locations, err := cache.Locations()
	if err != nil {
		return nil, err
	}

	for _, loc := range locations {
		if normalizeString(loc.Code) == normalizeString(code) {
			return &loc, nil
		}
	}

	return nil, fmt.Errorf("Can't find location by: %s", code)
}

func getServerModel(locationID int64, name string) (*scgo.ServerModelOption, error) {
	serverModels, err := cache.ServerModels(locationID)
	if err != nil {
		return nil, err
	}

	for _, sm := range serverModels {
		if normalizeString(sm.Name) == normalizeString(name) {
			return &sm, nil
		}
	}

	return nil, fmt.Errorf("Can't find server model by: %s", name)
}

func getDriveModel(locationID int64, serverModelID int64, name string) (*scgo.DriveModel, error) {
	driveModels, err := cache.DriveModels(locationID, serverModelID)
	if err != nil {
		return nil, err
	}

	for _, dm := range driveModels {
		if normalizeString(dm.Name) == normalizeString(name) {
			return &dm, nil
		}
	}

	return nil, fmt.Errorf("Can't find drive model by: %s", name)
}

func getOperatingSystem(locationID int64, serverModelID int64, name string) (*scgo.OperatingSystemOption, error) {
	operatingSystems, err := cache.OperatingSystems(locationID, serverModelID)
	if err != nil {
		return nil, err
	}

	for _, os := range operatingSystems {
		fullName := fmt.Sprintf("%s %s %s", os.Name, os.Version, os.Arch)

		if normalizeString(fullName) == normalizeString(name) {
			return &os, nil
		}
	}

	return nil, fmt.Errorf("Can't find operating system by: %s", name)
}

func getUplink(locationID int64, serverModelID int64, name string) (*scgo.UplinkOption, error) {
	uplinks, err := cache.Uplinks(locationID, serverModelID)
	if err != nil {
		return nil, err
	}

	for _, uplink := range uplinks {
		if normalizeString(name) == normalizeString(uplink.Name) {
			return &uplink, nil
		}
	}

	return nil, fmt.Errorf("Can't find uplink by: %s", name)
}

func getBandwidth(locationID int64, serverModelID int64, uplinkModelID int64, name string) (*scgo.BandwidthOption, error) {
	bandwidthList, err := cache.Bandwidth(locationID, serverModelID, uplinkModelID)
	if err != nil {
		return nil, err
	}

	for _, bandwidth := range bandwidthList {
		if normalizeString(name) == normalizeString(bandwidth.Name) {
			return &bandwidth, nil
		}
	}

	return nil, fmt.Errorf("Can't find bandwidth by: %s", name)
}

func getSlots(d *schema.ResourceData, locationID int64, serverModelID int64) ([]scgo.DedicatedServerSlotInput, error) {
	var slotsInput []scgo.DedicatedServerSlotInput

	if slotsList, ok := d.GetOk("slot"); ok {
		for _, slotSchema := range slotsList.([]interface{}) {
			slot := slotSchema.(map[string]interface{})

			var driveModelID *int64

			if value, ok := slot["drive_model"]; ok && len(value.(string)) != 0 {
				driveModel, err := getDriveModel(locationID, serverModelID, value.(string))
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
		for _, layoutSchema := range layoutsList.([]interface{}) {
			layout := layoutSchema.(map[string]interface{})

			currentLayout := scgo.DedicatedServerLayoutInput{}
			currentLayout.SlotPositions = expandIntList(layout["slot_positions"].([]interface{}))

			if len(currentLayout.SlotPositions) > 1 {
				raidLevel := layout["raid"].(int)
				currentLayout.Raid = &raidLevel
			}

			currentLayout.Partitions = []scgo.DedicatedServerLayoutPartitionInput{}

			partitionsList := layout["partition"].([]interface{})

			for _, partitionSchema := range partitionsList {
				partition := partitionSchema.(map[string]interface{})

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

func verifyLayouts(layouts []scgo.DedicatedServerLayoutInput) error {
	if len(layouts) == 0 {
		return fmt.Errorf("at least one layout must be specified")
	}

	return nil
}

func waitForDedicatedServerAttribute(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, timeoutKey string) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for dedicated server (%s) to have %s of %s",
		d.Id(), attribute, target,
	)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newDedicatedServerStateRefreshFunc(d, attribute, meta),
		Timeout:    d.Timeout(timeoutKey),
		Delay:      1 * time.Minute,
		MinTimeout: 30 * time.Second,
	}

	return stateConf.WaitForState()
}

func newDedicatedServerStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	return func() (interface{}, string, error) {
		err := resourceServerscomDedicatedServerRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		// See if we can access our attribute
		if attr, ok := d.GetOkExists(attribute); ok {
			dedicatedServer, err := client.Hosts.GetDedicatedServer(ctx, d.Id())

			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving dedicated server: %s", err)
			}

			switch attr.(type) {
			case bool:
				return &dedicatedServer, strconv.FormatBool(attr.(bool)), nil
			default:
				return &dedicatedServer, attr.(string), nil
			}
		}

		return nil, "", nil
	}
}
