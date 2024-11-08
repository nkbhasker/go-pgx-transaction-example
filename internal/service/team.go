package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/model"
)

type TeamService interface {
	Create(ctx context.Context, team *model.Team) (*model.Team, error)
	AddTeamMember(ctx context.Context, teamMember *model.TeamMember) error
}

type teamService struct {
	*baseService
}

func NewTeamService(baseService *baseService) TeamService {
	return &teamService{
		baseService,
	}
}

func (s *teamService) Create(ctx context.Context, team *model.Team) (*model.Team, error) {
	if team.ID == uuid.Nil {
		team.ID = uuid.New()
	}
	err := s.repo.Team.Create(ctx, team)

	return team, err
}

func (s *teamService) AddTeamMember(ctx context.Context, teamMember *model.TeamMember) error {
	if teamMember.ID == uuid.Nil {
		teamMember.ID = uuid.New()
	}
	return s.repo.Team.AddTeamMember(ctx, teamMember)
}
