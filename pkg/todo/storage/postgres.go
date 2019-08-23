package storage

import (
	"context"
	"errors"
	"log"

	"github.com/dikaeinstein/prototodo/pkg/todo"
	"github.com/jinzhu/gorm"
)

// ErrNotFound represents error when a todo item is not found in the postgres data store.
var ErrNotFound = errors.New("Todo item not found")

// PostgresStore represents the postgres db
type PostgresStore struct {
	*gorm.DB
}

// NewPostgresStore creates an instance of the PostgresStore with the db connection.
func NewPostgresStore(db *gorm.DB) *PostgresStore {
	return &PostgresStore{db}
}

// GetAll fetches all todo items from postgres data store.
func (p *PostgresStore) GetAll(ctx context.Context) (chan todo.ToDo, error) {
	rows, err := p.DB.Model(&todo.ToDo{}).Select("*").Rows()
	if err != nil {
		return nil, err
	}

	c := make(chan todo.ToDo)

	go func() {
		defer rows.Close()
		defer close(c)
		for rows.Next() {
			var t todo.ToDo
			p.DB.ScanRows(rows, &t)
			select {
			case <-ctx.Done():
				log.Println(ctx.Err())
				break
			default:
				c <- t
			}
		}
	}()

	return c, rows.Err()
}

// GetByID fetches one todo item from the postgres data store using its id.
func (p *PostgresStore) GetByID(ctx context.Context, id uint) (todo.ToDo, error) {
	var t todo.ToDo
	if err := p.DB.First(&t, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return t, ErrNotFound
		}
		return t, err
	}

	return t, nil
}

// Create saves the todo into the postgres data store.
func (p *PostgresStore) Create(ctx context.Context, t todo.ToDo) (todo.ToDo, error) {
	if err := p.DB.Create(&t).Error; err != nil {
		return t, err
	}

	return t, nil
}

// Delete removes a todo item from the postgres data store.
func (p *PostgresStore) Delete(ctx context.Context, id uint) (uint, error) {
	var t todo.ToDo
	if err := p.DB.First(&t, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return id, ErrNotFound
		}
		return id, err
	}

	if err := p.DB.Delete(&t).Error; err != nil {
		return id, err
	}

	return id, nil
}

// Update updates a todo item with attrs in the postgres data store
func (p *PostgresStore) Update(ctx context.Context, todoID uint, attrs todo.ToDo) (todo.ToDo, error) {
	var t todo.ToDo
	if err := p.DB.First(&t, todoID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return t, ErrNotFound
		}
		return t, err
	}

	if err := p.DB.Model(&t).Updates(attrs).Error; err != nil {
		return t, err
	}

	return t, nil
}
