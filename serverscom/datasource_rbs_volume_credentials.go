package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func dataSourceServerscomRBSVolumeCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerscomRBSVolumeCredentialsRead,
		Schema: map[string]*schema.Schema{
			"volume_id": {Type: schema.TypeString, Required: true},

			"username":   {Type: schema.TypeString, Computed: true},
			"password":   {Type: schema.TypeString, Computed: true},
			"target_iqn": {Type: schema.TypeString, Computed: true},
			"ip_address": {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceServerscomRBSVolumeCredentialsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)
	volumeID := d.Get("volume_id").(string)

	creds, err := client.RemoteBlockStorageVolumes.GetCredentials(ctx, volumeID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving RBS volume credentials: %s", err.Error()))
	}

	d.SetId(creds.VolumeID)
	d.Set("username", creds.Username)
	d.Set("password", creds.Password)
	if creds.TargetIQN != nil {
		_ = d.Set("target_iqn", *creds.TargetIQN)
	}
	if creds.IPAddress != nil {
		_ = d.Set("ip_address", *creds.IPAddress)
	}

	return nil
}
