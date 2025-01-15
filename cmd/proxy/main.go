package main

import (
	"log"

	"example.com/m/pkg/proxy"
)

func main() {
	starter := proxy.NewStarter("build/dummy")
	pool := proxy.NewServerPool(starter, 3)
	log.Fatal(proxy.NewProxy(4000, pool).Poll())
}
