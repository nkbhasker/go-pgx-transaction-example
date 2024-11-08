package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/service"
)

type API interface {
	Handler() http.Handler
}

type api struct {
	router chi.Router
}

func NewAPI(service *service.Service) API {
	router := chi.NewRouter()
	teamAPI := NewTeamAPI(service)
	userAPI := NewUserAPI(service)
	workspaceAPI := NewWorkspaceAPI(service)
	// Register routes
	router.Route("/api/workspaces", workspaceAPI.Routes())
	router.Route("/api/workspaces/{workspaceId}/teams", teamAPI.Route())
	router.Route("/api/workspaces/{workspaceId}/users", userAPI.Routes())

	return &api{
		router: router,
	}
}

func (a *api) Handler() http.Handler {
	return a.router
}

func parseWorkspaceId(r *http.Request) (*uuid.UUID, error) {
	workspaceId := chi.URLParam(r, "workspaceId")
	if workspaceId == "" {
		return nil, errors.New("workspace id is required")
	}
	id, err := uuid.Parse(workspaceId)
	if err != nil {
		return nil, errors.New("invalid workspace id")
	}

	return &id, nil
}
