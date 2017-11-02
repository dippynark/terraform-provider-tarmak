package tarmak

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/rpc"
	"os"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

type multiCloser struct {
	closers []io.Closer
}

func (mc multiCloser) Close() error {
	var err error
	for _, c := range mc.closers {
		if closeErr := c.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}

type procCloser struct {
	*os.Process
}

func (pc procCloser) Close() error {
	if pc.Process == nil {
		os.Exit(0)
		return nil
	}
	c := make(chan error, 1)
	go func() { _, err := pc.Process.Wait(); c <- err }()
	if err := pc.Process.Signal(os.Interrupt); err != nil {
		return err
	}
	select {
	case err := <-c:
		return err
	case <-time.After(1 * time.Second):
		return pc.Process.Kill()
	}
	return nil
}

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

	client := rpc.NewClient(struct {
		io.Reader
		io.Writer
		io.Closer
	}{os.Stdin, os.Stdout,
		multiCloser{[]io.Closer{os.Stdout, os.Stdin, procCloser{}}},
	})

	//var token faketarmak.InitToken

	cluster := d.Get("cluster").(string)
	//env := d.Get("environment").(string)
	//role := d.Get("role").(string)
	var token *string
	//args := &faketarmak.InitTokenArgs{
	//	Cluster: d.Get("cluster").(string),
	//	Env:     d.Get("environment").(string),
	//	Role:    d.Get("role").(string),
	//}
	if err := client.Call("Hook", cluster, &token); err != nil {
		return err
	}

	panic(*token)

	if err := d.Set("init_token", *token); err != nil {
		return err
	}

	h := sha1.New()
	if _, err := h.Write([]byte(*token)); err != nil {
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

//func getClient() *rpc.Client {
//	return rpc.NewClient(struct {
//		io.Reader
//		io.Writer
//		io.Closer
//	}{os.Stdin, os.Stdout,
//		multiCloser{[]io.Closer{os.Stdout, os.Stdin, procCloser{}}},
//	})
//	//for i := 0; i < 5; i++ {
//
//	//	client, err := rpc.DialHTTP("tcp", ":1234")
//	//	if err == nil {
//	//		return client, nil
//	//	}
//
//	//	time.Sleep(time.Second * time.Duration(i*2))
//	//}
//
//	//return nil, errors.New("Could not reslove rpc connection")
//}
