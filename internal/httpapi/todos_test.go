package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"challenge-backend-arancia/internal/application/todos"
	"challenge-backend-arancia/internal/storage/boltdb"
)

func TestTodos_CRUD(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := boltdb.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo, err := boltdb.NewTodoRepository(db)
	if err != nil {
		t.Fatalf("new repo: %v", err)
	}
	svc, err := todos.NewService(repo, todos.UUIDGenerator{})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	srv := NewRouter(RouterOptions{TodoService: svc})

	// create
	createBody, _ := json.Marshal(map[string]any{"title": "buy milk"})
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(createBody))
	req = req.WithContext(context.Background())
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	// list -> should have 1
	req = httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var list []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(list))
	}
	id, _ := list[0]["id"].(string)
	if id == "" {
		t.Fatalf("expected non-empty id, got %#v", list[0]["id"])
	}

	// update unknown -> 404
	updateBody, _ := json.Marshal(map[string]any{"title": "x", "completed": true})
	req = httptest.NewRequest(http.MethodPut, "/todos/missing", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}

	// update existing -> 200
	req = httptest.NewRequest(http.MethodPut, "/todos/"+id, bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	// delete -> 204
	req = httptest.NewRequest(http.MethodDelete, "/todos/"+id, nil)
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNoContent, rec.Code, rec.Body.String())
	}
}
