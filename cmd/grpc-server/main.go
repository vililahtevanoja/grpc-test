package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"grpc-test/protobuf/github.com/vililahtevanoja/grpc-test/reservationgrpc"
	"log"
	"math/rand"
	"net"
	"strings"
)

const port = ":5000"

type reservationServer struct {
	reservationgrpc.ReservationServer
}

func (s reservationServer) TryGetReservation(ctx context.Context, in *reservationgrpc.Credentials) (*reservationgrpc.TryGetReservationResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	println(fmt.Sprintf("API-key: %s", strings.Join(md["x-api-key"], "")))
	result := rand.Int() % 2 == 0
	if result {
		return &reservationgrpc.TryGetReservationResponse{Result: result, Msg: fmt.Sprintf("reservation made for %s", in.Username)}, nil
	} else {
		return &reservationgrpc.TryGetReservationResponse{Result: result, Msg: "could not mark job"}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reservationgrpc.RegisterReservationServer(s, reservationServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
