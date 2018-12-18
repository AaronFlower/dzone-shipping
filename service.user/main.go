package main

import (
	"fmt"
	"log"

	pb "github.com/aaronflower/dzone-shipping/service.user/proto/user"
	micro "github.com/micro/go-micro"
)

func main() {
	// create a database connection and handles
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// Automatically migrates the user struct into database columns/types etc.
	// This will check for changes and migrate them each itme this service is restrarted.
	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}
	tokenService := &TokenService{repo}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterUserServiceHandler(srv.Server(), &service{repo, tokenService})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
