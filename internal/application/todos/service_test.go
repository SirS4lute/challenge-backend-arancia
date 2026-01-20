package todos

import (
	"context"
	"errors"
	"testing"

	"challenge-backend-arancia/internal/domain"
	"challenge-backend-arancia/internal/ports"
)

type fakeIDGen struct{ id string }

func (f fakeIDGen) NewID() string { return f.id }

type fakeRepo struct {
	todos   map[string]domain.Todo
	creates int
	updates int
	deletes int
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{todos: map[string]domain.Todo{}}
}

func (r *fakeRepo) List(ctx context.Context) ([]domain.Todo, error) {
	out := make([]domain.Todo, 0, len(r.todos))
	for _, v := range r.todos {
		out = append(out, v)
	}
	return out, nil
}

func (r *fakeRepo) Get(ctx context.Context, id string) (domain.Todo, error) {
	td, ok := r.todos[id]
	if !ok {
		return domain.Todo{}, ports.ErrNotFound
	}
	return td, nil
}

func (r *fakeRepo) Create(ctx context.Context, todo domain.Todo) error {
	r.creates++
	if _, ok := r.todos[todo.ID]; ok {
		return ports.ErrConflict
	}
	r.todos[todo.ID] = todo
	return nil
}

func (r *fakeRepo) Update(ctx context.Context, todo domain.Todo) error {
	r.updates++
	if _, ok := r.todos[todo.ID]; !ok {
		return ports.ErrNotFound
	}
	r.todos[todo.ID] = todo
	return nil
}

func (r *fakeRepo) Delete(ctx context.Context, id string) error {
	r.deletes++
	if _, ok := r.todos[id]; !ok {
		return ports.ErrNotFound
	}
	delete(r.todos, id)
	return nil
}

func TestService_Create_ValidatesAndPersists(t *testing.T) {
	t.Parallel()

	repo := newFakeRepo()
	svc, err := NewService(repo, fakeIDGen{id: "id-1"})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	td, err := svc.Create(context.Background(), "buy milk")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if td.ID != "id-1" || td.Completed != false || td.Title != "buy milk" {
		t.Fatalf("unexpected todo: %+v", td)
	}
	if repo.creates != 1 {
		t.Fatalf("expected 1 create call, got %d", repo.creates)
	}
}

func TestService_Create_InvalidTitle_DoesNotCallRepo(t *testing.T) {
	t.Parallel()

	repo := newFakeRepo()
	svc, err := NewService(repo, fakeIDGen{id: "id-1"})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	_, err = svc.Create(context.Background(), "   ")
	if !errors.Is(err, domain.ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
	if repo.creates != 0 {
		t.Fatalf("expected 0 create calls, got %d", repo.creates)
	}
}

func TestService_Update_NotFound(t *testing.T) {
	t.Parallel()

	repo := newFakeRepo()
	svc, err := NewService(repo, fakeIDGen{id: "id-1"})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	_, err = svc.Update(context.Background(), "missing", "x", false)
	if !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestService_Delete_NotFound(t *testing.T) {
	t.Parallel()

	repo := newFakeRepo()
	svc, err := NewService(repo, fakeIDGen{id: "id-1"})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	err = svc.Delete(context.Background(), "missing")
	if !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
