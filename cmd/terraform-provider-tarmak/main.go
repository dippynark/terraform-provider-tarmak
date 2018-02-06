package main

import (
	//faketarmak "github.com/dippynark/terraform-provider-tarmak/pkg/faketarmak"
	tarmak "github.com/dippynark/terraform-provider-tarmak/pkg/tarmak/terraform/plugin"
	"github.com/hashicorp/terraform/plugin"
)

func main() {

	//go faketarmak.FakeTarmak()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tarmak.Provider})
}
