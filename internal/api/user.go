package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/model"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/service"
)

type userAPI struct {
	service *service.Service
}

type userCreateRequestPayload struct {
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	TeamID    *uuid.UUID `json:"teamId"`
}

func NewUserAPI(service *service.Service) *userAPI {
	return &userAPI{
		service: service,
	}
}

func (u *userAPI) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/users", u.createUser)
	}
}

func (u *userAPI) createUser(w http.ResponseWriter, r *http.Request) {
	var payload userCreateRequestPayload
	user, err := func() (*model.User, error) {
		_, err := parseWorkspaceId(r)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return nil, err
		}
		if payload.TeamID == nil {
			return nil, errors.New("team id is required")
		}

		return u.service.User.Create(r.Context(), *payload.TeamID, &model.User{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
		})
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	render.JSON(w, r, map[string]interface{}{"success": true, "data": user})
}
