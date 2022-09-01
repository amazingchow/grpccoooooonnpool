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
	opts := []gpool.PoolSettingsOption{
		gpool.WithAddr("127.0.0.1:18889"),
		gpool.WithDialFunc(gpool.DefaultDialWithInsecure),
		gpool.WithMaxIdles(8),
		gpool.WithMaxStreams(64),
	}
	p, err := gpool.NewGrpcClientConnPool(opts...)
	if err != nil {
		log.Fatalf("failed to create grpccoooooonnpool, err: %v\n", err)
	} else {
		log.Println("create grpccoooooonnpool")
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := p.PickOne(true /* wait */, 200 /* waitTime */)
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

		conn.Close()
	}

	http.HandleFunc("/performance", handler)
	log.Printf("run simple-proxy on ':18888'\n")
	http.ListenAndServe(":18888", nil)
}
