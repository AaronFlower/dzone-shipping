package main

import (
	"github.com/globalsign/mgo"

	pb "github.com/aaronflower/dzone-shipping/service.vessel/proto/vessel"
	"github.com/globalsign/mgo/bson"
)

const (
	dbName           = "shippy"
	vesselCollection = "vessel"
)

// Repository defines the interface.
type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
	Create(vessel *pb.Vessel) error
	Close()
}

// VesselRepository implements the interface.
type VesselRepository struct {
	session *mgo.Session
}

// FindAvailable checks a specification against a map of vessel.
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	var vessel *pb.Vessel

	// Here we define a more complex query than our consignment-microservice's
	// GetAll function. Here we'er asking for a vessel who's max weight and
	// capacity are greater than and equal to the given capacity and weight.
	// we're also using the `One` function has as that's all we want.
	err := repo.collection().Find(bson.M{
		"capacity":  bson.M{"$gte": spec.Capacity},
		"maxweight": bson.M{"$gte": spec.MaxWeight},
	}).One(&vessel)
	if err != nil {
		return nil, err
	}
	return vessel, nil
}

// Create creates a vessel
func (repo *VesselRepository) Create(vessel *pb.Vessel) error {
	return repo.collection().Insert(vessel)
}

// Close closes the session.
func (repo *VesselRepository) Close() {
	repo.session.Close()
}

func (repo *VesselRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(vesselCollection)
}
