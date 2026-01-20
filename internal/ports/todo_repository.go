package ports

import (
	"context"

	"challenge-backend-arancia/internal/domain"
)

// TodoRepository defines persistence operations for Todo entities.
type TodoRepository interface {
	List(ctx context.Context) ([]domain.Todo, error)
	Get(ctx context.Context, id string) (domain.Todo, error)
	Create(ctx context.Context, todo domain.Todo) error
	Update(ctx context.Context, todo domain.Todo) error
	Delete(ctx context.Context, id string) error
}
