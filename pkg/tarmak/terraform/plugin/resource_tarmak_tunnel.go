package tarmak

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTarmakTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceTarmakTunnelCreate,
		Read:   resourceTarmakTunnelRead,
		Delete: resourceTarmakTunnelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bind_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"bind_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"destination_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"ssh_config_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ssh_config_host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTarmakTunnelCreate(d *schema.ResourceData, meta interface{}) error {

	// build arguments
	bindAddress, err := getHostByName(d.Get("bind_address").(string))
	if err != nil {
		return err
	}

	bindPort, err := getUnusedPort(bindAddress)
	if err != nil {
		return err
	}
	d.Set("bindPort", bindPort)

	destinationAddress, err := getHostByName(d.Get("destination_address").(string))
	if err != nil {
		return err
	}

	destinationPort := d.Get("destination_port").(int)

	sshConfigPath := d.Get("ssh_config_path").(string)
	if _, err := os.Stat(sshConfigPath); err != nil {
		return err
	}

	sshConfigHost := d.Get("ssh_config_host").(string)

	// build command
	args := []string{
		"ssh",
		"-F",
		sshConfigPath,
		"-N",
		fmt.Sprintf("-L%s:%d:%s:%d", bindAddress, bindPort, destinationAddress, destinationPort),
		sshConfigHost,
	}
	cmd := exec.Command(args[0], args[1:len(args)]...)

	// start tunnel
	// TODO wait until tcp socket is reachable
	err = cmd.Start()
	if err != nil {
		return err
	}

	// set ID and return
	d.SetId(getIDFromProcess(cmd.Process))
	return nil
}

func resourceTarmakTunnelRead(d *schema.ResourceData, meta interface{}) error {

	// retreive PID
	process, err := getProcessFromID(d.Id())
	if err != nil {
		// asume dead tunnel
		d.SetId("")
		return nil
	}

	// test if process can be signalled
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		// assume dead tunnel
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTarmakTunnelDelete(d *schema.ResourceData, meta interface{}) error {

	// find process
	process, err := getProcessFromID(d.Id())
	if err != nil {
		// assume dead tunnel
		return nil
	}

	// signal process
	if err := process.Signal(syscall.SIGTERM); err != nil {
		// assume signal failed because of dead tunnel
		return nil
	}

	// wait on process
	if _, err := process.Wait(); err != nil {
		// assume process has already been waited on
		return nil
	}

	return nil
}

func getUnusedPort(bindAddress net.IP) (int, error) {

	l, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   bindAddress,
		Port: 0,
	})
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

func getHostByName(address string) (net.IP, error) {

	ip := net.ParseIP(address)
	if ip == nil {
		return ip, fmt.Errorf("could not parse address: %s", address)
	}

	ip = ip.To4()
	if ip == nil {
		return ip, fmt.Errorf("address is not an IPv4 address: %s", address)
	}

	return ip, nil
}

func getIDFromProcess(process *os.Process) string {
	id := strconv.Itoa(process.Pid)
	return id
}

func getProcessFromID(id string) (*os.Process, error) {

	// get PID from ID
	pid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	// find process
	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	return process, err
}
