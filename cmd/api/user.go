package main

import (
	"context"
	"encoding/json"
	appErrors "finance-manager/internal/errors"
	"finance-manager/internal/service"
	"log/slog"
	"net/http"
	"time"
)

type API struct {
	service service.Service
	logger  *slog.Logger
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewAPI(srv service.Service, logger *slog.Logger) *API {
	return &API{
		service: srv,
		logger:  logger,
	}
}

func (a *API) jsonResponse(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		a.logger.Error("Failed to encode JSON response", "error", err)
	}
}

func (a *API) httpError(w http.ResponseWriter, status int, message string) {
	a.jsonResponse(w, status, map[string]string{"error": message})
}

func (a *API) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := a.service.GetUsers(ctx)
	if err != nil {
		a.logger.Error("Failed to fetch users", "error", err)
		a.httpError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	a.jsonResponse(w, http.StatusOK, res)
}

func (a *API) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var body CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		a.httpError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	res, err := a.service.CreateUser(ctx, body.Name, body.Email)
	if err != nil {
		a.logger.Error("Failed to create user", "error", err)

		if appErrors.IsDuplicateEmail(err) {
			a.httpError(w, http.StatusConflict, "Email already exists")
			return
		}
		if appErrors.IsValidationError(err) {
			a.httpError(w, http.StatusBadRequest, err.Error())
			return
		}

		a.httpError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	a.jsonResponse(w, http.StatusCreated, res)
}
