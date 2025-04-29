package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomDedicatedServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerscomDedicatedServerRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rack_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"location_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operational_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"power_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration": {
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
			"lease_start_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduled_release_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"oob_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Flattened configuration_details fields
			"ram_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_model_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_model_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_uplink_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_uplink_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_uplink_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_uplink_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_system_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"operating_system_full_name": {
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

func dataSourceServerscomDedicatedServerRead(d *schema.ResourceData, meta any) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	serverID := d.Get("id").(string)

	server, err := client.Hosts.GetDedicatedServer(ctx, serverID)
	if err != nil {
		return fmt.Errorf("Error retrieving dedicated server: %s", err.Error())
	}

	d.SetId(server.ID)
	d.Set("rack_id", server.RackID)
	d.Set("title", server.Title)
	d.Set("location_id", server.LocationID)
	d.Set("location_code", server.LocationCode)
	d.Set("status", server.Status)
	d.Set("operational_status", server.OperationalStatus)
	d.Set("power_status", server.PowerStatus)
	d.Set("configuration", server.Configuration)
	d.Set("private_ipv4_address", server.PrivateIPv4Address)
	d.Set("public_ipv4_address", server.PublicIPv4Address)
	d.Set("lease_start_at", server.LeaseStart)
	d.Set("scheduled_release_at", server.ScheduledRelease)
	d.Set("type", server.Type)
	d.Set("oob_ipv4_address", server.OobIPv4Address)

	// Setting flattened configuration_details fields
	d.Set("ram_size", server.ConfigurationDetails.RAMSize)
	d.Set("server_model_id", server.ConfigurationDetails.ServerModelID)
	d.Set("server_model_name", server.ConfigurationDetails.ServerModelName)
	d.Set("bandwidth_id", server.ConfigurationDetails.BandwidthID)
	d.Set("bandwidth_name", server.ConfigurationDetails.BandwidthName)
	d.Set("private_uplink_id", server.ConfigurationDetails.PrivateUplinkID)
	d.Set("private_uplink_name", server.ConfigurationDetails.PrivateUplinkName)
	d.Set("public_uplink_id", server.ConfigurationDetails.PublicUplinkID)
	d.Set("public_uplink_name", server.ConfigurationDetails.PublicUplinkName)
	d.Set("operating_system_id", server.ConfigurationDetails.OperatingSystemID)
	d.Set("operating_system_full_name", server.ConfigurationDetails.OperatingSystemFullName)

	d.Set("labels", server.Labels)
	d.Set("created_at", server.Created)
	d.Set("updated_at", server.Updated)

	return nil
}
