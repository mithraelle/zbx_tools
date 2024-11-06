package agent

import (
	"context"
	pb "github.com/mithraelle/zbx_tools/go/pb/agent"
	"google.golang.org/grpc"
	"log"
	"net"
)

type AgentServer struct {
	pb.UnimplementedZbxAgentServer

	OutChan chan<- *pb.ZbxValue
}

func (z AgentServer) PushValue(ctx context.Context, value *pb.ZbxValue) (*pb.ZbxValueAck, error) {
	log.Println("AgentServer.SendValue: ", value.Key, value.Value, value.Ts)
	z.OutChan <- value
	return &pb.ZbxValueAck{Result: true}, nil
}

func (z AgentServer) PushValues(ctx context.Context, list *pb.ListZbxValue) (*pb.ZbxValueAck, error) {
	for _, value := range list.Values {
		log.Println("AgentServer.SendValue: ", value.Key, value.Value, value.Ts)
		z.OutChan <- value
	}
	return &pb.ZbxValueAck{Result: true}, nil
}

func RunAgent(ctx context.Context, lis net.Listener, opts []grpc.ServerOption, ch chan<- *pb.ZbxValue) {
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterZbxAgentServer(grpcServer, AgentServer{OutChan: ch})
	go grpcServer.Serve(lis)
	<-ctx.Done()
	grpcServer.GracefulStop()
}
