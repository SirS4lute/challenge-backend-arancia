package boltdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"challenge-backend-arancia/internal/domain"
	"challenge-backend-arancia/internal/ports"

	bolt "go.etcd.io/bbolt"
)

var todosBucket = []byte("todos")

type TodoRepository struct {
	db *bolt.DB
}

func NewTodoRepository(db *bolt.DB) (*TodoRepository, error) {
	if db == nil {
		return nil, errors.New("nil db")
	}
	r := &TodoRepository{db: db}
	if err := r.ensureBuckets(context.Background()); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *TodoRepository) ensureBuckets(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(todosBucket)
		return err
	})
}

func (r *TodoRepository) List(ctx context.Context) ([]domain.Todo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var out []domain.Todo
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(todosBucket)
		if b == nil {
			return fmt.Errorf("bucket %q not found", string(todosBucket))
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := ctx.Err(); err != nil {
				return err
			}
			var td domain.Todo
			if err := json.Unmarshal(v, &td); err != nil {
				return err
			}
			out = append(out, td)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *TodoRepository) Get(ctx context.Context, id string) (domain.Todo, error) {
	if err := ctx.Err(); err != nil {
		return domain.Todo{}, err
	}

	var out domain.Todo
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(todosBucket)
		if b == nil {
			return fmt.Errorf("bucket %q not found", string(todosBucket))
		}
		v := b.Get([]byte(id))
		if v == nil {
			return ports.ErrNotFound
		}
		return json.Unmarshal(v, &out)
	})
	if err != nil {
		return domain.Todo{}, err
	}
	return out, nil
}

func (r *TodoRepository) Create(ctx context.Context, todo domain.Todo) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if todo.ID == "" {
		return errors.New("missing id")
	}
	if err := todo.Validate(); err != nil {
		return err
	}

	payload, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(todosBucket)
		if b == nil {
			return fmt.Errorf("bucket %q not found", string(todosBucket))
		}
		k := []byte(todo.ID)
		if existing := b.Get(k); existing != nil {
			return ports.ErrConflict
		}
		return b.Put(k, payload)
	})
}

func (r *TodoRepository) Update(ctx context.Context, todo domain.Todo) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if todo.ID == "" {
		return errors.New("missing id")
	}
	if err := todo.Validate(); err != nil {
		return err
	}

	payload, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(todosBucket)
		if b == nil {
			return fmt.Errorf("bucket %q not found", string(todosBucket))
		}
		k := []byte(todo.ID)
		if existing := b.Get(k); existing == nil {
			return ports.ErrNotFound
		}
		return b.Put(k, payload)
	})
}

func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if id == "" {
		return errors.New("missing id")
	}

	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(todosBucket)
		if b == nil {
			return fmt.Errorf("bucket %q not found", string(todosBucket))
		}
		k := []byte(id)
		if existing := b.Get(k); existing == nil {
			return ports.ErrNotFound
		}
		return b.Delete(k)
	})
}

