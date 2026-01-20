package domain

import (
	"errors"
	"strings"
)

const (
	// MaxTitleLen is the maximum length allowed for a Todo title after trimming spaces.
	MaxTitleLen = 200
)

var (
	// ErrInvalidTitle indicates a title that is empty (after trim) or exceeds MaxTitleLen.
	ErrInvalidTitle = errors.New("invalid title")
)

// Todo is the core entity of the system.
type Todo struct {
	ID        string
	Title     string
	Completed bool
}

// Validate checks invariants for a Todo.
// At domain level we keep it minimal: only validates Title.
func (t Todo) Validate() error {
	title := strings.TrimSpace(t.Title)
	if title == "" || len(title) > MaxTitleLen {
		return ErrInvalidTitle
	}
	return nil
}
