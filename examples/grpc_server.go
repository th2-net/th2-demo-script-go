package main

import (
	"context"
	act_grpc "github.com/th2-net/th2-demo-script-go.git/examples/proto"

	"fmt"
	cg "github.com/th2-net/th2-common-go/proto" //common proto
	common_f "github.com/th2-net/th2-common-go/schema/factory"
	"github.com/th2-net/th2-common-go/schema/modules/grpcModule"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"time"
)

const (
	grpcJsonFileName = "grpc.json"
)

type server struct {
	act_grpc.UnimplementedActServer
}

func (s *server) PlaceOrderFIX(_ context.Context, in *act_grpc.PlaceMessageRequest) (*act_grpc.PlaceMessageResponse, error) {
	return &act_grpc.PlaceMessageResponse{Status: &cg.RequestStatus{
		Message: "the order has been received",
	},
		ResponseMessage: in.Message}, nil //return the request content as the response
}

func registerService(registrar grpc.ServiceRegistrar) {
	act_grpc.RegisterActServer(registrar, &server{})
}

const lifetimeSec = 30

func main() {
	newFactory := common_f.NewFactory()
	if err := newFactory.Register(grpcModule.NewGrpcModule); err != nil {
		panic(err)
	}

	module, err := grpcModule.ModuleID.GetModule(newFactory)
	if err != nil {
		panic("no module")
	} else {
		fmt.Println("module found", reflect.TypeOf(module))
	}

	grpcRouter := module.GrpcRouter
	stopServerFunc, err := grpcRouter.StartServerAsync(registerService)
	if err != nil {
		log.Fatalf(err.Error())
	}

	time.Sleep(time.Second * lifetimeSec)

	stopServerFunc()

}
