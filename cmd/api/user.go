package main

import (
	"encoding/json"
	"finance-manager/internal/service"
	"log/slog"
	"net/http"
)

type API struct {
	service service.Service
	logger  *slog.Logger
}

func NewAPI(srv service.Service, logger *slog.Logger) *API {
	return &API{
		service: srv,
		logger:  logger,
	}
}

func (a *API) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := a.service.GetUsers(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		a.logger.Error("Failed to encode users response", "error", err)
	}
}
