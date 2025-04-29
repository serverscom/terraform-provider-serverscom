package serverscom

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SERVERSCOM_API_TOKEN"}, nil),
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVERSCOM_API_URL", "https://api.servers.com/v1"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"serverscom_network_pool":                       datasourceServerscomNetworkPool(),
			"serverscom_dedicated_server":                   dataSourceServerscomDedicatedServer(),
			"serverscom_location":                           dataSourceServerscomLocation(),
			"serverscom_locations":                          dataSourceServerscomLocations(),
			"serverscom_server_model_order_option":          dataSourceServerscomServerModelOrderOption(),
			"serverscom_server_model_order_options":         dataSourceServerscomServerModelOrderOptions(),
			"serverscom_drive_model_order_option":           dataSourceServerscomDriveModelOrderOption(),
			"serverscom_drive_model_order_options":          dataSourceServerscomDriveModelOrderOptions(),
			"serverscom_operating_system_order_option":      dataSourceServerscomOperatingSystemOrderOption(),
			"serverscom_operating_system_order_options":     dataSourceServerscomOperatingSystemOrderOptions(),
			"serverscom_sbm_operating_system_order_option":  dataSourceServerscomSbmOperatingSystemOrderOption(),
			"serverscom_sbm_operating_system_order_options": dataSourceServerscomSbmOperatingSystemOrderOptions(),
			"serverscom_ram_order_option":                   dataSourceServerscomRamOrderOption(),
			"serverscom_ram_order_options":                  dataSourceServerscomRamOrderOptions(),
			"serverscom_uplink_model_order_option":          dataSourceServerscomUplinkModelOrderOption(),
			"serverscom_uplink_model_order_options":         dataSourceServerscomUplinkModelOrderOptions(),
			"serverscom_bandwidth_order_option":             dataSourceServerscomBandwidthOrderOption(),
			"serverscom_bandwidth_order_options":            dataSourceServerscomBandwidthOrderOptions(),
			"serverscom_sbm_flavor_order_option":            dataSourceServerscomSbmFlavorOrderOption(),
			"serverscom_sbm_flavor_order_options":           dataSourceServerscomSbmFlavorOrderOptions(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"serverscom_dedicated_server":         resourceServerscomDedicatedServer(),
			"serverscom_l2_segment":               resourceServerscomL2Segment(),
			"serverscom_cloud_computing_instance": resourceServerscomCloudComputingInstance(),
			"serverscom_ssh_key":                  resourceServerscomSSHKey(),
			"serverscom_subnetwork":               resourceServerscomSubnetwork(),
			"serverscom_sbm_server":               resourceServerscomSBM(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := scgo.NewClientWithEndpoint(
		d.Get("token").(string),
		d.Get("endpoint").(string),
	)

	client.SetupUserAgent("terraform-provider-serverscom")
	cache = NewCache(client)

	serverCollector = NewServerCollector(client)
	serverCollector.Run()

	return client, nil
}
