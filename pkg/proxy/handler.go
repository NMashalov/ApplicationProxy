package proxy

import (
	"log"
	"net"
)

func readConnection(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 40_000)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func communicationRound(clientConn, serverConn net.Conn) error {
	clientPayload, err := readConnection(clientConn)
	if err != nil {
		return err
	}
	if _, err := serverConn.Write(clientPayload); err != nil {
		return err
	}
	serverPayload, err := readConnection(serverConn)
	if err != nil {
		return err
	}
	if _, err := clientConn.Write(serverPayload); err != nil {
		return err
	}
	return nil
}

func (p *proxy) HandleConnection(clientConn net.Conn) {
	defer clientConn.Close()
	serverConn, err := p.serverPool.ProvideConnection()
	if err != nil {
		log.Print(err)
		return
	}
	defer serverConn.Close()
	if err := communicationRound(clientConn, serverConn); err != nil {
		log.Print(err)
	}
}
