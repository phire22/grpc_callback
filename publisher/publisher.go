package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	chat "callback/chat"
	hook "callback/hook"
	register "callback/register"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type Server struct {
	chat.UnimplementedChatServiceServer
	register.UnimplementedRegisterServiceServer
	hookAddr map[string]int
}

func (s *Server) SayHello(ctx context.Context, message *chat.Message) (*chat.Message, error) {
	log.Printf("Received message from client: %s", message.Body)
	for addr := range s.hookAddr {
		log.Printf("Notify subscriber %s", addr)
		callbackHook(addr, message.Body)
	}
	return &chat.Message{Body: "Hello from Server A!"}, nil
}

func (s *Server) Register(ctx context.Context, req *register.RegisterReq) (*register.RegisterResp, error) {
	verifyHook(&req.Body)
	log.Printf("Registred new subscriber %s.\n", req.Body)
	s.hookAddr[req.Body] = len(s.hookAddr)
	return &register.RegisterResp{}, nil
}

func (s *Server) List(ctx context.Context, req *register.ListReq) (*register.ListResp, error) {
	log.Printf("# : Hook")
	for hook, index := range s.hookAddr {
		log.Printf("%d : %s", index, hook)
	}
	log.Printf("")
	return &register.ListResp{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := Server{hookAddr: make(map[string]int)}

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)
	register.RegisterRegisterServiceServer(grpcServer, &s)
	log.Printf("server listening at %v\n", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc Server: %v", err)
	}
}

func verifyHook(addr *string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	c := hook.NewHookServiceClient(conn)
	_, err = c.Verify(context.Background(), &hook.VerifyMsg{})
	if err != nil {
		log.Fatalf("Error when calling Verify: %s", err)
	}
}

func callbackHook(addr string, msg string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	c := hook.NewHookServiceClient(conn)
	req := hook.CallbackReq{
		Body: msg,
	}
	_, err = c.Callback(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error when calling Callback: %s", err)
	}
}
