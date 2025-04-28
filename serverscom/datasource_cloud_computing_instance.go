package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomCloudInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomCloudInstanceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"openstack_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"gpn_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"backup_copies": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_port_blocked": {
				Type:     schema.TypeBool,
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

func dataSourceServerscomCloudInstanceRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	instanceID := d.Get("id").(string)

	instance, err := client.CloudComputingInstances.Get(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("Error retrieving cloud instance: %s", err.Error())
	}

	d.SetId(instance.ID)
	d.Set("openstack_uuid", instance.OpenstackUUID)
	d.Set("name", instance.Name)
	d.Set("status", instance.Status)
	d.Set("region_id", instance.RegionID)
	d.Set("region_code", instance.RegionCode)
	d.Set("flavor_id", instance.FlavorID)
	d.Set("flavor_name", instance.FlavorName)
	d.Set("image_id", instance.ImageID)
	d.Set("image_name", instance.ImageName)
	d.Set("public_ipv4_address", instance.PublicIPv4Address)
	d.Set("private_ipv4_address", instance.PrivateIPv4Address)
	d.Set("local_ipv4_address", instance.LocalIPv4Address)
	d.Set("public_ipv6_address", instance.PublicIPv6Address)
	d.Set("ipv6_enabled", instance.IPv6Enabled)
	d.Set("gpn_enabled", instance.GPNEnabled)
	d.Set("backup_copies", instance.BackupCopies)
	d.Set("public_port_blocked", instance.PublicPortBlocked)
	d.Set("labels", instance.Labels)
	d.Set("created_at", instance.Created)
	d.Set("updated_at", instance.Updated)

	return nil
}
