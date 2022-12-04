package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/MaoScut/grpc-demo/grpc-lb/app/appproto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l.Named("client"))
	addr := flag.String("server-addr", "", "")
	flag.Parse()
	conn, err := grpc.Dial(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`),
	)
	if err != nil {
		panic(err)
	}
	zap.L().Info("start", zap.Any("server-addr", addr))
	client := appproto.NewGreeterClient(conn)
	ctx := context.Background()
	sayHelloClient, err := client.SayHelloLoop(ctx)
	if err != nil {
		zap.L().Error("say hello", zap.Error(err))
		return
	}
	err = sayHelloClient.Send(&appproto.HelloLoopRequest{
		Name: "c1",
	})
	if err != nil {
		zap.L().Error("send", zap.Error(err))
		return
	}
	for {
		res, err := sayHelloClient.Recv()
		if err != nil {
			zap.L().Error("receive", zap.Error(err))
			zap.L().Info("reconnect")
			sayHelloClient, err = reconnect(ctx, client)
			if err != nil {
				zap.L().Error("reconnect", zap.Error(err))
			}
		} else {
			zap.L().Info("receive", zap.Any("res", res))
		}
		time.Sleep(time.Second * 3)
	}
}

func reconnect(parentCtx context.Context, client appproto.GreeterClient) (sayHelloClient appproto.Greeter_SayHelloLoopClient, err error) {
	// ctx, cancel := context.WithTimeout(parentCtx, time.Second*5)
	// defer cancel()
	ctx := parentCtx
	sayHelloClient, err = client.SayHelloLoop(ctx)
	if err != nil {
		err = fmt.Errorf("sayHelloLoop: %w", err)
		return
	}
	zap.L().Info("reconnect success")
	err = sayHelloClient.Send(&appproto.HelloLoopRequest{
		Name: "c1",
	})
	if err != nil {
		err = fmt.Errorf("send: %w", err)
		return
	}
	zap.L().Info("reconnect send success")
	return
}
