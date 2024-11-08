package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/model"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/service"
)

type workspaceAPI struct {
	service *service.Service
}

type workspaceCreateRequestPayload struct {
	Name  string                   `json:"name"`
	Owner userCreateRequestPayload `json:"owner"`
}

func NewWorkspaceAPI(service *service.Service) workspaceAPI {
	return workspaceAPI{
		service: service,
	}
}

func (a *workspaceAPI) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", a.CreateWorkspace)
	}
}

func (a *workspaceAPI) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var payload workspaceCreateRequestPayload
	workspace, err := func() (*model.Workspace, error) {
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return nil, err
		}
		return a.service.Workspace.Create(
			r.Context(),
			&model.Workspace{
				Name: payload.Name,
			},
			&model.User{
				FirstName: payload.Owner.FirstName,
				LastName:  payload.Owner.LastName,
				Email:     payload.Owner.Email,
			},
		)
	}()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"success":   true,
		"workspace": workspace,
	})
}
