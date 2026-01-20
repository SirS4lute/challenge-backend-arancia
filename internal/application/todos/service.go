package todos

import (
	"context"
	"errors"

	"challenge-backend-arancia/internal/domain"
	"challenge-backend-arancia/internal/ports"
)

type Service struct {
	repo  ports.TodoRepository
	idGen ports.IDGenerator
}

func NewService(repo ports.TodoRepository, idGen ports.IDGenerator) (*Service, error) {
	if repo == nil {
		return nil, errors.New("nil repo")
	}
	if idGen == nil {
		return nil, errors.New("nil id generator")
	}
	return &Service{repo: repo, idGen: idGen}, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Todo, error) {
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, title string) (domain.Todo, error) {
	td := domain.Todo{
		ID:        s.idGen.NewID(),
		Title:     title,
		Completed: false,
	}
	if err := td.Validate(); err != nil {
		return domain.Todo{}, err
	}
	if td.ID == "" {
		return domain.Todo{}, errors.New("generated empty id")
	}
	if err := s.repo.Create(ctx, td); err != nil {
		return domain.Todo{}, err
	}
	return td, nil
}

func (s *Service) Update(ctx context.Context, id string, title string, completed bool) (domain.Todo, error) {
	td := domain.Todo{
		ID:        id,
		Title:     title,
		Completed: completed,
	}
	if td.ID == "" {
		return domain.Todo{}, errors.New("missing id")
	}
	if err := td.Validate(); err != nil {
		return domain.Todo{}, err
	}
	if err := s.repo.Update(ctx, td); err != nil {
		return domain.Todo{}, err
	}
	return td, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("missing id")
	}
	return s.repo.Delete(ctx, id)
}
