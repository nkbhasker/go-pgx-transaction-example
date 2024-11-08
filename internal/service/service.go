package service

import (
	"context"

	"github.com/nkbhasker/go-pgx-transaction-example/internal/model"
)

// Validate that Service implements Transactioner
var _ model.Transactioner[*Service] = (*Service)(nil)

type baseService struct {
	repo *model.Repository
}

type Service struct {
	*baseService
	Team      TeamService
	User      UserService
	Workspace WorkspaceService
}

func NewBaseService(repo *model.Repository) *baseService {
	return &baseService{
		repo: repo,
	}
}

func NewService(baseService *baseService) *Service {
	teamService := NewTeamService(baseService)
	userService := NewUserService(baseService, teamService)
	workspaceService := NewWorkspaceService(baseService, teamService, userService)

	return &Service{
		baseService: baseService,
		Team:        teamService,
		User:        userService,
		Workspace:   workspaceService,
	}
}

func (s *baseService) Transact(ctx context.Context, fn func(ctx context.Context, service *Service) error) error {
	return s.repo.Transact(ctx, func(ctx context.Context, repo *model.Repository) error {
		return fn(ctx, NewService(NewBaseService(repo)))
	})
}
