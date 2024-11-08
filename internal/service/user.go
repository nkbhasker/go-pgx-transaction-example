package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/model"
)

type UserService interface {
	Get(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(ctx context.Context, teamId uuid.UUID, user *model.User) (*model.User, error)
	create(ctx context.Context, user *model.User) error
}

type userService struct {
	*baseService
	teamService TeamService
}

func NewUserService(baseService *baseService, teamService TeamService) UserService {
	return &userService{
		baseService,
		teamService,
	}
}

func (s *userService) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.repo.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *userService) Create(ctx context.Context, temaId uuid.UUID, user *model.User) (*model.User, error) {
	user.ID = uuid.New()
	err := s.Transact(ctx, func(ctx context.Context, service *Service) error {
		err := s.create(ctx, user)
		if err != nil {
			return err
		}
		teamMember := &model.TeamMember{
			TeamID: temaId,
			UserID: user.ID,
		}

		return s.teamService.AddTeamMember(ctx, teamMember)
	})

	return user, err
}

func (s *userService) create(ctx context.Context, user *model.User) error {
	return s.repo.User.Create(ctx, user)
}
