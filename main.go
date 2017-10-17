package main

import (
	"github.com/hashicorp/terraform/plugin"
	tarmak "github.com/jetstack/terraform-provider-tarmak/pkg/tarmak/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tarmak.Provider})
}
