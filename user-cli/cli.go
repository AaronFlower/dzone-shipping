package main

import (
	"context"
	"log"
	"os"

	pb "github.com/aaronflower/dzone-shipping/service.user/proto/user"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
)

func main() {
	cmd.Init()

	// create new greeter client
	client := pb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)

	// Define our flags
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "name",
				Usage: "You full name",
			},

			cli.StringFlag{
				Name:  "email",
				Usage: "You email",
			},

			cli.StringFlag{
				Name:  "password",
				Usage: "You password",
			},

			cli.StringFlag{
				Name:  "company",
				Usage: "You company",
			},
		))
	service.Init(
		micro.Action(
			func(c *cli.Context) {
				name := c.String("name")
				email := c.String("email")
				password := c.String("password")
				company := c.String("company")

				r, err := client.Create(context.TODO(), &pb.User{
					Name:     name,
					Email:    email,
					Password: password,
					Company:  company,
				})
				if err != nil {
					log.Fatalf("Could not create: %v", err)
				}

				log.Printf("Created: %v", r.User.Id)

				getAll, err := client.GetAll(context.Background(), &pb.Request{})
				if err != nil {
					log.Fatalf("Could not list users: %v", err)
				}
				for _, v := range getAll.Users {
					log.Println(v)
				}
				os.Exit(0)
			}),
	)
	// Run the server
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
