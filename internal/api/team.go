package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/model"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/service"
)

type teamAPI struct {
	service *service.Service
}

type teamCreateRequestPayload struct {
	Name string `json:"name"`
}

func NewTeamAPI(service *service.Service) *teamAPI {
	return &teamAPI{
		service: service,
	}
}

func (a *teamAPI) Route() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", a.create)
	}
}

func (a *teamAPI) create(w http.ResponseWriter, r *http.Request) {
	var payload teamCreateRequestPayload
	team, err := func() (*model.Team, error) {
		workspaceId, err := parseWorkspaceId(r)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return nil, err
		}
		return a.service.Team.Create(
			r.Context(),
			&model.Team{
				Name:        payload.Name,
				WorkspaceID: *workspaceId,
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
		"success": true,
		"team":    team,
	})
}
