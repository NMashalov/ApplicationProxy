package main

import (
	"context"
	"flag"
	"time"

	"example.com/m/pkg/dummy"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var port string
	flag.StringVar(&port, "port", ":8000", "running port")
	flag.Parse()
	go dummy.DummyGinServer().Run(port)
	<-ctx.Done()
	panic("I'm done")
}
