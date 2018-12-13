package main

import (
	"context"
	"log"
	"net"

	pb "github.com/aaronflower/dzone-shipping/service.consignment/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// IRepository defines Repository interface.
type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

// Repository - Dummy repository, this simulates the use oa a datastroe
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	consignements []*pb.Consignment
}

// Create creates a consignment
func (repo *Repository) Create(consignement *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignements, consignement)
	repo.consignements = updated
	return consignement, nil
}

// service should implement all the methods to satisfy the service we defined in our
// protobuf definition. You can check the interface in the generated code itself for
// the exact method signatures etc to give you a better idea.
type service struct {
	repo IRepository
}

// CreateConsignment - we created just one method on our service, which is a create method,
// which takes a context and a request as an argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Return matching the `Response` message we created in our protobuf definition.
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}

	// Set-up our gRPC server.
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our implementation into
	// the auto-generated interface code for our portobuf definition.
	pb.RegisterShippingServiceServer(s, &service{repo})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
