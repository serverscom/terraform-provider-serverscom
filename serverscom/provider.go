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
			"serverscom_network_pool": datasourceServerscomNetworkPool(),
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

	return client, nil
}
