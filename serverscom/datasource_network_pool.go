package serverscom

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func datasourceServerscomNetworkPool() *schema.Resource {
	recordSchema := networkPoolSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	searchFields := []string{"id", "cidr"}

	recordSchema["id"].ExactlyOneOf = searchFields
	recordSchema["id"].Optional = true

	recordSchema["cidr"].ExactlyOneOf = searchFields
	recordSchema["cidr"].Optional = true

	return &schema.Resource{
		ReadContext: datasourceServerscomNetworkPoolRead,
		Schema:      recordSchema,
	}
}

func datasourceServerscomNetworkPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scgo.Client)

	var foundNetworkPool *scgo.NetworkPool

	if id, ok := d.GetOk("id"); ok {
		networkPool, err := client.NetworkPools.Get(ctx, id.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		foundNetworkPool = networkPool
	} else if v, ok := d.GetOk("cidr"); ok {
		cidr := v.(string)
		networkPoolsList, err := findNetworkPools(ctx, client, cidr)
		if err != nil {
			return diag.FromErr(err)
		}

		networkPool, err := findNetworkPoolByCidr(networkPoolsList, cidr)
		if err != nil {
			return diag.FromErr(err)
		}

		foundNetworkPool = networkPool
	}

	flattenNetworkPool, err := flattenServerscomNetworkPool(foundNetworkPool, meta, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setResourceDataFromMap(d, flattenNetworkPool); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(foundNetworkPool.ID)

	return nil
}

func findNetworkPoolByCidr(pools []interface{}, cidr string) (*scgo.NetworkPool, error) {
	results := make([]scgo.NetworkPool, 0)

	for _, p := range pools {
		networkPool := p.(scgo.NetworkPool)

		if networkPool.CIDR == cidr {
			results = append(results, networkPool)
		}
	}

	if len(results) == 1 {
		return &results[0], nil
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no network pool found with cidr %s", cidr)
	}

	return nil, fmt.Errorf("too many network pools found with cidr %s (found %d, expected 1)", cidr, len(results))
}

func findNetworkPools(ctx context.Context, client *scgo.Client, searchPattern string) ([]interface{}, error) {
	networkPoolsList, err := client.NetworkPools.Collection().SetSearchPattern(searchPattern).Collect(ctx)
	if err != nil {
		return nil, err
	}

	var list = make([]interface{}, len(networkPoolsList))

	for i, p := range networkPoolsList {
		list[i] = p
	}

	return list, nil
}
