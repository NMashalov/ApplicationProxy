package main

import (
	"flag"

	"example.com/m/pkg/dummy"
)

func main() {
	var port string
	flag.StringVar(&port, "port", ":8000", "running port")
	flag.Parse()
	dummy.DummyGinServer().Run(port)
}
