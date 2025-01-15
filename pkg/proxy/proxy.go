package proxy

import (
	"fmt"
	"log"
	"net"
)

type Proxy interface {
	HandleConnection(conn net.Conn)
	Poll() error
}

type proxy struct {
	port    int
	starter Starter
}

func NewProxy(port int, starter Starter) Proxy {
	return &proxy{
		port,
		starter,
	}
}
func (p *proxy) Poll() error {
	servingPort := fmt.Sprintf(":%d", p.port)
	listener, err := net.Listen("tcp", servingPort)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go p.HandleConnection(conn)
	}
}
