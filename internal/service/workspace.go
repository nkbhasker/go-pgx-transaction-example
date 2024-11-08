package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/model"
)

type WorkspaceService interface {
	Create(ctx context.Context, workspace *model.Workspace, owner *model.User) (*model.Workspace, error)
}

type workspaceService struct {
	*baseService
	teamService TeamService
	userService UserService
}

func NewWorkspaceService(
	baseService *baseService,
	teamService TeamService,
	userService UserService,
) WorkspaceService {
	return &workspaceService{
		baseService,
		teamService,
		userService,
	}
}

func (s *workspaceService) Create(ctx context.Context, workspace *model.Workspace, owner *model.User) (*model.Workspace, error) {
	owner.ID = uuid.New()
	workspace.ID = uuid.New()
	// Set the owner of the workspace
	workspace.Owner = owner.ID
	team := &model.Team{
		ID:          uuid.New(),
		Name:        workspace.Name + "'s Team",
		WorkspaceID: workspace.ID,
	}
	err := s.Transact(ctx, func(ctx context.Context, service *Service) error {
		err := service.User.create(ctx, owner)
		if err != nil {
			return err
		}
		_, err = service.Team.Create(ctx, team)
		if err != nil {
			return err
		}
		err = service.Team.AddTeamMember(ctx, &model.TeamMember{
			ID:          uuid.New(),
			UserID:      owner.ID,
			TeamID:      team.ID,
			WorkspaceID: workspace.ID,
		})
		if err != nil {
			return err
		}

		return service.repo.Workspace.Create(ctx, workspace)
	})

	return workspace, err
}
