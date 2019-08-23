package service

import (
	"context"

	"github.com/dikaeinstein/prototodo/pkg/todo"
)

// Repository provides access to the todo data store.
type Repository interface {
	GetAll(ctx context.Context) (chan todo.ToDo, error)
	GetByID(ctx context.Context, id uint) (todo.ToDo, error)
	Create(ctx context.Context, t todo.ToDo) (todo.ToDo, error)
	Delete(ctx context.Context, id uint) (uint, error)
	Update(ctx context.Context, id uint, t todo.ToDo) (todo.ToDo, error)
}

// New creates a todo service with the necessary dependencies.
func New(r Repository) todo.Service {
	return &service{r}
}

type service struct {
	r Repository
}

func (s service) Create(ctx context.Context, t todo.ToDo) (todo.ToDo, error) {
	return s.r.Create(ctx, t)
}

func (s service) Delete(ctx context.Context, id uint) (uint, error) {
	return s.r.Delete(ctx, id)
}

func (s service) Read(ctx context.Context, id uint) (todo.ToDo, error) {
	return s.r.GetByID(ctx, id)
}

func (s service) ReadAll(ctx context.Context) (chan todo.ToDo, error) {
	c, err := s.r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return c, err
}

func (s service) Update(ctx context.Context, todoID uint, t todo.ToDo) (todo.ToDo, error) {
	return s.r.Update(ctx, todoID, t)
}
