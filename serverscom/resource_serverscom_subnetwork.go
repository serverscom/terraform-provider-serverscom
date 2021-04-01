package serverscom

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func resourceServerscomSubnetwork() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServerscomSubnetworkRead,
		Update: resourceServerscomSubnetworkUpdate,
		Delete: resourceServerscomSubnetworkDelete,
		Create: resourceServerscomSubnetworkCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: resourceServerscomSubnetworkCIDRDiffSupress,
			},
			"mask": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: resourceServerscomSubnetworkMaskDiffSupress,
			},
			"network_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServerscomSubnetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	networkPoolID := d.Get("network_pool_id").(string)

	subnetwork, err := client.NetworkPools.GetSubnetwork(ctx, networkPoolID, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom subnetwork (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving subnetwork: %s", err)
		}
	}

	_, ipv4Net, err := net.ParseCIDR(subnetwork.CIDR)
	if err != nil {
		return fmt.Errorf("Invalid cidr value: %s", err.Error())
	}

	mask, _ := ipv4Net.Mask.Size()

	d.Set("title", subnetwork.Title)
	d.Set("cidr", subnetwork.CIDR)
	d.Set("mask", mask)
	d.Set("network_pool_id", subnetwork.NetworkPoolID)

	return nil
}

func resourceServerscomSubnetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	networkPoolID := d.Get("network_pool_id").(string)

	var newTitle *string
	if v, ok := d.GetOk("title"); ok {
		title := v.(string)
		newTitle = &title
	} else {
		newTitle = nil
	}

	input := scgo.SubnetworkUpdateInput{}
	input.Title = newTitle

	if _, err := client.NetworkPools.UpdateSubnetwork(ctx, networkPoolID, d.Id(), input); err != nil {
		return err
	}

	return resourceServerscomSubnetworkRead(d, meta)
}

func resourceServerscomSubnetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	networkPoolID := d.Get("network_pool_id").(string)

	if _, err := client.NetworkPools.GetSubnetwork(ctx, networkPoolID, d.Id()); err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom subnetwork (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving subnetwork: %s", err.Error())
		}
	}

	return client.NetworkPools.DeleteSubnetwork(ctx, networkPoolID, d.Id())
}

func resourceServerscomSubnetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	networkPoolID := d.Get("network_pool_id").(string)

	var title *string
	if v, ok := d.GetOk("title"); ok {
		titleValue := v.(string)
		title = &titleValue
	} else {
		title = nil
	}

	input := scgo.SubnetworkCreateInput{}
	input.Title = title

	if cidr, cidrIsOk := d.GetOk("cidr"); cidrIsOk {
		cidrValue := cidr.(string)
		input.CIDR = &cidrValue
	} else if mask, maskIsOk := d.GetOk("mask"); maskIsOk {
		maskValue := mask.(int)
		input.Mask = &maskValue
	} else {
		return fmt.Errorf("mask or cidr must be set")
	}

	subnetwork, err := client.NetworkPools.CreateSubnetwork(ctx, networkPoolID, input)
	if err != nil {
		return err
	}

	d.SetId(subnetwork.ID)

	return resourceServerscomSubnetworkRead(d, meta)
}

func resourceServerscomSubnetworkCIDRDiffSupress(k, old, new string, d *schema.ResourceData) bool {
	if _, ok := d.GetOk("mask"); ok && new == "" && old != "" {
		return true
	}

	return false
}

func resourceServerscomSubnetworkMaskDiffSupress(k, old, new string, d *schema.ResourceData) bool {
	if _, ok := d.GetOk("cidr"); ok && new == "" && old != "" {
		return true
	}

	return false
}
