package main

import (
	"flag"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"callback/register"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

// type chatServiceClient struct {
// 	cc grpc.ClientConnInterface
// }

// func NewChatServiceClient(cc grpc.ClientConnInterface) chat.ChatServiceClient {
// 	return &chatServiceClient{cc}
// }

// func (c *chatServiceClient) SayHello(ctx context.Context, in *chat.Message, opts ...grpc.CallOption) (*chat.Message, error) {
// 	out := new(chat.Message)
// 	err := c.cc.Invoke(ctx, "/chat.ChatService/SayHello", in, out, opts...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return out, nil
// }


func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()
	
	c:= register.NewRegisterServiceClient(conn)
	_, err = c.List(context.Background(), &register.ListReq{})
	if err != nil {
		log.Fatalf("Error when calling List: %s", err)
	}
}