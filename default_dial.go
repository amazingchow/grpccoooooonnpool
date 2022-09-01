package grpccoooooonnpool

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// NOTE: Tune the default values to provide better system's throughput.
var (
	DefaultMaxSendMsgSize int = 8 * 1024 * 1024
	DefaultMaxRecvMsgSize int = 8 * 1024 * 1024

	DefaultInitialWindowSize     int32 = 128 * 1024 * 1024
	DefaultInitialConnWindowSize int32 = 128 * 1024 * 1024

	DefaultDialTimeout       = 5 * time.Second
	DefaultMinConnectTimeout = 3 * time.Second
	DefaultKeepAliveTime     = 10 * time.Second
	DefaultKeepAliveTimeout  = 3 * time.Second
)

// DefaultDialWithInsecure returns an insecure grpc client connection with default settings.
func DefaultDialWithInsecure(serverAddr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultDialTimeout)
	defer cancel()

	return grpc.DialContext(ctx, serverAddr,
		// Dial blocks until the underlying TCP connection is up.
		grpc.WithBlock(),
		// Dial disables transport security for the underlying TCP connection.
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(DefaultMaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(DefaultMaxRecvMsgSize)),
		// Dial sets the value for initial window size on a stream.
		grpc.WithInitialWindowSize(DefaultInitialWindowSize),
		// Dial sets the value for initial window size on a tcp connection.
		grpc.WithInitialConnWindowSize(DefaultInitialConnWindowSize),
		// Dial specifies the options for connection backoff.
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: DefaultMinConnectTimeout,
		}),
		// Dial specifies keepalive parameters for the client transport.
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                DefaultKeepAliveTime,
			Timeout:             DefaultKeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	)
}

// DefaultDial returns a secure grpc client connection with default settings.
func DefaultDial(serverAddr, serverName, certFile string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultDialTimeout)
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
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(DefaultMaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(DefaultMaxRecvMsgSize)),
		// Dial sets the value for initial window size on a stream.
		grpc.WithInitialWindowSize(DefaultInitialWindowSize),
		// Dial sets the value for initial window size on a tcp connection.
		grpc.WithInitialConnWindowSize(DefaultInitialConnWindowSize),
		// Dial specifies the options for connection backoff.
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: DefaultMinConnectTimeout,
		}),
		// Dial specifies keepalive parameters for the client transport.
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                DefaultKeepAliveTime,
			Timeout:             DefaultKeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	)
}
