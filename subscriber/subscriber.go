package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	hook "callback/hook"
	register "callback/register"
)

var (
	port = flag.Int("port", 50052, "The server port")
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type Server struct {
	hook.UnimplementedHookServiceServer
}

func (s *Server) Callback(ctx context.Context, resp *hook.CallbackReq) (*hook.CallbackResp, error) {
	log.Printf("Message from publisher: %s", resp.Body)
	return &hook.CallbackResp{}, nil
}

func (s *Server) Verify(ctx context.Context, resp *hook.VerifyMsg) (*hook.VerifyMsg, error) {
	return &hook.VerifyMsg{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := Server{}

	grpcServer := grpc.NewServer()

	hook.RegisterHookServiceServer(grpcServer, &s)
	log.Printf("server listening at %v\n", lis.Addr())

	go registerPublisherHook()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc Server: %v", err)
	}
}

func registerPublisherHook() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	c := register.NewRegisterServiceClient(conn)
	req := register.RegisterReq{
		Body: fmt.Sprintf(":%d", *port),
	}
	_, err = c.Register(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error when calling Register: %s", err)
	}
	log.Printf("Succesfully subscribed to publisher.")
}
