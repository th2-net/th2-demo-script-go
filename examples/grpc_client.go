package main

import (
	"context"
	ac "exactpro/th2/example/proto"
	"github.com/th2-net/th2-common-go/schema/factory"
	"github.com/th2-net/th2-demo-script-go.git/th2-common-go/schema/modules/grpcModule"
	//ac "exactpro/th2/example/proto" //act proto
	"fmt"
	"github.com/google/uuid"
	cg "github.com/th2-net/th2-common-go/proto" //common proto
	//"github.com/th2-net/th2-common-go/schema/factory"
	//"github.com/th2-net/th2-common-go/schema/modules/grpcModule"
	"log"
	"os"
	"reflect"
	"time"
)

func main() {
	newFactory := factory.NewFactory(os.Args)
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
	con, conErr := grpcRouter.GetConnection("actAttr")
	if conErr != nil {
		log.Fatalf(conErr.Error())
	}
	//defer con.Close() //doesnt have that func
	c := ac.NewActClient(con)

	// getting data ready for placing order

	eventID := cg.EventID{Id: uuid.New().String()}

	tradingPartyFields := map[string]*cg.Value{
		"NoPartyIDs": {Kind: &cg.Value_ListValue{ListValue: &cg.ListValue{
			Values: []*cg.Value{{Kind: &cg.Value_MessageValue{MessageValue: &cg.Message{
				Metadata: &cg.MessageMetadata{MessageType: "TradingParty_NoPartyIDs"},
				Fields: map[string]*cg.Value{
					"PartyID":       {Kind: &cg.Value_SimpleValue{SimpleValue: "Trader1"}},
					"PartyIDSource": {Kind: &cg.Value_SimpleValue{SimpleValue: "D"}},
					"PartyRole":     {Kind: &cg.Value_SimpleValue{SimpleValue: "76"}},
				},
			}}},
				{Kind: &cg.Value_MessageValue{MessageValue: &cg.Message{
					Metadata: &cg.MessageMetadata{MessageType: "TradingParty_NoPartyIDs"},
					Fields: map[string]*cg.Value{
						"PartyID":       {Kind: &cg.Value_SimpleValue{SimpleValue: "0"}},
						"PartyIDSource": {Kind: &cg.Value_SimpleValue{SimpleValue: "D"}},
						"PartyRole":     {Kind: &cg.Value_SimpleValue{SimpleValue: "3"}},
					},
				},
				}},
			}}},
		},
	}

	fields := map[string]*cg.Value{
		"Side":             {Kind: &cg.Value_SimpleValue{SimpleValue: "1"}},
		"SecurityID":       {Kind: &cg.Value_SimpleValue{SimpleValue: "INSTR1"}},
		"SecurityIDSource": {Kind: &cg.Value_SimpleValue{SimpleValue: "8"}},
		"OrdType":          {Kind: &cg.Value_SimpleValue{SimpleValue: "2"}},
		"AccountType":      {Kind: &cg.Value_SimpleValue{SimpleValue: "1"}},
		"OrderCapacity":    {Kind: &cg.Value_SimpleValue{SimpleValue: "A"}},
		"OrderQty":         {Kind: &cg.Value_SimpleValue{SimpleValue: "100"}},
		"Price":            {Kind: &cg.Value_SimpleValue{SimpleValue: "10"}},
		"ClOrdID":          {Kind: &cg.Value_SimpleValue{SimpleValue: "123"}}, //random in py
		"SecondaryClOrdID": {Kind: &cg.Value_SimpleValue{SimpleValue: "2"}},   //random in py
		"TransactTime":     {Kind: &cg.Value_SimpleValue{SimpleValue: time.Now().Format(time.RFC3339)}},
		"TradingParty":     {Kind: &cg.Value_MessageValue{MessageValue: &cg.Message{Fields: tradingPartyFields}}},
	}

	msg := cg.Message{Metadata: &cg.MessageMetadata{
		MessageType: "NewOrderSingle",
		Id:          &cg.MessageID{ConnectionId: &cg.ConnectionID{SessionAlias: "demo-conn1"}},
	},
		Fields: fields,
	}

	request := cg.PlaceMessageRequest{
		Message:       &msg,
		ParentEventId: &eventID,
		Description:   "User places an order.",
	}

	resp, err := c.PlaceOrderFIX(context.Background(), &request)
	if err != nil {
		log.Fatalf("could not send order: %v", err)
	}
	log.Printf("order sent. response: %s", resp.String())
}
