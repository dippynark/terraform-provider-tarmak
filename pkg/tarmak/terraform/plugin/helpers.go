package tarmak

import (
	"fmt"
	"io"
	"net"
	"net/rpc"
)

const (
	serverName   = "Tarmak"
	tarmakSocket = "tarmak.sock"
)

func newClient() (*rpc.Client, error) {

	conn, err := net.Dial("unix", tarmakSocket)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to socket: %s", err)
	}

	return rpc.NewClient(struct {
		io.Reader
		io.Writer
		io.Closer
	}{conn, conn, conn}), nil
}
