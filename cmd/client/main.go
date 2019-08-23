package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	pb "github.com/dikaeinstein/prototodo/pkg/proto"
	"github.com/golang/protobuf/ptypes"

	"google.golang.org/grpc"
)

var defaultTimeout = 5 * time.Second

func createToDo(client pb.ToDoServiceClient, t pb.ToDo) pb.ToDo {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.CreateToDo(ctx, &pb.CreateRequest{Todo: &t})
	if err != nil {
		log.Fatalf("%v.CreateToDo(_) = _, %v: ", client, err)
	}

	log.Println("CreateToDo result: ", resp.GetTodo())
	return *resp.Todo
}

func readToDo(client pb.ToDoServiceClient, id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.ReadToDo(ctx, &pb.ReadRequest{Id: id})
	if err != nil {
		log.Fatalf("%v.ReadToDo(_) = _, %v: ", client, err)
	}

	log.Println("ReadToDo result: ", resp.GetTodo())
}

func readAllToDos(client pb.ToDoServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.ReadAllToDos(ctx, &pb.ReadAllRequest{})
	if err != nil {
		log.Fatalf("%v.ReadAllToDos(_) = _, %v: ", client, err)
	}

	log.Println("ReadAllToDos result: ", resp.GetTodos())
}

func updateToDo(client pb.ToDoServiceClient, payload *pb.ToDo) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.UpdateToDo(ctx, &pb.UpdateRequest{Todo: payload})
	if err != nil {
		log.Fatalf("%v.UpdateToDo(_) = _, %v: ", client, err)
	}

	log.Println("UpdateToDo result: ", resp.GetUpdated())
}

func deleteToDo(client pb.ToDoServiceClient, todoID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	resp, err := client.DeleteToDo(ctx, &pb.DeleteRequest{Id: todoID})
	if err != nil {
		log.Fatalf("%v.DeleteToDo(_) = _, %v: ", client, err)
	}

	log.Println("DeleteToDo result: ", resp.GetDeleted())
}

func main() {
	creds, err := credentials.NewClientTLSFromFile("/Users/Dikaeinstein/Library/Application Support/mkcert/rootCA.pem", "")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewToDoServiceClient(conn)

	reminder := time.Now().Add(5 * time.Second).In(time.UTC)
	reminderProto, _ := ptypes.TimestampProto(reminder)
	t := pb.ToDo{
		Title:       "My first grpc todo item",
		Description: "Another first here.",
		Reminder:    reminderProto,
	}
	newToDo := createToDo(client, t)
	readToDo(client, newToDo.Id)
	readAllToDos(client)
	payload := &pb.ToDo{Id: newToDo.Id, Title: "My updated grpc todo item"}
	updateToDo(client, payload)
	deleteToDo(client, newToDo.Id)
}
