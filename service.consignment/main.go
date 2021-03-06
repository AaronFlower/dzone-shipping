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
	GetAll() []*pb.Consignment
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

// GetAll returns all consignments.
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignements
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

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	// consignments := s.repo.GetAll()
	consignments := s.repo.GetAll()
	return &pb.Response{Consignments: consignments}, nil
}

func main() {
	repo := &Repository{}

	// 创建一个监听接口。golang 是用核心库 net, net/http 来作网络通信的。所以端口都是这两个库来创建。
	// gRPC 是用 golang 实现的一个 RPC 库而已，没有监听端口的功能，只能为该端口提供服务。
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

	// gPRC 创建的服务指定在那个端口进行服务.
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
