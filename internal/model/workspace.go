package model

import (
	"context"

	"github.com/google/uuid"
)

type WorkspaceRepository interface {
	List(ctx context.Context) ([]*Workspace, error)
	Get(ctx context.Context, id uuid.UUID) (*Workspace, error)
	Create(ctx context.Context, workspace *Workspace) error
}

type Workspace struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Owner uuid.UUID `json:"owner"`
}

type workspaceRepository struct {
	baseRepository
}

func NewWorkspaceRepository(baseRepository *baseRepository) WorkspaceRepository {
	return &workspaceRepository{
		baseRepository: *baseRepository,
	}
}

func (r *workspaceRepository) List(ctx context.Context) ([]*Workspace, error) {
	workspaces := []*Workspace{}
	rows, err := r.db.Connection().Query(ctx, "SELECT id, name, owner FROM workspaces")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var workspace Workspace
		err := rows.Scan(&workspace.ID, &workspace.Name, &workspace.Owner)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, &workspace)
	}

	return workspaces, nil
}

func (r *workspaceRepository) Get(ctx context.Context, id uuid.UUID) (*Workspace, error) {
	var workspace Workspace
	stmt := `SELECT id, name, owner FROM workspaces WHERE id = $1`
	err := r.db.Connection().QueryRow(ctx, stmt, id).Scan(&workspace.ID, &workspace.Name, &workspace.Owner)
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (r *workspaceRepository) Create(ctx context.Context, workspace *Workspace) error {
	stmt := "INSERT INTO workspaces (id, name, owner) VALUES ($1, $2, $3)"
	_, err := r.db.Connection().Exec(ctx, stmt, workspace.ID, workspace.Name, workspace.Owner)
	if err != nil {
		return err
	}

	return nil
}
