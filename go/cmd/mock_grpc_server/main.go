package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/mithraelle/zbx_tools/go/pb/agent"
	"github.com/mithraelle/zbx_tools/go/pkg/agent"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	var port string
	flag.StringVar(&port, "p", "50000", "server port")
	flag.Parse()

	fmt.Printf("server port: %s\n", port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan *pb.ZbxValue, 100)

	lis, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go agent.RunAgent(ctx, lis, []grpc.ServerOption{}, ch)

	for v := range ch {
		fmt.Println(v)
	}
}
