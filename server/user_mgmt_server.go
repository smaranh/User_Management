package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "example.com/go-usermgmt-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":5051"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Unable to listen on %s. Error: %v", port, err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening on %v", lis.Addr())
	return s.Serve(lis)
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	readBytes, err := ioutil.ReadFile("users.json")

	var user_list *pb.UserList = &pb.UserList{}
	var user_id int32 = int32(rand.Intn(1000))
	created_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   user_id,
	}

	if err != nil {
		if os.IsNotExist(err) {
			log.Println("File not found. Creating a new file")
			user_list.Users = append(user_list.Users, created_user)
			writeJsonFile("users.json", user_list)
			return created_user, nil
		}
		log.Fatalf("Failed to read file. error: %v", err)
	}

	if err := protojson.Unmarshal(readBytes, user_list); err != nil {
		log.Fatalf("Failed to parse json file. error: %v", err)
	}
	user_list.Users = append(user_list.Users, created_user)
	writeJsonFile("users.json", user_list)

	return created_user, nil
}

func writeJsonFile(filename string, user_list *pb.UserList) {
	jsonBytes, err := protojson.Marshal(user_list)
	if err != nil {
		log.Fatalf("JSON marshalling failed. error: %v", err)
	}
	if err := ioutil.WriteFile(filename, jsonBytes, 0664); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	jsonBytes, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatalf("Failed to read from file. error: %v", err)
	}
	var user_list *pb.UserList = &pb.UserList{}
	if err := protojson.Unmarshal(jsonBytes, user_list); err != nil {
		log.Fatalf("Unmarshalling failed. error: %v", err)
	}

	return user_list, nil
}

func main() {
	user_mgmt_server := NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
