package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "github.com/thealanphipps-del/pqr/proto"
)

func main() {
	fmt.Println("Dialing pqr.info:443 via gRPC...")
	conn, err := grpc.Dial("pqr.info:443", grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewSwarmCommunicationClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	fmt.Println("Calling GetActiveShortcodes...")
	r, err := c.GetActiveShortcodes(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get shortcodes: %v", err)
	}
	fmt.Printf("Response: %v\n", r)
}
