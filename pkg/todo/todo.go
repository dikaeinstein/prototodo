package todo

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

// ToDo represents a todo item.
type ToDo struct {
	gorm.Model
	Title       string
	Description string
	Reminder    time.Time
}

// TableName sets ToDo table name to `todos`.
func (ToDo) TableName() string {
	return "todos"
}

// Service provides an interface to operate on ToDo items.
type Service interface {
	Create(ctx context.Context, t ToDo) (ToDo, error)
	Delete(ctx context.Context, id uint) (uint, error)
	Read(ctx context.Context, id uint) (ToDo, error)
	ReadAll(ctx context.Context) (chan ToDo, error)
	Update(ctx context.Context, todoID uint, t ToDo) (ToDo, error)
}
