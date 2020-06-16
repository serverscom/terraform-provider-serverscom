package serverscom

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func resourceServerscomSSHKey() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServerscomSSHKeyRead,
		Update: resourceServerscomSSHKeyUpdate,
		Delete: resourceServerscomSSHKeyDelete,
		Create: resourceServerscomSSHKeyCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
		},
	}
}

func resourceServerscomSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	sshKey, err := client.SSHKeys.Get(ctx, d.Id())
	if err != nil {
		return err
	}

	d.Set("name", sshKey.Name)
	d.Set("fingerprint", sshKey.Fingerprint)

	return nil
}

func resourceServerscomSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	if _, err := client.SSHKeys.Get(ctx, d.Id()); err != nil {
		return err
	}

	var newName string
	if v, ok := d.GetOk("name"); ok {
		newName = v.(string)
	}

	input := scgo.SSHKeyUpdateInput{}
	input.Name = newName

	if _, err := client.SSHKeys.Update(ctx, d.Id(), input); err != nil {
		return err
	}

	return resourceServerscomSSHKeyRead(d, meta)
}

func resourceServerscomSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	err := client.SSHKeys.Delete(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom ssh key (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving ssh key: %s", err.Error())
		}
	}

	return nil
}

func resourceServerscomSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	input := scgo.SSHKeyCreateInput{}
	input.PublicKey = d.Get("public_key").(string)
	input.Name = d.Get("name").(string)

	sshKey, err := client.SSHKeys.Create(ctx, input)
	if err != nil {
		return err
	}

	d.SetId(sshKey.Fingerprint)

	return resourceServerscomSSHKeyRead(d, meta)
}

func resourceServerscomSSHKeyPublicKeyDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}
