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
		Backoff: backoff.DefaultConfig,
	}
	return grpc.DialContext(ctx, addr, grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithConnectParams(cp),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}))
}
