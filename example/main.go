package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	common_proto "github.com/th2-net/th2-common-go/proto"
	common_f "github.com/th2-net/th2-common-go/schema/factory"
	"github.com/th2-net/th2-common-go/schema/modules/grpcModule"
	"github.com/th2-net/th2-common-go/schema/modules/mqModule"
	"github.com/th2-net/th2-demo-script-go.git/act"
	"github.com/th2-net/th2-demo-script-go.git/check1"
	"github.com/th2-net/th2-demo-script-go.git/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

func generateRandomClordID(n int) string {
	str := ""
	for i := 0; i < n; i++ {
		randDigit := rand.Intn(10)
		str += strconv.Itoa(randDigit)
	}

	return str
}

func genrateSecondaryRandomClordID(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	// 1) Initialize NewFactory instance and register modules for grpc and mq .
	factory := common_f.NewFactory()
	if err := factory.Register(mqModule.NewRabbitMQModule); err != nil {
		panic(err)
	}
	if err := factory.Register(grpcModule.NewGrpcModule); err != nil {
		panic(err)
	}

	rabbitModule, err := mqModule.ModuleID.GetModule(factory)
	if err != nil {
		panic("no module")
	} else {
		fmt.Println("module found", reflect.TypeOf(rabbitModule))
	}

	grpc_module, err := grpcModule.ModuleID.GetModule(factory)
	if err != nil {
		panic("no module")
	} else {
		fmt.Println("module found", reflect.TypeOf(grpc_module))
	}

	// 2) Initialize grpc router services to work with act and check1 boxes.
	grpc_router := grpc_module.GrpcRouter
	con, conErr := grpc_router.GetConnection()
	if conErr != nil {
		log.Fatalf(conErr.Error())
	}
	//defer con.Close() //doesnt have that func
	actCl := act.NewActClient(con)
	check := check1.NewCheck1Client(con)

	//# 3) Initialize mq router to work with estore.
	estore := rabbitModule.MqEventRouter

	// 4) Create root Event for report.

	current_timestamp := timestamppb.Now()
	log.Printf("current time: %v \n", current_timestamp)

	eventID := common_proto.EventID{Id: uuid.New().String()}
	event := common_proto.Event{
		Id:             &eventID,
		Name:           "Raw send example",
		Status:         common_proto.EventStatus_SUCCESS,
		StartTimestamp: current_timestamp,
	}
	fmt.Println(current_timestamp)

	// 5) Add this Event to EventBatch.
	eventBatch := common_proto.EventBatch{Events: []*common_proto.Event{&event}}
	fmt.Println(len(eventBatch.Events))

	// 6) Send EventBatch to estore.
	fail := estore.SendAll(&eventBatch, "group")
	log.Println(fail)

	//// 7) Create the NewOrderSingle Message.
	// ClOrdID stored separately for future use.
	clordid := generateRandomClordID(7)
	fmt.Println(clordid)
	secondaryClordid := genrateSecondaryRandomClordID(7)
	fmt.Println(secondaryClordid)

	tradingPartyFields := map[string]*common_proto.Value{
		"NoPartyIDs": {Kind: &common_proto.Value_ListValue{ListValue: &common_proto.ListValue{
			Values: []*common_proto.Value{{Kind: &common_proto.Value_MessageValue{MessageValue: &common_proto.Message{
				Metadata: &common_proto.MessageMetadata{MessageType: "TradingParty_NoPartyIDs"},
				Fields: map[string]*common_proto.Value{
					"PartyID":       {Kind: &common_proto.Value_SimpleValue{SimpleValue: "Trader1"}},
					"PartyIDSource": {Kind: &common_proto.Value_SimpleValue{SimpleValue: "D"}},
					"PartyRole":     {Kind: &common_proto.Value_SimpleValue{SimpleValue: "76"}},
				},
			}}},
				{Kind: &common_proto.Value_MessageValue{MessageValue: &common_proto.Message{
					Metadata: &common_proto.MessageMetadata{MessageType: "TradingParty_NoPartyIDs"},
					Fields: map[string]*common_proto.Value{
						"PartyID":       {Kind: &common_proto.Value_SimpleValue{SimpleValue: "0"}},
						"PartyIDSource": {Kind: &common_proto.Value_SimpleValue{SimpleValue: "D"}},
						"PartyRole":     {Kind: &common_proto.Value_SimpleValue{SimpleValue: "3"}},
					},
				},
				}},
			}}},
		},
	}

	fields := map[string]*common_proto.Value{
		"Side":             {Kind: &common_proto.Value_SimpleValue{SimpleValue: "1"}},
		"SecurityID":       {Kind: &common_proto.Value_SimpleValue{SimpleValue: "INSTR1"}},
		"SecurityIDSource": {Kind: &common_proto.Value_SimpleValue{SimpleValue: "8"}},
		"OrdType":          {Kind: &common_proto.Value_SimpleValue{SimpleValue: "2"}},
		"AccountType":      {Kind: &common_proto.Value_SimpleValue{SimpleValue: "1"}},
		"OrderCapacity":    {Kind: &common_proto.Value_SimpleValue{SimpleValue: "A"}},
		"OrderQty":         {Kind: &common_proto.Value_SimpleValue{SimpleValue: "100"}},
		"Price":            {Kind: &common_proto.Value_SimpleValue{SimpleValue: "10"}},
		"ClOrdID":          {Kind: &common_proto.Value_SimpleValue{SimpleValue: "123"}}, //random in py
		"SecondaryClOrdID": {Kind: &common_proto.Value_SimpleValue{SimpleValue: "2"}},   //random in py
		"TransactTime":     {Kind: &common_proto.Value_SimpleValue{SimpleValue: time.Now().Format(time.RFC3339)}},
		"TradingParty":     {Kind: &common_proto.Value_MessageValue{MessageValue: &common_proto.Message{Fields: tradingPartyFields}}},
	}

	msg := common_proto.Message{Metadata: &common_proto.MessageMetadata{
		MessageType: "NewOrderSingle",
		Id:          &common_proto.MessageID{ConnectionId: &common_proto.ConnectionID{SessionAlias: "demo-conn1"}},
	},
		Fields: fields,
	}

	//
	// 8) Create instance of PlaceMessageRequest class - grpc message object which used for calls to the act.

	request := act.PlaceMessageRequest{
		Message:       &msg,
		ParentEventId: &eventID,
		Description:   "User places an order.",
	}

	// 9) Call method placeOrderFix from the act interface.
	resp, failed := actCl.PlaceOrderFIX(context.Background(), &request)
	if failed != nil {
		log.Fatalf("could not send order: %v", err)
	}
	log.Printf("order sent. response: %s", resp)
	// 10) Create MessageFilter object - mask or pattern of message verification.

	message_filter := common_proto.MessageFilter{
		MessageType: "ExecutionReport",
		Fields: map[string]*common_proto.ValueFilter{
			"ClOrdID":   {Kind: &common_proto.ValueFilter_SimpleFilter{SimpleFilter: clordid}, Key: true},
			"Side":      {Kind: &common_proto.ValueFilter_SimpleFilter{SimpleFilter: "1"}},
			"Price":     {Operation: common_proto.FilterOperation_NOT_EMPTY},
			"LeavesQty": {Kind: &common_proto.ValueFilter_SimpleFilter{SimpleFilter: "0"}, Operation: common_proto.FilterOperation_NOT_EQUAL},
			"OrderID":   {Kind: &common_proto.ValueFilter_SimpleFilter{SimpleFilter: resp.ResponseMessage.Fields["OrderID"].GetSimpleValue()}},
		},
	}
	// 11) Create instance of CheckRuleRequest class - grpc message object which used for calls to the check1.
	check1_request := proto.CheckRuleRequest{
		ConnectivityId: &proto.ConnectionID{SessionAlias: "demo-conn1"},
		//filter = message_filter,
		//Checkpoint: response.checkpointId,
		Kind: &proto.CheckRuleRequest_Filter{Filter: &message_filter},

		//Kind: &proto.CheckRuleRequest_Filter{Filter: &message_filter},
		//Kind:          common_proto.MessageFilter{message_filter},
		Checkpoint:    resp.CheckpointId,
		Timeout:       3000,
		ParentEventId: &eventID,
		Description:   "User receives the ExecutionReport message.",
	}

	// 12) Call method submitCheckRule from the check1 interface.
	check1_response, checkErr := check.SubmitCheckRule(check1_request)
	if checkErr != nil {
		log.Fatal(checkErr)
	}
	log.Printf("check1_response : %v ", check1_response)

	// 13) Close CommonFactory.
	factory.Close()

}
