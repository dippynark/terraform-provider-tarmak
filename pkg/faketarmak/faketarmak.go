package faketarmak

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

type InitToken string

type Args struct {
	Env     string
	Cluster string
	Role    string
}

func (i *InitToken) TarmakInitToken(args *Args, resp *InitToken) error {
	*resp = InitToken(fmt.Sprintf("token-%s-%s-%s", args.Cluster, args.Env, args.Role))
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
