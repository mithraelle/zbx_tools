package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"github.com/mithraelle/zbx_tools/go/pkg/sender/grpcsender"
	"github.com/mithraelle/zbx_tools/go/pkg/sender/mocksender"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var port string
	flag.StringVar(&port, "p", "50000", "server port")
	flag.Parse()

	serverAddr := "localhost:" + port
	fmt.Println("Target agent: ", serverAddr)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	grpcSender, err := grpcsender.NewGRPCSender(serverAddr, opts)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan sender.Item)
	go mocksender.ThrowDice(ctx, ch, 5)

	iBuffer := sender.NewItemBuffer()
	iBuffer.Timeout = 15 * time.Second
	go iBuffer.Read(ctx, ch, grpcSender)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
