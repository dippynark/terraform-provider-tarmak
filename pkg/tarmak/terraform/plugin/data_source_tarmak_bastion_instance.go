package tarmak

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTarmakBastionInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTarmakBastionInstanceRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTarmakBastionInstanceRead(d *schema.ResourceData, meta interface{}) error {

	client, err := newClient()
	if err != nil {
		d.SetId("")
		return err
	}

	var args = ""
	var status string
	if err := client.Call(fmt.Sprintf("%s.BastionStatus", serverName), args, &status); err != nil {
		d.SetId("")
		return err
	}

	if err := d.Set("status", status); err != nil {
		return err
	}

	id := "bastioninstance"
	d.SetId(id)

	return nil
}
