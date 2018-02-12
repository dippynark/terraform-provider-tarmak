package tarmak

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTarmakVaultInstanceRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTarmakVaultInstanceRoleRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTarmakVaultInstanceRoleRead(d *schema.ResourceData, meta interface{}) error {

	client, err := newClient()
	if err != nil {
		return err
	}

	var args = ""
	var status string
	if err := client.Call(fmt.Sprintf("%s.VaultInstanceRole", serverName), args, &status); err != nil {
		return err
	}

	if err := d.Set("status", status); err != nil {
		return err
	}

	id := "vaultinstancerole"
	d.SetId(id)

	return nil
}
