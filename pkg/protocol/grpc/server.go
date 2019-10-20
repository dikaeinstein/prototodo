package grpc

import (
	"context"
	"fmt"

	pb "github.com/dikaeinstein/prototodo/pkg/proto"
	"github.com/dikaeinstein/prototodo/pkg/todo"
	"github.com/dikaeinstein/prototodo/pkg/todo/storage"
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errClientCancelled = status.Error(codes.Canceled, "Client cancelled, abandoning.")
)

type todoHandler struct {
	service todo.Service
}

// NewGRPCTodoHandler creates a new todoHandler
// which implements the pb.TodoServiceServer interface
func NewGRPCTodoHandler(s todo.Service) pb.TodoServiceServer {
	return &todoHandler{s}
}

func (h *todoHandler) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	t, err := makeTodo(req.Todo)
	if err != nil {
		return nil, err
	}

	// Check that there's still a client waiting for the response.
	if ctx.Err() == context.Canceled {
		return nil, errClientCancelled
	}

	newTodo, err := h.service.Create(ctx, *t)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to create todo: %v", err)
	}

	tProto, err := makeTodoProto(newTodo)
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Todo: tProto}, nil
}

func (h *todoHandler) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	// Check that there's still a client waiting for the response.
	if ctx.Err() == context.Canceled {
		return nil, errClientCancelled
	}

	id, err := h.service.Delete(ctx, uint(req.Id))
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "%v", storage.ErrNotFound)
		}
		return nil, status.Errorf(codes.Internal,
			"Failed to delete todo item: %v", err)
	}

	return &pb.DeleteResponse{Deleted: int64(id)}, nil
}

func (h *todoHandler) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	// Check that there's still a client waiting for the response.
	if ctx.Err() == context.Canceled {
		return nil, errClientCancelled
	}

	t, err := h.service.Read(ctx, uint(req.Id))
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "%v", storage.ErrNotFound)
		}
		return nil, status.Errorf(codes.Internal,
			"Failed to fetch todo item: %v", err)
	}

	tProto, err := makeTodoProto(t)
	if err != nil {
		return nil, err
	}

	return &pb.ReadResponse{Todo: tProto}, nil
}

func (h *todoHandler) ReadAll(ctx context.Context, req *pb.ReadAllRequest) (*pb.ReadAllResponse, error) {
	// Check that there's still a client waiting for the response.
	if ctx.Err() == context.Canceled {
		return nil, errClientCancelled
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := h.service.ReadAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"Failed to fetch todo items: %v", err)
	}

	ttProto := make([]*pb.Todo, 0)

	for t := range ch {
		tProto, err := makeTodoProto(t)
		if err != nil {
			return nil, err
		}
		ttProto = append(ttProto, tProto)
	}

	return &pb.ReadAllResponse{Todos: ttProto}, err
}

func (h *todoHandler) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	t, err := makeTodo(req.Todo)
	if err != nil {
		return nil, err
	}

	// Check that there's still a client waiting for the response.
	if ctx.Err() == context.Canceled {
		return nil, errClientCancelled
	}

	updated, err := h.service.Update(ctx, uint(req.Todo.Id), *t)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "%v", storage.ErrNotFound)
		}
		return nil, status.Errorf(codes.Internal,
			"Failed to update todo item: %v", err)
	}

	tProto, err := makeTodoProto(updated)
	return &pb.UpdateResponse{Updated: tProto}, nil
}

func makeParseTimeStampErrorMsg(field string, err error) string {
	return fmt.Sprintf("failed to convert %s to a google.protobuf.Timestamp proto."+
		"Resulting Timestamp is invalid: %v", field, err)
}

func makeTodoProto(t todo.Todo) (*pb.Todo, error) {
	var deletedAtProto *tspb.Timestamp
	if t.DeletedAt != nil {
		var err error
		deletedAtProto, err = ptypes.TimestampProto(*t.DeletedAt)
		if err != nil {
			return nil, status.Error(codes.Internal,
				makeParseTimeStampErrorMsg("DeletedAt", err))
		}
	}

	reminderProto, err := ptypes.TimestampProto(t.Reminder)
	if err != nil {
		return nil, status.Error(codes.Internal,
			makeParseTimeStampErrorMsg("Reminder", err))
	}
	createdAtProto, err := ptypes.TimestampProto(t.CreatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal,
			makeParseTimeStampErrorMsg("CreatedAt", err))
	}
	updatedAtProto, err := ptypes.TimestampProto(t.UpdatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal,
			makeParseTimeStampErrorMsg("UpdatedAt", err))
	}

	return &pb.Todo{
		Id:          int64(t.ID),
		Description: t.Description,
		Title:       t.Title,
		Reminder:    reminderProto,
		CreatedAt:   createdAtProto,
		UpdatedAt:   updatedAtProto,
		DeletedAt:   deletedAtProto,
	}, nil
}

func makeTodo(tProto *pb.Todo) (*todo.Todo, error) {
	var t todo.Todo
	r := tProto.GetReminder()
	if r != nil {
		reminder, err := ptypes.Timestamp(r)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"Request field todo.Reminder is invalid: %v", err)
		}
		t.Reminder = reminder
	}

	t.Description = tProto.GetDescription()
	t.Title = tProto.GetTitle()

	return &t, nil
}
