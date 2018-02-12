package tarmak

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTarmakVaultCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTarmakVaultClusterRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTarmakVaultClusterRead(d *schema.ResourceData, meta interface{}) error {

	client, err := newClient()
	if err != nil {
		d.SetId("")
		return err
	}

	var args = ""
	var status string
	if err := client.Call(fmt.Sprintf("%s.VaultStatus", serverName), args, &status); err != nil {
		d.SetId("")
		return err
	}

	if err := d.Set("status", status); err != nil {
		return err
	}

	id := "vaultcluster"
	d.SetId(id)

	return nil
}
