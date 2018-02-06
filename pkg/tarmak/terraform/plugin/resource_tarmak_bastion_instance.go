package tarmak

import (
	"fmt"
	"io"
	"net"
	"net/rpc"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTarmakBastionInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTarmakBastionInstanceCreate,
		Read:   resourceTarmakBastionInstanceRead,
		Delete: resourceTarmakBastionInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTarmakBastionInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client, err := newClient()
	if err != nil {
		return err
	}

	var args = ""
	var status string
	if err := client.Call("Tarmak.BastionStatus", args, &status); err != nil {
		return err
	}

	if err := d.Set("status", status); err != nil {
		return err
	}

	id := "id"
	d.SetId(id)

	return nil
}

func resourceTarmakBastionInstanceRead(d *schema.ResourceData, meta interface{}) error {

	client, err := newClient()
	if err != nil {
		d.SetId("")
		return err
	}

	var args = ""
	var status string
	if err := client.Call("Tarmak.BastionStatus", args, &status); err != nil {
		d.SetId("")
		return err
	}

	if err := d.Set("status", status); err != nil {
		return err
	}

	return nil
}

func resourceTarmakBastionInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func newClient() (*rpc.Client, error) {

	conn, err := net.Dial("unix", "tarmak.sock")
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to socket: %s", err)
	}

	return rpc.NewClient(struct {
		io.Reader
		io.Writer
		io.Closer
	}{conn, conn, conn}), nil
}
