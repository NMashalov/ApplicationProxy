package main

import (
	"log"
	"net"
	"time"
)

// slow echo server
func main() {
	listener, err := net.Listen("tcp", ":9009")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go func() {
			defer conn.Close()
			payloadBuffer := make([]byte, 4096)
			n, err := conn.Read(payloadBuffer)
			if err != nil {
				log.Println(err)
				return
			}
			clientPayload := payloadBuffer[:n]
			time.Sleep(5 * time.Second)
			conn.Write(clientPayload)
		}()
	}
}
