package boltdb

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"challenge-backend-arancia/internal/domain"
	"challenge-backend-arancia/internal/ports"
)

func TestTodoRepository_CRUD(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo, err := NewTodoRepository(db)
	if err != nil {
		t.Fatalf("new repo: %v", err)
	}

	ctx := context.Background()
	td := domain.Todo{ID: "1", Title: "buy milk", Completed: false}

	// create
	if err := repo.Create(ctx, td); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := repo.Create(ctx, td); !errors.Is(err, ports.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}

	// get
	got, err := repo.Get(ctx, "1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != td.ID || got.Title != td.Title || got.Completed != td.Completed {
		t.Fatalf("unexpected todo: %+v", got)
	}
	if _, err := repo.Get(ctx, "nope"); !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// list
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(list))
	}

	// update
	td.Completed = true
	if err := repo.Update(ctx, td); err != nil {
		t.Fatalf("update: %v", err)
	}
	updated, err := repo.Get(ctx, "1")
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if updated.Completed != true {
		t.Fatalf("expected completed=true, got %+v", updated)
	}
	if err := repo.Update(ctx, domain.Todo{ID: "missing", Title: "x"}); !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// delete
	if err := repo.Delete(ctx, "1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := repo.Delete(ctx, "1"); !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

