package tarmak

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/rpc"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jetstack/terraform-provider-tarmak/pkg/faketarmak"
)

func resourceVaultInitToken() *schema.Resource {

	r := &schema.Resource{
		Create: resourceVaultInitTokenCreateOrUpdate,
		Read:   resourceVaultInitTokenRead,
		Update: resourceVaultInitTokenCreateOrUpdate,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			//"name": {
			//	Type:     schema.TypeString,
			//	ForceNew: true,
			//},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"init_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return r
}

func resourceVaultInitTokenCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {

	client, err := waitForConnection()
	if err != nil {
		return err
	}

	var token faketarmak.InitToken

	args := &faketarmak.InitTokenArgs{
		Cluster: d.Get("cluster").(string),
		Env:     d.Get("environment").(string),
		Role:    d.Get("role").(string),
	}
	if err := client.Call("InitToken.TarmakInitToken", args, &token); err != nil {
		return err
	}

	if err := d.Set("init_token", string(token)); err != nil {
		return err
	}

	h := sha1.New()
	if _, err := h.Write([]byte(token)); err != nil {
		return fmt.Errorf("failed to hash init token: %v", token)
	}

	id := fmt.Sprintf("%x", h.Sum(nil))
	d.SetId(id)

	return nil
}

func resourceVaultInitTokenRead(d *schema.ResourceData, meta interface{}) error {
	role := d.Get("role").(string)
	if role == "" {
		return fmt.Errorf("Role must be configured for the Tarmak provider")
	}
	return nil
}

func resourceVaultInitTokenDelete(d *schema.ResourceData, meta interface{}) error {
	role := d.Get("role").(string)
	return fmt.Errorf("not implemented: role=%s", role)
}

func waitForConnection() (client *rpc.Client, err error) {
	for i := 0; i < 5; i++ {
		client, err := rpc.DialHTTP("tcp", ":1234")
		if err == nil {
			return client, nil
		}

		time.Sleep(time.Second * time.Duration(i*2))
	}

	return nil, errors.New("Could not reslove rpc connection")
}
