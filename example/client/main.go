package main

import (
	"context"
	"github.com/2se/dolphin-sdk/mock"
	"github.com/golang/protobuf/ptypes"

	"github.com/2se/dolphin-sdk/pb"

	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	//resource: MockUser
	//action: GetUser version:v2
	//input param:GetUserRequest
	//output param:User

	ctx1, _ := context.WithTimeout(context.Background(), time.Second*3)
	//defer cel()
	//conn, err := grpc.DialContext(ctx1, address, grpc.WithBlock(), grpc.WithInsecure())
	conn, err := grpc.DialContext(ctx1, "127.0.0.1:8848", grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Println("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewAppServeClient(conn)
	p := &mock.GetUserRequest{
		UserId: 1,
	}

	object, err := ptypes.MarshalAny(p)
	if err != nil {
		log.Println(err)
		return
	}

	req := &pb.ClientComRequest{
		Meta: &pb.ClientComMeta{
			Resource: "MockUser",
			Revision: "v2",
			Action:   "GetUser",
		},
		Params: object,
	}
	res, err := c.Request(context.Background(), req)
	if err != nil {
		log.Println(err)
		return
	}
	if res.Code == 200 {
		pmu := &mock.User{}
		err = ptypes.UnmarshalAny(res.Body, pmu)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(pmu)
	}

}