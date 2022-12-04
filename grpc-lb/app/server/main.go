package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

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

func (s *GreeterService) SayHelloLoop(stream appproto.Greeter_SayHelloLoopServer) (err error) {
	counter := 0
	var req *appproto.HelloLoopRequest
	req, err = stream.Recv()
	if err != nil {
		err = fmt.Errorf("stream.Recv: %w", err)
		return
	}
	zap.L().Info("say hello loop, receive", zap.Any("req", req))
	for {
		err = stream.Send(&appproto.HelloLoopReply{
			Message: fmt.Sprintf("hi, %d", counter),
		})
		if err != nil {
			zap.L().Error("stream.Send", zap.Error(err))
		}
		time.Sleep(time.Second * 3)
		counter = counter + 1
	}
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
