package model

import (
	"context"

	"github.com/google/uuid"
)

type TeamRepository interface {
	Create(ctx context.Context, team *Team) error
	AddTeamMember(ctx context.Context, teamMember *TeamMember) error
}

type Team struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
}

type TeamMember struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	TeamID      uuid.UUID `json:"team_id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
}

type teamRepository struct {
	*baseRepository
}

func NewTeamRepository(baseRepository *baseRepository) TeamRepository {
	return &teamRepository{
		baseRepository: baseRepository,
	}
}

func (r *teamRepository) Create(ctx context.Context, team *Team) error {
	stmt := `INSERT INTO teams (
		id, name, workspace_id
	) VALUES ($1, $2, $3)`
	_, err := r.db.Connection().Exec(ctx, stmt, team.ID, team.Name, team.WorkspaceID)
	if err != nil {
		return err
	}

	return nil
}

func (r *teamRepository) AddTeamMember(ctx context.Context, teamMember *TeamMember) error {
	stmt := `INSERT INTO team_members (
		id, user_id, team_id, workspace_id
	) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Connection().Exec(ctx, stmt, teamMember.ID, teamMember.UserID, teamMember.TeamID, teamMember.WorkspaceID)
	if err != nil {
		return err
	}

	return nil
}
