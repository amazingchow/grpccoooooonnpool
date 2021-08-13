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
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (srv *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	time.Sleep(50 * time.Millisecond)
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s", req.Name)}, nil
}

func main() {
	l, err := net.Listen("tcp", ":18889")
	if err != nil {
		log.Fatalf("failed to create tcp listener, err: %v\n", err)
	}

	srv := grpc.NewServer()
	pb.RegisterGreeterServer(srv, &server{})
	reflection.Register(srv)

	quitCh := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err = srv.Serve(l); err != nil {
			log.Printf("failed to serve, err: %v\n", err)
		}
		quitCh <- struct{}{}
	}()

	<-sigCh
	srv.GracefulStop()
	<-quitCh
}
