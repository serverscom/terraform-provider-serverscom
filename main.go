package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/serverscom/terraform-provider-serverscom/serverscom"
)

var (
	version string = "dev"
	commit  string = ""
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: serverscom.Provider})
}
