package tarmak

import (
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
		Delete: resourceVaultInitTokenDelete,
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

type Args struct {
	Env     string
	Cluster string
	Role    string
}

func resourceVaultInitTokenCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {

	client, err := waitForConnection()
	if err != nil {
		return err
	}

	var token faketarmak.InitToken

	args := &faketarmak.Args{
		Cluster: d.Get("cluster").(string),
		Env:     d.Get("environment").(string),
		Role:    d.Get("role").(string),
	}
	if err := client.Call("InitToken.TarmakInitToken", args, &token); err != nil {
		return err
	}

	d.Partial(true)
	if err := d.Set("init_token", string(token)); err != nil {
		return err
	}
	d.SetPartial("init_token")
	d.Partial(false)

	//role := d.Get("role").(string)
	//return fmt.Errorf("not implemented: role=%s", role)
	return nil
}

func resourceVaultInitTokenRead(d *schema.ResourceData, meta interface{}) error {
	role := d.Get("role").(string)
	return fmt.Errorf("not implemented: role=%s", role)
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
