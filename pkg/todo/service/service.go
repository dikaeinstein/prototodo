package service

import (
	"context"

	"github.com/dikaeinstein/prototodo/pkg/protocol/grpc"
	"github.com/dikaeinstein/prototodo/pkg/todo"
)

// Repository provides access to the todo data store.
type Repository interface {
	GetAll(ctx context.Context) (chan todo.Todo, error)
	GetByID(ctx context.Context, id uint) (todo.Todo, error)
	Create(ctx context.Context, t todo.Todo) (todo.Todo, error)
	Delete(ctx context.Context, id uint) (uint, error)
	Update(ctx context.Context, id uint, t todo.Todo) (todo.Todo, error)
}

// New creates a todo service with the necessary dependencies.
// This contains the core business logic to operate on todo items.
func New(r Repository) grpc.Service {
	return &service{r}
}

type service struct {
	r Repository
}

func (s service) Create(ctx context.Context, t todo.Todo) (todo.Todo, error) {
	return s.r.Create(ctx, t)
}

func (s service) Delete(ctx context.Context, id uint) (uint, error) {
	return s.r.Delete(ctx, id)
}

func (s service) Read(ctx context.Context, id uint) (todo.Todo, error) {
	return s.r.GetByID(ctx, id)
}

func (s service) ReadAll(ctx context.Context) (chan todo.Todo, error) {
	c, err := s.r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return c, err
}

func (s service) Update(ctx context.Context, todoID uint, t todo.Todo) (todo.Todo, error) {
	return s.r.Update(ctx, todoID, t)
}
