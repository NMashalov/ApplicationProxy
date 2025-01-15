package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

type Proxy interface {
	HandleConnection(conn net.Conn)
	Poll() error
}

type proxy struct {
	port       int
	serverPool ServerPool
}

func NewProxy(port int, serverPool ServerPool) Proxy {
	return &proxy{
		port,
		serverPool,
	}
}
func (p *proxy) Poll() error {
	servingPort := fmt.Sprintf(":%d", p.port)
	listener, err := net.Listen("tcp", servingPort)
	if err != nil {
		return err
	}
	go func() {
		for range time.NewTicker(time.Second).C {
			p.serverPool.Refresh(context.Background())
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go p.HandleConnection(conn)
	}
}
