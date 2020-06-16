package serverscom

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func resourceServerscomCloudComputingInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServerscomCloudComputingInstanceRead,
		Update: resourceServerscomCloudComputingInstanceUpdate,
		Delete: resourceServerscomCloudComputingInstanceDelete,
		Create: resourceServerscomCloudComputingInstanceCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"image": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"flavor": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"gpn_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ipv6_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"backup_copies": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"ssh_key_fingerprint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"openstack_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServerscomCloudComputingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)

	ctx := context.TODO()

	cloudInstance, err := client.CloudComputingInstances.Get(ctx, d.Id())
	if err != nil {
		return err
	}

	d.Set("status", cloudInstance.Status)
	d.Set("name", cloudInstance.Name)
	d.Set("image", cloudInstance.ImageName)
	d.Set("flavor", cloudInstance.FlavorName)
	d.Set("private_ipv4_address", cloudInstance.PrivateIPv4Address)
	d.Set("public_ipv4_address", cloudInstance.PublicIPv4Address)
	d.Set("public_ipv6_address", cloudInstance.PublicIPv6Address)
	d.Set("ipv6_enabled", cloudInstance.PublicIPv6Address)
	d.Set("gpn_enabled", cloudInstance.GPNEnabled)
	d.Set("openstack_uuid", cloudInstance.OpenstackUUID)

	if cloudInstance.PublicIPv4Address != nil {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": *cloudInstance.PublicIPv4Address,
		})
	}

	return nil
}

func resourceServerscomCloudComputingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error

	client := meta.(*scgo.Client)

	input := scgo.CloudComputingInstanceUpdateInput{}

	name := d.Get("name").(string)
	input.Name = &name

	if d.HasChange("backup_copies") {
		backupCopies := d.Get("backup_copies").(int)
		input.BackupCopies = &backupCopies
	}

	if d.HasChange("ipv6_enabled") {
		ipv6Enabled := d.Get("ipv6_enabled").(bool)
		input.IPv6Enabled = &ipv6Enabled
	}

	if d.HasChange("gpn_enabled") {
		gpnEnabled := d.Get("gpn_enabled").(bool)
		input.GPNEnabled = &gpnEnabled
	}

	ctx := context.TODO()

	_, err = client.CloudComputingInstances.Update(ctx, d.Id(), input)
	if err != nil {
		return err
	}

	return resourceServerscomCloudComputingInstanceRead(d, meta)
}

func resourceServerscomCloudComputingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)

	ctx := context.TODO()

	cloudInstance, err := client.CloudComputingInstances.Get(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom cloud computing instance (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving cloud computing instance: %s", err.Error())
		}
	}

	if cloudInstance.Status == "DELETING" {
		log.Printf("[WARN] Serverscom cloud computing instance (%s) already scheduled to delete", d.Id())
		d.SetId("")
		return nil
	}

	return client.CloudComputingInstances.Delete(ctx, d.Id())
}

func resourceServerscomCloudComputingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)

	input := scgo.CloudComputingInstanceCreateInput{}
	input.Name = d.Get("name").(string)

	region, err := getRegion(d.Get("region").(string))
	if err != nil {
		return err
	}

	input.RegionID = region.ID

	flavor, err := getFlavor(region.ID, d.Get("flavor").(string))
	if err != nil {
		return err
	}

	input.FlavorID = flavor.ID

	image, err := getImage(region.ID, d.Get("image").(string))
	if err != nil {
		return err
	}

	input.ImageID = image.ID

	if v, ok := d.GetOk("gpn_enabled"); ok {
		gpnEnabled := v.(bool)
		input.GPNEnabled = &gpnEnabled
	}

	if v, ok := d.GetOk("ipv6_enabled"); ok {
		ipv6Enabled := v.(bool)
		input.IPv6Enabled = &ipv6Enabled
	}

	if v, ok := d.GetOk("backup_copies"); ok {
		backupCopies := v.(int)
		input.BackupCopies = &backupCopies
	}

	sshKeyFp := d.Get("ssh_key_fingerprint").(string)
	input.SSHKeyFingerprint = &sshKeyFp

	ctx := context.TODO()

	cloudInstance, err := client.CloudComputingInstances.Create(ctx, input)
	if err != nil {
		return err
	}

	d.SetId(cloudInstance.ID)

	_, err = waitForCloudComputingInstanceAttribute(d, "ACTIVE", []string{"PROVISIONING", "BUILDING", "REBOOTING"}, "status", meta)
	if err != nil {
		return fmt.Errorf("Error waiting for cloud computing instance (%s) to become active: %s", d.Id(), err)
	}

	return resourceServerscomCloudComputingInstanceRead(d, meta)
}

func waitForCloudComputingInstanceAttribute(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for cloud computing instance (%s) to have %s of %s",
		d.Id(), attribute, target,
	)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newCloudComputingInstanceStateRefreshFunc(d, attribute, meta),
		Timeout:    5 * time.Minute,
		Delay:      1 * time.Minute,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForState()
}

func newCloudComputingInstanceStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	return func() (interface{}, string, error) {
		err := resourceServerscomCloudComputingInstanceRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		// See if we can access our attribute
		if attr, ok := d.GetOkExists(attribute); ok {
			cloudInstance, err := client.CloudComputingInstances.Get(ctx, d.Id())

			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving cloud instance: %s", err)
			}

			switch attr.(type) {
			case bool:
				return &cloudInstance, strconv.FormatBool(attr.(bool)), nil
			default:
				return &cloudInstance, attr.(string), nil
			}
		}

		return nil, "", nil
	}
}

func getRegion(code string) (*scgo.CloudComputingRegion, error) {
	regions, err := cache.CloudComputingRegions()
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		if region.Code == code {
			return &region, nil
		}
	}

	return nil, fmt.Errorf("Can't find cloud computing region by: %s", code)
}

func getFlavor(regionID int64, name string) (*scgo.CloudComputingFlavor, error) {
	flavors, err := cache.CloudComputingFlavors(regionID)
	if err != nil {
		return nil, err
	}

	for _, flavor := range flavors {
		if flavor.Name == name {
			return &flavor, nil
		}
	}

	return nil, fmt.Errorf("Can't find cloud computing flavor by: %s", name)
}

func getImage(regionID int64, name string) (*scgo.CloudComputingImage, error) {
	images, err := cache.CloudComputingImages(regionID)
	if err != nil {
		return nil, err
	}

	for _, image := range images {
		if image.Name == name {
			return &image, nil
		}
	}

	return nil, fmt.Errorf("Can't find cloud computing image by: %s", name)
}
