package model

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
}

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type userRepository struct {
	*baseRepository
}

func NewUserRepository(baseRepository *baseRepository) UserRepository {
	return &userRepository{
		baseRepository: baseRepository,
	}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	stmt := `INSERT INTO users 
		(id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.Connection().Exec(ctx, stmt, user.ID, user.FirstName, user.LastName, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	user := &User{}
	err := r.db.Connection().QueryRow(ctx, "SELECT id, first_name, last_name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
