package serverscom

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func resourceServerscomSSHKey() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceServerscomSSHKeyRead,
		UpdateContext: resourceServerscomSSHKeyUpdate,
		DeleteContext: resourceServerscomSSHKeyDelete,
		CreateContext: resourceServerscomSSHKeyCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"public_key": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validation.NoZeroValues,
				DiffSuppressFunc: resourceServerscomSSHKeyPublicKeyDiffSuppress,
			},
			"fingerprint": {
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

func resourceServerscomSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	sshKey, err := client.SSHKeys.Get(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom ssh key (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.Errorf("error retrieving ssh key: %s", err)
		}
	}

	d.Set("name", sshKey.Name)
	d.Set("fingerprint", sshKey.Fingerprint)
	d.Set("labels", sshKey.Labels)

	return nil
}

func resourceServerscomSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	var newName string
	if v, ok := d.GetOk("name"); ok {
		newName = v.(string)
	}

	input := scgo.SSHKeyUpdateInput{}
	input.Name = newName

	if d.HasChange("labels") {
		if labelsRaw, ok := d.GetOk("labels"); ok {
			labels := labelsRaw.(map[string]any)
			stringLabels := make(map[string]string)
			for k, v := range labels {
				stringLabels[k] = v.(string)
			}
			input.Labels = stringLabels
		}
	}

	if _, err := client.SSHKeys.Update(ctx, d.Id(), input); err != nil {
		return diag.FromErr(err)
	}

	return resourceServerscomSSHKeyRead(ctx, d, meta)
}

func resourceServerscomSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	err := client.SSHKeys.Delete(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom ssh key (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.Errorf("error retrieving ssh key: %s", err.Error())
		}
	}

	return nil
}

func resourceServerscomSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	input := scgo.SSHKeyCreateInput{}
	input.PublicKey = d.Get("public_key").(string)
	input.Name = d.Get("name").(string)

	if labelsRaw, ok := d.GetOk("labels"); ok {
		labels := labelsRaw.(map[string]any)
		stringLabels := make(map[string]string)
		for k, v := range labels {
			stringLabels[k] = v.(string)
		}
		input.Labels = stringLabels
	}

	sshKey, err := client.SSHKeys.Create(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sshKey.Fingerprint)

	return resourceServerscomSSHKeyRead(ctx, d, meta)
}

func resourceServerscomSSHKeyPublicKeyDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}
