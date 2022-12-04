package main

import (
	"context"
	"flag"
	"net"

	"github.com/MaoScut/grpc-demo/grpc-lb/app/appproto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GreeterService struct {
	appproto.UnimplementedGreeterServer
}

func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

func (s *GreeterService) SayHello(ctx context.Context, req *appproto.HelloRequest) (res *appproto.HelloReply, err error) {
	zap.L().Info("receive", zap.Any("req", req))
	res = &appproto.HelloReply{
		Message: "h1",
	}
	return
}

func main() {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	port := flag.String("prort", ":9100", "server port")
	flag.Parse()
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	appproto.RegisterGreeterServer(grpcServer, NewGreeterService())
	grpcServer.Serve(lis)
}
