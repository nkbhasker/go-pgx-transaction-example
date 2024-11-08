package model

import (
	"context"
	"log"

	"github.com/nkbhasker/go-pgx-transaction-example/internal/db"
)

// Validate that Repository implements Transactioner
var _ Transactioner[*Repository] = (*Repository)(nil)

type Transactioner[T interface{}] interface {
	Transact(context.Context, func(context.Context, T) error) error
}

type baseRepository struct {
	db db.DB
}

type Repository struct {
	*baseRepository
	Team      TeamRepository
	User      UserRepository
	Workspace WorkspaceRepository
}

func NewBaseRepository(db db.DB) *baseRepository {
	return &baseRepository{
		db: db,
	}
}

func NewRepository(baseRepository *baseRepository) *Repository {
	return &Repository{
		baseRepository: baseRepository,
		Team:           NewTeamRepository(baseRepository),
		User:           NewUserRepository(baseRepository),
		Workspace:      NewWorkspaceRepository(baseRepository),
	}
}

func (r *baseRepository) Transact(ctx context.Context, fn func(ctx context.Context, repo *Repository) error) error {
	tx, db, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	baseRepository := NewBaseRepository(db)
	err = fn(ctx, NewRepository(baseRepository))
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			log.Printf("Error rolling back transaction: %v", rbErr)
		}

		return err
	}

	return tx.Commit(ctx)
}
