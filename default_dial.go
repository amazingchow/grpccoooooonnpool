package gpool

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/keepalive"
)

// tune them to provide better system's throughput
const (
	DialTimeout           = 5 * time.Second
	MinConnectTimeout     = 3 * time.Second
	MaxSendMsgSize        = 8 * 1024 * 1024
	MaxRecvMsgSize        = 8 * 1024 * 1024
	InitialWindowSize     = 128 * 1024 * 1024
	InitialConnWindowSize = 128 * 1024 * 1024
	KeepAliveTime         = 10 * time.Second
	KeepAliveTimeout      = 3 * time.Second
)

// DefaultDialWithInsecure returns a insecure grpc connection with default settings.
func DefaultDialWithInsecure(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()

	cp := grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: MinConnectTimeout,
	}
	return grpc.DialContext(ctx, addr,
		// Dial blocks until the underlying TCP connection is up.
		grpc.WithBlock(),
		// Dial disables transport security for the underlying TCP connection.
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		// Dial sets the value for initial window size on a stream.
		grpc.WithInitialWindowSize(InitialWindowSize),
		// Dial sets the value for initial window size on on a TCP connection.
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		// Dial specifies the options for connection backoff.
		grpc.WithConnectParams(cp),
		// Dial specifies keepalive parameters for the client transport.
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}))
}
