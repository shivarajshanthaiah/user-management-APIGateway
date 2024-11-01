package user

import (
	"log"

	"github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/config"
	pb "github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/user/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ClientDial(cfg config.Config) (pb.UserServiceClient, error) {
	// address := "user-service:" + cfg.USERPORT
	address := "localhost:" + cfg.USERPORT

	grpcConn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Error dialing to gRPC client at address %s: %s", address, err.Error())
		return nil, err
	}

	log.Printf("Successfully connected to gRPC client at address: %s", address)
	return pb.NewUserServiceClient(grpcConn), nil
}
