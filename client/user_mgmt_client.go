package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example.com/go-usermgmt-grpc/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:5051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect %v", err)
	}
	defer conn.Close()

	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var new_users = make(map[string]int32)
	new_users["Janice"] = 35
	new_users["Chandler"] = 33

	for name, age := range new_users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{
			Name: name,
			Age:  age,
		})
		if err != nil {
			log.Fatalf("could not create new user: %v", err)
		}

		log.Printf(`User Details: 
		Name: %s,
		Age: %d,
		Id: %d
		`, r.GetName(), r.GetAge(), r.GetId())
	}

	params := &pb.GetUsersParams{}
	res, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("failed to retrive users: %v", err)
	}

	log.Println("\nList of Users")
	fmt.Printf("GetUsers(): %v", res.GetUsers())
}
