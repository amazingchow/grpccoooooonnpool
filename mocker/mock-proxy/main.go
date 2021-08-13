package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	gpool "github.com/amazingchow/photon-dance-gpool"
)

func main() {
	p, err := gpool.NewGrpcConnPool("localhost:18889", gpool.PoolOptions{
		Dial:                 gpool.DefaultDialWithInsecure,
		MaxIdles:             32,
		MaxActives:           64,
		MaxConcurrentStreams: 100,
		Reuse:                true,
	})
	if err != nil {
		log.Fatalf("failed to create gpool, err: %v\n", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := p.Get()
		if err != nil {
			log.Printf("failed to fetch conn from gpool, err: %v\n", err)
			w.WriteHeader(500)
			return
		}
		cli := pb.NewGreeterClient(conn.Underlay())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := cli.SayHello(ctx, &pb.HelloRequest{Name: "Grpc"})
		if err != nil {
			log.Printf("failed to do great, err: %v\n", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(resp.Message))
	}

	http.HandleFunc("/performance", handler)
	_ = http.ListenAndServe("localhost:18888", nil)
}
