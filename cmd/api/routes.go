package main

import (
	"net/http"
)

func (a *API) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/users", a.GetUsersHandler)
	mux.HandleFunc("GET /api/v1/transactions", a.GetTransactionsHandler)

	return mux
}
