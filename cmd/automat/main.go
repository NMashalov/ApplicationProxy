package main

import (
	"log"

	"example.com/m/pkg/proxy"
)

func main() {
	starter := proxy.NewStarter("build/dummy")
	log.Fatal(proxy.NewProxy(4000, starter).Poll())
}
