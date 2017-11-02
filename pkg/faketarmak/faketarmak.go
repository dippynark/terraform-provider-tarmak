package faketarmak

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"time"
)

type InitToken string

type InitTokenArgs struct {
	Env     string
	Cluster string
	Role    string
}

type HookRPC interface {
	TarmakInitToken(cluster, env, role string, token *string) error
}

type hookRPCWrap struct {
	HookRPC
}

type initTokenRPC struct{}

func (i *initTokenRPC) TarmakInitToken(cluster, env, role string, token *string) error {
	*token = fmt.Sprintf("token-%s-%s-%s", cluster, env, role)
	return nil
}

func NewServerServe(receiver HookRPC) {
	s := rpc.NewServer()
	s.RegisterName("Hook", receiver)
	s.ServeConn(struct {
		io.Reader
		io.Writer
		io.Closer
	}{os.Stdin, os.Stdout,
		multiCloser{[]io.Closer{os.Stdout, os.Stdin, procCloser{}}},
	})
}

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

// Close the underlying process by killing it.
// If there's no such process, then commit suicide.
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

func FakeTarmak() {
	l := log.New(os.Stderr, "", 0)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			l.Printf("done: %s\n", sig)
			os.Exit(1)
		}
	}()

	NewServerServe(hookRPCWrap{&initTokenRPC{}})
}

//func FakeTarmak() {
//	initToken := new(InitToken)
//	err := rpc.Register(initToken)
//	if err != nil {
//		panic(err)
//	}
//	rpc.HandleHTTP()
//
//	listener, err := net.Listen("tcp", ":1234")
//	if err != nil {
//		panic(err)
//	}
//
//	err = http.Serve(listener, nil)
//	if err != nil {
//		panic(err)
//	}
//}
