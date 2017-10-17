package tarmak

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVaultInitToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceVaultInitTokenCreateOrUpdate,
		Read:   resourceVaultInitTokenRead,
		Update: resourceVaultInitTokenCreateOrUpdate,
		Delete: resourceVaultInitTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVaultInitTokenCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	role := d.Get("role").(string)
	return fmt.Errorf("not implemented: role=%s", role)
}

func resourceVaultInitTokenRead(d *schema.ResourceData, meta interface{}) error {
	role := d.Get("role").(string)
	return fmt.Errorf("not implemented: role=%s", role)
}

func resourceVaultInitTokenDelete(d *schema.ResourceData, meta interface{}) error {
	role := d.Get("role").(string)
	return fmt.Errorf("not implemented: role=%s", role)
}
