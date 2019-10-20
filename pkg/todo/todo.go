package todo

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

// Todo represents a todo item.
type Todo struct {
	gorm.Model
	Title       string
	Description string
	Reminder    time.Time
}

// TableName sets Todo table name to `todos`.
func (Todo) TableName() string {
	return "todos"
}

// Service provides an interface to operate on Todo items.
type Service interface {
	Create(ctx context.Context, t Todo) (Todo, error)
	Delete(ctx context.Context, id uint) (uint, error)
	Read(ctx context.Context, id uint) (Todo, error)
	ReadAll(ctx context.Context) (chan Todo, error)
	Update(ctx context.Context, todoID uint, t Todo) (Todo, error)
}
