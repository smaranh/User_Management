package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "example.com/go-usermgmt-grpc/proto"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const (
	port = ":5051"
)

type UserManagementServer struct {
	conn *pgx.Conn
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
	createSql := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name text,
		age int
	);`

	_, err := s.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}

	created_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
	}
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Failed to being transaction. error: %v", err)
	}
	_, err = tx.Exec(context.Background(), "insert into users(name, age) values ($1, $2)", created_user.Name, created_user.Age)
	if err != nil {
		log.Fatalf("Failed to execute. error: %v", err)
	}
	tx.Commit(context.Background())
	return created_user, nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var user_list *pb.UserList = &pb.UserList{}

	rows, err := s.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		user_list.Users = append(user_list.Users, &user)
	}

	return user_list, nil
}

func main() {
	database_url := "postgres://smaranh:password@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())

	user_mgmt_server := NewUserManagementServer()
	user_mgmt_server.conn = conn
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
