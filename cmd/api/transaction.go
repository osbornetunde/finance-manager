package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	appErrors "finance-manager/internal/errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CreateTransactionRequest struct {
	UserID      int64           `json:"user_id"`
	Amount      int64           `json:"amount"`
	CategoryID  int64           `json:"category_id"`
	Description string          `json:"description"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
}

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

func (a *API) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var body CreateTransactionRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		a.httpError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		a.httpError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	requestKey := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
	if requestKey == "" {
		requestKey = "fp:" + buildTransactionFingerprint(body)
	} else {
		requestKey = "idem:" + requestKey
	}
	if a.idem.seenRecently(requestKey, 60*time.Second) {
		a.logger.Debug("Duplicate transaction request blocked", "key", requestKey)
		a.httpError(w, http.StatusConflict, "Duplicate transaction request")
		return
	}

	res, err := a.service.CreateTransaction(ctx, body.UserID, body.Amount, body.CategoryID, body.Description, body.Metadata, body.Tags)
	if err != nil {
		if appErrors.IsValidationError(err) {
			a.logger.Debug("Validation failed", "error", err)
			a.httpError(w, http.StatusBadRequest, err.Error())
			return
		}
		if appErrors.IsInvalidUserReference(err) || appErrors.IsInvalidCategoryReference(err) {
			a.logger.Debug("Invalid transaction reference", "error", err)
			a.httpError(w, http.StatusBadRequest, err.Error())
			return
		}
		a.logger.Error("Failed to create transaction", "error", err)
		a.httpError(w, http.StatusInternalServerError, "Failed to create transaction")
		return
	}
	a.jsonResponse(w, http.StatusCreated, res)
}

func buildTransactionFingerprint(body CreateTransactionRequest) string {
	metadata := strings.TrimSpace(string(body.Metadata))
	if len(body.Metadata) > 0 {
		var compact bytes.Buffer
		if err := json.Compact(&compact, body.Metadata); err == nil {
			metadata = compact.String()
		}
	}
	payload := fmt.Sprintf("%d|%d|%d|%s|%s|%s",
		body.UserID,
		body.Amount,
		body.CategoryID,
		strings.TrimSpace(body.Description),
		metadata,
		strings.Join(body.Tags, ","),
	)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
