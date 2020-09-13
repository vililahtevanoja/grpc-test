package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-test/protobuf/github.com/vililahtevanoja/grpc-test/reservationgrpc"
	"log"
	"time"
)

type ApiKeyCredentials struct {

}

func (creds *ApiKeyCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	credentialMetaData :=  map[string]string{
		"X-API-KEY": "test-api-key",
	}
	return credentialMetaData, nil
}

func (creds *ApiKeyCredentials) RequireTransportSecurity() bool {
	return false
}

func main() {

	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithPerRPCCredentials(&ApiKeyCredentials{}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := reservationgrpc.NewReservationClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < 10; i++ {
		r, err := c.TryGetReservation(ctx, &reservationgrpc.Credentials{Username: "username", Password: "password"})
		if err != nil {
			log.Fatal(err)
		}
		if r.Result {
			println(fmt.Sprintf("SUCCESS: %s", r.Msg))
		} else {
			println(fmt.Sprintf("FAILURE: %s", r.Msg))
		}
		time.Sleep(time.Second)
	}
}