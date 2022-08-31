package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	helloworldpb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

type server struct {
	helloworldpb.UnimplementedGreeterServer
}

func (srv *server) SayHello(ctx context.Context, req *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	time.Sleep(50 * time.Millisecond)
	return &helloworldpb.HelloReply{Message: fmt.Sprintf("Hello %s", req.Name)}, nil
}

func main() {
	l, err := net.Listen("tcp", ":18889")
	if err != nil {
		log.Fatalf("failed to create tcp listener, err: %v\n", err)
	}

	srv := grpc.NewServer()
	helloworldpb.RegisterGreeterServer(srv, &server{})
	reflection.Register(srv)

	done := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err = srv.Serve(l); err != nil {
			log.Printf("failed to serve, err: %v\n", err)
		}
		done <- struct{}{}
	}()

	<-sigCh
	srv.GracefulStop()
	<-done
}
