package zbxgrpcagent

import (
	"context"
	"github.com/mithraelle/zbx_tools/go/pkg/zbxgrpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

type ZbxGrpcAgent struct {
	zbxgrpc.UnimplementedZbxSenderServer

	OutChan chan<- *zbxgrpc.ZbxValue
}

func (z ZbxGrpcAgent) PushValue(ctx context.Context, value *zbxgrpc.ZbxValue) (*zbxgrpc.ZbxValueAck, error) {
	log.Println("ZbxGrpcAgent.SendValue: ", value.Key, value.Value, value.Ts)
	z.OutChan <- value
	return &zbxgrpc.ZbxValueAck{Result: true}, nil
}

func (z ZbxGrpcAgent) PushValues(ctx context.Context, list *zbxgrpc.ListZbxValue) (*zbxgrpc.ZbxValueAck, error) {
	for _, value := range list.Values {
		log.Println("ZbxGrpcAgent.SendValue: ", value.Key, value.Value, value.Ts)
		z.OutChan <- value
	}
	return &zbxgrpc.ZbxValueAck{Result: true}, nil
}

func RunZBXGrpcAgent(ctx context.Context, lis net.Listener, opts []grpc.ServerOption, ch chan<- *zbxgrpc.ZbxValue) {
	grpcServer := grpc.NewServer(opts...)
	zbxgrpc.RegisterZbxSenderServer(grpcServer, ZbxGrpcAgent{OutChan: ch})
	go grpcServer.Serve(lis)
	<-ctx.Done()
	grpcServer.GracefulStop()
}
