package proxy

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"os/exec"
)

type StopServer func() error

type Starter interface {
	StartServer(ctx context.Context) (net.Conn, StopServer, error)
}

type starter struct {
	binaryPath string
}

func NewStarter(binaryPath string) Starter {
	return &starter{binaryPath}
}

func randomPort() int {
	r, _ := rand.Int(rand.Reader, big.NewInt(1000))
	return int(r.Int64()) + 60_000
}

func (s *starter) StartServer(ctx context.Context) (net.Conn, StopServer, error) {
	runningPort := fmt.Sprintf("%d:", randomPort)
	cmd := exec.CommandContext(
		ctx, s.binaryPath,
		"--port", runningPort)
	err := cmd.Start()
	if err != nil {
		return nil, nil, err
	}
	conn, err := net.Dial("tcp", runningPort)
	if err != nil {
		cmd.Cancel()
		return nil, nil, err
	}
	return conn, cmd.Cancel, nil
}
