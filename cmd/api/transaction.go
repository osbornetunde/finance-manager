package main

import (
	"net/http"
)

func (a *API) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := a.service.GetTransactions(ctx)
	if err != nil {
		a.logger.Error("Failed to fetch transactions", "error", err)
		a.httpError(w, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	a.jsonResponse(w, http.StatusOK, res)
}
