package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaoScut/grpc-demo/grpc-lb/app/appproto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GreeterService struct {
	appproto.UnimplementedGreeterServer
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

func (s *GreeterService) SayHello(ctx context.Context, req *appproto.HelloRequest) (res *appproto.HelloReply, err error) {
	zap.L().Info("receive", zap.Any("req", req))
	res = &appproto.HelloReply{
		Message: "h1",
	}
	return
}

func (s *GreeterService) SayHelloLoop(stream appproto.Greeter_SayHelloLoopServer) (err error) {
	if !s.avaliable {
		err = status.Error(codes.Unavailable, "server unavaliable")
		return
	}
	counter := 0
	var req *appproto.HelloLoopRequest
	req, err = stream.Recv()
	if err != nil {
		err = fmt.Errorf("stream.Recv: %w", err)
		return
	}
	zap.L().Info("receive", zap.Any("req", req))
	for {
		err = stream.Send(&appproto.HelloLoopReply{
			Message: fmt.Sprintf("hi, %d", counter),
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.Unavailable {
					err = nil
					return
				}
			}
			zap.L().Error("stream.Send", zap.Error(err))
		} else {
			zap.L().Info("send msg")
		}
		time.Sleep(time.Second * 3)
		counter = counter + 1
	}
}

func main() {
	exitCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer stop()
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l.Named("server"))
	port := flag.String("grpc-port", ":9100", "server grpc port")
	httpPort := flag.String("http-port", ":9101", "server http port")
	flag.Parse()
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}
	httpList, err := net.Listen("tcp", *httpPort)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	greeterService := NewGreeterService()
	appproto.RegisterGreeterServer(grpcServer, greeterService)
	go func() {
		zap.L().Info("start grpc server", zap.Any("addr", port))
		err := grpcServer.Serve(lis)
		if err != nil {
			zap.L().Error("serve", zap.Error(err))
		}
	}()
	http.HandleFunc("/enable", func(w http.ResponseWriter, r *http.Request) {
		greeterService.Enable()
		w.Write([]byte("enable done"))
	})
	http.HandleFunc("/disable", func(w http.ResponseWriter, r *http.Request) {
		greeterService.Disable()
		w.Write([]byte("disable done"))
	})
	go func() {
		zap.L().Info("start http server", zap.Any("addr", httpPort))
		err := http.Serve(httpList, nil)
		if err != nil {
			zap.L().Error("http serve", zap.Error(err))
		}
	}()
	<-exitCtx.Done()
}
