package grpcsender

import (
	"context"
	pb "github.com/mithraelle/zbx_tools/go/pb/agent"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"google.golang.org/grpc"
	"time"
)

const GRPCTimeout = 10 * time.Second

type GRPCSender struct {
	client pb.ZbxAgentClient
	conn   *grpc.ClientConn

	Timeout time.Duration
}

func NewGRPCSender(serverAddr string, opts []grpc.DialOption) (*GRPCSender, error) {
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, err
	}
	client := pb.NewZbxAgentClient(conn)

	return &GRPCSender{client: client, conn: conn, Timeout: GRPCTimeout}, nil
}

func (s *GRPCSender) Send(items []sender.Item, errorSink chan<- sender.ItemSendError) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	grpcVals := pb.ListZbxValue{}
	for _, item := range items {
		grpcVals.Values = append(grpcVals.Values, &pb.ZbxValue{Key: item.Key, Value: item.Value, Ts: int32(item.Clock)})
	}

	_, err := s.client.PushValues(ctx, &grpcVals)
	if err != nil && errorSink != nil {
		errorSink <- sender.ItemSendError{Items: items, Err: err}
	}
}
