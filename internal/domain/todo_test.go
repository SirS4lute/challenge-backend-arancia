package domain

import "testing"

func TestTodoValidate_TitleRequired(t *testing.T) {
	t.Parallel()

	td := Todo{ID: "x", Title: "   ", Completed: false}
	if err := td.Validate(); err != ErrInvalidTitle {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestTodoValidate_TitleMaxLen(t *testing.T) {
	t.Parallel()

	long := make([]byte, MaxTitleLen+1)
	for i := range long {
		long[i] = 'a'
	}

	td := Todo{ID: "x", Title: string(long)}
	if err := td.Validate(); err != ErrInvalidTitle {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestTodoValidate_OK(t *testing.T) {
	t.Parallel()

	td := Todo{ID: "x", Title: "buy milk", Completed: true}
	if err := td.Validate(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
