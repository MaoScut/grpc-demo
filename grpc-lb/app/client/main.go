package main

import (
	"context"
	"flag"

	"github.com/MaoScut/grpc-demo/grpc-lb/app/appproto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	addr := flag.String("server-addr", "", "")
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := appproto.NewGreeterClient(conn)
	ctx := context.Background()
	res, err := client.SayHello(ctx, &appproto.HelloRequest{
		Name: "demo",
	})
	if err != nil {
		zap.L().Error("say hello", zap.Error(err))
		return
	}
	zap.L().Info("say hello", zap.Any("res", res))
}
