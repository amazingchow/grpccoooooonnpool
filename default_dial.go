package gpool

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
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

// DefaultDialWithInsecure returns an insecure grpc connection with default settings.
func DefaultDialWithInsecure(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()

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
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: MinConnectTimeout,
		}),
		// Dial specifies keepalive parameters for the client transport.
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}))
}

// DefaultDial returns a secure grpc connection with default settings.
func DefaultDial(serverAddr, serverName, certFile string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()

	creds, err := credentials.NewClientTLSFromFile(certFile, serverName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create TLS using root-cert '%s'", certFile)
	}

	return grpc.DialContext(ctx, serverAddr,
		// Dial blocks until the underlying TCP connection is up.
		grpc.WithBlock(),
		// Dial configures a connection level security credentials (e.g., TLS/SSL).
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		// Dial sets the value for initial window size on a stream.
		grpc.WithInitialWindowSize(InitialWindowSize),
		// Dial sets the value for initial window size on on a TCP connection.
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		// Dial specifies the options for connection backoff.
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: MinConnectTimeout,
		}),
		// Dial specifies keepalive parameters for the client transport.
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}))
}
