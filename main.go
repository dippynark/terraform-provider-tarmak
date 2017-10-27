package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/jetstack/terraform-provider-tarmak/pkg/faketarmak"
	tarmak "github.com/jetstack/terraform-provider-tarmak/pkg/tarmak/terraform/plugin"
)

func main() {
	go faketarmak.FakeTarmak()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tarmak.Provider})
}
