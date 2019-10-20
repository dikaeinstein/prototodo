package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dikaeinstein/prototodo/pkg/config"
	pb "github.com/dikaeinstein/prototodo/pkg/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var defaultTimeout = 5 * time.Second

func createTodo(client pb.TodoServiceClient, t pb.Todo) pb.Todo {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.Create(ctx, &pb.CreateRequest{Todo: &t})
	if err != nil {
		log.Fatalf("%v.Create(_) = _, %v: ", client, err)
	}

	log.Println("Create result: ", resp.GetTodo())
	return *resp.Todo
}

func readTodo(client pb.TodoServiceClient, id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.Read(ctx, &pb.ReadRequest{Id: id})
	if err != nil {
		log.Fatalf("%v.Read(_) = _, %v: ", client, err)
	}

	log.Println("Read result: ", resp.GetTodo())
}

func readAllTodos(client pb.TodoServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.ReadAll(ctx, &pb.ReadAllRequest{})
	if err != nil {
		log.Fatalf("%v.ReadAll(_) = _, %v: ", client, err)
	}

	log.Println("ReadAll result: ", resp.GetTodos())
}

func updateTodo(client pb.TodoServiceClient, payload *pb.Todo) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.Update(ctx, &pb.UpdateRequest{Todo: payload})
	if err != nil {
		log.Fatalf("%v.Update(_) = _, %v: ", client, err)
	}

	log.Println("Update result: ", resp.GetUpdated())
}

func deleteTodo(client pb.TodoServiceClient, todoID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.Delete(ctx, &pb.DeleteRequest{Id: todoID})
	if err != nil {
		log.Fatalf("%v.Delete(_) = _, %v: ", client, err)
	}

	log.Println("Delete result: ", resp.GetDeleted())
}

func main() {
	cfg := config.New()
	creds, err := credentials.NewClientTLSFromFile(cfg.RootCert, "")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	var conn *grpc.ClientConn
	if cfg.TLS {
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	} else {
		conn, err = grpc.Dial(addr, grpc.WithInsecure())
	}
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	client := pb.NewTodoServiceClient(conn)

	reminder := time.Now().Add(5 * time.Second).In(time.UTC)
	reminderProto, _ := ptypes.TimestampProto(reminder)
	t := pb.Todo{
		Title:       "My first grpc todo item",
		Description: "Another first here.",
		Reminder:    reminderProto,
	}

	checkHealth(healthClient)
	newTodo := createTodo(client, t)
	readTodo(client, newTodo.Id)
	readAllTodos(client)
	payload := &pb.Todo{Id: newTodo.Id, Title: "My updated grpc todo item"}
	updateTodo(client, payload)
	deleteTodo(client, newTodo.Id)
	checkHealth(healthClient)
}

func checkHealth(c grpc_health_v1.HealthClient) {
	resp, err := c.Check(context.TODO(),
		&grpc_health_v1.HealthCheckRequest{Service: "TodoService"})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.Unimplemented {
			log.Println("the server doesn't implement the grpc health protocol")
		} else {
			log.Printf("gRPC health check failed: %v", err)
		}
	}
	log.Printf("Health status: %s", resp.GetStatus().String())
}
