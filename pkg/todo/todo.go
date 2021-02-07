package todo

import (
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
