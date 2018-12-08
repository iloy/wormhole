package server

import (
	"net"
	"time"

	proxyproto "github.com/armon/go-proxyproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/iloy/wormhole/protocol"
	"github.com/iloy/wormhole/topicmanager"
)

const (
	useProxyProtocolV1 = false
	pingTime           = 30
)

var (
	topicManager = topicmanager.New()
)

// StreamingServer :
type StreamingServer struct {
	server *grpc.Server
}

// rpcServer :
type rpcServer struct{}

// Serve :
func (ss *StreamingServer) Serve() error {
	lis := func() net.Listener {
		for {
			lis, err := net.Listen("tcp", ":18554")
			if err != nil {
				time.Sleep(time.Millisecond * time.Duration(100))
			} else {
				return lis
			}
		}
	}()

	if useProxyProtocolV1 {
		lis = &proxyproto.Listener{
			Listener:           lis,
			ProxyHeaderTimeout: time.Second * 1,
		}
	}

	kp := keepalive.ServerParameters{}
	kep := keepalive.EnforcementPolicy{}

	if pingTime > 20 {
		kp.MaxConnectionIdle = time.Second * 4
		//kp.MaxConnectionAgeGrace = time.Second * 4

		kp.Time = time.Second * pingTime
		kp.Timeout = time.Second * 20 // gRPC-go default value
		kep.MinTime = time.Second * (pingTime - 5)
		kep.PermitWithoutStream = true
	}

	server := grpc.NewServer(
		grpc.KeepaliveParams(kp),
		grpc.KeepaliveEnforcementPolicy(kep),
	)
	ss.server = server

	protocol.RegisterWormholeServer(server, &rpcServer{})
	go func() {
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return nil
}

// Stop :
func (ss *StreamingServer) Stop() {
	ss.server.GracefulStop()
}
