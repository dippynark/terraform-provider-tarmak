package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

type tarmakRPC struct{}

func (i *tarmakRPC) BastionStatus(args string, reply *string) error {
	fmt.Printf("BastionStatus called\n")
	*reply = "running"
	return nil
}

func main() {
	log.Println("Starting fake Tarmak")
	ln, err := net.Listen("unix", "tarmak.sock")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go accept(fd)
	}
}

func accept(conn net.Conn) {

	s := rpc.NewServer()
	s.RegisterName("Tarmak", &tarmakRPC{})

	fmt.Printf("Connection made\n")

	s.ServeConn(struct {
		io.Reader
		io.Writer
		io.Closer
	}{conn, conn, conn})

}
