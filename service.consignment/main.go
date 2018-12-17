package main

import (
	"context"
	"log"

	pb "github.com/aaronflower/dzone-shipping/service.consignment/proto/consignment"
	vesselProto "github.com/aaronflower/dzone-shipping/service.vessel/proto/vessel"
	micro "github.com/micro/go-micro"
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
	repo         IRepository
	vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - we created just one method on our service, which is a create method,
// which takes a context and a request as an argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})

	log.Printf("Found vessel: %s", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	req.VesselId = vesselResponse.Vessel.Id

	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our protobuf definition.
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	// consignments := s.repo.GetAll()
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition.
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	// Init will parse the command line flags.
	srv.Init()

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// Register our service with the gRPC server, this will tie our implementation into
	// the auto-generated interface code for our portobuf definition.
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
