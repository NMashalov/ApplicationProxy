package proxy

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"os/exec"
	"time"

	"github.com/samber/lo"
)

type EphemeralServer struct {
	port          int
	isLiveChannel chan int
}

type Starter interface {
	StartServer(ctx context.Context) (*EphemeralServer, error)
}

type starter struct {
	binaryPath string
}

func NewStarter(binaryPath string) Starter {
	return &starter{binaryPath}
}

func (s *starter) StartServer(ctx context.Context) (*EphemeralServer, error) {
	port := randomPort()
	runningPort := fmt.Sprintf(":%d", port)
	cmd := exec.CommandContext(
		ctx, s.binaryPath,
		"--port", runningPort)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	isLiveChan := make(chan int, 0)
	go func() {
		cmd.Wait()
		close(isLiveChan)
	}()
	time.Sleep(time.Second)
	return &EphemeralServer{
		port,
		isLiveChan,
	}, nil
}

type ServerPool interface {
	Refresh(ctx context.Context)
	ProvideConnection() (net.Conn, error)
}
type serverPool struct {
	starter    Starter
	maxRunning uint
	servers    []*EphemeralServer
}

func NewServerPool(starter Starter,
	maxRunning uint) ServerPool {
	return &serverPool{
		starter, maxRunning, []*EphemeralServer{},
	}
}

func (sp *serverPool) Refresh(ctx context.Context) {
	actualServers := lo.Filter(
		sp.servers,
		func(server *EphemeralServer, _ int) bool {
			select {
			case <-server.isLiveChannel:
				return false
			default:
				return true
			}
		},
	)
	newServers := lo.FilterMap(
		lo.Range(
			int(sp.maxRunning)-len(sp.servers),
		),
		func(_ int, _ int) (*EphemeralServer, bool) {
			server, err := sp.starter.StartServer(ctx)
			if err != nil {
				return nil, false
			}
			return server, true
		},
	)
	sp.servers = append(actualServers, newServers...)
}

func randomPort() int {
	r, _ := rand.Int(rand.Reader, big.NewInt(1000))
	return int(r.Int64()) + 60_000
}

func (sp *serverPool) randomBalancer() *EphemeralServer {
	r, _ := rand.Int(
		rand.Reader,
		big.NewInt(int64(len(sp.servers))),
	)
	return sp.servers[r.Int64()]
}

func (sp *serverPool) ProvideConnection() (net.Conn, error) {
	selectedJobs := sp.randomBalancer()
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", selectedJobs.port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
