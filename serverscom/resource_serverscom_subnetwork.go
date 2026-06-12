package serverscom

import (
	"context"
	"log"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func resourceServerscomSubnetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceServerscomSubnetworkRead,
		UpdateContext: resourceServerscomSubnetworkUpdate,
		DeleteContext: resourceServerscomSubnetworkDelete,
		CreateContext: resourceServerscomSubnetworkCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceServerscomSubnetworkRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	networkPoolID := d.Get("network_pool_id").(string)

	subnetwork, err := client.NetworkPools.GetSubnetwork(ctx, networkPoolID, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom subnetwork (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.Errorf("error retrieving subnetwork: %s", err)
		}
	}

	_, ipv4Net, err := net.ParseCIDR(subnetwork.CIDR)
	if err != nil {
		return diag.Errorf("invalid cidr value: %s", err.Error())
	}

	mask, _ := ipv4Net.Mask.Size()

	d.Set("title", subnetwork.Title)
	d.Set("cidr", subnetwork.CIDR)
	d.Set("mask", mask)
	d.Set("network_pool_id", subnetwork.NetworkPoolID)

	return nil
}

func resourceServerscomSubnetworkUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	networkPoolID := d.Get("network_pool_id").(string)

	if d.HasChange("title") {
		var newTitle *string
		if v, ok := d.GetOk("title"); ok {
			title := v.(string)
			newTitle = &title
		}

		input := scgo.SubnetworkUpdateInput{}
		input.Title = newTitle

		if _, err := client.NetworkPools.UpdateSubnetwork(ctx, networkPoolID, d.Id(), input); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceServerscomSubnetworkRead(ctx, d, meta)
}

func resourceServerscomSubnetworkDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

	networkPoolID := d.Get("network_pool_id").(string)

	if _, err := client.NetworkPools.GetSubnetwork(ctx, networkPoolID, d.Id()); err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom subnetwork (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.Errorf("error retrieving subnetwork: %s", err.Error())
		}
	}

	if err := client.NetworkPools.DeleteSubnetwork(ctx, networkPoolID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServerscomSubnetworkCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*scgo.Client)

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
		return diag.Errorf("mask or cidr must be set")
	}

	subnetwork, err := client.NetworkPools.CreateSubnetwork(ctx, networkPoolID, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(subnetwork.ID)

	return resourceServerscomSubnetworkRead(ctx, d, meta)
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
