package grpccoooooonnpool

import (
	"google.golang.org/grpc"
)

// PoolSettings includes pool settings for pool initialization.
type PoolSettings struct {
	// Grpc server address to serve client connections.
	Addr string
	// An app-supplied-function for connection creation and configuration.
	Dial func(addr string) (conn *grpc.ClientConn, err error)
	// Maximum number of idle connections inside the pool, should be the power of 2.
	MaxIdles uint32
	// Maximum number of connections can be allocated by the pool at a given time, should be the power of 2.
	MaxActives uint32
	// Maximum number of concurrent streams attached to a single tcp connection, should be the power of 2.
	MaxStreams uint32
}

type PoolSettingsOption func(settings *PoolSettings)

func WithAddr(addr string) PoolSettingsOption {
	return func(settings *PoolSettings) {
		settings.Addr = addr
	}
}

func WithDialFunc(f func(addr string) (conn *grpc.ClientConn, err error)) PoolSettingsOption {
	return func(settings *PoolSettings) {
		settings.Dial = f
	}
}

func WithMaxIdles(n uint32) PoolSettingsOption {
	return func(settings *PoolSettings) {
		settings.MaxIdles = n
	}
}

func WithMaxActives(n uint32) PoolSettingsOption {
	return func(settings *PoolSettings) {
		settings.MaxActives = n
	}
}

func WithMaxStreams(n uint32) PoolSettingsOption {
	return func(settings *PoolSettings) {
		settings.MaxStreams = n
	}
}
