package main

import (
	//"github.com/th2-net/th2-common-go/schema/factory"
	"github.com/google/uuid"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	//"th2-demo-script-go/proto"
	"github.com/th2-net/th2-demo-script-go.git/proto"
)

func main() {
	//# 1) Initialize CommonFactory instance with paths to json configs as parameters.
	//	factory = CommonFactory(
	//	grpc_router_config_filepath="./configs/grpc.json",
	//	rabbit_mq_config_filepath="./configs/rabbit.json",
	//	mq_router_config_filepath="./configs/mq.json")
	//
	//# 2) Initialize grpc router services to work with act and check1 boxes.
	//	grpc_router = factory.grpc_router
	//act = grpc_router.get_service(ActService)
	//check = grpc_router.get_service(Check1Service)
	//
	//# 3) Initialize mq router to work with estore.
	//	estore = factory.event_batch_router

	// 4) Create root Event for report.

	current_timestamp := timestamppb.Now()
	log.Printf("current time: %v \n", current_timestamp)

	eventID := proto.EventID{Id: uuid.New().String()}
	event := proto.Event{
		Id:             &eventID,
		Name:           "Raw send example",
		Status:         proto.EventStatus_SUCCESS,
		StartTimestamp: current_timestamp,
	}

	// 5) Add this Event to EventBatch.
	eventBatch := proto.EventBatch{Events: []*proto.Event{&event}}

	// 6) Send EventBatch to estore.
	//	estore.SendAll(eventBatch)

	// 7) Create the NewOrderSingle Message.

	tradingPartyFields := map[string]*proto.Value{
		"NoPartyIDs": {Kind: &proto.Value_ListValue{ListValue: &proto.ListValue{
			Values: []*proto.Value{{Kind: &proto.Value_MessageValue{MessageValue: &proto.Message{
				Metadata: &proto.MessageMetadata{MessageType: "TradingParty_NoPartyIDs"},
				Fields: map[string]*proto.Value{
					"PartyID":       {Kind: &proto.Value_SimpleValue{SimpleValue: "Trader1"}},
					"PartyIDSource": {Kind: &proto.Value_SimpleValue{SimpleValue: "D"}},
					"PartyRole":     {Kind: &proto.Value_SimpleValue{SimpleValue: "76"}},
				},
			}}},
				{Kind: &proto.Value_MessageValue{MessageValue: &proto.Message{
					Metadata: &proto.MessageMetadata{MessageType: "TradingParty_NoPartyIDs"},
					Fields: map[string]*proto.Value{
						"PartyID":       {Kind: &proto.Value_SimpleValue{SimpleValue: "0"}},
						"PartyIDSource": {Kind: &proto.Value_SimpleValue{SimpleValue: "D"}},
						"PartyRole":     {Kind: &proto.Value_SimpleValue{SimpleValue: "3"}},
					},
				},
				}},
			}}},
		},
	}
	fields := map[string]*proto.Value{
		"Side":             {Kind: &proto.Value_SimpleValue{SimpleValue: "1"}},
		"SecurityID":       {Kind: &proto.Value_SimpleValue{SimpleValue: "INSTR1"}},
		"SecurityIDSource": {Kind: &proto.Value_SimpleValue{SimpleValue: "8"}},
		"OrdType":          {Kind: &proto.Value_SimpleValue{SimpleValue: "2"}},
		"AccountType":      {Kind: &proto.Value_SimpleValue{SimpleValue: "1"}},
		"OrderCapacity":    {Kind: &proto.Value_SimpleValue{SimpleValue: "A"}},
		"OrderQty":         {Kind: &proto.Value_SimpleValue{SimpleValue: "100"}},
		"Price":            {Kind: &proto.Value_SimpleValue{SimpleValue: "10"}},
		"ClOrdID":          {Kind: &proto.Value_SimpleValue{SimpleValue: "123"}}, //random in py
		"SecondaryClOrdID": {Kind: &proto.Value_SimpleValue{SimpleValue: "2"}},   //random in py
		"TransactTime":     {Kind: &proto.Value_SimpleValue{SimpleValue: time.Now().Format(time.RFC3339)}},
		"TradingParty":     {Kind: &proto.Value_MessageValue{MessageValue: &proto.Message{Fields: tradingPartyFields}}},
	}
	message := proto.Message{
		Metadata: &proto.MessageMetadata{
			MessageType: "NewOrderSingle",
			Id:          &proto.MessageID{ConnectionId: &proto.ConnectionID{SessionAlias: "demo-conn1"}},
		},
		Fields: fields,
	}

	// 8) Create instance of PlaceMessageRequest class - grpc message object which used for calls to the act.
	request := proto.PlaceMessageRequest{
		Message:       &message,
		ParentEventId: &eventID,
		ConnectionId:  &proto.ConnectionID{SessionAlias: "demo-conn1"},
		Description:   "User places an order.",
	}
	log.Println(request)

	// 9) Call method placeOrderFix from the act interface.
	//response = proto.placeOrderFIX(request)

	//// 10) Create MessageFilter object - mask or pattern of message verification.
	//	message_filter := proto.MessageFilter{
	//	MessageType: "ExecutionReport",
	//	Fields: map[string]*proto.ValueFilter{"ClOrdID": proto.ValueFilter{proto., key=True},
	//	'Side': ValueFilter(simple_filter='1'),
	//	'Price': ValueFilter(operation=FilterOperation.NOT_EMPTY),
	//	'LeavesQty': ValueFilter(simple_filter='0', operation=FilterOperation.NOT_EQUAL),
	//	'OrderID': ValueFilter(simple_filter=response.response_message.fields['OrderID'].simple_value)}}

	// 11) Create instance of CheckRuleRequest class - grpc message object which used for calls to the check1.
	//	check1_request := proto.CheckRuleRequest{
	//		ConnectivityId: &proto.ConnectionID{SessionAlias: "demo-conn1"},
	//		//filter = message_filter,
	//		//Checkpoint: response.checkpointId,
	//		Timeout: 3000,
	//		ParentEventId:  &eventID,
	//		Description: "User receives the ExecutionReport message."
	//	}

	// 12) Call method submitCheckRule from the check1 interface.
	//	check1_response = check.submitCheckRule(check1_request)

	//# 13) Close CommonFactory.
	//	factory.close()

}
