package main

import (
	"encoding/json"
	"net/http"
)

func (a *API) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := a.service.GetTransactions(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		a.logger.Error("Failed to encode transactions response", "error", err)
	}
}
