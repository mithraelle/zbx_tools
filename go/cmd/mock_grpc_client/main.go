package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/mithraelle/zbx_tools/go/pb/agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func pushValue(client pb.ZbxAgentClient, value *pb.ZbxValue) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.PushValue(ctx, value)
	if err != nil {
		log.Fatalf("client.SendValue failed: %v", err)
	}
	log.Println(resp)
}

func main() {
	var key, value, port string
	flag.StringVar(&port, "p", "50000", "server port")
	flag.Parse()

	serverAddr := "localhost:" + port
	fmt.Println("Target agent: ", serverAddr)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewZbxAgentClient(conn)

	fmt.Println("Enter key value pairs")
	for {
		_, err := fmt.Scanf("%v %v", &key, &value)
		if err != nil {
			fmt.Println("Error reading input: ", err.Error())
		} else {
			fmt.Printf("Key: %v, Value: %v\n", key, value)
			_, err = client.PushValue(context.Background(), &pb.ZbxValue{Key: key, Value: value, Ts: int32(time.Now().Unix())})
			if err != nil {
				log.Fatalf("client.PushValue failed: %v", err)
			}
		}
	}
}
