package tarmak

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTarmakVaultCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTarmakVaultClusterRead,

		Schema: map[string]*schema.Schema{
			"instances": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTarmakVaultClusterRead(d *schema.ResourceData, meta interface{}) error {

	instanceInterfaces := d.Get("instances").([]interface{})
	instances := []string{}
	for _, instancesInterface := range instanceInterfaces {
		instances = append(instances, instancesInterface.(string))
	}

	client, err := newClient()
	if err != nil {
		d.SetId("")
		return err
	}

	var status string
	if err := client.Call(fmt.Sprintf("%s.VaultClusterStatus", serverName), instances, &status); err != nil {
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
