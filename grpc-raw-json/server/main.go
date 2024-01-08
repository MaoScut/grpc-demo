package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/MaoScut/grpc-demo/grpc-raw-json/proto/greeter"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

type GreeterService struct {
	greeter.UnimplementedGreeterServer
	avaliable bool
}

func NewGreeterService() *GreeterService {
	return &GreeterService{
		avaliable: true,
	}
}

func (s *GreeterService) Enable() {
	s.avaliable = true
}

func (s *GreeterService) Disable() {
	s.avaliable = false
}

func (s *GreeterService) SayHello(ctx context.Context, req *greeter.HelloRequest) (res *greeter.HelloReply, err error) {
	zap.L().Info("receive", zap.Any("req", req))
	rawBytes := `{"a": {"b": {"c": "a"}}}`
	st := &structpb.Struct{}
	err = json.Unmarshal([]byte(rawBytes), st)
	if err != nil {
		err = fmt.Errorf("json unmarshal: %w", err)
		return
	}
	res = &greeter.HelloReply{
		Message: st,
	}
	return
}

func main() {
	// exitCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	// defer stop()
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l.Named("server"))
	port := flag.String("grpc-port", ":9100", "server grpc port")
	httpPort := flag.String("http-port", ":9101", "server http port")
	flag.Parse()
	l.Info("start")
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}
	// httpList, err := net.Listen("tcp", *httpPort)
	// if err != nil {
	// 	panic(err)
	// }
	grpcServer := grpc.NewServer()
	greeterService := NewGreeterService()
	greeter.RegisterGreeterServer(grpcServer, greeterService)
	go func() {
		zap.L().Info("start grpc server", zap.Any("addr", port))
		err := grpcServer.Serve(lis)
		if err != nil {
			zap.L().Error("serve", zap.Error(err))
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		*port,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = greeter.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    *httpPort,
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on", *httpPort)
	log.Fatalln(gwServer.ListenAndServe())
}
