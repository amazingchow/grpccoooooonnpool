package main

import (
	"context"
	"log"
	"net/http"
	"time"

	helloworldpb "google.golang.org/grpc/examples/helloworld/helloworld"

	gpool "github.com/amazingchow/grpccoooooonnpool"
)

func main() {
	p, err := gpool.NewGrpcConnPool("localhost:18889", gpool.PoolOptions{
		Dial:                 gpool.DefaultDialWithInsecure,
		MaxIdles:             32,
		MaxConcurrentStreams: 100,
	})
	if err != nil {
		log.Fatalf("failed to create grpccoooooonnpool, err: %v\n", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := p.PickOne(true /* wait */)
		if err != nil {
			log.Printf("failed to fetch conn from grpccoooooonnpool, err: %v\n", err)
			w.WriteHeader(500)
			return
		}

		cli := helloworldpb.NewGreeterClient(conn.Underlay())
		resp, err := cli.SayHello(ctx, &helloworldpb.HelloRequest{Name: "grpccoooooonnpool"})
		if err != nil {
			log.Printf("failed to do great, err: %v\n", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(resp.Message))

		_ = conn.Close()
	}

	http.HandleFunc("/performance", handler)
	_ = http.ListenAndServe("localhost:18888", nil)
}
