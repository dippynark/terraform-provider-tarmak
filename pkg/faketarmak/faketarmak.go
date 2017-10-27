package faketarmak

import (
	"net"
	"net/http"
	"net/rpc"
)

type InitToken string
type Result string

type Args struct {
	Env     string
	Cluster string
	Role    string
}

func (i *InitToken) TarmakInitToken(args *Args, resp *Result) error {
	*resp = "This is an init token"
	return nil
}

func FakeTarmak() {
	initToken := new(InitToken)
	err := rpc.Register(initToken)
	if err != nil {
		panic(err)
	}
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}
