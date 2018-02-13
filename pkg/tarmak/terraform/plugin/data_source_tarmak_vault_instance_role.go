package tarmak

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTarmakVaultInstanceRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTarmakVaultInstanceRoleRead,

		Schema: map[string]*schema.Schema{
			"vault_cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"init_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTarmakVaultInstanceRoleRead(d *schema.ResourceData, meta interface{}) error {

	vaultClusterName := d.Get("vault_cluster_name").(string)
	roleName := d.Get("role_name").(string)

	client, err := newClient()
	if err != nil {
		return err
	}

	var args = [2]string{vaultClusterName, roleName}
	var initToken string
	if err := client.Call(fmt.Sprintf("%s.VaultInstanceRoleStatus", serverName), args, &initToken); err != nil {
		return err
	}

	if err := d.Set("init_token", initToken); err != nil {
		return err
	}

	d.SetId(initToken)

	return nil
}
