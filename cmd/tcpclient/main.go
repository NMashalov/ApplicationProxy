package main

import (
	"log"
	"net"
	"time"
)

// note that tcp things like mangle algo corrupts
// eny attempt to send small amount of bytes
// https://en.wikipedia.org/wiki/Nagle%27s_algorithm
func main() {
	conn, err := net.Dial("tcp", "localhost:9009")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := conn.Write([]byte("Why you leave me!")); err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	if _, err := conn.Write([]byte("Tell me!")); err != nil {
		log.Fatal(err)
	}
	payloadBuffer := make([]byte, 4096)
	n, err := conn.Read(payloadBuffer)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(payloadBuffer[:n]))
}
